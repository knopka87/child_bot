# 19: Мини-игры (Brain Breaks)

> Фаза 3 | Приоритет: P3 | Сложность: Очень высокая | Срок: 2-3 недели

## Цель

Короткие образовательные мини-игры для разнообразия и отдыха между задачами.

## Игры

| Игра | Описание | Навык |
|------|----------|-------|
| `catch_numbers` | Ловить падающие числа | Арифметика |
| `word_builder` | Собрать слово из букв | Грамотность |
| `memory_cards` | Найти пары | Память |
| `bubble_pop` | Лопать пузыри с ответами | Скорость |

## Архитектура

Игры работают преимущественно на клиенте. Backend только:
- Сохраняет результаты
- Выдаёт награды
- Генерирует данные для игр

## Миграция

```sql
-- migrations/040_mini_games.up.sql

CREATE TABLE mini_game_session (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES "user"(chat_id),
    game_type TEXT NOT NULL,
    score INT NOT NULL,
    duration_seconds INT,
    difficulty TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE mini_game_highscore (
    user_id BIGINT REFERENCES "user"(chat_id),
    game_type TEXT NOT NULL,
    best_score INT NOT NULL,
    achieved_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, game_type)
);
```

## API Endpoints

```
GET  /api/v1/games                     # Список доступных игр
GET  /api/v1/games/{type}/data         # Данные для игры
POST /api/v1/games/{type}/complete     # Завершение игры
GET  /api/v1/games/{type}/highscores   # Таблица рекордов
```

## Генерация данных для игр

```go
// internal/service/games/generator.go

func (g *Generator) GenerateCatchNumbers(difficulty string) *CatchNumbersData {
    var maxNum int
    switch difficulty {
    case "easy":
        maxNum = 10
    case "medium":
        maxNum = 50
    case "hard":
        maxNum = 100
    }

    // Генерируем 10 примеров
    var items []CatchNumberItem
    for i := 0; i < 10; i++ {
        a := rand.Intn(maxNum)
        b := rand.Intn(maxNum)
        items = append(items, CatchNumberItem{
            Expression: fmt.Sprintf("%d + %d", a, b),
            Answer:     a + b,
        })
    }

    return &CatchNumbersData{Items: items}
}

func (g *Generator) GenerateWordBuilder(grade int) *WordBuilderData {
    // Выбираем слово из словаря для класса
    word := g.wordDict.RandomWord(grade)
    letters := shuffleString(word)

    return &WordBuilderData{
        Letters: letters,
        Hint:    getWordHint(word),
    }
}

func (g *Generator) GenerateMemoryCards(topic string) *MemoryCardsData {
    // Пары: число - слово, операция - результат
    pairs := []MemoryPair{
        {A: "5 + 3", B: "8"},
        {A: "2 × 4", B: "8"},
        {A: "10 - 2", B: "8"},
        // ...
    }
    return &MemoryCardsData{Pairs: pairs[:6]}
}
```

## Сохранение результатов

```go
func (s *Store) SaveGameResult(ctx context.Context, userID int64, gameType string, score int, duration int) (*GameResult, error) {
    tx, _ := s.DB.BeginTx(ctx, nil)
    defer tx.Rollback()

    // Сохраняем сессию
    tx.ExecContext(ctx, `
        INSERT INTO mini_game_session (user_id, game_type, score, duration_seconds)
        VALUES ($1, $2, $3, $4)
    `, userID, gameType, score, duration)

    // Обновляем рекорд если нужно
    var isNewRecord bool
    err := tx.QueryRowContext(ctx, `
        INSERT INTO mini_game_highscore (user_id, game_type, best_score, achieved_at)
        VALUES ($1, $2, $3, NOW())
        ON CONFLICT (user_id, game_type) DO UPDATE
        SET best_score = GREATEST(mini_game_highscore.best_score, $3),
            achieved_at = CASE WHEN $3 > mini_game_highscore.best_score THEN NOW() ELSE mini_game_highscore.achieved_at END
        RETURNING best_score = $3
    `, userID, gameType, score).Scan(&isNewRecord)

    // XP за игру
    xpGain := score / 10 // 1 XP за каждые 10 очков
    tx.ExecContext(ctx, `UPDATE "user" SET xp = xp + $2 WHERE chat_id = $1`, userID, xpGain)

    return &GameResult{
        Score:       score,
        XPGained:    xpGain,
        IsNewRecord: isNewRecord,
    }, tx.Commit()
}
```

## Frontend (пример: Catch Numbers)

```typescript
// Игра на Canvas
class CatchNumbersGame {
  private canvas: HTMLCanvasElement;
  private ctx: CanvasRenderingContext2D;
  private score: number = 0;
  private items: FallingItem[] = [];
  private currentTarget: number;

  constructor(canvas: HTMLCanvasElement, data: CatchNumbersData) {
    this.canvas = canvas;
    this.ctx = canvas.getContext('2d')!;
    this.data = data;
    this.nextTarget();
  }

  private nextTarget() {
    const item = this.data.items.shift();
    if (!item) {
      this.endGame();
      return;
    }
    this.currentTarget = item.answer;
    this.showExpression(item.expression);
    this.spawnNumbers();
  }

  private spawnNumbers() {
    // Спавним правильный ответ и несколько неправильных
    const numbers = [this.currentTarget];
    while (numbers.length < 5) {
      const wrong = this.currentTarget + (Math.random() > 0.5 ? 1 : -1) * Math.floor(Math.random() * 10);
      if (!numbers.includes(wrong) && wrong > 0) {
        numbers.push(wrong);
      }
    }

    // Создаём падающие элементы
    shuffle(numbers).forEach((num, i) => {
      this.items.push({
        value: num,
        x: (i + 0.5) * (this.canvas.width / 5),
        y: -50,
        speed: 2 + Math.random(),
      });
    });
  }

  private update() {
    // Обновляем позиции
    this.items.forEach(item => {
      item.y += item.speed;
    });

    // Удаляем упавшие
    this.items = this.items.filter(item => item.y < this.canvas.height);

    this.draw();
    requestAnimationFrame(() => this.update());
  }

  private onTap(x: number, y: number) {
    const tapped = this.items.find(item =>
      Math.abs(item.x - x) < 30 && Math.abs(item.y - y) < 30
    );

    if (tapped) {
      if (tapped.value === this.currentTarget) {
        this.score += 10;
        this.showSuccess();
      } else {
        this.score -= 5;
        this.showError();
      }
      this.nextTarget();
    }
  }

  private async endGame() {
    const result = await api.completeGame('catch_numbers', this.score);
    this.showResults(result);
  }
}
```

## Чек-лист

- [ ] Миграция `040_mini_games.up.sql`
- [ ] Backend: генератор данных для игр
- [ ] REST API endpoints
- [ ] Сохранение результатов и highscores
- [ ] Frontend: реализация игр
- [ ] Интеграция наград (XP, achievements)
- [ ] Баланс сложности
- [ ] Тестирование UX

---

[← Voice Mode](./18-voice-mode.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Seasonal Events →](./20-seasonal-events.md)
