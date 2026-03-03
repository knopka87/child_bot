package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"child-bot/api/internal/llmclient"
	"child-bot/api/internal/service"
	"child-bot/api/internal/store"
	"child-bot/api/internal/v2/telegram"
)

const (
	tasksDir = "testdata/tasks" // Фото заданий
)

// TestE2E_HintFlow tests the complete hint flow for all task images:
// 1. Set grade
// 2. Upload task image
// 3. Confirm parse result (click "parse_yes")
// 4. Request hints (click "hint_next") until no more available
// 5. Save all results to file
func TestE2E_HintFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode (requires real LLM)")
	}

	// Load config from environment
	cfg := loadTestConfig(t)

	// Find all task images
	taskImages := findImagesInDir(t, tasksDir)
	if len(taskImages) == 0 {
		t.Skip("No task images found in testdata/tasks/. Please add task images (*.jpg, *.png)")
	}

	t.Logf("Found %d task images for hint flow testing", len(taskImages))

	// Setup real LLM client (shared across all subtests)
	llmClient := llmclient.New(cfg.LLMProxyURL)

	// Setup real store (test database)
	st := setupTestStore(t, cfg.DatabaseURL)

	// Configure templates directory for tests (relative to api/test/e2e/)
	telegram.SetTemplatesDir("../../internal/v2/templates")
	telegram.ResetTemplatesCache()

	// Run test for each image
	for i, imagePath := range taskImages {
		imageName := filepath.Base(imagePath)
		chatID := int64(255509524)
		userID := int64(88888001 + i)

		t.Run(imageName, func(t *testing.T) {
			runHintFlowForImage(t, cfg, llmClient, st, imagePath, chatID, userID)
		})
	}
}

// runHintFlowForImage runs the hint flow test for a single image
func runHintFlowForImage(t *testing.T, cfg *TestConfig, llmClient *llmclient.Client, st *store.Store, imagePath string, chatID, userID int64) {
	mockBot := NewMockBot()

	router := &telegram.Router{
		Bot:        mockBot,
		LlmManager: service.NewLlmManager(cfg.LLMName),
		LLMClient:  llmClient,
		Store:      st,
	}

	testImage := loadTestImage(t, imagePath)
	imageName := filepath.Base(imagePath)

	result := NewTestResult(chatID)
	result.ImagePath = imagePath

	// Step 0: Initialize - trigger grade selection and select grade 3
	router.HandleUpdate(makePhotoUpdate(chatID, userID, "dummy"), cfg.LLMName) // triggers grade selection
	waitForTextMessage(t, mockBot, 5*time.Second)
	mockBot.ClearMessages()

	msgID := mockBot.MessageID
	router.HandleUpdate(makeGradeCallback(chatID, userID, msgID, "grade3"), cfg.LLMName)
	waitForTextMessage(t, mockBot, 5*time.Second)
	t.Logf("Grade set, response: %s", truncateText(mockBot.LastTextMessage().Text, 100))
	mockBot.ClearMessages()

	// Step 1: Send photo
	t.Logf("Step 1: Sending task photo: %s", imageName)
	fileID := fmt.Sprintf("test-file-%d", chatID)
	mockBot.AddFile(fileID, testImage, "photos/"+imageName)

	router.HandleUpdate(makePhotoUpdate(chatID, userID, fileID), cfg.LLMName)

	// Wait for debounce + LLM processing (Detect + Parse)
	t.Log("Waiting for Detect + Parse (may take 30-60 seconds)...")
	parseMsg := waitForMessageWithLogging(t, mockBot, 3*time.Minute, func(msg *CapturedMessage) bool {
		return containsButton(msg, "parse_yes")
	})
	t.Log("Parse completed, found parse_yes button")

	// Extract parsed text for logging
	if parseMsg != nil {
		t.Logf("Parsed text: %s", truncateText(parseMsg.Text, 200))
	}

	// Step 2: Confirm parse
	t.Log("Step 2: Confirming parse result")
	msgID = mockBot.MessageID
	router.HandleUpdate(makeCallbackUpdate(chatID, userID, msgID, "parse_yes"), cfg.LLMName)

	// Wait for first hint (LLM Hint call)
	t.Log("Waiting for Hint L1...")
	hintMsg := waitForMessage(t, mockBot, 3*time.Minute, func(msg *CapturedMessage) bool {
		return strings.Contains(msg.Text, "подсказка") || strings.Contains(msg.Text, "Подсказка")
	})
	t.Logf("Hint 1: %s", truncateText(hintMsg.Text, 150))

	// Step 4: Request hints until exhausted
	hintCount := 1
	for i := 0; i < 10; i++ { // safety limit
		lastMsg := mockBot.LastTextMessage()
		if lastMsg == nil {
			t.Log("No text message found")
			break
		}

		if !containsButton(lastMsg, "hint_next") {
			t.Logf("No more hints available after %d hints", hintCount)
			break
		}

		t.Logf("Step 4.%d: Requesting next hint", i+1)
		msgID = mockBot.MessageID
		router.HandleUpdate(makeCallbackUpdate(chatID, userID, msgID, "hint_next"), cfg.LLMName)

		// Wait for hint response
		hintMsg = waitForMessage(t, mockBot, 3*time.Minute, func(msg *CapturedMessage) bool {
			return strings.Contains(msg.Text, "подсказка") ||
				strings.Contains(msg.Text, "Подсказка") ||
				strings.Contains(msg.Text, "показаны") ||
				strings.Contains(msg.Text, "закончились")
		})
		hintCount++
		t.Logf("Hint %d: %s", hintCount, truncateText(hintMsg.Text, 150))
	}

	result.EndTime = time.Now()
	result.PipelineOK = true
	result.HintsCount = hintCount

	// Fetch timeline events from database
	fetchTimelineEvents(t, st, chatID, result)

	// Save results with image name in filename
	safeImageName := strings.TrimSuffix(imageName, filepath.Ext(imageName))
	resultsPath := filepath.Join("results", fmt.Sprintf("%s_%d.json", safeImageName, time.Now().Unix()))
	saveTestResults(t, resultsPath, result)

	t.Logf("Test completed for %s!", imageName)
	t.Logf("  - Total hints received: %d", hintCount)
	t.Logf("  - Total duration: %v", result.Duration())
	t.Logf("  - Results saved to: %s", resultsPath)
}

