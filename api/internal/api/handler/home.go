package handler

import (
	"net/http"

	"child-bot/api/internal/api/response"
	"child-bot/api/internal/api/validation"
	"child-bot/api/internal/service"
)

// HomeHandler обрабатывает запросы главного экрана
type HomeHandler struct {
	service *service.HomeService
}

// NewHomeHandler создает новый HomeHandler
func NewHomeHandler(service *service.HomeService) *HomeHandler {
	return &HomeHandler{service: service}
}

// HomeData структура данных для главного экрана
type HomeData struct {
	Profile struct {
		ID                      string `json:"id"`
		DisplayName             string `json:"displayName"`
		Level                   int    `json:"level"`
		LevelProgress           int    `json:"levelProgress"` // 0-100
		CoinsBalance            int    `json:"coinsBalance"`
		TasksSolvedCorrectCount int    `json:"tasksSolvedCorrectCount"`
	} `json:"profile"`
	Mascot struct {
		ID       string `json:"id"`
		State    string `json:"state"` // idle, happy, thinking, celebrating
		ImageURL string `json:"imageUrl"`
		Message  string `json:"message"`
	} `json:"mascot"`
	Villain           *VillainInfo    `json:"villain"`
	UnfinishedAttempt *AttemptInfo    `json:"unfinishedAttempt"`
	RecentAttempts    []RecentAttempt `json:"recentAttempts"`
	Achievements      struct {
		UnlockedCount int `json:"unlockedCount"`
		TotalCount    int `json:"totalCount"`
	} `json:"achievements"`
}

type VillainInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ImageURL      string `json:"imageUrl"`
	HealthPercent int    `json:"healthPercent"` // 0-100
	IsDefeated    bool   `json:"isDefeated"`
}

type AttemptInfo struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type RecentAttempt struct {
	ID            string `json:"id"`
	Mode          string `json:"mode"`   // help or check
	Status        string `json:"status"` // success, error, in_progress
	CreatedAt     string `json:"createdAt"`
	Thumbnail     string `json:"thumbnail,omitempty"`
	ResultSummary string `json:"resultSummary,omitempty"`
}

// GetHomeData получает все данные для главного экрана
// GET /home/{childProfileId}
func (h *HomeHandler) GetHomeData(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.PathValue("childProfileId")
	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	// Получение данных через service layer
	serviceData, err := h.service.GetHomeData(r.Context(), childProfileID)
	if err != nil {
		response.InternalError(w, "Failed to get home data")
		return
	}

	// Преобразуем в response format
	data := HomeData{}
	data.Profile.ID = serviceData.Profile.ID
	data.Profile.DisplayName = serviceData.Profile.DisplayName
	data.Profile.Level = serviceData.Profile.Level
	data.Profile.LevelProgress = serviceData.Profile.LevelProgress
	data.Profile.CoinsBalance = serviceData.Profile.CoinsBalance
	data.Profile.TasksSolvedCorrectCount = serviceData.Profile.TasksSolvedCorrectCount

	data.Mascot.ID = serviceData.Mascot.ID
	data.Mascot.State = serviceData.Mascot.State
	data.Mascot.ImageURL = serviceData.Mascot.ImageURL
	data.Mascot.Message = serviceData.Mascot.Message

	// Преобразуем данные villain
	if serviceData.Villain != nil {
		healthPercent := 0
		if serviceData.Villain.MaxHP > 0 {
			healthPercent = (serviceData.Villain.HP * 100) / serviceData.Villain.MaxHP
		}
		data.Villain = &VillainInfo{
			ID:            serviceData.Villain.ID,
			Name:          serviceData.Villain.Name,
			ImageURL:      serviceData.Villain.ImageURL,
			HealthPercent: healthPercent,
			IsDefeated:    serviceData.Villain.IsDefeated,
		}
	}

	// Преобразуем unfinished attempt
	if serviceData.UnfinishedAttempt != nil {
		data.UnfinishedAttempt = &AttemptInfo{
			ID:        serviceData.UnfinishedAttempt.ID,
			Type:      serviceData.UnfinishedAttempt.Type,
			Status:    serviceData.UnfinishedAttempt.Status,
			CreatedAt: serviceData.UnfinishedAttempt.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	data.RecentAttempts = []RecentAttempt{}

	data.Achievements.UnlockedCount = serviceData.Achievements.UnlockedCount
	data.Achievements.TotalCount = serviceData.Achievements.TotalCount

	response.OK(w, data)
}
