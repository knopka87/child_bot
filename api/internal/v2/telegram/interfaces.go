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

// maxFileDownloadSize — максимальный размер скачиваемого файла (20MB)
// Защита от OOM при скачивании больших файлов
const maxFileDownloadSize = 20 * 1024 * 1024

// DownloadFile downloads file bytes using the Telegram API
// Ограничивает размер скачиваемого файла для защиты от OOM
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
		// Ограничиваем чтение тела ошибки
		limitedBody := io.LimitReader(resp.Body, 4096)
		body, _ := io.ReadAll(limitedBody)
		return nil, fmt.Errorf("download failed: status %d: %s", resp.StatusCode, string(body))
	}

	// Ограничиваем размер скачиваемого файла для защиты от OOM
	limitedReader := io.LimitReader(resp.Body, maxFileDownloadSize+1)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if len(data) > maxFileDownloadSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed (%d bytes)", maxFileDownloadSize)
	}

	return data, nil
}
