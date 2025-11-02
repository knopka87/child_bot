package service

import (
	"sync"
)

type LlmManager struct {
	engine string
	m      sync.Map // chatID -> llmName
}

func NewLlmManager(defaultEngine string) *LlmManager {
	return &LlmManager{engine: defaultEngine}
}

func (m *LlmManager) Get(chatID int64) string {
	if v, ok := m.m.Load(chatID); ok {
		return v.(string)
	}
	return m.engine
}
func (m *LlmManager) Set(chatID int64, llmName string) {
	m.m.Store(chatID, llmName)
}
