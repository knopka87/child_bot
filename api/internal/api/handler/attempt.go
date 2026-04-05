package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"child-bot/api/internal/api/middleware"
	"child-bot/api/internal/api/response"
	"child-bot/api/internal/api/validation"
	"child-bot/api/internal/domain"
	"child-bot/api/internal/llm/types"
	"child-bot/api/internal/service"
)

// AttemptServiceInterface определяет интерфейс для AttemptService
type AttemptServiceInterface interface {
	CreateAttempt(ctx context.Context, childProfileID, attemptType string) (string, error)
	UploadImage(ctx context.Context, attemptID, imageType, imageData string) (string, error)
	ProcessHelp(ctx context.Context, attemptID, imageBase64 string) error
	ProcessCheck(ctx context.Context, attemptID, childProfileID, taskImageBase64, answerImageBase64 string) error
	GetAttemptResult(ctx context.Context, attemptID string) (*service.AttemptData, error)
	GetNextHint(ctx context.Context, attemptID string) (*domain.HelpResult, error)
	DeleteAttempt(ctx context.Context, attemptID string) error
	GetUnfinishedAttempt(ctx context.Context, childProfileID string) (*service.AttemptData, error)
	GetRecentAttempts(ctx context.Context, childProfileID string, limit int) ([]service.AttemptData, error)
}

// AttemptHandler обрабатывает запросы, связанные с попытками
type AttemptHandler struct {
	service        AttemptServiceInterface
	profileService *service.ProfileService
}

// NewAttemptHandler создает новый AttemptHandler
func NewAttemptHandler(service AttemptServiceInterface) *AttemptHandler {
	return &AttemptHandler{
		service: service,
	}
}

// NewAttemptHandlerWithService создает AttemptHandler с конкретным service типом
func NewAttemptHandlerWithService(attemptSvc *service.AttemptService, profileSvc *service.ProfileService) *AttemptHandler {
	return &AttemptHandler{
		service:        attemptSvc,
		profileService: profileSvc,
	}
}

