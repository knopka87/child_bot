# 11: Связь родитель-ребёнок

> Фаза 1 | Приоритет: P2 | Сложность: Средняя | Срок: 3-4 дня

## Цель

Реализовать связывание аккаунтов родителя и ребёнка для функций родительского портала.

## Функциональность

- Родитель привязывает ребёнка по коду
- Просмотр статистики ребёнка
- Отправка похвал
- Parental Gate (математический замок)

## Миграция

```sql
-- migrations/033_parent_child.up.sql

CREATE TABLE parent_child (
    id SERIAL PRIMARY KEY,
    parent_id BIGINT NOT NULL REFERENCES "user"(chat_id) ON DELETE CASCADE,
    child_id BIGINT NOT NULL REFERENCES "user"(chat_id) ON DELETE CASCADE,
    link_code TEXT,
    linked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(parent_id, child_id)
);

CREATE TABLE praise (
    id SERIAL PRIMARY KEY,
    from_user_id BIGINT NOT NULL REFERENCES "user"(chat_id),
    to_user_id BIGINT NOT NULL REFERENCES "user"(chat_id),
    message TEXT,
    sticker_type TEXT NOT NULL,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_praise_to_user ON praise(to_user_id, created_at DESC);
```

## Store методы

```go
// internal/store/parent.go

func (s *Store) GenerateLinkCode(ctx context.Context, childID int64) (string, error) {
    code := generateRandomCode(6) // e.g., "ABC123"

    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO parent_child (parent_id, child_id, link_code)
        VALUES (0, $1, $2)
        ON CONFLICT (parent_id, child_id) WHERE parent_id = 0
        DO UPDATE SET link_code = $2, created_at = NOW()
    `, childID, code)

    return code, err
}

func (s *Store) LinkParentChild(ctx context.Context, parentID int64, code string) (int64, error) {
    var childID int64

    err := s.DB.QueryRowContext(ctx, `
        UPDATE parent_child
        SET parent_id = $1, linked_at = NOW(), link_code = NULL
        WHERE link_code = $2 AND parent_id = 0
          AND created_at > NOW() - INTERVAL '24 hours'
        RETURNING child_id
    `, parentID, code).Scan(&childID)

    if err == sql.ErrNoRows {
        return 0, fmt.Errorf("invalid or expired code")
    }

    // Update parent role
    s.DB.ExecContext(ctx, `UPDATE "user" SET role = 'parent' WHERE chat_id = $1`, parentID)

    return childID, err
}

func (s *Store) GetChildren(ctx context.Context, parentID int64) ([]ChildInfo, error) {
    rows, err := s.DB.QueryContext(ctx, `
        SELECT pc.child_id, c.first_name, c.username, u.grade,
               us.current_streak, u.xp, u.level
        FROM parent_child pc
        JOIN chat c ON pc.child_id = c.id
        JOIN "user" u ON pc.child_id = u.chat_id
        LEFT JOIN user_streak us ON pc.child_id = us.user_id
        WHERE pc.parent_id = $1 AND pc.linked_at IS NOT NULL
    `, parentID)
    // ...
}

func (s *Store) SendPraise(ctx context.Context, fromID, toID int64, message, stickerType string) error {
    _, err := s.DB.ExecContext(ctx, `
        INSERT INTO praise (from_user_id, to_user_id, message, sticker_type)
        VALUES ($1, $2, $3, $4)
    `, fromID, toID, message, stickerType)
    return err
}

func (s *Store) GetUnreadPraises(ctx context.Context, userID int64) ([]Praise, error) {
    // ...
}

func (s *Store) MarkPraiseRead(ctx context.Context, praiseID int64, userID int64) error {
    _, err := s.DB.ExecContext(ctx, `
        UPDATE praise SET read_at = NOW()
        WHERE id = $1 AND to_user_id = $2
    `, praiseID, userID)
    return err
}
```

## REST API

```go
// internal/api/handlers/parent.go

// GET /api/v1/me/link-code - генерация кода для привязки
func (h *ParentHandler) GenerateLinkCode(w http.ResponseWriter, r *http.Request) {
    user := middleware.GetUser(r.Context())

    code, err := h.store.GenerateLinkCode(r.Context(), user.UserID)
    // ...

    json.NewEncoder(w).Encode(map[string]string{"code": code})
}

// POST /api/v1/parent/link - привязка ребёнка по коду
func (h *ParentHandler) LinkChild(w http.ResponseWriter, r *http.Request) {
    user := middleware.GetUser(r.Context())

    var req struct {
        Code string `json:"code"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    childID, err := h.store.LinkParentChild(r.Context(), user.UserID, req.Code)
    // ...
}

// GET /api/v1/parent/children - список детей
func (h *ParentHandler) GetChildren(w http.ResponseWriter, r *http.Request) {
    user := middleware.GetUser(r.Context())

    children, err := h.store.GetChildren(r.Context(), user.UserID)
    // ...
}

// POST /api/v1/parent/child/{childID}/praise - отправка похвалы
func (h *ParentHandler) SendPraise(w http.ResponseWriter, r *http.Request) {
    user := middleware.GetUser(r.Context())
    childID, _ := strconv.ParseInt(chi.URLParam(r, "childID"), 10, 64)

    // Verify parent-child relationship
    // ...

    var req struct {
        Message     string `json:"message"`
        StickerType string `json:"sticker_type"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    err := h.store.SendPraise(r.Context(), user.UserID, childID, req.Message, req.StickerType)
    // ...

    // Notify via WebSocket
    h.eventBus.Publish(events.Event{
        Type:   events.EventPraiseReceived,
        UserID: childID,
        Payload: map[string]any{
            "from_name": user.FirstName,
            "message":   req.Message,
            "sticker":   req.StickerType,
        },
    })
}
```

## Parental Gate

```go
// Простая математическая проверка для доступа к родительским функциям
func (h *ParentHandler) VerifyParentalGate(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Answer int `json:"answer"`
        Question string `json:"question"` // e.g., "7+8"
    }
    json.NewDecoder(r.Body).Decode(&req)

    // Verify answer
    expected := evaluateQuestion(req.Question) // 15
    if req.Answer != expected {
        http.Error(w, `{"error": "incorrect answer"}`, http.StatusForbidden)
        return
    }

    // Issue parent session token
    // ...
}
```

## Чек-лист

- [ ] Миграция `033_parent_child.up.sql`
- [ ] Store методы для связи и похвал
- [ ] REST API endpoints
- [ ] WebSocket уведомления о похвалах
- [ ] Parental Gate
- [ ] Unit-тесты

---

[← Task History](./10-task-history.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Pet Evolution →](./12-pet-evolution.md)
