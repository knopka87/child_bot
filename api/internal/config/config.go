package config

import (
	"log"
	"os"
)

type Config struct {
	Port string

	// LLMClient
	DefaultLLM   string
	LLMServerURL string // например: https://llm.example.com  (без хвоста /)

	// CORS
	AllowedOrigins string

	// App
	AppURL string // базовый URL приложения для реферальных ссылок
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
		Port: getEnv("PORT", "8080"),

		DefaultLLM:   getEnv("DEFAULT_LLM", "gemini"),
		LLMServerURL: mustEnv("LLM_SERVER_URL"),

		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000"),
		AppURL:         getEnv("APP_URL", "http://localhost:5173"),
	}
}
