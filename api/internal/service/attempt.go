package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"child-bot/api/internal/domain"
	"child-bot/api/internal/llm"
	"child-bot/api/internal/llm/types"
	"child-bot/api/internal/store"

	"github.com/google/uuid"
)

// AttemptService бизнес-логика для работы с попытками
type AttemptService struct {
	store          *store.Store
	llmClient      *llm.Client
	defaultLLM     string
	profileService *ProfileService
}

// NewAttemptService создает новый AttemptService
func NewAttemptService(store *store.Store, llmClient *llm.Client, defaultLLM string) *AttemptService {
	return &AttemptService{
		store:      store,
		llmClient:  llmClient,
		defaultLLM: defaultLLM,
	}
}

// SetProfileService устанавливает ProfileService (для избежания циклических зависимостей)
func (s *AttemptService) SetProfileService(profileService *ProfileService) {
	s.profileService = profileService
}

// AttemptData внутренняя структура для хранения данных попытки
type AttemptData struct {
	ID             string
	ChildProfileID string
	Type           string // help or check
	Status         string // created, processing, completed, failed
	TaskImageData  string // base64
	AnswerImageData string // base64
	ParseResult    *types.ParseResponse
	DetectResult   *types.DetectResponse
	HintsResult    *types.HintResponse
	CheckResult    *types.CheckResponse
	CurrentHint    int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// CreateAttempt создает новую попытку
func (s *AttemptService) CreateAttempt(ctx context.Context, childProfileID, attemptType string) (string, error) {
	// Валидация типа
	if attemptType != "help" && attemptType != "check" {
		return "", domain.ErrInvalidInput
	}

	// Парсим UUID
	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		return "", fmt.Errorf("invalid child_profile_id: %w", err)
	}

	// Создаём в БД
	attemptID, err := s.store.Attempts.CreateAttempt(ctx, profileUUID, attemptType)
	if err != nil {
		return "", fmt.Errorf("failed to create attempt: %w", err)
	}

	log.Printf("[AttemptService] Created attempt: id=%s, child_profile_id=%s, type=%s",
		attemptID, childProfileID, attemptType)

	return attemptID.String(), nil
}

// UploadImage загружает изображение для попытки
func (s *AttemptService) UploadImage(ctx context.Context, attemptID, imageType, imageData string) (string, error) {
	// Валидация типа изображения
	if imageType != "task" && imageType != "answer" {
		return "", domain.ErrInvalidInput
	}

	// Парсим UUID
	id, err := uuid.Parse(attemptID)
	if err != nil {
		return "", fmt.Errorf("invalid attempt_id: %w", err)
	}

	// Сохраняем data URI в БД (в колонке task_image_url/answer_image_url)
	// TODO: Phase 6 - загрузка в S3 и сохранение URL
	if imageType == "task" {
		err = s.store.Attempts.UpdateTaskImage(ctx, id, imageData)
	} else {
		err = s.store.Attempts.UpdateAnswerImage(ctx, id, imageData)
	}

	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	log.Printf("[AttemptService] Uploaded %s image for attempt: %s", imageType, attemptID)

	// Возвращаем data URI как "URL"
	return imageData, nil
}