// Create создает новую попытку
// POST /attempts
func (h *AttemptHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAttemptRequest
	if err := validation.DecodeJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Валидация
	if err := validation.ValidateAttemptType(req.Type); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	if err := validation.ValidateUUID(req.ChildProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	// Создание attempt через service layer
	attemptID, err := h.service.CreateAttempt(r.Context(), req.ChildProfileID, req.Type)
	if err != nil {
		response.InternalError(w, "Failed to create attempt")
		return
	}

	// Активируем реферала если это его первая попытка
	if h.profileService != nil {
		log.Printf("[AttemptHandler] Attempting to activate referral for child: %s", req.ChildProfileID)
		if err := h.profileService.ActivateReferral(r.Context(), req.ChildProfileID); err != nil {
			// Логируем ошибку, но не блокируем создание попытки
			log.Printf("[AttemptHandler] Failed to activate referral for child %s: %v", req.ChildProfileID, err)
		} else {
			log.Printf("[AttemptHandler] Successfully activated referral for child: %s", req.ChildProfileID)
		}
	} else {
		log.Printf("[AttemptHandler] ProfileService is nil, cannot activate referral")
	}

	response.Created(w, CreateAttemptResponse{
		AttemptID: attemptID,
		Status:    "created",
	})
}

// UploadImage загружает изображение для попытки
// POST /attempts/{id}/images
func (h *AttemptHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	attemptID := r.PathValue("id")
	if err := validation.ValidateUUID(attemptID); err != nil {
		response.BadRequest(w, "invalid attempt_id: "+err.Error())
		return
	}

	var req UploadImageRequest
	if err := validation.DecodeJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Валидация
	if req.ImageType != "task" && req.ImageType != "answer" {
		response.BadRequest(w, "image_type must be 'task' or 'answer'")
		return
	}

	if err := validation.ValidateBase64Image(req.ImageData); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// Загрузка изображения через service layer
	imageURL, err := h.service.UploadImage(r.Context(), attemptID, req.ImageType, req.ImageData)
	if err != nil {
		response.InternalError(w, "Failed to upload image")
		return
	}

	response.OK(w, UploadImageResponse{
		ImageURL: imageURL,
		Message:  "Image uploaded successfully",
	})
}

// Process начинает обработку попытки через LLM
// POST /attempts/{id}/process
func (h *AttemptHandler) Process(w http.ResponseWriter, r *http.Request) {
	attemptID := r.PathValue("id")
	if err := validation.ValidateUUID(attemptID); err != nil {
		response.BadRequest(w, "invalid attempt_id: "+err.Error())
		return
	}

	// Получаем child_profile_id из context (добавлен middleware.Auth)
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// Получаем данные попытки
	attemptData, err := h.service.GetAttemptResult(r.Context(), attemptID)
	if err != nil {
		log.Printf("[AttemptHandler] Failed to get attempt: %v", err)
		response.InternalError(w, "Failed to get attempt")
		return
	}

	// Проверяем что попытка принадлежит пользователю
	if attemptData.ChildProfileID != childProfileID {
		response.Forbidden(w, "Attempt belongs to another user")
		return
	}

	// Проверяем что есть изображение
	if attemptData.TaskImageData == "" {
		response.BadRequest(w, "No image uploaded")
		return
	}

	// Обрабатываем в зависимости от типа
	if attemptData.Type == "help" {
		// Запускаем обработку в goroutine (не блокируем запрос)
		go func() {
			ctx := context.Background() // новый context для фоновой задачи
			err := h.service.ProcessHelp(ctx, attemptID, attemptData.TaskImageData)
			if err != nil {
				log.Printf("[AttemptHandler] ProcessHelp failed for attempt %s: %v", attemptID, err)
			}
		}()
	} else if attemptData.Type == "check" {
		if attemptData.AnswerImageData == "" {
			response.BadRequest(w, "No answer image uploaded")
			return
		}

		// Запускаем обработку в goroutine
		go func() {
			ctx := context.Background()
			err := h.service.ProcessCheck(ctx, attemptID, childProfileID,
				attemptData.TaskImageData, attemptData.AnswerImageData)
			if err != nil {
				log.Printf("[AttemptHandler] ProcessCheck failed for attempt %s: %v", attemptID, err)
			}
		}()
	}

	response.OK(w, ProcessAttemptResponse{
		Status:  "processing",
		Message: "Attempt is being processed",
	})
}

// GetResult получает результат попытки
// GET /attempts/{id}/result
func (h *AttemptHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	attemptID := r.PathValue("id")
	if err := validation.ValidateUUID(attemptID); err != nil {
		response.BadRequest(w, "invalid attempt_id: "+err.Error())
		return
	}

	// Получаем результат через service layer
	attemptData, err := h.service.GetAttemptResult(r.Context(), attemptID)
	if err != nil {
		log.Printf("[AttemptHandler] Failed to get attempt result: %v", err)
		response.InternalError(w, "Failed to get attempt result")
		return
	}

	// Формируем ответ
	resultData := make(map[string]interface{})

	// Для help - добавляем hints
	if attemptData.Type == "help" && attemptData.HintsResult != nil {
		hintsArray := make([]map[string]interface{}, 0)
		hintOrder := 1
		for _, item := range attemptData.HintsResult.Items {
			// Каждая подсказка из item - отдельный элемент
			for _, hint := range item.Hints {
				level := 1
				switch hint.Level {
				case types.HintL1:
					level = 1
				case types.HintL2:
					level = 2
				case types.HintL3:
					level = 3
				}

				hintsArray = append(hintsArray, map[string]interface{}{
					"id":      fmt.Sprintf("%d", hintOrder),
					"level":   level,
					"title":   fmt.Sprintf("Подсказка %d уровня", level),
					"content": hint.HintText,
					"order":   hintOrder,
				})
				hintOrder++
			}
		}
		resultData["hints"] = hintsArray
		log.Printf("[AttemptHandler] Formatted %d hints for attempt %s", len(hintsArray), attemptID)
	}

	// Для check - добавляем ошибки
	if attemptData.Type == "check" && attemptData.CheckResult != nil {
		resultData["is_correct"] = attemptData.CheckResult.Decision == types.CheckDecisionCorrect
		if attemptData.CheckResult.Feedback != "" {
			resultData["feedback"] = attemptData.CheckResult.Feedback
		}
	}

	response.OK(w, GetResultResponse{
		AttemptID: attemptID,
		Type:      attemptData.Type,
		Status:    attemptData.Status,
		Result:    resultData,
		CreatedAt: attemptData.CreatedAt,
		UpdatedAt: attemptData.UpdatedAt,
	})
}

// NextHint получает следующую подсказку
// POST /attempts/{id}/next-hint
func (h *AttemptHandler) NextHint(w http.ResponseWriter, r *http.Request) {
	attemptID := r.PathValue("id")
	if err := validation.ValidateUUID(attemptID); err != nil {
		response.BadRequest(w, "invalid attempt_id: "+err.Error())
		return
	}

	hint, err := h.service.GetNextHint(r.Context(), attemptID)
	if err != nil {
		log.Printf("[AttemptHandler] Failed to get next hint: %v", err)
		response.InternalError(w, "Failed to get next hint")
		return
	}

	response.OK(w, NextHintResponse{
		Hint:        hint.Hints[hint.CurrentHint],
		HintIndex:   hint.CurrentHint,
		TotalHints:  hint.TotalHints,
		HasMoreHint: hint.CurrentHint < hint.TotalHints-1,
	})
}

// Delete удаляет попытку
// DELETE /attempts/{id}
func (h *AttemptHandler) Delete(w http.ResponseWriter, r *http.Request) {
	attemptID := r.PathValue("id")
	if err := validation.ValidateUUID(attemptID); err != nil {
		response.BadRequest(w, "invalid attempt_id: "+err.Error())
		return
	}

	err := h.service.DeleteAttempt(r.Context(), attemptID)
	if err != nil {
		log.Printf("[AttemptHandler] Failed to delete attempt: %v", err)
		response.InternalError(w, "Failed to delete attempt")
		return
	}

	response.NoContent(w)
}

// GetUnfinished получает незавершенную попытку
// GET /attempts/unfinished
func (h *AttemptHandler) GetUnfinished(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.URL.Query().Get("childProfileId")
	if childProfileID == "" {
		childProfileID = middleware.GetChildProfileID(r.Context())
	}

	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	attempt, err := h.service.GetUnfinishedAttempt(r.Context(), childProfileID)
	if err != nil {
		log.Printf("[AttemptHandler] Failed to get unfinished attempt: %v", err)
		response.InternalError(w, "Failed to get unfinished attempt")
		return
	}

	// Возвращаем null если нет незавершенных
	if attempt == nil {
		response.OK(w, nil)
		return
	}

	// Формируем ответ
	response.OK(w, map[string]interface{}{
		"attempt_id": attempt.ID,
		"type":       attempt.Type,
		"status":     attempt.Status,
		"created_at": attempt.CreatedAt,
	})
}

// GetRecent получает последние попытки
// GET /attempts/recent
func (h *AttemptHandler) GetRecent(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.URL.Query().Get("childProfileId")
	if childProfileID == "" {
		childProfileID = middleware.GetChildProfileID(r.Context())
	}

	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	// Parse limit parameter
	limit := 3
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && parsedLimit == 1 {
			if limit < 1 || limit > 100 {
				limit = 3
			}
		}
	}

	attempts, err := h.service.GetRecentAttempts(r.Context(), childProfileID, limit)
	if err != nil {
		log.Printf("[AttemptHandler] Failed to get recent attempts: %v", err)
		response.InternalError(w, "Failed to get recent attempts")
		return
	}

	// Формируем массив ответов
	recent := make([]map[string]interface{}, 0, len(attempts))
	for _, attempt := range attempts {
		recent = append(recent, map[string]interface{}{
			"attempt_id": attempt.ID,
			"type":       attempt.Type,
			"status":     attempt.Status,
			"created_at": attempt.CreatedAt,
		})
	}

	response.OK(w, recent)
}