// findImagesInDir finds all images in a directory
func findImagesInDir(t *testing.T, dir string) []string {
	t.Helper()

	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	var images []string

	patterns := []string{
		filepath.Join(dir, "*.jpg"),
		filepath.Join(dir, "*.jpeg"),
		filepath.Join(dir, "*.png"),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			t.Logf("Warning: glob error for %s: %v", pattern, err)
			continue
		}
		images = append(images, matches...)
	}

	// Sort for consistent order
	for i := 0; i < len(images); i++ {
		for j := i + 1; j < len(images); j++ {
			if images[i] > images[j] {
				images[i], images[j] = images[j], images[i]
			}
		}
	}

	for _, img := range images {
		t.Logf("Found image: %s", img)
	}

	return images
}

// setupTestStore creates a store connected to the test database
func setupTestStore(t *testing.T, databaseURL string) *store.Store {
	t.Helper()

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return store.NewStore(db)
}

// waitForTextMessage waits for any text message
func waitForTextMessage(t *testing.T, bot *MockBot, timeout time.Duration) *CapturedMessage {
	t.Helper()

	return waitForMessage(t, bot, timeout, func(msg *CapturedMessage) bool {
		return msg.Type == "text"
	})
}

// truncateText truncates text to maxLen with ellipsis
func truncateText(text string, maxLen int) string {
	text = strings.ReplaceAll(text, "\n", " ")
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

// TestE2E_GradeSelection tests just the grade selection flow
func TestE2E_GradeSelection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	cfg := loadTestConfig(t)
	llmClient := llmclient.New(cfg.LLMProxyURL)
	mockBot := NewMockBot()
	st := setupTestStore(t, cfg.DatabaseURL)

	router := &telegram.Router{
		Bot:        mockBot,
		LlmManager: service.NewLlmManager(cfg.LLMName),
		LLMClient:  llmClient,
		Store:      st,
	}

	chatID := int64(99999002)
	userID := int64(88888002)

	grades := []string{"grade1", "grade2", "grade3", "grade4"}

	for _, grade := range grades {
		mockBot.ClearMessages()
		msgID := mockBot.MessageID

		t.Logf("Testing grade selection: %s", grade)
		router.HandleUpdate(makeGradeCallback(chatID, userID, msgID, grade), cfg.LLMName)

		msg := waitForTextMessage(t, mockBot, 5*time.Second)
		if msg == nil {
			t.Errorf("No response for grade %s", grade)
			continue
		}

		t.Logf("Response for %s: %s", grade, truncateText(msg.Text, 80))
	}
}

// TestE2E_StartCommand tests the /start command
func TestE2E_StartCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	cfg := loadTestConfig(t)
	llmClient := llmclient.New(cfg.LLMProxyURL)
	mockBot := NewMockBot()
	st := setupTestStore(t, cfg.DatabaseURL)

	router := &telegram.Router{
		Bot:        mockBot,
		LlmManager: service.NewLlmManager(cfg.LLMName),
		LLMClient:  llmClient,
		Store:      st,
	}

	chatID := int64(99999003)
	userID := int64(88888003)

	t.Log("Testing /start command")
	router.HandleUpdate(makeCommandUpdate(chatID, userID, "start"), cfg.LLMName)

	msg := waitForTextMessage(t, mockBot, 5*time.Second)
	if msg == nil {
		t.Error("No response for /start command")
		return
	}

	t.Logf("/start response: %s", truncateText(msg.Text, 80))

	if !containsButton(msg, "grade1") && !containsButton(msg, "grade3") {
		t.Logf("Warning: Grade buttons not found in /start response")
	}
}

// fetchTimelineEvents loads timeline events from DB and adds them to result
func fetchTimelineEvents(t *testing.T, st *store.Store, chatID int64, result *TestResult) {
	t.Helper()

	ctx := context.Background()

	// Get session by chatID
	session, err := st.FindSession(ctx, chatID)
	if err != nil {
		t.Logf("Warning: failed to find session for chat %d: %v", chatID, err)
		return
	}
	if session.SessionID == "" {
		t.Logf("Warning: no session found for chat %d", chatID)
		return
	}

	result.SessionID = session.SessionID

	// Fetch timeline events by session_id
	events, err := st.FindALLHistoryBySID(ctx, session.SessionID)
	if err != nil {
		t.Logf("Warning: failed to fetch timeline events: %v", err)
		return
	}

	// Convert to TimelineRecord
	var records []TimelineRecord
	for _, e := range events {
		record := TimelineRecord{
			TaskSessionID: e.TaskSessionID,
			Direction:     e.Direction,
			EventType:     e.EventType,
			Provider:      e.Provider,
			OK:            e.OK,
			LatencyMS:     e.LatencyMS,
			Text:          e.Text,
			InputPayload:  e.InputPayload,
			OutputPayload: e.OutputPayload,
			CreatedAt:     e.CreatedAt,
		}
		if e.Error != nil {
			record.Error = e.Error.Error()
		}
		records = append(records, record)
	}

	// P2.3: Sanitize to remove base64 images from logs
	result.TimelineEvents = SanitizeTimelineEvents(records)

	t.Logf("Loaded %d timeline events from database", len(events))
}
