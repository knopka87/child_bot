# 12: Эволюция питомца

> Фаза 2 | Приоритет: P1 | Сложность: Средняя | Срок: 3-4 дня

## Цель

Расширить систему питомцев: 4 стадии эволюции, выбор из нескольких видов.

## Стадии эволюции

| Стадия | Название | XP требование | Визуал |
|--------|----------|---------------|--------|
| 1 | Яйцо | 0 | Яйцо с узором |
| 2 | Малыш | 100 XP | Маленький дракончик |
| 3 | Подросток | 500 XP | Средний дракон |
| 4 | Взрослый | 2000 XP | Полноразмерный дракон |

## Виды питомцев

| Тип | Название | Цвет | Особенность |
|-----|----------|------|-------------|
| `dragon` | Дракон Знаний | Синий | +10% XP за математику |
| `phoenix` | Феникс Мудрости | Оранжевый | +10% XP за русский |
| `owl` | Сова Эрудиции | Фиолетовый | +5% ко всему |
| `fox` | Лис Хитрости | Рыжий | Двойной streak бонус |

## Миграции

```sql
-- migrations/034_pet_evolution.up.sql

CREATE TABLE pet_type (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    color TEXT,
    bonus_type TEXT,     -- xp_math, xp_russian, xp_all, streak
    bonus_value INT      -- процент бонуса
);

INSERT INTO pet_type VALUES
('dragon', 'Дракон Знаний', 'Любит математику', '#3B82F6', 'xp_math', 10),
('phoenix', 'Феникс Мудрости', 'Мастер русского языка', '#F97316', 'xp_russian', 10),
('owl', 'Сова Эрудиции', 'Знает всё понемногу', '#8B5CF6', 'xp_all', 5),
('fox', 'Лис Хитрости', 'Хитрый и быстрый', '#EA580C', 'streak', 100);

ALTER TABLE user_pet ADD COLUMN pet_xp INT NOT NULL DEFAULT 0;
ALTER TABLE user_pet ADD COLUMN can_evolve BOOLEAN NOT NULL DEFAULT false;
```

## Store методы

```go
func (s *Store) AddPetXP(ctx context.Context, userID int64, xp int) (*PetXPResult, error) {
    tx, err := s.DB.BeginTx(ctx, nil)
    // ...

    var pet UserPet
    err = tx.QueryRowContext(ctx, `
        UPDATE user_pet
        SET pet_xp = pet_xp + $2
        WHERE user_id = $1
        RETURNING pet_xp, evolution_stage
    `, userID, xp).Scan(&pet.PetXP, &pet.EvolutionStage)

    // Check evolution threshold
    thresholds := []int{0, 100, 500, 2000}
    canEvolve := pet.EvolutionStage < 4 && pet.PetXP >= thresholds[pet.EvolutionStage]

    if canEvolve {
        tx.ExecContext(ctx, `UPDATE user_pet SET can_evolve = true WHERE user_id = $1`, userID)
    }

    return &PetXPResult{
        NewXP:      pet.PetXP,
        CanEvolve:  canEvolve,
        NextStage:  pet.EvolutionStage + 1,
        XPToNext:   thresholds[min(pet.EvolutionStage, 3)] - pet.PetXP,
    }, tx.Commit()
}

func (s *Store) EvolvePet(ctx context.Context, userID int64) (*UserPet, error) {
    var pet UserPet
    err := s.DB.QueryRowContext(ctx, `
        UPDATE user_pet
        SET evolution_stage = evolution_stage + 1, can_evolve = false
        WHERE user_id = $1 AND can_evolve = true AND evolution_stage < 4
        RETURNING *
    `, userID).Scan(/* ... */)
    return &pet, err
}

func (s *Store) ChangePetType(ctx context.Context, userID int64, newType string) error {
    // Можно менять только на стадии 1 (яйцо)
    _, err := s.DB.ExecContext(ctx, `
        UPDATE user_pet SET pet_type = $2
        WHERE user_id = $1 AND evolution_stage = 1
    `, userID, newType)
    return err
}
```

## API Endpoints

```
GET  /api/v1/pet                    # Текущее состояние
POST /api/v1/pet/evolve             # Эволюционировать
POST /api/v1/pet/type               # Сменить вид (только яйцо)
GET  /api/v1/pet/types              # Доступные виды
```

## Чек-лист

- [ ] Миграция `034_pet_evolution.up.sql`
- [ ] Store методы для эволюции
- [ ] Накопление pet_xp при решении задач
- [ ] REST API endpoints
- [ ] Бонусы от типа питомца
- [ ] Анимации эволюции (ТЗ дизайнеру)

---

[← Parent-Child](./11-parent-child.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Customization →](./13-customization.md)
