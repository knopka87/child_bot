package handler

import (
	"log"
	"net/http"

	"child-bot/api/internal/api/middleware"
	"child-bot/api/internal/api/response"
	"child-bot/api/internal/api/validation"
	"child-bot/api/internal/service"
)

// ProfileHandler обрабатывает запросы профиля
type ProfileHandler struct {
	service *service.ProfileService
}

// NewProfileHandler создает новый ProfileHandler
func NewProfileHandler(service *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

// ProfileData структура данных профиля
type ProfileData struct {
	ID           string           `json:"id"`
	DisplayName  string           `json:"display_name"`
	AvatarID     string           `json:"avatar_id"`
	AvatarURL    string           `json:"avatar_url"`
	Grade        int              `json:"grade"`
	Subscription SubscriptionData `json:"subscription"`
}

type SubscriptionData struct {
	Status             string `json:"status"` // trial, active, expired, cancelled
	PlanID             string `json:"plan_id,omitempty"`
	PlanName           string `json:"plan_name,omitempty"`
	TrialDaysRemaining int    `json:"trial_days_remaining,omitempty"`
	ExpiresAt          string `json:"expires_at,omitempty"`
}

type CreateChildProfileRequest struct {
	ParentUserID string `json:"parentUserId"`
	Grade        int    `json:"grade"`
	AvatarID     string `json:"avatarId"`
	DisplayName  string `json:"displayName"`
	ReferralCode string `json:"referralCode,omitempty"`
}

type CreateChildProfileResponse struct {
	ChildProfileID string `json:"childProfileId"`
}

type UpdateProfileRequest struct {
	DisplayName string `json:"display_name,omitempty"`
	AvatarID    string `json:"avatar_id,omitempty"`
	Grade       int    `json:"grade,omitempty"`
}

type HistoryAttempt struct {
	ID           string         `json:"id"`
	Mode         string         `json:"mode"`                    // help or check
	Status       string         `json:"status"`                  // success, error, in_progress
	ScenarioType string         `json:"scenario_type,omitempty"` // single_photo, two_photo
	CreatedAt    string         `json:"created_at"`
	CompletedAt  string         `json:"completed_at,omitempty"`
	Images       []HistoryImage `json:"images"`
	Result       *HistoryResult `json:"result,omitempty"`
	HintsUsed    int            `json:"hints_used,omitempty"`
}

type HistoryImage struct {
	ID           string `json:"id"`
	Role         string `json:"role"` // task, answer, single
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type HistoryResult struct {
	Status     string          `json:"status"` // correct, has_errors, processing
	ErrorCount int             `json:"error_count,omitempty"`
	Feedback   []ErrorFeedback `json:"feedback,omitempty"`
	Summary    string          `json:"summary,omitempty"`
}

type ErrorFeedback struct {
	ID            string `json:"id"`
	StepNumber    int    `json:"step_number,omitempty"`
	LineReference string `json:"line_reference,omitempty"`
	Description   string `json:"description"`
	LocationType  string `json:"location_type"` // line, step, general
}

type ProfileStats struct {
	TotalAttempts      int     `json:"total_attempts"`
	SuccessfulAttempts int     `json:"successful_attempts"`
	ErrorsFixed        int     `json:"errors_fixed"`
	StreakDays         int     `json:"streak_days"`
	AverageAccuracy    float64 `json:"average_accuracy"`
	TotalHintsUsed     int     `json:"total_hints_used"`
}

// CreateChild создает профиль ребенка
// POST /profiles/child
func (h *ProfileHandler) CreateChild(w http.ResponseWriter, r *http.Request) {
	// Получаем platformID из middleware
	platformID := middleware.GetPlatformID(r.Context())
	if platformID == "" {
		response.Unauthorized(w, "Missing platform ID")
		return
	}

	var req CreateChildProfileRequest
	if err := validation.DecodeJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Валидация
	if req.ParentUserID == "" {
		response.BadRequest(w, "parentUserId is required")
		return
	}
	if req.Grade < 1 || req.Grade > 4 {
		response.BadRequest(w, "grade must be between 1 and 4")
		return
	}
	if req.AvatarID == "" {
		response.BadRequest(w, "avatarId is required")
		return
	}
	if req.DisplayName == "" {
		response.BadRequest(w, "displayName is required")
		return
	}
	if err := validation.ValidateMaxLength(req.DisplayName, "displayName", 50); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Создание профиля через service layer
	childProfileID, err := h.service.CreateChildProfile(r.Context(), req.ParentUserID, req.DisplayName, req.AvatarID, platformID, req.Grade)
	if err != nil {
		response.InternalError(w, "Failed to create child profile")
		return
	}

	// Если указан реферальный код, создаём реферальную связь
	if req.ReferralCode != "" {
		if err := h.service.ProcessReferral(r.Context(), childProfileID, req.ReferralCode); err != nil {
			// Логируем ошибку, но не блокируем создание профиля
			log.Printf("Failed to process referral code %s for child %s: %v", req.ReferralCode, childProfileID, err)
		}
	}

	response.Created(w, CreateChildProfileResponse{
		ChildProfileID: childProfileID,
	})
}

// Get получает профиль пользователя
// GET /profile
func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// Получение профиля через service layer
	serviceProfile, err := h.service.GetProfile(r.Context(), childProfileID)
	if err != nil {
		response.InternalError(w, "Failed to get profile")
		return
	}

	// Преобразуем в response format
	profile := ProfileData{
		ID:          serviceProfile.ID,
		DisplayName: serviceProfile.DisplayName,
		AvatarID:    serviceProfile.AvatarID,
		AvatarURL:   serviceProfile.AvatarURL,
		Grade:       serviceProfile.Grade,
		Subscription: SubscriptionData{
			Status:             serviceProfile.Subscription.Status,
			PlanID:             serviceProfile.Subscription.PlanID,
			PlanName:           serviceProfile.Subscription.PlanName,
			TrialDaysRemaining: serviceProfile.Subscription.TrialDaysRemaining,
		},
	}

	response.OK(w, profile)
}

// Update обновляет профиль пользователя
// PUT /profile
func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	var req UpdateProfileRequest
	if err := validation.DecodeJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Валидация
	if req.DisplayName != "" {
		if err := validation.ValidateMaxLength(req.DisplayName, "display_name", 50); err != nil {
			response.BadRequest(w, err.Error())
			return
		}
	}

	if req.Grade < 1 || req.Grade > 4 {
		response.BadRequest(w, "grade must be between 1 and 4")
		return
	}

	// TODO: Phase 4 - обновление через service layer
	// err := h.service.UpdateProfile(r.Context(), childProfileID, req)

	response.OK(w, map[string]string{"message": "Profile updated successfully"})
}

