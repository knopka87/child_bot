package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"child-bot/api/internal/api/router"
	"child-bot/api/internal/config"
	"child-bot/api/internal/llm"
	"child-bot/api/internal/store"
)

// Test configuration
type E2ETestConfig struct {
	DatabaseURL  string
	LLMProxyURL  string
	LLMName      string
	UseRealLLM   bool // If false, use mock LLM
	TestPlatform string
}

// loadE2EConfig loads configuration for E2E tests
func loadE2EConfig(t *testing.T) *E2ETestConfig {
	t.Helper()

	cfg := &E2ETestConfig{
		DatabaseURL:  os.Getenv("TEST_DATABASE_URL"),
		LLMProxyURL:  os.Getenv("LLM_PROXY_URL"),
		LLMName:      getEnvOrDefault("LLM_NAME", "gpt-4"),
		UseRealLLM:   os.Getenv("USE_REAL_LLM") == "true",
		TestPlatform: "test",
	}

	if cfg.DatabaseURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// setupE2EServer creates a test server with all dependencies
func setupE2EServer(t *testing.T, cfg *E2ETestConfig) (*httptest.Server, *sql.DB) {
	t.Helper()

	// Connect to database
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	// Create store
	st := store.NewStore(db)

	// Create LLM client
	var llmClient *llm.Client
	if cfg.UseRealLLM {
		llmClient = llm.NewClient(cfg.LLMProxyURL)
	} else {
		// TODO: Use mock LLM client for faster tests
		llmClient = llm.NewClient(cfg.LLMProxyURL)
	}

	// Create router
	r := router.New(&router.Dependencies{
		Store:      st,
		LLMClient:  llmClient,
		Config:     &config.Config{},
		DefaultLLM: cfg.LLMName,
	})

	// Create test server
	server := httptest.NewServer(r)
	t.Cleanup(func() {
		server.Close()
	})

	return server, db
}

// createTestProfile creates a test profile and returns ID
func createE2ETestProfile(t *testing.T, db *sql.DB, platformID string) string {
	t.Helper()

	var profileID string
	query := `
		INSERT INTO child_profiles (display_name, platform_id, platform_user_id, grade, avatar_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	userID := fmt.Sprintf("e2e_user_%d", time.Now().UnixNano())
	err := db.QueryRow(query, "E2E Test User", platformID, userID, 5, "avatar1").Scan(&profileID)
	if err != nil {
		t.Fatalf("failed to create test profile: %v", err)
	}

	t.Cleanup(func() {
		db.Exec("DELETE FROM child_profiles WHERE id = $1", profileID)
	})

	return profileID
}

// makeE2ERequest makes an HTTP request to the test server
func makeE2ERequest(t *testing.T, server *httptest.Server, method, path string, body interface{}, platformID, profileID string) *http.Response {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, server.URL+path, reqBody)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if platformID != "" {
		req.Header.Set("X-Platform-ID", platformID)
	}
	if profileID != "" {
		req.Header.Set("X-Child-Profile-ID", profileID)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	return resp
}

// decodeE2EResponse decodes JSON response
func decodeE2EResponse(t *testing.T, resp *http.Response, target interface{}) {
	t.Helper()

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("failed to decode response: %v\nBody: %s", err, string(body))
	}
}

// TestE2E_HealthCheck tests the health check endpoint
func TestE2E_HealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	cfg := loadE2EConfig(t)
	server, _ := setupE2EServer(t, cfg)

	resp := makeE2ERequest(t, server, http.MethodGet, "/health", nil, "", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	decodeE2EResponse(t, resp, &result)

	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", result["status"])
	}

	t.Log("Health check passed")
}

// TestE2E_AttemptFlow_Help tests the complete help attempt flow
func TestE2E_AttemptFlow_Help(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	cfg := loadE2EConfig(t)
	server, db := setupE2EServer(t, cfg)

	// Create test profile
	profileID := createE2ETestProfile(t, db, cfg.TestPlatform)

	// Step 1: Create attempt
	t.Log("Step 1: Creating help attempt")
	createReq := map[string]string{
		"child_profile_id": profileID,
		"type":             "help",
	}
	resp := makeE2ERequest(t, server, http.MethodPost, "/attempts", createReq, cfg.TestPlatform, profileID)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create attempt failed with status %d", resp.StatusCode)
	}

	var createResp map[string]interface{}
	decodeE2EResponse(t, resp, &createResp)

	attemptID, ok := createResp["attempt_id"].(string)
	if !ok || attemptID == "" {
		t.Fatalf("expected attempt_id in response, got %v", createResp)
	}
	t.Logf("Created attempt: %s", attemptID)

	// Cleanup
	defer db.Exec("DELETE FROM attempts WHERE id = $1", attemptID)

	// Step 2: Upload task image
	t.Log("Step 2: Uploading task image")
	sampleImage := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
	uploadReq := map[string]string{
		"image_type":   "task",
		"image_base64": sampleImage,
	}
	resp = makeE2ERequest(t, server, http.MethodPost, "/attempts/"+attemptID+"/images", uploadReq, cfg.TestPlatform, profileID)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("upload image failed with status %d", resp.StatusCode)
	}
	t.Log("Image uploaded successfully")

	// Step 3: Process attempt (skip if not using real LLM)
	if cfg.UseRealLLM {
		t.Log("Step 3: Processing attempt with LLM")
		resp = makeE2ERequest(t, server, http.MethodPost, "/attempts/"+attemptID+"/process", nil, cfg.TestPlatform, profileID)
		if resp.StatusCode != http.StatusAccepted {
			t.Fatalf("process attempt failed with status %d", resp.StatusCode)
		}

		// Wait for processing (with timeout)
		t.Log("Waiting for processing to complete...")
		time.Sleep(30 * time.Second)

		// Step 4: Get result
		t.Log("Step 4: Getting result")
		resp = makeE2ERequest(t, server, http.MethodGet, "/attempts/"+attemptID+"/result", nil, cfg.TestPlatform, profileID)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("get result failed with status %d", resp.StatusCode)
		}

		var result map[string]interface{}
		decodeE2EResponse(t, resp, &result)
		t.Logf("Result status: %v", result["status"])

		// Step 5: Get first hint
		t.Log("Step 5: Getting first hint")
		resp = makeE2ERequest(t, server, http.MethodPost, "/attempts/"+attemptID+"/next-hint", nil, cfg.TestPlatform, profileID)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("get hint failed with status %d", resp.StatusCode)
		}

		var hintResp map[string]interface{}
		decodeE2EResponse(t, resp, &hintResp)
		t.Logf("Hint: %v", hintResp)
	} else {
		t.Log("Skipping LLM processing (USE_REAL_LLM not set)")
	}

	// Step 6: Delete attempt
	t.Log("Step 6: Deleting attempt")
	resp = makeE2ERequest(t, server, http.MethodDelete, "/attempts/"+attemptID, nil, cfg.TestPlatform, profileID)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete attempt failed with status %d", resp.StatusCode)
	}

	t.Log("Help attempt flow completed successfully")
}

// TestE2E_AttemptFlow_Check tests the complete check attempt flow
func TestE2E_AttemptFlow_Check(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	cfg := loadE2EConfig(t)
	server, db := setupE2EServer(t, cfg)

	profileID := createE2ETestProfile(t, db, cfg.TestPlatform)

	// Step 1: Create check attempt
	t.Log("Step 1: Creating check attempt")
	createReq := map[string]string{
		"child_profile_id": profileID,
		"type":             "check",
	}
	resp := makeE2ERequest(t, server, http.MethodPost, "/attempts", createReq, cfg.TestPlatform, profileID)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create attempt failed with status %d", resp.StatusCode)
	}

	var createResp map[string]interface{}
	decodeE2EResponse(t, resp, &createResp)

	attemptID := createResp["attempt_id"].(string)
	t.Logf("Created attempt: %s", attemptID)

	defer db.Exec("DELETE FROM attempts WHERE id = $1", attemptID)

	// Step 2: Upload task image
	t.Log("Step 2: Uploading task image")
	sampleImage := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
	uploadReq := map[string]string{
		"image_type":   "task",
		"image_base64": sampleImage,
	}
	resp = makeE2ERequest(t, server, http.MethodPost, "/attempts/"+attemptID+"/images", uploadReq, cfg.TestPlatform, profileID)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("upload task image failed with status %d", resp.StatusCode)
	}

	// Step 3: Upload answer image
	t.Log("Step 3: Uploading answer image")
	uploadReq["image_type"] = "answer"
	resp = makeE2ERequest(t, server, http.MethodPost, "/attempts/"+attemptID+"/images", uploadReq, cfg.TestPlatform, profileID)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("upload answer image failed with status %d", resp.StatusCode)
	}

	// Step 4: Process check (skip if not using real LLM)
	if cfg.UseRealLLM {
		t.Log("Step 4: Processing check with LLM")
		resp = makeE2ERequest(t, server, http.MethodPost, "/attempts/"+attemptID+"/process", nil, cfg.TestPlatform, profileID)
		if resp.StatusCode != http.StatusAccepted {
			t.Fatalf("process attempt failed with status %d", resp.StatusCode)
		}

		// Wait for processing
		t.Log("Waiting for processing to complete...")
		time.Sleep(30 * time.Second)

		// Get result
		t.Log("Step 5: Getting check result")
		resp = makeE2ERequest(t, server, http.MethodGet, "/attempts/"+attemptID+"/result", nil, cfg.TestPlatform, profileID)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("get result failed with status %d", resp.StatusCode)
		}

		var result map[string]interface{}
		decodeE2EResponse(t, resp, &result)
		t.Logf("Check result: %v", result)
	} else {
		t.Log("Skipping LLM processing (USE_REAL_LLM not set)")
	}

	t.Log("Check attempt flow completed successfully")
}

// TestE2E_ErrorHandling tests error scenarios
func TestE2E_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	cfg := loadE2EConfig(t)
	server, db := setupE2EServer(t, cfg)

	profileID := createE2ETestProfile(t, db, cfg.TestPlatform)

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}{
		{
			name:   "missing platform ID",
			method: http.MethodPost,
			path:   "/attempts",
			body: map[string]string{
				"child_profile_id": profileID,
				"type":             "help",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid attempt type",
			method: http.MethodPost,
			path:   "/attempts",
			body: map[string]string{
				"child_profile_id": profileID,
				"type":             "invalid",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "attempt not found",
			method:         http.MethodGet,
			path:           "/attempts/00000000-0000-0000-0000-000000000000/result",
			body:           nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid UUID",
			method:         http.MethodGet,
			path:           "/attempts/invalid-uuid/result",
			body:           nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var platformID string
			if tt.name != "missing platform ID" {
				platformID = cfg.TestPlatform
			}

			resp := makeE2ERequest(t, server, tt.method, tt.path, tt.body, platformID, profileID)
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("expected status %d, got %d\nBody: %s", tt.expectedStatus, resp.StatusCode, string(body))
			}
		})
	}
}

// TestE2E_ConcurrentRequests tests handling of concurrent requests
func TestE2E_ConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	cfg := loadE2EConfig(t)
	server, db := setupE2EServer(t, cfg)

	profileID := createE2ETestProfile(t, db, cfg.TestPlatform)

	// Create multiple attempts concurrently
	concurrency := 10
	results := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			createReq := map[string]string{
				"child_profile_id": profileID,
				"type":             "help",
			}
			resp := makeE2ERequest(t, server, http.MethodPost, "/attempts", createReq, cfg.TestPlatform, profileID)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				results <- fmt.Errorf("unexpected status: %d", resp.StatusCode)
				return
			}

			var createResp map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
				results <- err
				return
			}

			attemptID := createResp["attempt_id"].(string)
			defer db.Exec("DELETE FROM attempts WHERE id = $1", attemptID)

			results <- nil
		}()
	}

	// Wait for all goroutines
	for i := 0; i < concurrency; i++ {
		if err := <-results; err != nil {
			t.Errorf("concurrent request %d failed: %v", i, err)
		}
	}

	t.Logf("Successfully handled %d concurrent requests", concurrency)
}
