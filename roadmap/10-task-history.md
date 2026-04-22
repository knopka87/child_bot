# 10: История задач (Task History)

> Фаза 1 | Приоритет: P2 | Сложность: Низкая | Срок: 1-2 дня

## Цель

Предоставить API для просмотра истории решённых задач с фильтрацией и пагинацией.

## Функциональность

- Список всех задач пользователя
- Фильтрация по предмету, статусу, дате
- Пагинация
- Детали задачи (текст, подсказки, результат)

## Store методы

```go
// internal/store/history.go
package store

type TaskHistoryItem struct {
    SessionID   string    `json:"session_id"`
    Subject     string    `json:"subject"`
    TaskText    string    `json:"task_text"`
    Grade       int       `json:"grade"`
    Status      string    `json:"status"` // "correct", "incorrect", "in_progress"
    HintsUsed   int       `json:"hints_used"`
    CreatedAt   time.Time `json:"created_at"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type HistoryFilter struct {
    Subject   string
    Status    string
    DateFrom  *time.Time
    DateTo    *time.Time
}

func (s *Store) GetUserTaskHistory(ctx context.Context, userID int64, filter HistoryFilter, limit, offset int) ([]TaskHistoryItem, int, error) {
    // Build query with filters
    query := `
        SELECT ts.session_id, pt.subject, pt.task_text_clean, u.grade,
               CASE
                   WHEN te.ok = true THEN 'correct'
                   WHEN te.ok = false THEN 'incorrect'
                   ELSE 'in_progress'
               END as status,
               COALESCE(hc.hints_count, 0) as hints_used,
               ts.created_at,
               te.created_at as completed_at
        FROM task_sessions ts
        LEFT JOIN parsed_tasks pt ON ts.session_id = pt.session_id
        LEFT JOIN "user" u ON ts.chat_id = u.chat_id
        LEFT JOIN LATERAL (
            SELECT ok, created_at FROM timeline_events
            WHERE task_session_id = ts.session_id AND event_type = 'api_check'
            ORDER BY created_at DESC LIMIT 1
        ) te ON true
        LEFT JOIN LATERAL (
            SELECT COUNT(*) as hints_count FROM hints_cache
            WHERE session_id = ts.session_id
        ) hc ON true
        WHERE ts.chat_id = $1
    `

    args := []any{userID}
    argIdx := 2

    if filter.Subject != "" {
        query += fmt.Sprintf(" AND pt.subject = $%d", argIdx)
        args = append(args, filter.Subject)
        argIdx++
    }

    if filter.Status != "" {
        switch filter.Status {
        case "correct":
            query += " AND te.ok = true"
        case "incorrect":
            query += " AND te.ok = false"
        case "in_progress":
            query += " AND te.ok IS NULL"
        }
    }

    if filter.DateFrom != nil {
        query += fmt.Sprintf(" AND ts.created_at >= $%d", argIdx)
        args = append(args, *filter.DateFrom)
        argIdx++
    }

    if filter.DateTo != nil {
        query += fmt.Sprintf(" AND ts.created_at <= $%d", argIdx)
        args = append(args, *filter.DateTo)
        argIdx++
    }

    // Count total
    countQuery := "SELECT COUNT(*) FROM (" + query + ") sub"
    var total int
    if err := s.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
        return nil, 0, err
    }

    // Add pagination
    query += fmt.Sprintf(" ORDER BY ts.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
    args = append(args, limit, offset)

    rows, err := s.DB.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var items []TaskHistoryItem
    for rows.Next() {
        var item TaskHistoryItem
        if err := rows.Scan(
            &item.SessionID, &item.Subject, &item.TaskText, &item.Grade,
            &item.Status, &item.HintsUsed, &item.CreatedAt, &item.CompletedAt,
        ); err != nil {
            return nil, 0, err
        }
        items = append(items, item)
    }

    return items, total, rows.Err()
}

func (s *Store) GetTaskDetails(ctx context.Context, userID int64, sessionID string) (*TaskDetails, error) {
    // Verify ownership
    var ownerID int64
    err := s.DB.QueryRowContext(ctx, `
        SELECT chat_id FROM task_sessions WHERE session_id = $1
    `, sessionID).Scan(&ownerID)
    if err != nil {
        return nil, err
    }
    if ownerID != userID {
        return nil, fmt.Errorf("not found")
    }

    // Get full details
    // ... implementation
}
```

## REST API

```go
// internal/api/handlers/history.go
package handlers

type HistoryHandler struct {
    store *store.Store
}

func (h *HistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
    user := middleware.GetUser(r.Context())

    // Parse query params
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if limit <= 0 || limit > 50 {
        limit = 20
    }

    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

    filter := store.HistoryFilter{
        Subject: r.URL.Query().Get("subject"),
        Status:  r.URL.Query().Get("status"),
    }

    if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
        if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
            filter.DateFrom = &t
        }
    }

    items, total, err := h.store.GetUserTaskHistory(r.Context(), user.UserID, filter, limit, offset)
    if err != nil {
        http.Error(w, `{"error": "failed to get history"}`, http.StatusInternalServerError)
        return
    }

    resp := dto.PaginatedResponse[store.TaskHistoryItem]{
        Data:    items,
        Total:   total,
        HasMore: offset+limit < total,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
```

## API Endpoints

```
GET /api/v1/tasks/history
    ?limit=20
    &offset=0
    &subject=math
    &status=correct
    &date_from=2024-01-01
    &date_to=2024-01-31

GET /api/v1/tasks/{sessionID}
```

## Чек-лист

- [ ] Реализовать `store/history.go`
- [ ] REST API handler
- [ ] Фильтрация и пагинация
- [ ] Детальный просмотр задачи
- [ ] Unit-тесты

---

[← WebSocket](./09-websocket.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Parent-Child →](./11-parent-child.md)
