package e2e

import (
	"os"
	"testing"
)

// TestConfig holds configuration for E2E tests
type TestConfig struct {
	LLMProxyURL string // e.g., "http://138.124.55.145:80" or real proxy
	DatabaseURL string // PostgreSQL connection string
	LLMName     string // LLM name: "gpt", "gemini", etc.
}

// loadTestConfig loads test configuration from environment variables
func loadTestConfig(t *testing.T) *TestConfig {
	t.Helper()

	cfg := &TestConfig{
		LLMProxyURL: getEnvOrDefault("TEST_LLM_PROXY_URL", "http://138.124.55.145:80"),
		DatabaseURL: getEnvOrDefault("TEST_DATABASE_URL", "postgres://localhost:5432/child_bot_test?sslmode=disable"),
		LLMName:     getEnvOrDefault("TEST_LLM_NAME", "gpt"),
	}

	// Validate required configuration
	if cfg.LLMProxyURL == "" {
		t.Skip("TEST_LLM_PROXY_URL not set, skipping E2E test")
	}

	return cfg
}

// getTemplatesDir returns the path to templates directory
func getTemplatesDir() string {
	// Try relative path from test directory
	paths := []string{
		"../../internal/v2/templates",
		"../../../internal/v2/templates",
		"api/internal/v2/templates",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return "../../internal/v2/templates" // default
}
