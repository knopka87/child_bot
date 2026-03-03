package telegram

import (
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotSender abstracts Telegram bot operations for testing
type BotSender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	GetFile(config tgbotapi.FileConfig) (tgbotapi.File, error)
	GetFileDirectURL(fileID string) (string, error)
	GetToken() string
	// DownloadFile downloads file bytes by file ID
	// This allows for mocking file downloads in tests
	DownloadFile(fileID string) ([]byte, error)
}

// BotAPIWrapper wraps *tgbotapi.BotAPI to implement BotSender interface
type BotAPIWrapper struct {
	*tgbotapi.BotAPI
}

// NewBotAPIWrapper creates a new wrapper for tgbotapi.BotAPI
func NewBotAPIWrapper(bot *tgbotapi.BotAPI) *BotAPIWrapper {
	return &BotAPIWrapper{BotAPI: bot}
}

// GetToken returns the bot token
func (w *BotAPIWrapper) GetToken() string {
	return w.Token
}

// DownloadFile downloads file bytes using the Telegram API
func (w *BotAPIWrapper) DownloadFile(fileID string) ([]byte, error) {
	url, err := w.GetFileDirectURL(fileID)
	if err != nil {
		return nil, fmt.Errorf("get file URL: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed: status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}
