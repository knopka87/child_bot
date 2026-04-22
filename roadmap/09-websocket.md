# 09: WebSocket (Real-time Updates)

> Фаза 1 | Приоритет: P1 | Сложность: Средняя | Срок: 2-3 дня

## Цель

Добавить WebSocket для real-time обновлений в Mini App: уведомления о похвалах, достижениях, изменениях состояния босса.

## Use Cases

1. **Похвала от родителя** — мгновенное уведомление ребёнку
2. **Achievement unlocked** — анимация в Mini App
3. **Босс побеждён** — уведомление всем участникам
4. **Pet state changed** — обновление UI без refresh

## Архитектура

```
┌─────────────────┐         ┌─────────────────┐
│   Mini App      │◄───WS───│   API Server    │
│   (Frontend)    │         │   (Go)          │
└─────────────────┘         └────────┬────────┘
                                     │
                            ┌────────▼────────┐
                            │   Event Bus     │
                            └────────┬────────┘
                                     │
              ┌──────────────────────┼──────────────────────┐
              │                      │                      │
    ┌─────────▼─────────┐ ┌─────────▼─────────┐ ┌─────────▼─────────┐
    │ Gamification Svc  │ │ Parent Svc        │ │ Other Services    │
    └───────────────────┘ └───────────────────┘ └───────────────────┘
```

## Зависимости

```go
// go.mod
require (
    github.com/gorilla/websocket v1.5.1
)
```

## WebSocket Hub

```go
// internal/api/ws/hub.go
package ws

import (
    "encoding/json"
    "sync"
)

type Hub struct {
    // Registered clients by user ID
    clients map[int64]map[*Client]bool
    mu      sync.RWMutex

    // Broadcast to all clients
    broadcast chan Message

    // Register/unregister channels
    register   chan *Client
    unregister chan *Client
}

type Client struct {
    hub    *Hub
    userID int64
    conn   *websocket.Conn
    send   chan []byte
}

type Message struct {
    Type    string `json:"type"`
    Payload any    `json:"payload"`
    UserID  int64  `json:"-"` // Target user (0 = broadcast to all)
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[int64]map[*Client]bool),
        broadcast:  make(chan Message, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            if h.clients[client.userID] == nil {
                h.clients[client.userID] = make(map[*Client]bool)
            }
            h.clients[client.userID][client] = true
            h.mu.Unlock()

        case client := <-h.unregister:
            h.mu.Lock()
            if clients, ok := h.clients[client.userID]; ok {
                if _, ok := clients[client]; ok {
                    delete(clients, client)
                    close(client.send)
                    if len(clients) == 0 {
                        delete(h.clients, client.userID)
                    }
                }
            }
            h.mu.Unlock()

        case message := <-h.broadcast:
            data, _ := json.Marshal(message)

            h.mu.RLock()
            if message.UserID == 0 {
                // Broadcast to all
                for _, clients := range h.clients {
                    for client := range clients {
                        select {
                        case client.send <- data:
                        default:
                            // Buffer full, skip
                        }
                    }
                }
            } else {
                // Send to specific user
                if clients, ok := h.clients[message.UserID]; ok {
                    for client := range clients {
                        select {
                        case client.send <- data:
                        default:
                        }
                    }
                }
            }
            h.mu.RUnlock()
        }
    }
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID int64, msgType string, payload any) {
    h.broadcast <- Message{
        Type:    msgType,
        Payload: payload,
        UserID:  userID,
    }
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(msgType string, payload any) {
    h.broadcast <- Message{
        Type:    msgType,
        Payload: payload,
        UserID:  0,
    }
}
```

## WebSocket Client

```go
// internal/api/ws/client.go
package ws

import (
    "log"
    "time"

    "github.com/gorilla/websocket"
)

const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 512
)

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, _, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("ws error: %v", err)
            }
            break
        }
        // We don't expect messages from client in this implementation
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            // Batch queued messages
            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write([]byte{'\n'})
                w.Write(<-c.send)
            }

            if err := w.Close(); err != nil {
                return
            }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

## WebSocket Handler

```go
// internal/api/ws/handler.go
package ws

import (
    "net/http"

    "github.com/gorilla/websocket"

    mw "child_bot/api/internal/api/middleware"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // TODO: Validate origin in production
        return true
    },
}

func ServeWs(hub *Hub) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get user from JWT (passed as query param for WS)
        token := r.URL.Query().Get("token")
        if token == "" {
            http.Error(w, "missing token", http.StatusUnauthorized)
            return
        }

        user, err := mw.ValidateToken(token)
        if err != nil {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }

        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Printf("ws upgrade error: %v", err)
            return
        }

        client := &Client{
            hub:    hub,
            userID: user.UserID,
            conn:   conn,
            send:   make(chan []byte, 256),
        }

        hub.register <- client

        go client.writePump()
        go client.readPump()
    }
}
```

## Интеграция с Event Bus

```go
// internal/service/notification/service.go
package notification

import (
    "child_bot/api/internal/api/ws"
    "child_bot/api/internal/service/events"
)

type NotificationService struct {
    hub *ws.Hub
}

