# Roadmap: Миграция на Mini App

> Переход от Telegram Bot к Telegram Mini App с геймификацией

## Обзор фаз

| Фаза | Название | Срок | Цель |
|------|----------|------|------|
| **0** | Quick Wins | 1-2 недели | WOW-эффект в текущем боте без Mini App |
| **1** | Mini App MVP | 3-4 недели | Базовая геймификация, REST API |
| **2** | Расширение | 2-3 месяца | Социальные функции, кастомизация |
| **3** | Продвинутые фичи | 6+ месяцев | Голос, мини-игры, сезоны |

---

## Фаза 0: Quick Wins (1-2 недели)

Минимальные изменения backend для максимального эффекта в текущем Telegram боте.

| # | Шаг | Файл | Приоритет | Сложность |
|---|-----|------|-----------|-----------|
| 01 | Streak-система | [01-streak-system.md](./01-streak-system.md) | P0 | Низкая |
| 02 | Система достижений | [02-achievements-system.md](./02-achievements-system.md) | P0 | Низкая |
| 03 | Ежедневные отчёты | [03-daily-reports.md](./03-daily-reports.md) | P1 | Низкая |

---

## Фаза 1: Mini App MVP (3-4 недели)

Создание REST API и базовой геймификации для Mini App.

| # | Шаг | Файл | Приоритет | Сложность |
|---|-----|------|-----------|-----------|
| 04 | REST API Layer | [04-api-layer.md](./04-api-layer.md) | P0 | Высокая |
| 05 | Авторизация (Init Data + JWT) | [05-auth-system.md](./05-auth-system.md) | P0 | Средняя |
| 06 | Рефакторинг сервисного слоя | [06-service-refactoring.md](./06-service-refactoring.md) | P0 | Высокая |
| 07 | Питомец-компаньон | [07-pet-system.md](./07-pet-system.md) | P1 | Низкая |
| 08 | Босс недели | [08-boss-system.md](./08-boss-system.md) | P1 | Средняя |
| 09 | WebSocket (real-time) | [09-websocket.md](./09-websocket.md) | P1 | Средняя |
| 10 | История задач | [10-task-history.md](./10-task-history.md) | P2 | Низкая |
| 11 | Связь родитель-ребёнок | [11-parent-child.md](./11-parent-child.md) | P2 | Средняя |

---

## Фаза 2: Расширение (2-3 месяца)

Социальные функции и продвинутая геймификация.

| # | Шаг | Файл | Приоритет | Сложность |
|---|-----|------|-----------|-----------|
| 12 | Эволюция питомца | [12-pet-evolution.md](./12-pet-evolution.md) | P1 | Средняя |
| 13 | Кастомизация маскота | [13-customization.md](./13-customization.md) | P2 | Высокая |
| 14 | Карта знаний | [14-knowledge-map.md](./14-knowledge-map.md) | P1 | Высокая |
| 15 | Spaced Repetition | [15-spaced-repetition.md](./15-spaced-repetition.md) | P2 | Средняя |
| 16 | Таблица лидеров | [16-leaderboard.md](./16-leaderboard.md) | P2 | Низкая |
| 17 | Семейные квесты | [17-family-quests.md](./17-family-quests.md) | P3 | Средняя |

---

## Фаза 3: Продвинутые фичи (6+ месяцев)

Инновационные функции для долгосрочного retention.

| # | Шаг | Файл | Приоритет | Сложность |
|---|-----|------|-----------|-----------|
| 18 | Голосовой режим | [18-voice-mode.md](./18-voice-mode.md) | P2 | Высокая |
| 19 | Мини-игры | [19-mini-games.md](./19-mini-games.md) | P3 | Очень высокая |
| 20 | Сезонные события | [20-seasonal-events.md](./20-seasonal-events.md) | P3 | Средняя |

---

## Архитектурная диаграмма

