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

// EnsureJSONRaw makes sure the blob is valid JSON. If it's plain text (not JSON),
// it wraps it as a properly quoted JSON string. If empty, returns nil.
func EnsureJSONRaw(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(raw, &v); err == nil {
		return raw
	}
	// treat raw as plain text; quote it
	b, _ := json.Marshal(string(raw))
	return json.RawMessage(b)
}
