package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
)

func StripCodeFences(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

// SHA256Hex возвращает SHA-256 хэш входных данных в виде шестнадцатеричной строки (нижний регистр).
func SHA256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func PrettyJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func NullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
