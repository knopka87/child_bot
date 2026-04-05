# 17: Семейные квесты (Family Quests)

> Фаза 2 | Приоритет: P3 | Сложность: Средняя | Срок: 3-4 дня

## Цель

Совместные задания для родителя и ребёнка. Укрепление связи через обучение.

## Типы квестов

| Тип | Описание | Награда |
|-----|----------|---------|
| `weekend_challenge` | Решить 10 задач за выходные | Семейный бейдж |
| `topic_master` | Освоить тему вместе | XP x2 |
| `streak_together` | Оба не пропускают дни | Редкий предмет |
| `teach_me` | Родитель объясняет задачу | Бонус XP |

## Миграция

```sql
-- migrations/039_family_quests.up.sql

CREATE TABLE family_quest (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    quest_type TEXT NOT NULL,
    target_value INT NOT NULL,
    reward_type TEXT,        -- xp, achievement, cosmetic
    reward_value TEXT,
    duration_days INT
);

CREATE TABLE family_quest_progress (
    id SERIAL PRIMARY KEY,
    parent_id BIGINT REFERENCES "user"(chat_id),
    child_id BIGINT REFERENCES "user"(chat_id),
    quest_id TEXT REFERENCES family_quest(id),
    current_value INT DEFAULT 0,
    started_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    reward_claimed BOOLEAN DEFAULT false,
    UNIQUE(parent_id, child_id, quest_id, started_at::date)
);
```

## Store методы

```go
func (s *Store) GetActiveQuests(ctx context.Context, parentID, childID int64) ([]QuestWithProgress, error) {
    // ...
}

func (s *Store) StartQuest(ctx context.Context, parentID, childID int64, questID string) error {
    // ...
}

func (s *Store) UpdateQuestProgress(ctx context.Context, progressID int64, delta int) (*QuestProgress, error) {
    // Возвращает completed=true если цель достигнута
}

func (s *Store) ClaimQuestReward(ctx context.Context, progressID int64) error {
    // Выдаёт награду обоим участникам
}
```

## API Endpoints

```
GET  /api/v1/family/quests                    # Доступные квесты
POST /api/v1/family/quests/{id}/start         # Начать квест
GET  /api/v1/family/quests/active             # Активные квесты
POST /api/v1/family/quests/{progressId}/claim # Забрать награду
```

## Чек-лист

- [ ] Миграция `039_family_quests.up.sql`
- [ ] Store методы
- [ ] REST API
- [ ] Интеграция с task completion
- [ ] Уведомления о прогрессе

---

[← Leaderboard](./16-leaderboard.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Voice Mode →](./18-voice-mode.md)
