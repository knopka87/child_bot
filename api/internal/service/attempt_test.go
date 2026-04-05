package service

import (
	"context"
	"encoding/base64"
	"testing"
	"time"
)

func TestAttemptService_CreateAttempt(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	st := setupTestStore(t)

	// Create test profile
	profileID := createTestProfile(t, db, "test", "user123", 5)

	// Create service
	mockLLM := &mockLLMClient{}
	service := NewAttemptService(st, mockLLM, "gpt-4")

	tests := []struct {
		name        string
		profileID   string
		attemptType string
		expectError bool
	}{
		{
			name:        "success - help attempt",
			profileID:   profileID,
			attemptType: "help",
			expectError: false,
		},
		{
			name:        "success - check attempt",
			profileID:   profileID,
			attemptType: "check",
			expectError: false,
		},
		{
			name:        "error - invalid type",
			profileID:   profileID,
			attemptType: "invalid",
			expectError: true,
		},
		{
			name:        "error - invalid profile ID",
			profileID:   "invalid-uuid",
			attemptType: "help",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			attemptID, err := service.CreateAttempt(ctx, tt.profileID, tt.attemptType)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if attemptID == "" {
				t.Error("expected non-empty attempt ID")
			}

			// Verify attempt was created in database
			var status, attemptType string
			query := "SELECT status, attempt_type FROM attempts WHERE id = $1"
			err = db.QueryRow(query, attemptID).Scan(&status, &attemptType)
			if err != nil {
				t.Fatalf("failed to query attempt: %v", err)
			}

			if status != "created" {
				t.Errorf("expected status 'created', got %q", status)
			}

			if attemptType != tt.attemptType {
				t.Errorf("expected type %q, got %q", tt.attemptType, attemptType)
			}

			// Cleanup
			cleanupTestData(t, db, "attempts", "id = $1", attemptID)
		})
	}
}

func TestAttemptService_UploadImage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	st := setupTestStore(t)

	// Create test profile and attempt
	profileID := createTestProfile(t, db, "test", "user123", 5)
	attemptID := createTestAttempt(t, db, profileID, "help")

	// Create service
	mockLLM := &mockLLMClient{}
	service := NewAttemptService(st, mockLLM, "gpt-4")

	// Sample base64 image (1x1 transparent PNG)
	sampleImage := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="

	tests := []struct {
		name        string
		attemptID   string
		imageType   string
		imageBase64 string
		expectError bool
	}{
		{
			name:        "success - task image",
			attemptID:   attemptID,
			imageType:   "task",
			imageBase64: sampleImage,
			expectError: false,
		},
		{
			name:        "success - answer image",
			attemptID:   attemptID,
			imageType:   "answer",
			imageBase64: sampleImage,
			expectError: false,
		},
		{
			name:        "error - invalid image type",
			attemptID:   attemptID,
			imageType:   "invalid",
			imageBase64: sampleImage,
			expectError: true,
		},
		{
			name:        "error - invalid base64",
			attemptID:   attemptID,
			imageType:   "task",
			imageBase64: "not-base64!@#$",
			expectError: true,
		},
		{
			name:        "error - attempt not found",
			attemptID:   "00000000-0000-0000-0000-000000000000",
			imageType:   "task",
			imageBase64: sampleImage,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			err := service.UploadImage(ctx, tt.attemptID, tt.imageType, tt.imageBase64)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify image was saved in database
			var taskImageURL, answerImageURL *string
			query := "SELECT task_image_url, answer_image_url FROM attempts WHERE id = $1"
			err = db.QueryRow(query, tt.attemptID).Scan(&taskImageURL, &answerImageURL)
			if err != nil {
				t.Fatalf("failed to query attempt: %v", err)
			}

			if tt.imageType == "task" && taskImageURL == nil {
				t.Error("expected task_image_url to be set")
			}

			if tt.imageType == "answer" && answerImageURL == nil {
				t.Error("expected answer_image_url to be set")
			}
		})
	}
}

func TestAttemptService_ProcessHelp_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	st := setupTestStore(t)

	// Create test profile and attempt
	profileID := createTestProfile(t, db, "test", "user123", 5)
	attemptID := createTestAttempt(t, db, profileID, "help")

	// Upload task image
	sampleImage := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
	_, err := db.Exec("UPDATE attempts SET task_image_url = $1 WHERE id = $2", "data:image/png;base64,"+sampleImage, attemptID)
	if err != nil {
		t.Fatalf("failed to set task image: %v", err)
	}

	// Create service with mock LLM
	mockLLM := &mockLLMClient{
		detectFunc: func(ctx context.Context, llmName string, req interface{}) (interface{}, error) {
			return map[string]interface{}{
				"schema_version": "1.0",
				"classification": map[string]interface{}{
					"subject_candidate": "math",
					"confidence":        0.95,
				},
				"quality": map[string]interface{}{
					"is_acceptable": true,
					"score":         0.9,
				},
			}, nil
		},
		parseFunc: func(ctx context.Context, llmName string, req interface{}) (interface{}, error) {
			return map[string]interface{}{
				"task": map[string]interface{}{
					"id":   attemptID,
					"text": "Solve 2+2",
				},
				"items": []interface{}{
					map[string]interface{}{
						"id":   "item1",
						"text": "2+2=?",
					},
				},
			}, nil
		},
		hintFunc: func(ctx context.Context, llmName string, req interface{}) (interface{}, error) {
			return map[string]interface{}{
				"hints": []interface{}{
					map[string]interface{}{
						"text":  "First hint: add the numbers",
						"level": 1,
					},
					map[string]interface{}{
						"text":  "Second hint: 2+2=4",
						"level": 2,
					},
				},
			}, nil
		},
	}

	service := NewAttemptService(st, mockLLM, "gpt-4")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Process help
	err = service.ProcessHelp(ctx, attemptID)
	if err != nil {
		t.Fatalf("ProcessHelp failed: %v", err)
	}

	// Verify results in database
	var status string
	var detectResult, parseResult, hintsResult []byte
	query := `
		SELECT status, detect_result, parse_result, hints_result
		FROM attempts WHERE id = $1
	`
	err = db.QueryRow(query, attemptID).Scan(&status, &detectResult, &parseResult, &hintsResult)
	if err != nil {
		t.Fatalf("failed to query attempt: %v", err)
	}

	if status != "completed" {
		t.Errorf("expected status 'completed', got %q", status)
	}

	if len(detectResult) == 0 {
		t.Error("expected detect_result to be populated")
	}

	if len(parseResult) == 0 {
		t.Error("expected parse_result to be populated")
	}

	if len(hintsResult) == 0 {
		t.Error("expected hints_result to be populated")
	}

	t.Logf("ProcessHelp completed successfully for attempt %s", attemptID)
}

