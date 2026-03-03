package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"child-bot/api/internal/llmclient"
	"child-bot/api/internal/service"
	"child-bot/api/internal/store"
	"child-bot/api/internal/v2/telegram"
)

const (
	answersDir = "testdata/answers" // Пары task_X + answer_X
)

// TaskAnswerPair holds a pair of task and answer images
type TaskAnswerPair struct {
	Name       string // e.g., "homework1" (without prefix)
	TaskPath   string // path to task_homework1.jpg
	AnswerPath string // path to answer_homework1.jpg
}

// CheckAnswerResult extends TestResult with check-specific fields
type CheckAnswerResult struct {
	*TestResult
	TaskImagePath   string `json:"task_image_path"`
	AnswerImagePath string `json:"answer_image_path"`
	CheckResult     string `json:"check_result"` // "correct", "incorrect", "error"
	CheckFeedback   string `json:"check_feedback,omitempty"`
}

// TestE2E_CheckAnswer tests the complete answer checking flow:
// 1. Set grade
// 2. Upload task image
// 3. Confirm parse result (click "parse_yes")
// 4. Click "ready_solution" to enter answer mode
// 5. Upload answer image
// 6. Get verification result (correct/incorrect)
func TestE2E_CheckAnswer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode (requires real LLM)")
	}

	// Load config from environment
	cfg := loadTestConfig(t)

	// Find all task-answer pairs
	pairs := findTaskAnswerPairs(t, answersDir)
	if len(pairs) == 0 {
		t.Skip("No task-answer pairs found in testdata/answers/. Please add pairs (task_X.jpg + answer_X.jpg)")
	}

	t.Logf("Found %d task-answer pairs for check flow testing", len(pairs))

	// Setup real LLM client (shared across all subtests)
	llmClient := llmclient.New(cfg.LLMProxyURL)

	// Setup real store (test database)
	st := setupTestStore(t, cfg.DatabaseURL)

	// Configure templates directory for tests
	telegram.SetTemplatesDir("../../internal/v2/templates")
	telegram.ResetTemplatesCache()

	// Run test for each pair
	for i, pair := range pairs {
		chatID := int64(255509624 + i) // Different from hint flow tests
		userID := int64(88889001 + i)

		t.Run(pair.Name, func(t *testing.T) {
			runCheckAnswerForPair(t, cfg, llmClient, st, pair, chatID, userID)
		})
	}
}

// findTaskAnswerPairs finds matching task_X and answer_X pairs in the directory
func findTaskAnswerPairs(t *testing.T, dir string) []TaskAnswerPair {
	t.Helper()

	// Find all images
	images := findImagesInDir(t, dir)

	// Separate task and answer images
	tasks := make(map[string]string)   // name -> path
	answers := make(map[string]string) // name -> path

	for _, img := range images {
		base := filepath.Base(img)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)

		if strings.HasPrefix(name, "task_") {
			key := strings.TrimPrefix(name, "task_")
			tasks[key] = img
		} else if strings.HasPrefix(name, "answer_") {
			key := strings.TrimPrefix(name, "answer_")
			answers[key] = img
		}
	}

	// Find matching pairs
	var pairs []TaskAnswerPair
	for name, taskPath := range tasks {
		if answerPath, ok := answers[name]; ok {
			pairs = append(pairs, TaskAnswerPair{
				Name:       name,
				TaskPath:   taskPath,
				AnswerPath: answerPath,
			})
			t.Logf("Found pair: task_%s + answer_%s", name, name)
		} else {
			t.Logf("Warning: task_%s has no matching answer_%s", name, name)
		}
	}

	// Check for orphaned answers
	for name := range answers {
		if _, ok := tasks[name]; !ok {
			t.Logf("Warning: answer_%s has no matching task_%s", name, name)
		}
	}

	// Sort by name for consistent order
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Name < pairs[j].Name
	})

	return pairs
}