func NewNotificationService(hub *ws.Hub, eventBus *events.EventBus) *NotificationService {
    svc := &NotificationService{hub: hub}

    // Subscribe to events
    eventBus.Subscribe(events.EventAchievementUnlocked, svc.onAchievement)
    eventBus.Subscribe(events.EventPraiseReceived, svc.onPraise)
    eventBus.Subscribe(events.EventBossDefeated, svc.onBossDefeated)
    eventBus.Subscribe(events.EventPetHappy, svc.onPetHappy)

    return svc
}

func (s *NotificationService) onAchievement(event events.Event) {
    achievement := event.Payload["achievement"]

    s.hub.SendToUser(event.UserID, "achievement_unlocked", map[string]any{
        "achievement": achievement,
    })
}

func (s *NotificationService) onPraise(event events.Event) {
    s.hub.SendToUser(event.UserID, "praise_received", map[string]any{
        "from":    event.Payload["from_name"],
        "message": event.Payload["message"],
        "sticker": event.Payload["sticker_type"],
    })
}

func (s *NotificationService) onBossDefeated(event events.Event) {
    // Broadcast to all connected clients
    s.hub.Broadcast("boss_defeated", map[string]any{
        "boss_name": event.Payload["boss_name"],
    })
}

func (s *NotificationService) onPetHappy(event events.Event) {
    s.hub.SendToUser(event.UserID, "pet_state_changed", map[string]any{
        "state": "happy",
    })
}
```

## Интеграция в Router

```go
// internal/api/router.go

func NewRouter(cfg Config, ...) http.Handler {
    r := chi.NewRouter()

    // ... middleware ...

    // WebSocket hub
    hub := ws.NewHub()
    go hub.Run()

    // WebSocket endpoint
    r.Get("/api/v1/ws", ws.ServeWs(hub))

    // Pass hub to services that need it
    notificationSvc := notification.NewNotificationService(hub, eventBus)

    // ... other routes ...

    return r
}
```

## Message Types

```typescript
// Frontend TypeScript types

interface WSMessage {
  type: string;
  payload: any;
}

// Achievement unlocked
interface AchievementUnlockedPayload {
  achievement: {
    id: string;
    name: string;
    description: string;
    icon_key: string;
    rarity: string;
    xp_reward: number;
  };
}

// Praise received
interface PraiseReceivedPayload {
  from: string;       // Parent name
  message: string;
  sticker: string;    // Sticker type
}

// Boss defeated
interface BossDefeatedPayload {
  boss_name: string;
}

// Pet state changed
interface PetStateChangedPayload {
  state: 'hungry' | 'fed' | 'happy';
}
```

## Frontend Integration

```typescript
// Mini App frontend

class WebSocketClient {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;

  constructor(private token: string, private onMessage: (msg: WSMessage) => void) {}

  connect() {
    const wsUrl = `wss://api.example.com/api/v1/ws?token=${this.token}`;
    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.reconnectAttempts = 0;
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data) as WSMessage;
        this.onMessage(message);
      } catch (e) {
        console.error('Failed to parse WS message', e);
      }
    };

    this.ws.onclose = () => {
      console.log('WebSocket closed');
      this.reconnect();
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error', error);
    };
  }

  private reconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnect attempts reached');
      return;
    }

    this.reconnectAttempts++;
    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);

    setTimeout(() => {
      console.log(`Reconnecting... attempt ${this.reconnectAttempts}`);
      this.connect();
    }, delay);
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}

// Usage
const wsClient = new WebSocketClient(jwtToken, (message) => {
  switch (message.type) {
    case 'achievement_unlocked':
      showAchievementAnimation(message.payload);
      break;
    case 'praise_received':
      showPraiseNotification(message.payload);
      break;
    case 'boss_defeated':
      showBossDefeatedCelebration(message.payload);
      break;
    case 'pet_state_changed':
      updatePetUI(message.payload);
      break;
  }
});

wsClient.connect();
```

## Тестирование

```go
// internal/api/ws/hub_test.go
func TestHub(t *testing.T) {
    hub := NewHub()
    go hub.Run()

    // Create mock client
    client := &Client{
        hub:    hub,
        userID: 123,
        send:   make(chan []byte, 256),
    }
    hub.register <- client

    // Wait for registration
    time.Sleep(10 * time.Millisecond)

    // Send message
    hub.SendToUser(123, "test", map[string]string{"foo": "bar"})

    // Check received
    select {
    case msg := <-client.send:
        var m Message
        json.Unmarshal(msg, &m)
        if m.Type != "test" {
            t.Errorf("expected type 'test', got %s", m.Type)
        }
    case <-time.After(100 * time.Millisecond):
        t.Error("message not received")
    }

    // Unregister
    hub.unregister <- client
}
```

## Чек-лист

- [ ] Добавить gorilla/websocket в go.mod
- [ ] Создать `internal/api/ws/` package
- [ ] Реализовать Hub и Client
- [ ] WebSocket handler с JWT валидацией
- [ ] Интегрировать с event bus
- [ ] Создать NotificationService
- [ ] Документировать message types для frontend
- [ ] Unit-тесты для Hub
- [ ] Интеграционные тесты
- [ ] Load testing (много подключений)

## Связанные шаги

- [05-auth-system.md](./05-auth-system.md) — JWT для WS аутентификации
- [11-parent-child.md](./11-parent-child.md) — отправка похвал через WS

---

[← Boss System](./08-boss-system.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Task History →](./10-task-history.md)
