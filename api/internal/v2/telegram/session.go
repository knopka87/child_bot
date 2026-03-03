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

// ensureSessionForNewTask создаёт новую сессию только если нет активного batch
// Для альбомов фото (MediaGroupID != "") создаётся одна сессия на весь альбом
func (r *Router) ensureSessionForNewTask(cid int64, mediaGroupID string) string {
	// Если это часть альбома, проверяем есть ли уже session для этого batch
	if mediaGroupID != "" {
		batchKey := "grp:" + mediaGroupID
		if v, ok := batchSessionKeys.Load(cid); ok {
			if existingKey, ok := v.(string); ok && existingKey == batchKey {
				// Уже есть session для этого batch, не создаём новую
				if sid, ok := r.getSession(cid); ok {
					return sid
				}
			}
		}
		// Создаём новую session и запоминаем batch key
		sid := uuid.NewString()
		r.setSession(cid, sid)
		batchSessionKeys.Store(cid, batchKey)
		return sid
	}

	// Для одиночного фото (без MediaGroupID) всегда создаём новую session
	sid := uuid.NewString()
	r.setSession(cid, sid)
	batchSessionKeys.Delete(cid) // очищаем связь с предыдущим batch
	return sid
}