func TestAttemptService_NextHint(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	st := setupTestStore(t)

	// Create test profile and attempt
	profileID := createTestProfile(t, db, "test", "user123", 5)
	attemptID := createTestAttempt(t, db, profileID, "help")

	// Set up hints_result with 3 hints
	hintsJSON := `{
		"hints": [
			{"text": "Hint 1", "level": 1},
			{"text": "Hint 2", "level": 2},
			{"text": "Hint 3", "level": 3}
		]
	}`
	_, err := db.Exec("UPDATE attempts SET hints_result = $1, current_hint_index = 0 WHERE id = $2", hintsJSON, attemptID)
	if err != nil {
		t.Fatalf("failed to set hints: %v", err)
	}

	// Create service
	mockLLM := &mockLLMClient{}
	service := NewAttemptService(st, mockLLM, "gpt-4")

	ctx := context.Background()

	// Test getting hints sequentially
	for i := 0; i < 3; i++ {
		t.Run("hint_"+string(rune('0'+i)), func(t *testing.T) {
			result, err := service.NextHint(ctx, attemptID)
			if err != nil {
				t.Fatalf("NextHint failed: %v", err)
			}

			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatalf("expected map result, got %T", result)
			}

			hintIndex, ok := resultMap["hint_index"].(int)
			if !ok || hintIndex != i {
				t.Errorf("expected hint_index %d, got %v", i, resultMap["hint_index"])
			}

			hasMore, ok := resultMap["has_more"].(bool)
			if !ok {
				t.Error("expected has_more field")
			}

			expectedHasMore := i < 2
			if hasMore != expectedHasMore {
				t.Errorf("expected has_more=%v, got %v", expectedHasMore, hasMore)
			}
		})
	}

	// Test getting hint when all hints exhausted
	t.Run("no_more_hints", func(t *testing.T) {
		result, err := service.NextHint(ctx, attemptID)
		if err != nil {
			t.Fatalf("NextHint failed: %v", err)
		}

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("expected map result, got %T", result)
		}

		hasMore, ok := resultMap["has_more"].(bool)
		if !ok || hasMore {
			t.Error("expected has_more=false when hints exhausted")
		}
	})
}

func TestAttemptService_DeleteAttempt(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	st := setupTestStore(t)

	// Create test profile and attempt
	profileID := createTestProfile(t, db, "test", "user123", 5)
	attemptID := createTestAttempt(t, db, profileID, "help")

	// Create service
	mockLLM := &mockLLMClient{}
	service := NewAttemptService(st, mockLLM, "gpt-4")

	ctx := context.Background()

	// Delete attempt
	err := service.DeleteAttempt(ctx, attemptID)
	if err != nil {
		t.Fatalf("DeleteAttempt failed: %v", err)
	}

	// Verify attempt was deleted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM attempts WHERE id = $1", attemptID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query attempts: %v", err)
	}

	if count != 0 {
		t.Errorf("expected attempt to be deleted, but found %d records", count)
	}
}

// Benchmark tests

func BenchmarkAttemptService_CreateAttempt(b *testing.B) {
	db := setupTestDB(&testing.T{})
	st := setupTestStore(&testing.T{})

	profileID := createTestProfile(&testing.T{}, db, "test", "bench_user", 5)

	mockLLM := &mockLLMClient{}
	service := NewAttemptService(st, mockLLM, "gpt-4")

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attemptID, err := service.CreateAttempt(ctx, profileID, "help")
		if err != nil {
			b.Fatalf("CreateAttempt failed: %v", err)
		}

		// Cleanup
		db.Exec("DELETE FROM attempts WHERE id = $1", attemptID)
	}
}

func BenchmarkAttemptService_UploadImage(b *testing.B) {
	db := setupTestDB(&testing.T{})
	st := setupTestStore(&testing.T{})

	profileID := createTestProfile(&testing.T{}, db, "test", "bench_user", 5)
	attemptID := createTestAttempt(&testing.T{}, db, profileID, "help")

	mockLLM := &mockLLMClient{}
	service := NewAttemptService(st, mockLLM, "gpt-4")

	sampleImage := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := service.UploadImage(ctx, attemptID, "task", sampleImage)
		if err != nil {
			b.Fatalf("UploadImage failed: %v", err)
		}
	}
}
