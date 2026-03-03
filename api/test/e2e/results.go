package e2e

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestResult holds the results of an E2E test run
type TestResult struct {
	ChatID         int64            `json:"chat_id"`
	SessionID      string           `json:"session_id,omitempty"`
	ImagePath      string           `json:"image_path,omitempty"`
	StartTime      time.Time        `json:"start_time"`
	EndTime        time.Time        `json:"end_time"`
	HintsCount     int              `json:"hints_count,omitempty"`
	Success        bool             `json:"success"`
	Error          string           `json:"error,omitempty"`
	TimelineEvents []TimelineRecord `json:"timeline_events,omitempty"`
}

// TimelineRecord is a JSON-friendly version of store.TimelineEvent
type TimelineRecord struct {
	TaskSessionID string      `json:"task_session_id"`
	Direction     string      `json:"direction"`
	EventType     string      `json:"event_type"`
	Provider      string      `json:"provider"`
	OK            bool        `json:"ok"`
	LatencyMS     *int64      `json:"latency_ms,omitempty"`
	Text          string      `json:"text,omitempty"`
	InputPayload  interface{} `json:"input_payload,omitempty"`
	OutputPayload interface{} `json:"output_payload,omitempty"`
	Error         string      `json:"error,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
}

// NewTestResult creates a new test result
func NewTestResult(chatID int64) *TestResult {
	return &TestResult{
		ChatID:    chatID,
		StartTime: time.Now(),
	}
}

// Duration returns the total test duration
func (r *TestResult) Duration() time.Duration {
	if r.EndTime.IsZero() {
		return time.Since(r.StartTime)
	}
	return r.EndTime.Sub(r.StartTime)
}

// Finish marks the test as finished
func (r *TestResult) Finish(success bool, err error) {
	r.EndTime = time.Now()
	r.Success = success
	if err != nil {
		r.Error = err.Error()
	}
}

// ToJSON converts the result to JSON string
func (r *TestResult) ToJSON() string {
	data, _ := json.MarshalIndent(r, "", "  ")
	return string(data)
}

// saveTestResults saves test results to a JSON file
func saveTestResults(t *testing.T, path string, result *TestResult) {
	t.Helper()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Logf("Warning: failed to create results directory: %v", err)
		return
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Logf("Warning: failed to marshal test results: %v", err)
		return
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Logf("Warning: failed to write test results: %v", err)
		return
	}

	t.Logf("Test results saved to: %s", path)
}