// GetHistory получает историю попыток
// GET /profile/history
func (h *ProfileHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// Query parameters для фильтрации (пока не используются)
	filters := map[string]string{
		"mode":      r.URL.Query().Get("mode"),
		"status":    r.URL.Query().Get("status"),
		"date_from": r.URL.Query().Get("date_from"),
		"date_to":   r.URL.Query().Get("date_to"),
	}

	// Получение истории через service layer
	serviceHistory, err := h.service.GetHistory(r.Context(), childProfileID, filters)
	if err != nil {
		log.Printf("GetHistory error: %v", err)
		response.InternalError(w, "Failed to get history")
		return
	}

	// Преобразуем в response format
	history := make([]HistoryAttempt, 0, len(serviceHistory))
	for _, sh := range serviceHistory {
		ha := HistoryAttempt{
			ID:           sh.ID,
			Mode:         sh.Mode,
			Status:       sh.Status,
			ScenarioType: sh.ScenarioType,
			CreatedAt:    sh.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			HintsUsed:    sh.HintsUsed,
			Images:       make([]HistoryImage, 0, len(sh.Images)),
		}

		if sh.CompletedAt != nil {
			completedAt := sh.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
			ha.CompletedAt = completedAt
		}

		for _, img := range sh.Images {
			ha.Images = append(ha.Images, HistoryImage{
				ID:           img.ID,
				Role:         img.Role,
				URL:          img.URL,
				ThumbnailURL: img.ThumbnailURL,
			})
		}

		history = append(history, ha)
	}

	response.OK(w, history)
}

// GetStats получает статистику профиля
// GET /profile/stats
func (h *ProfileHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// TODO: Phase 4 - получение статистики через service layer
	// stats, err := h.service.GetProfileStats(r.Context(), childProfileID)

	// Placeholder
	stats := ProfileStats{
		TotalAttempts:      50,
		SuccessfulAttempts: 42,
		ErrorsFixed:        35,
		StreakDays:         7,
		AverageAccuracy:    84.0,
		TotalHintsUsed:     15,
	}

	response.OK(w, stats)
}

// GetByPlatform получает профиль по platform_id и platform_user_id
// GET /profiles/by-platform?platform_id=vk&platform_user_id=12345
func (h *ProfileHandler) GetByPlatform(w http.ResponseWriter, r *http.Request) {
	platformID := r.URL.Query().Get("platform_id")
	platformUserID := r.URL.Query().Get("platform_user_id")

	if platformID == "" {
		response.BadRequest(w, "platform_id is required")
		return
	}
	if platformUserID == "" {
		response.BadRequest(w, "platform_user_id is required")
		return
	}

	// Валидация platform_id
	validPlatforms := map[string]bool{
		"vk":       true,
		"telegram": true,
		"max":      true,
		"web":      true,
	}
	if !validPlatforms[platformID] {
		response.BadRequest(w, "invalid platform_id, must be one of: vk, telegram, max, web")
		return
	}

	childProfileID, err := h.service.GetProfileByPlatform(r.Context(), platformID, platformUserID)
	if err != nil {
		log.Printf("GetByPlatform error: %v", err)
		response.NotFound(w, "Profile not found")
		return
	}

	response.OK(w, map[string]string{
		"child_profile_id": childProfileID,
	})
}
