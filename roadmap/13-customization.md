# 13: Кастомизация маскота

> Фаза 2 | Приоритет: P2 | Сложность: Высокая | Срок: 5-7 дней

## Цель

Система кастомизации маскота: гардероб, аксессуары, магазин за очки.

## Компоненты

### Категории предметов

| Категория | Примеры | Слот |
|-----------|---------|------|
| Головные уборы | Шапка, корона, ушки | head |
| Очки | Солнцезащитные, учёные | eyes |
| Аксессуары | Шарф, галстук, медаль | accessory |
| Фон | Космос, лес, школа | background |

### Источники предметов

- **Достижения** — награда за разблокировку
- **Магазин** — покупка за XP
- **Сезоны** — ограниченные предметы

## Миграция

```sql
-- migrations/035_customization.up.sql

CREATE TABLE cosmetic_item (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,      -- head, eyes, accessory, background
    rarity TEXT NOT NULL,        -- common, rare, legendary
    price_xp INT,                -- NULL = не продаётся
    achievement_id TEXT REFERENCES achievement(id),
    season TEXT,                 -- NULL = всегда доступен
    asset_key TEXT NOT NULL,     -- ключ для frontend
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_cosmetic (
    user_id BIGINT REFERENCES "user"(chat_id) ON DELETE CASCADE,
    item_id TEXT REFERENCES cosmetic_item(id),
    obtained_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, item_id)
);

CREATE TABLE user_equipped (
    user_id BIGINT PRIMARY KEY REFERENCES "user"(chat_id) ON DELETE CASCADE,
    head_item TEXT REFERENCES cosmetic_item(id),
    eyes_item TEXT REFERENCES cosmetic_item(id),
    accessory_item TEXT REFERENCES cosmetic_item(id),
    background_item TEXT REFERENCES cosmetic_item(id),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Начальные предметы
INSERT INTO cosmetic_item (id, name, category, rarity, price_xp, asset_key) VALUES
('hat_student', 'Шапка выпускника', 'head', 'common', 100, 'hat_student'),
('glasses_nerd', 'Очки отличника', 'eyes', 'common', 50, 'glasses_nerd'),
('medal_gold', 'Золотая медаль', 'accessory', 'rare', 500, 'medal_gold'),
('bg_space', 'Космос', 'background', 'rare', 300, 'bg_space');
```

## Store методы

```go
func (s *Store) GetUserCosmetics(ctx context.Context, userID int64) ([]CosmeticItem, error) {
    // Все предметы пользователя
}

func (s *Store) GetShopItems(ctx context.Context, userID int64) ([]ShopItem, error) {
    // Предметы в магазине, которых у пользователя нет
}

func (s *Store) PurchaseItem(ctx context.Context, userID int64, itemID string) error {
    // Покупка за XP
    tx, _ := s.DB.BeginTx(ctx, nil)
    defer tx.Rollback()

    // Проверяем XP и цену
    var userXP, price int
    tx.QueryRowContext(ctx, `SELECT xp FROM "user" WHERE chat_id = $1 FOR UPDATE`, userID).Scan(&userXP)
    tx.QueryRowContext(ctx, `SELECT price_xp FROM cosmetic_item WHERE id = $1`, itemID).Scan(&price)

    if userXP < price {
        return fmt.Errorf("insufficient XP")
    }

    // Списываем XP и добавляем предмет
    tx.ExecContext(ctx, `UPDATE "user" SET xp = xp - $2 WHERE chat_id = $1`, userID, price)
    tx.ExecContext(ctx, `INSERT INTO user_cosmetic VALUES ($1, $2)`, userID, itemID)

    return tx.Commit()
}

func (s *Store) EquipItem(ctx context.Context, userID int64, itemID string) error {
    // Надеть предмет
}

func (s *Store) GetEquippedItems(ctx context.Context, userID int64) (*EquippedItems, error) {
    // Текущий наряд
}
```

## API Endpoints

```
GET  /api/v1/cosmetics              # Инвентарь пользователя
GET  /api/v1/cosmetics/shop         # Магазин
POST /api/v1/cosmetics/purchase     # Купить
POST /api/v1/cosmetics/equip        # Надеть
GET  /api/v1/cosmetics/equipped     # Текущий наряд
```

## Чек-лист

- [ ] Миграция `035_customization.up.sql`
- [ ] Store методы
- [ ] REST API
- [ ] Интеграция с достижениями (выдача предметов)
- [ ] Сезонные предметы (Фаза 3)
- [ ] ТЗ на ассеты для дизайнера

---

[← Pet Evolution](./12-pet-evolution.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Knowledge Map →](./14-knowledge-map.md)