// ProcessHelp обрабатывает help попытку через LLM
func (s *AttemptService) ProcessHelp(ctx context.Context, attemptID string, imageBase64 string) error {
	// Парсим UUID
	id, err := uuid.Parse(attemptID)
	if err != nil {
		return fmt.Errorf("invalid attempt_id: %w", err)
	}

	// Обновляем статус на processing
	err = s.store.Attempts.UpdateStatus(ctx, id, "processing")
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	log.Printf("[AttemptService] Processing help attempt: %s", attemptID)

	// 1. Detect - определить предмет и качество
	detectReq := types.DetectRequest{
		Image:  imageBase64,
		Locale: "ru-RU",
	}

	detectResp, err := s.llmClient.Detect(ctx, s.defaultLLM, detectReq)
	if err != nil {
		// Обновляем статус на failed
		_ = s.store.Attempts.UpdateStatus(ctx, id, "failed")
		return fmt.Errorf("detect failed: %w", err)
	}

	// Сохраняем результат Detect
	err = s.store.Attempts.SaveDetectResult(ctx, id, &detectResp)
	if err != nil {
		log.Printf("[AttemptService] Failed to save detect result: %v", err)
	}

	log.Printf("[AttemptService] Detect completed: subject=%s, confidence=%.2f",
		detectResp.Classification.SubjectCandidate, detectResp.Classification.Confidence)

	// 2. Parse - распарсить задачу
	parseReq := types.ParseRequest{
		Image:             imageBase64,
		TaskId:            attemptID,
		Grade:             5, // TODO: получать из профиля пользователя
		SubjectCandidate:  string(detectResp.Classification.SubjectCandidate),
		SubjectConfidence: fmt.Sprintf("%.2f", detectResp.Classification.Confidence),
		Locale:            "ru-RU",
	}

	parseResp, err := s.llmClient.Parse(ctx, s.defaultLLM, parseReq)
	if err != nil {
		_ = s.store.Attempts.UpdateStatus(ctx, id, "failed")
		return fmt.Errorf("parse failed: %w", err)
	}

	// Сохраняем результат Parse
	err = s.store.Attempts.SaveParseResult(ctx, id, &parseResp)
	if err != nil {
		log.Printf("[AttemptService] Failed to save parse result: %v", err)
	}

	log.Printf("[AttemptService] Parse completed: task_text=%s", parseResp.Task.TaskTextClean)

	// 3. Hint - сгенерировать подсказки
	hintReq := types.HintRequest{
		Task:  parseResp.Task,
		Mode:  "learn",
		Items: parseResp.Items,
		// TODO: правильно заполнить AppliedPolicy и Template
	}

	hintResp, err := s.llmClient.Hint(ctx, s.defaultLLM, hintReq)
	if err != nil {
		_ = s.store.Attempts.UpdateStatus(ctx, id, "failed")
		return fmt.Errorf("hint generation failed: %w", err)
	}

	// Сохраняем результат Hints (и обновляем статус на completed)
	err = s.store.Attempts.SaveHintsResult(ctx, id, &hintResp)
	if err != nil {
		log.Printf("[AttemptService] Failed to save hints result: %v", err)
		_ = s.store.Attempts.UpdateStatus(ctx, id, "failed")
		return fmt.Errorf("failed to save hints result: %w", err)
	}

	log.Printf("[AttemptService] Hints generated successfully: %d items", len(hintResp.Items))

	return nil
}

