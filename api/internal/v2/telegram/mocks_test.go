package telegram

import (
	"context"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
)

// MockBotSender implements BotSender for testing
type MockBotSender struct {
	mu            sync.Mutex
	SentMessages  []tgbotapi.Chattable
	Requests      []tgbotapi.Chattable
	Token         string
	DownloadBytes []byte
	DownloadError error
	FileURL       string
}

func NewMockBotSender() *MockBotSender {
	return &MockBotSender{
		Token: "test-token",
	}
}

func (m *MockBotSender) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SentMessages = append(m.SentMessages, c)
	return tgbotapi.Message{MessageID: len(m.SentMessages)}, nil
}

func (m *MockBotSender) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Requests = append(m.Requests, c)
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func (m *MockBotSender) GetFile(config tgbotapi.FileConfig) (tgbotapi.File, error) {
	return tgbotapi.File{FileID: config.FileID, FilePath: "test/path"}, nil
}

func (m *MockBotSender) GetFileDirectURL(fileID string) (string, error) {
	if m.FileURL != "" {
		return m.FileURL, nil
	}
	return "https://api.telegram.org/file/test/" + fileID, nil
}

func (m *MockBotSender) GetToken() string {
	return m.Token
}

func (m *MockBotSender) DownloadFile(fileID string) ([]byte, error) {
	if m.DownloadError != nil {
		return nil, m.DownloadError
	}
	if m.DownloadBytes != nil {
		return m.DownloadBytes, nil
	}
	// Return JPEG magic bytes by default
	return []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46}, nil
}

func (m *MockBotSender) GetSentMessages() []tgbotapi.Chattable {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]tgbotapi.Chattable, len(m.SentMessages))
	copy(result, m.SentMessages)
	return result
}

func (m *MockBotSender) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SentMessages = nil
	m.Requests = nil
}

// MockStore implements store methods needed for testing
type MockStore struct {
	mu               sync.Mutex
	Users            map[int64]store.User
	Chats            map[int64]store.Chat
	Sessions         map[int64]store.TaskSession
	ParsedTasks      map[string]*store.ParsedTasks
	HintCache        map[string]store.HintCache
	TimelineEvents   []store.TimelineEvent
	MetricEvents     []store.MetricEvent
	InsertEventCalls int
}

func NewMockStore() *MockStore {
	return &MockStore{
		Users:       make(map[int64]store.User),
		Chats:       make(map[int64]store.Chat),
		Sessions:    make(map[int64]store.TaskSession),
		ParsedTasks: make(map[string]*store.ParsedTasks),
		HintCache:   make(map[string]store.HintCache),
	}
}

func (m *MockStore) FindUserByChatID(ctx context.Context, chatID int64) (store.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if user, ok := m.Users[chatID]; ok {
		return user, nil
	}
	return store.User{}, nil
}

func (m *MockStore) UpsertUser(ctx context.Context, user store.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Users[user.ID] = user
	return nil
}

func (m *MockStore) FindChatByID(ctx context.Context, id int64) (store.Chat, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if chat, ok := m.Chats[id]; ok {
		return chat, nil
	}
	return store.Chat{}, nil
}

func (m *MockStore) UpsertChat(ctx context.Context, chat store.Chat) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Chats[chat.ID] = chat
	return nil
}

func (m *MockStore) FindSession(ctx context.Context, chatID int64) (store.TaskSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if session, ok := m.Sessions[chatID]; ok {
		return session, nil
	}
	return store.TaskSession{}, nil
}

func (m *MockStore) UpsertSession(ctx context.Context, session store.TaskSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Sessions[session.ChatID] = session
	return nil
}

func (m *MockStore) UpdateSessionState(ctx context.Context, chatID int64, state, mode *string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if session, ok := m.Sessions[chatID]; ok {
		if state != nil {
			session.CurrentState = state
		}
		if mode != nil {
			session.ChatMode = mode
		}
		m.Sessions[chatID] = session
	}
	return nil
}

func (m *MockStore) InsertHistory(ctx context.Context, event store.TimelineEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TimelineEvents = append(m.TimelineEvents, event)
	return nil
}

func (m *MockStore) InsertEvent(ctx context.Context, event store.MetricEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MetricEvents = append(m.MetricEvents, event)
	m.InsertEventCalls++
	return nil
}

func (m *MockStore) FindLastConfirmedParse(ctx context.Context, sid string) (*store.ParsedTasks, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if pt, ok := m.ParsedTasks[sid]; ok {
		return pt, true
	}
	return nil, false
}

func (m *MockStore) UpsertParse(ctx context.Context, pt store.ParsedTasks) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ParsedTasks[pt.SessionID] = &pt
	return nil
}

func (m *MockStore) MarkAcceptedParseBySID(ctx context.Context, sessionID, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if pt, ok := m.ParsedTasks[sessionID]; ok {
		pt.Accepted = true
		pt.AcceptReason = reason
	}
	return nil
}

func (m *MockStore) UpsertHint(ctx context.Context, hc store.HintCache) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := hc.SessionID + "_" + hc.Level
	m.HintCache[key] = hc
	return nil
}

func (m *MockStore) GetMetricEvents() []store.MetricEvent {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]store.MetricEvent, len(m.MetricEvents))
	copy(result, m.MetricEvents)
	return result
}

func (m *MockStore) GetTimelineEvents() []store.TimelineEvent {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]store.TimelineEvent, len(m.TimelineEvents))
	copy(result, m.TimelineEvents)
	return result
}

func (m *MockStore) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Users = make(map[int64]store.User)
	m.Chats = make(map[int64]store.Chat)
	m.Sessions = make(map[int64]store.TaskSession)
	m.ParsedTasks = make(map[string]*store.ParsedTasks)
	m.HintCache = make(map[string]store.HintCache)
	m.TimelineEvents = nil
	m.MetricEvents = nil
	m.InsertEventCalls = 0
}

// SetUser adds a user to the mock store
func (m *MockStore) SetUser(chatID int64, grade int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Users[chatID] = store.User{
		ID:    chatID,
		Grade: &grade,
	}
}
