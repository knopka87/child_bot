package telegram

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"child-bot/api/internal/store"
)

var sessionByChat sync.Map // chatID(int64) -> string (UUID)

// helpers
func (r *Router) setSession(cid int64, sid string) {
	sessionByChat.Store(cid, sid)
	_ = r.Session.Insert(context.Background(), store.TaskSession{
		ChatID:    cid,
		SessionID: sid,
	})
}
func (r *Router) getSession(cid int64) (string, bool) {
	if v, ok := sessionByChat.Load(cid); ok {
		return v.(string), true
	}

	if s, err := r.Session.Find(context.Background(), cid); err == nil && s.SessionID != "" {
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
