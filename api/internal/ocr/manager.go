package ocr

import (
	"sync"
)

type Manager struct {
	engine string
	m      sync.Map // chatID -> llmName
}

func NewManager(defaultEngine string) *Manager {
	return &Manager{engine: defaultEngine}
}

func (m *Manager) Get(chatID int64) string {
	if v, ok := m.m.Load(chatID); ok {
		return v.(string)
	}
	return m.engine
}
func (m *Manager) Set(chatID int64, llmName string) {
	m.m.Store(chatID, llmName)
}
