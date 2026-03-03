package e2e

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	PipelineOK     bool             `json:"pipeline_ok"` // P2.1: renamed from "success"
	Error          string           `json:"error,omitempty"`
	TimelineEvents []TimelineRecord `json:"timeline_events,omitempty"`
}

// TimelineRecord is a JSON-friendly version of store.TimelineEvent
// P2.3: base64 images are stripped and replaced with metadata
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

// sanitizePayload removes base64 images from payloads and replaces them with metadata
// P2.3: Don't store base64 in logs
func sanitizePayload(payload interface{}) interface{} {
	if payload == nil {
		return nil
	}

	switch p := payload.(type) {
	case map[string]interface{}:
		sanitized := make(map[string]interface{})
		for k, v := range p {
			if k == "image" || k == "answer_image" || k == "task_image" {
				if str, ok := v.(string); ok && len(str) > 1000 {
					// Replace base64 with metadata
					sanitized[k+"_metadata"] = map[string]interface{}{
						"sha256":     fmt.Sprintf("%x", sha256.Sum256([]byte(str)))[:16],
						"bytes_size": len(str),
						"truncated":  true,
					}
					continue
				}
			}
			sanitized[k] = sanitizePayload(v)
		}
		return sanitized
	case []interface{}:
		sanitized := make([]interface{}, len(p))
		for i, v := range p {
			sanitized[i] = sanitizePayload(v)
		}
		return sanitized
	case string:
		// Truncate very long strings (likely base64)
		if len(p) > 1000 && !strings.Contains(p, " ") {
			return fmt.Sprintf("[truncated: %d bytes, sha256: %s]", len(p), fmt.Sprintf("%x", sha256.Sum256([]byte(p)))[:16])
		}
		return p
	default:
		return payload
	}
}

// SanitizeTimelineEvents removes base64 data from timeline events
func SanitizeTimelineEvents(events []TimelineRecord) []TimelineRecord {
	sanitized := make([]TimelineRecord, len(events))
	for i, e := range events {
		sanitized[i] = TimelineRecord{
			TaskSessionID: e.TaskSessionID,
			Direction:     e.Direction,
			EventType:     e.EventType,
			Provider:      e.Provider,
			OK:            e.OK,
			LatencyMS:     e.LatencyMS,
			Text:          e.Text,
			InputPayload:  sanitizePayload(e.InputPayload),
			OutputPayload: sanitizePayload(e.OutputPayload),
			Error:         e.Error,
			CreatedAt:     e.CreatedAt,
		}
	}
	return sanitized
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
func (r *TestResult) Finish(pipelineOK bool, err error) {
	r.EndTime = time.Now()
	r.PipelineOK = pipelineOK
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
