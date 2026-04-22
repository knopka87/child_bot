package service

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"child-bot/api/internal/store"
)

// Test helpers

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// setupTestStore creates a store connected to test database
func setupTestStore(t *testing.T) *store.Store {
	t.Helper()
	return store.NewStore(setupTestDB(t))
}

// cleanupTestData removes test data from database
func cleanupTestData(t *testing.T, db *sql.DB, tableName string, whereClause string, args ...interface{}) {
	t.Helper()

	query := "DELETE FROM " + tableName
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	_, err := db.Exec(query, args...)
	if err != nil {
		t.Logf("Warning: failed to cleanup test data from %s: %v", tableName, err)
	}
}

// createTestProfile creates a test child_profile and returns its ID
func createTestProfile(t *testing.T, db *sql.DB, platformID, platformUserID string, grade int) string {
	t.Helper()

	var profileID string
	query := `
		INSERT INTO child_profiles (display_name, platform_id, platform_user_id, grade, avatar_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := db.QueryRow(query, "Test User", platformID, platformUserID, grade, "avatar1").Scan(&profileID)
	if err != nil {
		t.Fatalf("failed to create test profile: %v", err)
	}

	t.Cleanup(func() {
		cleanupTestData(t, db, "child_profiles", "id = $1", profileID)
	})

	return profileID
}

// createTestAttempt creates a test attempt and returns its ID
func createTestAttempt(t *testing.T, db *sql.DB, childProfileID, attemptType string) string {
	t.Helper()

	var attemptID string
	query := `
		INSERT INTO attempts (child_profile_id, attempt_type, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := db.QueryRow(query, childProfileID, attemptType, "created").Scan(&attemptID)
	if err != nil {
		t.Fatalf("failed to create test attempt: %v", err)
	}

	t.Cleanup(func() {
		cleanupTestData(t, db, "attempts", "id = $1", attemptID)
	})

	return attemptID
}

// Mock LLM Client for testing

// mockLLMClient is a mock LLM client that returns predefined responses
type mockLLMClient struct {
	detectFunc func(ctx context.Context, llmName string, req interface{}) (interface{}, error)
	parseFunc  func(ctx context.Context, llmName string, req interface{}) (interface{}, error)
	hintFunc   func(ctx context.Context, llmName string, req interface{}) (interface{}, error)
	checkFunc  func(ctx context.Context, llmName string, req interface{}) (interface{}, error)
}

func (m *mockLLMClient) Detect(ctx context.Context, llmName string, req interface{}) (interface{}, error) {
	if m.detectFunc != nil {
		return m.detectFunc(ctx, llmName, req)
	}
	// Return default mock response
	return map[string]interface{}{
		"classification": map[string]interface{}{
			"subject_candidate": "math",
			"confidence":        0.95,
		},
		"quality": map[string]interface{}{
			"is_acceptable": true,
		},
	}, nil
}

func (m *mockLLMClient) Parse(ctx context.Context, llmName string, req interface{}) (interface{}, error) {
	if m.parseFunc != nil {
		return m.parseFunc(ctx, llmName, req)
	}
	// Return default mock response
	return map[string]interface{}{
		"task": map[string]interface{}{
			"text": "Solve 2+2",
		},
		"items": []interface{}{
			map[string]interface{}{
				"text": "2+2=?",
			},
		},
	}, nil
}

func (m *mockLLMClient) Hint(ctx context.Context, llmName string, req interface{}) (interface{}, error) {
	if m.hintFunc != nil {
		return m.hintFunc(ctx, llmName, req)
	}
	// Return default mock response
	return map[string]interface{}{
		"hints": []interface{}{
			map[string]interface{}{
				"text":  "First hint",
				"level": 1,
			},
			map[string]interface{}{
				"text":  "Second hint",
				"level": 2,
			},
		},
	}, nil
}

func (m *mockLLMClient) CheckSolution(ctx context.Context, llmName string, req interface{}) (interface{}, error) {
	if m.checkFunc != nil {
		return m.checkFunc(ctx, llmName, req)
	}
	// Return default mock response
	return map[string]interface{}{
		"is_correct": true,
		"feedback":   "Correct answer!",
	}, nil
}
