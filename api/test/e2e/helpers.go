package e2e

import (
	"os"
	"strings"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// makePhotoUpdate creates a tgbotapi.Update simulating a photo message
func makePhotoUpdate(chatID int64, userID int64, fileID string) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			MessageID: int(time.Now().UnixNano() % 1000000),
			From: &tgbotapi.User{
				ID:        userID,
				FirstName: "TestUser",
			},
			Chat: &tgbotapi.Chat{
				ID:   chatID,
				Type: "private",
			},
			Date: int(time.Now().Unix()),
			Photo: []tgbotapi.PhotoSize{
				{
					FileID:       fileID + "_small",
					FileUniqueID: fileID + "_small_unique",
					Width:        100,
					Height:       100,
				},
				{
					FileID:       fileID,
					FileUniqueID: fileID + "_unique",
					Width:        800,
					Height:       600,
				},
			},
		},
	}
}

// makeCallbackUpdate creates a tgbotapi.Update simulating a callback button click
func makeCallbackUpdate(chatID int64, userID int64, msgID int, data string) tgbotapi.Update {
	return tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			ID: "callback_" + data,
			From: &tgbotapi.User{
				ID:        userID,
				FirstName: "TestUser",
			},
			Message: &tgbotapi.Message{
				MessageID: msgID,
				Chat: &tgbotapi.Chat{
					ID:   chatID,
					Type: "private",
				},
			},
			ChatInstance: "test_instance",
			Data:         data,
		},
	}
}

// makeGradeCallback creates a callback update for grade selection
func makeGradeCallback(chatID int64, userID int64, msgID int, grade string) tgbotapi.Update {
	return makeCallbackUpdate(chatID, userID, msgID, grade)
}

// loadTestImage loads a test image from the testdata directory
func loadTestImage(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to load test image %s: %v", path, err)
	}
	return data
}

// containsButton checks if a message contains a button with the given callback data
func containsButton(msg *CapturedMessage, data string) bool {
	if msg == nil || msg.Buttons == nil {
		return false
	}
	for _, row := range msg.Buttons {
		for _, btn := range row {
			if btn.Data == data {
				return true
			}
		}
	}
	return false
}

// containsText checks if a message contains the given text
func containsText(msg *CapturedMessage, text string) bool {
	if msg == nil {
		return false
	}
	return strings.Contains(msg.Text, text)
}

// WaitResult holds the result of waiting for a message
type WaitResult struct {
	Message *CapturedMessage
	Index   int
	Err     error
}

// waitForMessage waits for a message matching the predicate
func waitForMessage(t *testing.T, bot *MockBot, timeout time.Duration, predicate func(*CapturedMessage) bool) *CapturedMessage {
	t.Helper()

	deadline := time.Now().Add(timeout)
	lastCount := 0

	for time.Now().Before(deadline) {
		messages := bot.GetMessages()

		// Check new messages since last check
		for i := lastCount; i < len(messages); i++ {
			if predicate(&messages[i]) {
				return &messages[i]
			}
		}
		lastCount = len(messages)

		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("Timeout waiting for message (waited %v)", timeout)
	return nil
}

// waitForMessageWithLogging is like waitForMessage but logs all received messages for debugging
func waitForMessageWithLogging(t *testing.T, bot *MockBot, timeout time.Duration, predicate func(*CapturedMessage) bool) *CapturedMessage {
	t.Helper()

	deadline := time.Now().Add(timeout)
	lastCount := 0
	logInterval := 10 * time.Second
	lastLog := time.Now()

	for time.Now().Before(deadline) {
		messages := bot.GetMessages()

		// Log periodically
		if time.Since(lastLog) > logInterval {
			t.Logf("DEBUG: Total messages so far: %d", len(messages))
			for i, msg := range messages {
				t.Logf("DEBUG: Message[%d] type=%s text=%.100s buttons=%v", i, msg.Type, msg.Text, len(msg.Buttons) > 0)
			}
			lastLog = time.Now()
		}

		// Check new messages since last check
		for i := lastCount; i < len(messages); i++ {
			t.Logf("DEBUG: New message[%d] type=%s text=%.150s", i, messages[i].Type, messages[i].Text)
			if predicate(&messages[i]) {
				return &messages[i]
			}
		}
		lastCount = len(messages)

		time.Sleep(100 * time.Millisecond)
	}

	// Final debug dump before failing
	messages := bot.GetMessages()
	t.Logf("DEBUG FINAL: Total messages: %d", len(messages))
	for i, msg := range messages {
		t.Logf("DEBUG FINAL: Message[%d] type=%s text=%.200s", i, msg.Type, msg.Text)
	}

	t.Fatalf("Timeout waiting for message (waited %v)", timeout)
	return nil
}

// waitForMessageCount waits until at least N messages are received
func waitForMessageCount(t *testing.T, bot *MockBot, count int, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if bot.GetMessageCount() >= count {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("Timeout waiting for %d messages (got %d)", count, bot.GetMessageCount())
}

// getEnvOrDefault returns environment variable or default value
func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// findButtonData finds a button by its data and returns it
func findButtonData(msg *CapturedMessage, data string) *ButtonInfo {
	if msg == nil || msg.Buttons == nil {
		return nil
	}
	for _, row := range msg.Buttons {
		for _, btn := range row {
			if btn.Data == data {
				return &btn
			}
		}
	}
	return nil
}

// getLastMessageID returns the last message ID from the mock bot
func getLastMessageID(bot *MockBot) int {
	bot.mu.Lock()
	defer bot.mu.Unlock()
	return bot.MessageID
}

// makeTextUpdate creates a tgbotapi.Update simulating a text message
func makeTextUpdate(chatID int64, userID int64, text string) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			MessageID: int(time.Now().UnixNano() % 1000000),
			From: &tgbotapi.User{
				ID:        userID,
				FirstName: "TestUser",
			},
			Chat: &tgbotapi.Chat{
				ID:   chatID,
				Type: "private",
			},
			Date: int(time.Now().Unix()),
			Text: text,
		},
	}
}

// makeCommandUpdate creates a tgbotapi.Update simulating a command
func makeCommandUpdate(chatID int64, userID int64, command string) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			MessageID: int(time.Now().UnixNano() % 1000000),
			From: &tgbotapi.User{
				ID:        userID,
				FirstName: "TestUser",
			},
			Chat: &tgbotapi.Chat{
				ID:   chatID,
				Type: "private",
			},
			Date: int(time.Now().Unix()),
			Text: "/" + command,
			Entities: []tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: len(command) + 1,
				},
			},
		},
	}
}
