package config

import (
	"log"
	"os"
)

type Config struct {
	Port               string
	WebhookURL         string
	TelegramBotToken   string
	TelegramBotVersion string

	// LLMClient
	DefaultLLM   string
	LLMServerURL string // например: https://llm.example.com  (без хвоста /)
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing required env %s", k)
	}
	return v
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		WebhookURL:         getEnv("WEBHOOK_URL", ""),
		TelegramBotToken:   mustEnv("TELEGRAM_BOT_TOKEN"),
		TelegramBotVersion: getEnv("TELEGRAM_BOT_VERSION", "v1"),

		DefaultLLM:   getEnv("DEFAULT_LLM", "gemini"),
		LLMServerURL: mustEnv("LLM_SERVER_URL"),
	}
}