// ProcessCheck обрабатывает check попытку через LLM
func (s *AttemptService) ProcessCheck(ctx context.Context, attemptID, childProfileID, taskImageBase64, answerImageBase64 string) error {
	// Парсим UUID
	id, err := uuid.Parse(attemptID)
	if err != nil {
		return fmt.Errorf("invalid attempt_id: %w", err)
	}

	// Обновляем статус на processing
	err = s.store.Attempts.UpdateStatus(ctx, id, "processing")
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	log.Printf("[AttemptService] Processing check attempt: %s", attemptID)

	// 1. Detect + Parse задачу
	detectReq := types.DetectRequest{
		Image:  taskImageBase64,
		Locale: "ru-RU",
	}

	detectResp, err := s.llmClient.Detect(ctx, s.defaultLLM, detectReq)
	if err != nil {
		_ = s.store.Attempts.UpdateStatus(ctx, id, "failed")
		return fmt.Errorf("detect failed: %w", err)
	}

	// Сохраняем результат Detect
	err = s.store.Attempts.SaveDetectResult(ctx, id, &detectResp)
	if err != nil {
		log.Printf("[AttemptService] Failed to save detect result: %v", err)
	}

	log.Printf("[AttemptService] Detect completed: subject=%s, confidence=%.2f",
		detectResp.Classification.SubjectCandidate, detectResp.Classification.Confidence)

	parseReq := types.ParseRequest{
		Image:             taskImageBase64,
		TaskId:            attemptID,
		Grade:             5, // TODO: получать из профиля
		SubjectCandidate:  string(detectResp.Classification.SubjectCandidate),
		SubjectConfidence: fmt.Sprintf("%.2f", detectResp.Classification.Confidence),
		Locale:            "ru-RU",
	}

	parseResp, err := s.llmClient.Parse(ctx, s.defaultLLM, parseReq)
	if err != nil {
		_ = s.store.Attempts.UpdateStatus(ctx, id, "failed")
		return fmt.Errorf("parse failed: %w", err)
	}

	// Сохраняем результат Parse
	err = s.store.Attempts.SaveParseResult(ctx, id, &parseResp)
	if err != nil {
		log.Printf("[AttemptService] Failed to save parse result: %v", err)
	}

	log.Printf("[AttemptService] Parse completed: task_text=%s", parseResp.Task.TaskTextClean)

	// 2. CheckSolution - проверить решение
	checkReq := types.CheckRequest{
		Image: answerImageBase64,
		TaskStruct: types.TaskStructCheck{
			TaskTextClean:   parseResp.Task.TaskTextClean,
			VisualReasoning: parseResp.Task.VisualReasoning,
			VisualFacts:     parseResp.Task.VisualFacts,
			QualityFlags:    parseResp.Task.Quality,
			Items:           parseResp.Items,
		},
		RawTaskText: parseResp.Task.TaskTextClean,
		Student: types.StudentCheck{
			Grade:   parseResp.Task.Grade,
			Subject: string(parseResp.Task.Subject),
			Locale:  "ru-RU",
		},
		PhotoQualityHint: "", // TODO: передавать качество фото
	}

	checkResp, err := s.llmClient.CheckSolution(ctx, s.defaultLLM, checkReq)
	if err != nil {
		_ = s.store.Attempts.UpdateStatus(ctx, id, "failed")
		return fmt.Errorf("check solution failed: %w", err)
	}

	// 3. Проверяем результат и начисляем монеты
	if checkResp.Decision == types.CheckDecisionCorrect && s.profileService != nil {
		// Начисляем 5 монет за правильное решение
		err := s.profileService.AddCoins(ctx, childProfileID, 5)
		if err != nil {
			// Логируем ошибку, но не блокируем основной флоу
			log.Printf("[AttemptService] Failed to add coins for child %s: %v", childProfileID, err)
		}
	}

	// 4. Сохраняем результат Check (и обновляем статус на completed)
	err = s.store.Attempts.SaveCheckResult(ctx, id, &checkResp)
	if err != nil {
		log.Printf("[AttemptService] Failed to save check result: %v", err)
		_ = s.store.Attempts.UpdateStatus(ctx, id, "failed")
		return fmt.Errorf("failed to save check result: %w", err)
	}

	log.Printf("[AttemptService] Check completed successfully: decision=%s", checkResp.Decision)

	return nil
}

// GetNextHint получает следующую подсказку
func (s *AttemptService) GetNextHint(ctx context.Context, attemptID string) (*domain.HelpResult, error) {
	// TODO: Phase 5 - загрузка текущего состояния из БД
	// Пока возвращаем placeholder

	// Пример: загружаем сохраненные Parse и Detect результаты
	// и запрашиваем следующий уровень подсказки

	result := &domain.HelpResult{
		Subject:     "Математика",
		TaskText:    "Решить уравнение",
		Hints:       []string{"Первая подсказка", "Вторая подсказка", "Третья подсказка"},
		CurrentHint: 1,
		TotalHints:  3,
	}

	return result, nil
}

// GetAttemptResult получает результат попытки
func (s *AttemptService) GetAttemptResult(ctx context.Context, attemptID string) (*AttemptData, error) {
	id, err := uuid.Parse(attemptID)
	if err != nil {
		return nil, fmt.Errorf("invalid attempt_id: %w", err)
	}

	attempt, err := s.store.Attempts.GetAttempt(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get attempt: %w", err)
	}

	return s.convertToAttemptData(attempt), nil
}