```
┌─────────────────────────────────────────────────────────────┐
│                        CLIENTS                               │
├─────────────────────────────┬───────────────────────────────┤
│     Telegram Bot (v2)       │        Mini App (Web)         │
│   (polling/webhook)         │      (REST + WebSocket)       │
└─────────────┬───────────────┴───────────────┬───────────────┘
              │                               │
              ▼                               ▼
┌─────────────────────────────────────────────────────────────┐
│                      API GATEWAY                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ TG Handler  │  │ REST API    │  │ WebSocket Server    │  │
│  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘  │
└─────────┼────────────────┼────────────────────┼──────────────┘
          │                │                    │
          ▼                ▼                    ▼
┌─────────────────────────────────────────────────────────────┐
│                    SERVICE LAYER                             │
│  ┌──────────┐  ┌──────────────┐  ┌────────────────────────┐ │
│  │ TaskSvc  │  │ GamificationSvc│ │ NotificationSvc       │ │
│  │ (PARSE,  │  │ (Streak,     │  │ (Push, WS broadcast)  │ │
│  │  HINT,   │  │  Pet, Boss,  │  │                       │ │
│  │  CHECK)  │  │  Achievements)│ │                       │ │
│  └────┬─────┘  └───────┬──────┘  └───────────┬───────────┘ │
└───────┼────────────────┼─────────────────────┼──────────────┘
        │                │                     │
        ▼                ▼                     ▼
┌─────────────────────────────────────────────────────────────┐
│                      EVENT BUS                               │
│         (task_solved, hint_used, check_correct, ...)        │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    DATA LAYER                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐   │
│  │  PostgreSQL  │  │  Redis       │  │  LLM Server      │   │
│  │  (sessions,  │  │  (cache,     │  │  (DETECT, PARSE, │   │
│  │   users,     │  │   rate limit)│  │   HINT, CHECK)   │   │
│  │   gamification)│ │              │  │                  │   │
│  └──────────────┘  └──────────────┘  └──────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## Зависимости между шагами

```
Фаза 0:
  01-streak ──┐
  02-achievements ──┼──▶ (можно параллельно)
  03-reports ──┘

Фаза 1:
  04-api-layer ──▶ 05-auth ──▶ 06-service-refactoring
                                      │
                    ┌─────────────────┼─────────────────┐
                    ▼                 ▼                 ▼
              07-pet-system     08-boss-system    09-websocket
                    │                 │                 │
                    └─────────────────┼─────────────────┘
                                      ▼
                              10-task-history
                                      │
                                      ▼
                              11-parent-child

Фаза 2:
  07-pet-system ──▶ 12-pet-evolution
  06-service ──▶ 13-customization
  06-service ──▶ 14-knowledge-map ──▶ 15-spaced-repetition
  02-achievements ──▶ 16-leaderboard
  11-parent-child ──▶ 17-family-quests

Фаза 3:
  06-service ──▶ 18-voice-mode
  06-service ──▶ 19-mini-games
  13-customization ──▶ 20-seasonal-events
```

---

## Метрики успеха

| Фаза | Метрика | Цель |
|------|---------|------|
| 0 | DAU retention D7 | +15% |
| 1 | Конверсия Bot → Mini App | >50% |
| 1 | Среднее время сессии | +30% |
| 2 | MAU retention | +25% |
| 3 | LTV | +40% |

---

## Технический стек

| Компонент | Технология |
|-----------|------------|
| Backend | Go 1.22+ |
| Database | PostgreSQL 15 |
| Cache | Redis (Фаза 1+) |
| API | REST + WebSocket |
| Auth | Telegram Init Data + JWT |
| LLM | Существующий LLM Server |
| Hosting | Docker, K8s |

---

## Навигация

- [Фаза 0: Streak-система](./01-streak-system.md)
- [Полное описание проекта](../Full_Project_Description_Explainer.md)
- [ТЗ для дизайнера](../TZ_Designer_Obiasnyatel_DZ.md)
