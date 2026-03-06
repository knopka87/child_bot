package telegram

import (
	"context"

	"github.com/google/uuid"

	"child-bot/api/internal/store"
)

// sessionByChat хранит сессии с TTL (24 часа)
var sessionByChat = NewTTLCache("sessionByChat", UserDataTTL)

// helpers
func (r *Router) setSession(cid int64, sid string) {
	sessionByChat.Store(cid, sid)
	_ = r.Store.UpsertSession(context.Background(), store.TaskSession{
		ChatID:    cid,
		SessionID: sid,
	})
}
func (r *Router) getSession(cid int64) (string, bool) {
	if v, ok := sessionByChat.Load(cid); ok {
		if sid, ok := v.(string); ok {
			return sid, true
		}
	}

	if s, err := r.Store.FindSession(context.Background(), cid); err == nil && s.SessionID != "" {
		// Кешируем результат из БД для избежания повторных запросов
		sessionByChat.Store(cid, s.SessionID)
		return s.SessionID, true
	}

	return "", false
}
func (r *Router) clearSession(cid int64) { sessionByChat.Delete(cid) }
func (r *Router) ensureSession(cid int64) string {
	if sid, ok := r.getSession(cid); ok && sid != "" {
		return sid
	}
	sid := uuid.NewString()
	r.setSession(cid, sid)
	return sid
}

// batchSessionKeys хранит ключи batch -> session ID для корректной обработки альбомов
var batchSessionKeys = NewTTLCache("batchSessionKeys", PendingTTL)

// ensureSessionForNewTask возвращает существующую сессию или создаёт новую
// Для альбомов фото (MediaGroupID != "") создаётся одна сессия на весь альбом
func (r *Router) ensureSessionForNewTask(cid int64, mediaGroupID string) string {
	// Если сессия уже есть — используем её
	if sid, ok := r.getSession(cid); ok && sid != "" {
		if mediaGroupID != "" {
			// Для альбома сохраняем batch key
			batchKey := "grp:" + mediaGroupID
			batchSessionKeys.Store(cid, batchKey)
		}
		return sid
	}

	// Сессии нет — создаём новую
	if mediaGroupID != "" {
		batchKey := "grp:" + mediaGroupID
		sid := uuid.NewString()
		r.setSession(cid, sid)
		batchSessionKeys.Store(cid, batchKey)
		return sid
	}

	sid := uuid.NewString()
	r.setSession(cid, sid)
	batchSessionKeys.Delete(cid)
	return sid
}
