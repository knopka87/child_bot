package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"child-bot/api/internal/domain"
	"child-bot/api/internal/service"
)

// Test Helpers

// mockResponseWriter captures response for testing
type mockResponseWriter struct {
	statusCode int
	body       *bytes.Buffer
	header     http.Header
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		statusCode: http.StatusOK,
		body:       new(bytes.Buffer),
		header:     make(http.Header),
	}
}

func (m *mockResponseWriter) Header() http.Header {
	return m.header
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	return m.body.Write(b)
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

// makeRequest creates a test HTTP request
func makeRequest(t *testing.T, method, path string, body interface{}) *http.Request {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Extract path parameters (e.g., /attempts/{id})
	req.SetPathValue("id", extractPathParam(path, "id"))
	req.SetPathValue("childProfileId", extractPathParam(path, "childProfileId"))

	return req
}

// extractPathParam extracts path parameter from URL for testing
func extractPathParam(path, param string) string {
	// Simple extraction: /attempts/attempt-123/images -> "attempt-123"
	segments := []string{}
	current := ""
	for _, ch := range path {
		if ch == '/' {
			if current != "" {
				segments = append(segments, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		segments = append(segments, current)
	}

	// For /attempts/{id} patterns, the ID is usually the 2nd segment
	if len(segments) >= 2 && param == "id" {
		return segments[1]
	}

	// For /home/{childProfileId}, it's the 2nd segment
	if len(segments) >= 2 && param == "childProfileId" {
		return segments[1]
	}

	return ""
}

// makeAuthRequest creates a test HTTP request with auth headers
func makeAuthRequest(t *testing.T, method, path string, body interface{}, platformID, childProfileID string) *http.Request {
	t.Helper()

	req := makeRequest(t, method, path, body)
	req.Header.Set("X-Platform-ID", platformID)
	req.Header.Set("X-Child-Profile-ID", childProfileID)
	return req
}

// decodeResponse decodes JSON response
func decodeResponse(t *testing.T, w *mockResponseWriter, target interface{}) {
	t.Helper()

	if err := json.Unmarshal(w.body.Bytes(), target); err != nil {
		t.Fatalf("failed to decode response: %v\nBody: %s", err, w.body.String())
	}
}

// assertStatus checks response status code
func assertStatus(t *testing.T, w *mockResponseWriter, expectedStatus int) {
	t.Helper()

	if w.statusCode != expectedStatus {
		t.Errorf("expected status %d, got %d\nBody: %s",
			expectedStatus, w.statusCode, w.body.String())
	}
}

// assertError checks that response contains error
func assertError(t *testing.T, w *mockResponseWriter, expectedStatus int, expectedMessage string) {
	t.Helper()

	assertStatus(t, w, expectedStatus)

	var resp map[string]interface{}
	decodeResponse(t, w, &resp)

	errMsg, ok := resp["error"].(string)
	if !ok {
		t.Errorf("expected error field in response, got: %v", resp)
		return
	}

	if errMsg != expectedMessage {
		t.Errorf("expected error message %q, got %q", expectedMessage, errMsg)
	}
}

// Mock Services

// mockAttemptService is a mock implementation of AttemptService for testing
type mockAttemptService struct {
	createFunc            func(ctx context.Context, childProfileID, attemptType string) (string, error)
	uploadImageFunc       func(ctx context.Context, attemptID, imageType, imageData string) (string, error)
	processHelpFunc       func(ctx context.Context, attemptID, imageBase64 string) error
	processCheckFunc      func(ctx context.Context, attemptID, taskImageBase64, answerImageBase64 string) error
	getAttemptResultFunc  func(ctx context.Context, attemptID string) (*service.AttemptData, error)
	getNextHintFunc       func(ctx context.Context, attemptID string) (*domain.HelpResult, error)
	deleteFunc            func(ctx context.Context, attemptID string) error
	getUnfinishedFunc     func(ctx context.Context, childProfileID string) (*service.AttemptData, error)
	getRecentAttemptsFunc func(ctx context.Context, childProfileID string, limit int) ([]service.AttemptData, error)
}

func (m *mockAttemptService) CreateAttempt(ctx context.Context, childProfileID, attemptType string) (string, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, childProfileID, attemptType)
	}
	return "", errors.New("not implemented")
}

func (m *mockAttemptService) UploadImage(ctx context.Context, attemptID, imageType, imageData string) (string, error) {
	if m.uploadImageFunc != nil {
		return m.uploadImageFunc(ctx, attemptID, imageType, imageData)
	}
	return "", errors.New("not implemented")
}

func (m *mockAttemptService) ProcessHelp(ctx context.Context, attemptID, imageBase64 string) error {
	if m.processHelpFunc != nil {
		return m.processHelpFunc(ctx, attemptID, imageBase64)
	}
	return errors.New("not implemented")
}

func (m *mockAttemptService) ProcessCheck(ctx context.Context, attemptID, taskImageBase64, answerImageBase64 string) error {
	if m.processCheckFunc != nil {
		return m.processCheckFunc(ctx, attemptID, taskImageBase64, answerImageBase64)
	}
	return errors.New("not implemented")
}

func (m *mockAttemptService) GetAttemptResult(ctx context.Context, attemptID string) (*service.AttemptData, error) {
	if m.getAttemptResultFunc != nil {
		return m.getAttemptResultFunc(ctx, attemptID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAttemptService) GetNextHint(ctx context.Context, attemptID string) (*domain.HelpResult, error) {
	if m.getNextHintFunc != nil {
		return m.getNextHintFunc(ctx, attemptID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAttemptService) DeleteAttempt(ctx context.Context, attemptID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, attemptID)
	}
	return errors.New("not implemented")
}

func (m *mockAttemptService) GetUnfinishedAttempt(ctx context.Context, childProfileID string) (*service.AttemptData, error) {
	if m.getUnfinishedFunc != nil {
		return m.getUnfinishedFunc(ctx, childProfileID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAttemptService) GetRecentAttempts(ctx context.Context, childProfileID string, limit int) ([]service.AttemptData, error) {
	if m.getRecentAttemptsFunc != nil {
		return m.getRecentAttemptsFunc(ctx, childProfileID, limit)
	}
	return nil, errors.New("not implemented")
}
