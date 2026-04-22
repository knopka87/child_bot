package handler

import "time"

// CreateAttemptRequest запрос на создание попытки
type CreateAttemptRequest struct {
	Type           string `json:"type"`             // "help" или "check"
	ChildProfileID string `json:"child_profile_id"` // UUID профиля ребенка
}

// CreateAttemptResponse ответ на создание попытки
type CreateAttemptResponse struct {
	AttemptID string `json:"attempt_id"`
	Status    string `json:"status"`
}

// UploadImageRequest запрос на загрузку изображения
type UploadImageRequest struct {
	ImageType string `json:"image_type"` // "task" или "answer"
	ImageData string `json:"image_data"` // base64 encoded image
}

// UploadImageResponse ответ на загрузку изображения
type UploadImageResponse struct {
	ImageURL string `json:"image_url"`
	Message  string `json:"message"`
}

// ProcessAttemptResponse ответ на обработку попытки
type ProcessAttemptResponse struct {
	Status  string `json:"status"` // "processing", "completed", "failed"
	Message string `json:"message,omitempty"`
}

// GetResultResponse ответ с результатом попытки
type GetResultResponse struct {
	AttemptID string                 `json:"attempt_id"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"`
	Result    map[string]interface{} `json:"result,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// NextHintResponse ответ на запрос следующей подсказки
type NextHintResponse struct {
	Hint        string `json:"hint"`
	HintIndex   int    `json:"hint_index"`
	TotalHints  int    `json:"total_hints"`
	HasMoreHint bool   `json:"has_more_hints"`
	Completed   bool   `json:"completed"` // true если все подсказки просмотрены и попытка завершена
}

// ErrorResponse стандартный формат ошибки
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}
