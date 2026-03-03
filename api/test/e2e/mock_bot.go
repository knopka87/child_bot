package e2e

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MockBot implements telegram.BotSender for testing
type MockBot struct {
	mu           sync.Mutex
	Messages     []CapturedMessage
	MessageID    int
	Files        map[string]FileData // fileID -> file info
	CallbackAcks []string
	Token        string
}

// FileData holds file information for the mock
type FileData struct {
	Data     []byte
	FilePath string
}

// CapturedMessage represents a message captured by the mock bot
type CapturedMessage struct {
	ChatID    int64          `json:"chat_id"`
	Text      string         `json:"text"`
	ParseMode string         `json:"parse_mode,omitempty"`
	Buttons   [][]ButtonInfo `json:"buttons,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
	Type      string         `json:"type"` // "text", "edit_markup", "callback_ack"
}

// ButtonInfo represents a button in a message
type ButtonInfo struct {
	Text string `json:"text"`
	Data string `json:"data"`
}

// NewMockBot creates a new mock bot for testing
func NewMockBot() *MockBot {
	return &MockBot{
		Messages:     make([]CapturedMessage, 0),
		MessageID:    1000,
		Files:        make(map[string]FileData),
		CallbackAcks: make([]string, 0),
		Token:        "test-bot-token",
	}
}

// Send implements BotSender.Send
func (m *MockBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.MessageID++
	msg := tgbotapi.Message{
		MessageID: m.MessageID,
	}

	switch v := c.(type) {
	case tgbotapi.MessageConfig:
		msg.Chat = &tgbotapi.Chat{ID: v.ChatID}
		msg.Text = v.Text

		captured := CapturedMessage{
			ChatID:    v.ChatID,
			Text:      v.Text,
			ParseMode: v.ParseMode,
			Timestamp: time.Now(),
			Type:      "text",
		}

		if v.ReplyMarkup != nil {
			if kb, ok := v.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup); ok {
				captured.Buttons = extractButtons(kb.InlineKeyboard)
			}
		}

		m.Messages = append(m.Messages, captured)

	case tgbotapi.EditMessageReplyMarkupConfig:
		captured := CapturedMessage{
			ChatID:    v.ChatID,
			Timestamp: time.Now(),
			Type:      "edit_markup",
		}
		if v.ReplyMarkup != nil {
			captured.Buttons = extractButtons(v.ReplyMarkup.InlineKeyboard)
		}
		m.Messages = append(m.Messages, captured)

	case tgbotapi.DocumentConfig:
		msg.Chat = &tgbotapi.Chat{ID: v.ChatID}
		captured := CapturedMessage{
			ChatID:    v.ChatID,
			Text:      "[Document]",
			Timestamp: time.Now(),
			Type:      "document",
		}
		m.Messages = append(m.Messages, captured)
	}

	return msg, nil
}

// Request implements BotSender.Request
func (m *MockBot) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch v := c.(type) {
	case tgbotapi.CallbackConfig:
		m.CallbackAcks = append(m.CallbackAcks, v.CallbackQueryID)
	case tgbotapi.ChatActionConfig:
		// Ignore typing indicators
	case tgbotapi.DeleteMessageConfig:
		// Ignore delete requests
	}

	return &tgbotapi.APIResponse{Ok: true}, nil
}

// GetFile implements BotSender.GetFile
func (m *MockBot) GetFile(config tgbotapi.FileConfig) (tgbotapi.File, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if fd, ok := m.Files[config.FileID]; ok {
		return tgbotapi.File{
			FileID:   config.FileID,
			FilePath: fd.FilePath,
		}, nil
	}

	return tgbotapi.File{}, fmt.Errorf("file not found: %s", config.FileID)
}

// GetFileDirectURL implements BotSender.GetFileDirectURL
func (m *MockBot) GetFileDirectURL(fileID string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if fd, ok := m.Files[fileID]; ok {
		return fmt.Sprintf("http://mock-telegram/file/bot%s/%s", m.Token, fd.FilePath), nil
	}

	return "", fmt.Errorf("file not found: %s", fileID)
}

// GetToken implements BotSender.GetToken
func (m *MockBot) GetToken() string {
	return m.Token
}

// DownloadFile implements BotSender.DownloadFile
func (m *MockBot) DownloadFile(fileID string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if fd, ok := m.Files[fileID]; ok {
		return fd.Data, nil
	}

	return nil, fmt.Errorf("file not found: %s", fileID)
}

// AddFile adds a file to the mock for test setup
func (m *MockBot) AddFile(fileID string, data []byte, filePath string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Files[fileID] = FileData{
		Data:     data,
		FilePath: filePath,
	}
}

// LastMessage returns the most recent captured message
func (m *MockBot) LastMessage() *CapturedMessage {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.Messages) == 0 {
		return nil
	}
	return &m.Messages[len(m.Messages)-1]
}

// LastTextMessage returns the most recent text message (not edit_markup)
func (m *MockBot) LastTextMessage() *CapturedMessage {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := len(m.Messages) - 1; i >= 0; i-- {
		if m.Messages[i].Type == "text" {
			return &m.Messages[i]
		}
	}
	return nil
}

// GetMessages returns a copy of all captured messages
func (m *MockBot) GetMessages() []CapturedMessage {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]CapturedMessage, len(m.Messages))
	copy(result, m.Messages)
	return result
}

// GetMessageCount returns the number of captured messages
func (m *MockBot) GetMessageCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.Messages)
}

// ClearMessages clears all captured messages
func (m *MockBot) ClearMessages() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Messages = make([]CapturedMessage, 0)
}

// ToJSON serializes messages to JSON for debugging
func (m *MockBot) ToJSON() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, _ := json.MarshalIndent(m.Messages, "", "  ")
	return string(data)
}

func extractButtons(keyboard [][]tgbotapi.InlineKeyboardButton) [][]ButtonInfo {
	result := make([][]ButtonInfo, len(keyboard))
	for i, row := range keyboard {
		result[i] = make([]ButtonInfo, len(row))
		for j, btn := range row {
			data := ""
			if btn.CallbackData != nil {
				data = *btn.CallbackData
			}
			result[i][j] = ButtonInfo{
				Text: btn.Text,
				Data: data,
			}
		}
	}
	return result
}
