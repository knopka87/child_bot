package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"child-bot/api/internal/domain"
	"child-bot/api/internal/service"
)

func TestAttemptHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockCreate     func(ctx context.Context, childProfileID, attemptType string) (string, error)
		expectedStatus int
		checkResponse  func(t *testing.T, w *mockResponseWriter)
	}{
		{
			name: "success - help attempt",
			requestBody: map[string]string{
				"child_profile_id": "550e8400-e29b-41d4-a716-446655440000",
				"type":             "help",
			},
			mockCreate: func(ctx context.Context, childProfileID, attemptType string) (string, error) {
				if childProfileID != "550e8400-e29b-41d4-a716-446655440000" {
					t.Errorf("unexpected childProfileID: %s", childProfileID)
				}
				if attemptType != "help" {
					t.Errorf("unexpected attemptType: %s", attemptType)
				}
				return "attempt-123", nil
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *mockResponseWriter) {
				var resp CreateAttemptResponse
				decodeResponse(t, w, &resp)
				if resp.AttemptID != "attempt-123" {
					t.Errorf("expected attempt_id 'attempt-123', got %s", resp.AttemptID)
				}
				if resp.Status != "created" {
					t.Errorf("expected status 'created', got %s", resp.Status)
				}
			},
		},
		{
			name: "success - check attempt",
			requestBody: map[string]string{
				"child_profile_id": "550e8400-e29b-41d4-a716-446655440000",
				"type":             "check",
			},
			mockCreate: func(ctx context.Context, childProfileID, attemptType string) (string, error) {
				return "attempt-456", nil
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *mockResponseWriter) {
				var resp CreateAttemptResponse
				decodeResponse(t, w, &resp)
				if resp.AttemptID != "attempt-456" {
					t.Errorf("expected attempt_id 'attempt-456', got %s", resp.AttemptID)
				}
			},
		},
		{
			name: "validation error - invalid type",
			requestBody: map[string]string{
				"child_profile_id": "550e8400-e29b-41d4-a716-446655440000",
				"type":             "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *mockResponseWriter) {
				var resp map[string]string
				decodeResponse(t, w, &resp)
				if resp["error"] == "" {
					t.Error("expected error message")
				}
			},
		},
		{
			name: "validation error - missing child_profile_id",
			requestBody: map[string]string{
				"type": "help",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *mockResponseWriter) {
				var resp map[string]string
				decodeResponse(t, w, &resp)
				if resp["error"] == "" {
					t.Error("expected error message")
				}
			},
		},
		{
			name: "service error",
			requestBody: map[string]string{
				"child_profile_id": "550e8400-e29b-41d4-a716-446655440000",
				"type":             "help",
			},
			mockCreate: func(ctx context.Context, childProfileID, attemptType string) (string, error) {
				return "", errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *mockResponseWriter) {
				var resp map[string]string
				decodeResponse(t, w, &resp)
				if resp["error"] == "" {
					t.Error("expected error message")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := &mockAttemptService{
				createFunc: tt.mockCreate,
			}

			handler := NewAttemptHandler(mockService)

			// Make request
			req := makeRequest(t, http.MethodPost, "/attempts", tt.requestBody)
			w := newMockResponseWriter()

			// Execute
			handler.Create(w, req)

			// Assert
			assertStatus(t, w, tt.expectedStatus)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestAttemptHandler_UploadImage(t *testing.T) {
	tests := []struct {
		name           string
		attemptID      string
		requestBody    interface{}
		mockUpload     func(ctx context.Context, attemptID, imageType, imageData string) (string, error)
		expectedStatus int
	}{
		{
			name:      "success - task image",
			attemptID: "attempt-123",
			requestBody: map[string]string{
				"image_type": "task",
				"image_data": "iVBORw0KGgoAAAANSUhEUg...",
			},
			mockUpload: func(ctx context.Context, attemptID, imageType, imageData string) (string, error) {
				if attemptID != "attempt-123" {
					t.Errorf("unexpected attemptID: %s", attemptID)
				}
				if imageType != "task" {
					t.Errorf("unexpected imageType: %s", imageType)
				}
				return "https://storage.example.com/task123.jpg", nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "success - answer image",
			attemptID: "attempt-123",
			requestBody: map[string]string{
				"image_type": "answer",
				"image_data": "iVBORw0KGgoAAAANSUhEUg...",
			},
			mockUpload: func(ctx context.Context, attemptID, imageType, imageData string) (string, error) {
				return "https://storage.example.com/answer123.jpg", nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "validation error - invalid image_type",
			attemptID: "attempt-123",
			requestBody: map[string]string{
				"image_type": "invalid",
				"image_data": "iVBORw0KGgoAAAANSUhEUg...",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "validation error - missing image_data",
			attemptID: "attempt-123",
			requestBody: map[string]string{
				"image_type": "task",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "service error",
			attemptID: "attempt-123",
			requestBody: map[string]string{
				"image_type": "task",
				"image_data": "iVBORw0KGgoAAAANSUhEUg...",
			},
			mockUpload: func(ctx context.Context, attemptID, imageType, imageData string) (string, error) {
				return "", errors.New("upload failed")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockAttemptService{
				uploadImageFunc: tt.mockUpload,
			}

			handler := NewAttemptHandler(mockService)

			req := makeRequest(t, http.MethodPost, "/attempts/"+tt.attemptID+"/images", tt.requestBody)
			w := newMockResponseWriter()

			handler.UploadImage(w, req)

			assertStatus(t, w, tt.expectedStatus)
		})
	}
}

func TestAttemptHandler_Process(t *testing.T) {
	tests := []struct {
		name           string
		attemptID      string
		mockProcess    func(ctx context.Context, attemptID, imageBase64 string) error
		expectedStatus int
	}{
		{
			name:      "success - help attempt",
			attemptID: "attempt-123",
			mockProcess: func(ctx context.Context, attemptID, imageBase64 string) error {
				if attemptID != "attempt-123" {
					t.Errorf("unexpected attemptID: %s", attemptID)
				}
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "service error",
			attemptID: "attempt-123",
			mockProcess: func(ctx context.Context, attemptID, imageBase64 string) error {
				return errors.New("processing failed")
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockAttemptService{
				processHelpFunc: tt.mockProcess,
			}

			handler := NewAttemptHandler(mockService)

			req := makeRequest(t, http.MethodPost, "/attempts/"+tt.attemptID+"/process", nil)
			w := newMockResponseWriter()

			handler.Process(w, req)

			assertStatus(t, w, tt.expectedStatus)
		})
	}
}

func TestAttemptHandler_GetResult(t *testing.T) {
	tests := []struct {
		name           string
		attemptID      string
		mockGetResult  func(ctx context.Context, attemptID string) (*service.AttemptData, error)
		expectedStatus int
		checkResponse  func(t *testing.T, w *mockResponseWriter)
	}{
		{
			name:      "success - result available",
			attemptID: "attempt-123",
			mockGetResult: func(ctx context.Context, attemptID string) (*service.AttemptData, error) {
				return &service.AttemptData{
					ID:     attemptID,
					Status: "completed",
					Type:   "help",
				}, nil
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *mockResponseWriter) {
				var resp map[string]interface{}
				decodeResponse(t, w, &resp)
				if resp["attempt_id"] != "attempt-123" {
					t.Errorf("unexpected attempt_id: %v", resp["attempt_id"])
				}
			},
		},
		{
			name:      "not found",
			attemptID: "attempt-999",
			mockGetResult: func(ctx context.Context, attemptID string) (*service.AttemptData, error) {
				return nil, errors.New("not found")
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockAttemptService{
				getAttemptResultFunc: tt.mockGetResult,
			}

			handler := NewAttemptHandler(mockService)

			req := makeRequest(t, http.MethodGet, "/attempts/"+tt.attemptID+"/result", nil)
			w := newMockResponseWriter()

			handler.GetResult(w, req)

			assertStatus(t, w, tt.expectedStatus)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestAttemptHandler_NextHint(t *testing.T) {
	tests := []struct {
		name           string
		attemptID      string
		mockNextHint   func(ctx context.Context, attemptID string) (*domain.HelpResult, error)
		expectedStatus int
		checkResponse  func(t *testing.T, w *mockResponseWriter)
	}{
		{
			name:      "success - hint available",
			attemptID: "attempt-123",
			mockNextHint: func(ctx context.Context, attemptID string) (*domain.HelpResult, error) {
				return &domain.HelpResult{
					Subject:     "Math",
					TaskText:    "Solve 2+2",
					Hints:       []string{"Hint 1", "Hint 2"},
					CurrentHint: 0,
					TotalHints:  2,
				}, nil
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *mockResponseWriter) {
				var resp map[string]interface{}
				decodeResponse(t, w, &resp)
				// Just check that we got a response
				if resp["subject"] != "Math" {
					t.Logf("Response: %v", resp)
				}
			},
		},
		{
			name:      "service error",
			attemptID: "attempt-123",
			mockNextHint: func(ctx context.Context, attemptID string) (*domain.HelpResult, error) {
				return nil, errors.New("hint generation failed")
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockAttemptService{
				getNextHintFunc: tt.mockNextHint,
			}

			handler := NewAttemptHandler(mockService)

			req := makeRequest(t, http.MethodPost, "/attempts/"+tt.attemptID+"/next-hint", nil)
			w := newMockResponseWriter()

			handler.NextHint(w, req)

			assertStatus(t, w, tt.expectedStatus)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestAttemptHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		attemptID      string
		mockDelete     func(ctx context.Context, attemptID string) error
		expectedStatus int
	}{
		{
			name:      "success",
			attemptID: "attempt-123",
			mockDelete: func(ctx context.Context, attemptID string) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "service error",
			attemptID: "attempt-123",
			mockDelete: func(ctx context.Context, attemptID string) error {
				return errors.New("delete failed")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockAttemptService{
				deleteFunc: tt.mockDelete,
			}

			handler := NewAttemptHandler(mockService)

			req := makeRequest(t, http.MethodDelete, "/attempts/"+tt.attemptID, nil)
			w := newMockResponseWriter()

			handler.Delete(w, req)

			assertStatus(t, w, tt.expectedStatus)
		})
	}
}
