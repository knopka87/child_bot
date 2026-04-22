package domain

import "time"

// AttemptType тип попытки
type AttemptType string

const (
	AttemptTypeHelp  AttemptType = "help"
	AttemptTypeCheck AttemptType = "check"
)

// AttemptStatus статус попытки
type AttemptStatus string

const (
	AttemptStatusCreated    AttemptStatus = "created"
	AttemptStatusProcessing AttemptStatus = "processing"
	AttemptStatusCompleted  AttemptStatus = "completed"
	AttemptStatusFailed     AttemptStatus = "failed"
)

// Attempt представляет попытку решения задачи
type Attempt struct {
	ID             string        `json:"id"`
	ChildProfileID string        `json:"child_profile_id"`
	Type           AttemptType   `json:"type"`
	Status         AttemptStatus `json:"status"`
	TaskImageURL   string        `json:"task_image_url,omitempty"`
	AnswerImageURL string        `json:"answer_image_url,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// HelpResult результат обработки help запроса
type HelpResult struct {
	Subject     string   `json:"subject"`
	TaskText    string   `json:"task_text"`
	Hints       []string `json:"hints"`
	CurrentHint int      `json:"current_hint"`
	TotalHints  int      `json:"total_hints"`
}

// CheckResult результат проверки решения
type CheckResult struct {
	IsCorrect   bool   `json:"is_correct"`
	Decision    string `json:"decision"`
	Explanation string `json:"explanation"`
	Score       int    `json:"score,omitempty"`
}