// DeleteAttempt удаляет попытку
func (s *AttemptService) DeleteAttempt(ctx context.Context, attemptID string) error {
	id, err := uuid.Parse(attemptID)
	if err != nil {
		return fmt.Errorf("invalid attempt_id: %w", err)
	}

	err = s.store.Attempts.DeleteAttempt(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete attempt: %w", err)
	}

	log.Printf("[AttemptService] Deleted attempt: %s", attemptID)
	return nil
}

// GetUnfinishedAttempt получает незавершенную попытку для профиля
func (s *AttemptService) GetUnfinishedAttempt(ctx context.Context, childProfileID string) (*AttemptData, error) {
	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		return nil, fmt.Errorf("invalid child_profile_id: %w", err)
	}

	attempt, err := s.store.Attempts.GetUnfinishedAttempt(ctx, profileUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unfinished attempt: %w", err)
	}

	if attempt == nil {
		return nil, nil // Нет незавершенных попыток
	}

	// Конвертируем в AttemptData
	return s.convertToAttemptData(attempt), nil
}

// GetRecentAttempts получает последние попытки
func (s *AttemptService) GetRecentAttempts(ctx context.Context, childProfileID string, limit int) ([]AttemptData, error) {
	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		return nil, fmt.Errorf("invalid child_profile_id: %w", err)
	}

	attempts, err := s.store.Attempts.GetRecentAttempts(ctx, profileUUID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent attempts: %w", err)
	}

	// Конвертируем в []AttemptData
	result := make([]AttemptData, 0, len(attempts))
	for _, attempt := range attempts {
		result = append(result, *s.convertToAttemptData(attempt))
	}

	return result, nil
}

// convertToAttemptData конвертирует store.Attempt в service.AttemptData
func (s *AttemptService) convertToAttemptData(attempt *store.Attempt) *AttemptData {
	var parseResult *types.ParseResponse
	if len(attempt.ParseResult) > 0 {
		if err := json.Unmarshal(attempt.ParseResult, &parseResult); err != nil {
			log.Printf("[AttemptService] Failed to unmarshal parse result: %v", err)
		}
	}

	var detectResult *types.DetectResponse
	if len(attempt.DetectResult) > 0 {
		if err := json.Unmarshal(attempt.DetectResult, &detectResult); err != nil {
			log.Printf("[AttemptService] Failed to unmarshal detect result: %v", err)
		}
	}

	var hintsResult *types.HintResponse
	if len(attempt.HintsResult) > 0 {
		if err := json.Unmarshal(attempt.HintsResult, &hintsResult); err != nil {
			log.Printf("[AttemptService] Failed to unmarshal hints result: %v", err)
		}
	}

	var checkResult *types.CheckResponse
	if len(attempt.CheckResult) > 0 {
		if err := json.Unmarshal(attempt.CheckResult, &checkResult); err != nil {
			log.Printf("[AttemptService] Failed to unmarshal check result: %v", err)
		}
	}

	taskImage := ""
	if attempt.TaskImageURL.Valid {
		taskImage = attempt.TaskImageURL.String
	}

	answerImage := ""
	if attempt.AnswerImageURL.Valid {
		answerImage = attempt.AnswerImageURL.String
	}

	return &AttemptData{
		ID:              attempt.ID.String(),
		ChildProfileID:  attempt.ChildProfileID.String(),
		Type:            attempt.AttemptType,
		Status:          attempt.Status,
		TaskImageData:   taskImage,
		AnswerImageData: answerImage,
		ParseResult:     parseResult,
		DetectResult:    detectResult,
		HintsResult:     hintsResult,
		CheckResult:     checkResult,
		CurrentHint:     attempt.CurrentHintIndex,
		CreatedAt:       attempt.CreatedAt,
		UpdatedAt:       attempt.UpdatedAt,
	}
}

// serializeJSON helper для сериализации в JSON
func serializeJSON(v interface{}) (json.RawMessage, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(data), nil
}