// runCheckAnswerForPair runs the check answer test for a single task-answer pair
func runCheckAnswerForPair(t *testing.T, cfg *TestConfig, llmClient *llmclient.Client, st *store.Store, pair TaskAnswerPair, chatID, userID int64) {
	mockBot := NewMockBot()

	router := &telegram.Router{
		Bot:        mockBot,
		LlmManager: service.NewLlmManager(cfg.LLMName),
		LLMClient:  llmClient,
		Store:      st,
	}

	taskImage := loadTestImage(t, pair.TaskPath)
	answerImage := loadTestImage(t, pair.AnswerPath)

	result := &CheckAnswerResult{
		TestResult:      NewTestResult(chatID),
		TaskImagePath:   pair.TaskPath,
		AnswerImagePath: pair.AnswerPath,
	}

	// Step 0: Initialize - trigger grade selection and select grade 3
	router.HandleUpdate(makePhotoUpdate(chatID, userID, "dummy"), cfg.LLMName)
	waitForTextMessage(t, mockBot, 5*time.Second)
	mockBot.ClearMessages()

	msgID := mockBot.MessageID
	router.HandleUpdate(makeGradeCallback(chatID, userID, msgID, "grade3"), cfg.LLMName)
	waitForTextMessage(t, mockBot, 5*time.Second)
	t.Logf("Grade set, response: %s", truncateText(mockBot.LastTextMessage().Text, 100))
	mockBot.ClearMessages()

	// Step 1: Send task photo
	t.Logf("Step 1: Sending task photo: %s", filepath.Base(pair.TaskPath))
	taskFileID := fmt.Sprintf("test-task-%d", chatID)
	mockBot.AddFile(taskFileID, taskImage, "photos/task_"+pair.Name+".jpg")

	router.HandleUpdate(makePhotoUpdate(chatID, userID, taskFileID), cfg.LLMName)

	// Wait for parse result
	t.Log("Waiting for Detect + Parse...")
	parseMsg := waitForMessageWithLogging(t, mockBot, 3*time.Minute, func(msg *CapturedMessage) bool {
		return containsButton(msg, "parse_yes")
	})
	t.Logf("Parse completed: %s", truncateText(parseMsg.Text, 150))

	// Step 2: Confirm parse
	t.Log("Step 2: Confirming parse result")
	msgID = mockBot.MessageID
	router.HandleUpdate(makeCallbackUpdate(chatID, userID, msgID, "parse_yes"), cfg.LLMName)

	// Wait for hint (we need to get past hints to ready_solution)
	t.Log("Waiting for first hint...")
	waitForMessage(t, mockBot, 3*time.Minute, func(msg *CapturedMessage) bool {
		return strings.Contains(msg.Text, "подсказка") || strings.Contains(msg.Text, "Подсказка")
	})

	// Step 3: Click "ready_solution" (Проверь мой ответ)
	t.Log("Step 3: Clicking 'ready_solution'")
	lastMsg := mockBot.LastTextMessage()
	if !containsButton(lastMsg, "ready_solution") {
		t.Fatalf("ready_solution button not found in message: %s", truncateText(lastMsg.Text, 200))
	}

	msgID = mockBot.MessageID
	router.HandleUpdate(makeCallbackUpdate(chatID, userID, msgID, "ready_solution"), cfg.LLMName)

	// Wait for "send your solution" prompt
	t.Log("Waiting for solution prompt...")
	solutionPrompt := waitForMessage(t, mockBot, 30*time.Second, func(msg *CapturedMessage) bool {
		return strings.Contains(msg.Text, "фото") ||
			strings.Contains(msg.Text, "решени") ||
			strings.Contains(msg.Text, "ответ")
	})
	t.Logf("Solution prompt: %s", truncateText(solutionPrompt.Text, 100))

	// Step 4: Send answer photo
	t.Logf("Step 4: Sending answer photo: %s", filepath.Base(pair.AnswerPath))
	answerFileID := fmt.Sprintf("test-answer-%d", chatID)
	mockBot.AddFile(answerFileID, answerImage, "photos/answer_"+pair.Name+".jpg")

	// Clear messages before sending answer to only capture check result
	mockBot.ClearMessages()

	router.HandleUpdate(makePhotoUpdate(chatID, userID, answerFileID), cfg.LLMName)

	// Step 5: Wait for check result
	t.Log("Step 5: Waiting for check result (may take 30-90 seconds)...")
	checkMsg := waitForMessageWithLogging(t, mockBot, 3*time.Minute, func(msg *CapturedMessage) bool {
		text := strings.ToLower(msg.Text)
		return strings.Contains(text, "верно") ||
			strings.Contains(text, "правильно") ||
			strings.Contains(text, "молодец") ||
			strings.Contains(text, "ошибк") ||
			strings.Contains(text, "неверно") ||
			strings.Contains(text, "поправ") ||
			strings.Contains(text, "отлично") ||
			strings.Contains(text, "не видно") ||
			strings.Contains(text, "не удалось")
	})

	// Analyze check result
	// IMPORTANT: Check negative patterns FIRST because "неверно" contains "верно"
	checkText := strings.ToLower(checkMsg.Text)
	if strings.Contains(checkText, "не видно") || strings.Contains(checkText, "не удалось") {
		result.CheckResult = "no_visible_answer"
	} else if strings.Contains(checkText, "неверно") ||
		strings.Contains(checkText, "неправильно") ||
		strings.Contains(checkText, "ошибк") ||
		strings.Contains(checkText, "поправ") ||
		strings.Contains(checkText, "почти получилось") ||
		strings.Contains(checkText, "что можно поправить") {
		result.CheckResult = "incorrect"
	} else if strings.Contains(checkText, "всё верно") ||
		strings.Contains(checkText, "правильно") ||
		strings.Contains(checkText, "молодец") ||
		strings.Contains(checkText, "отлично") {
		result.CheckResult = "correct"
	} else {
		result.CheckResult = "unknown"
	}
	result.CheckFeedback = checkMsg.Text

	t.Logf("Check result: %s", result.CheckResult)
	t.Logf("Check feedback: %s", truncateText(checkMsg.Text, 200))

	result.EndTime = time.Now()
	result.Success = true

	// Fetch timeline events from database
	fetchTimelineEvents(t, st, chatID, result.TestResult)

	// Save results
	resultsPath := filepath.Join("results", fmt.Sprintf("check_%s_%d.json", pair.Name, time.Now().Unix()))
	saveCheckResults(t, resultsPath, result)

	t.Logf("Test completed for pair %s!", pair.Name)
	t.Logf("  - Check result: %s", result.CheckResult)
	t.Logf("  - Total duration: %v", result.Duration())
	t.Logf("  - Results saved to: %s", resultsPath)
}

// saveCheckResults saves check answer results to a JSON file
func saveCheckResults(t *testing.T, path string, result *CheckAnswerResult) {
	t.Helper()

	// Use the same save function but with extended result
	// We can marshal CheckAnswerResult directly since it embeds TestResult
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
