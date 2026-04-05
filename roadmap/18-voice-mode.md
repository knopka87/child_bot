# 18: Голосовой режим

> Фаза 3 | Приоритет: P2 | Сложность: Высокая | Срок: 7-10 дней

## Цель

Добавить голосовой ввод для ответов ребёнка. Особенно полезно для младших классов.

## Архитектура

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Mini App      │────▶│   API Server    │────▶│   STT Service   │
│   (WebRTC)      │◀────│   (Go)          │◀────│   (Whisper)     │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

## Варианты STT

| Сервис | Плюсы | Минусы |
|--------|-------|--------|
| OpenAI Whisper API | Качество, простота | Цена |
| Google Cloud STT | Streaming | Сложность |
| Self-hosted Whisper | Контроль | Инфра |

## API Endpoints

```
POST /api/v1/voice/transcribe
    Content-Type: multipart/form-data
    audio: <audio file>

Response:
{
    "text": "двадцать пять",
    "confidence": 0.95,
    "normalized": "25"
}
```

## Backend Service

```go
// internal/service/voice/transcriber.go
package voice

import (
    "context"
    "io"

    "github.com/sashabaranov/go-openai"
)

type Transcriber struct {
    client *openai.Client
}

func NewTranscriber(apiKey string) *Transcriber {
    return &Transcriber{
        client: openai.NewClient(apiKey),
    }
}

func (t *Transcriber) Transcribe(ctx context.Context, audio io.Reader, filename string) (*TranscriptionResult, error) {
    resp, err := t.client.CreateTranscription(ctx, openai.AudioRequest{
        Model:    openai.Whisper1,
        Reader:   audio,
        FilePath: filename,
        Language: "ru",
    })
    if err != nil {
        return nil, err
    }

    return &TranscriptionResult{
        Text:       resp.Text,
        Normalized: normalizeAnswer(resp.Text),
    }, nil
}

func normalizeAnswer(text string) string {
    // "двадцать пять" -> "25"
    // "икс равно семь" -> "x = 7"
    // ... нормализация чисел и математических выражений
}
```

## Handler

```go
// internal/api/handlers/voice.go
func (h *VoiceHandler) Transcribe(w http.ResponseWriter, r *http.Request) {
    // Max 10MB
    r.ParseMultipartForm(10 << 20)

    file, header, err := r.FormFile("audio")
    if err != nil {
        http.Error(w, `{"error": "no audio file"}`, http.StatusBadRequest)
        return
    }
    defer file.Close()

    result, err := h.transcriber.Transcribe(r.Context(), file, header.Filename)
    if err != nil {
        http.Error(w, `{"error": "transcription failed"}`, http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(result)
}
```

## Frontend Integration

```typescript
// Mini App
class VoiceRecorder {
  private mediaRecorder: MediaRecorder | null = null;
  private chunks: Blob[] = [];

  async start() {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    this.mediaRecorder = new MediaRecorder(stream);

    this.mediaRecorder.ondataavailable = (e) => {
      this.chunks.push(e.data);
    };

    this.mediaRecorder.start();
  }

  async stop(): Promise<string> {
    return new Promise((resolve) => {
      this.mediaRecorder!.onstop = async () => {
        const blob = new Blob(this.chunks, { type: 'audio/webm' });
        this.chunks = [];

        const formData = new FormData();
        formData.append('audio', blob, 'audio.webm');

        const response = await fetch('/api/v1/voice/transcribe', {
          method: 'POST',
          body: formData,
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        });

        const result = await response.json();
        resolve(result.normalized);
      };

      this.mediaRecorder!.stop();
    });
  }
}
```

## Нормализация ответов

```go
var numberWords = map[string]int{
    "ноль": 0, "один": 1, "два": 2, "три": 3, "четыре": 4,
    "пять": 5, "шесть": 6, "семь": 7, "восемь": 8, "девять": 9,
    "десять": 10, "одиннадцать": 11, "двенадцать": 12,
    // ...
    "двадцать": 20, "тридцать": 30, "сорок": 40, "пятьдесят": 50,
    // ...
}

func normalizeAnswer(text string) string {
    text = strings.ToLower(text)

    // Replace number words
    for word, num := range numberWords {
        text = strings.ReplaceAll(text, word, strconv.Itoa(num))
    }

    // Handle compound numbers: "двадцать пять" -> "25"
    // ...

    // Clean up
    text = strings.TrimSpace(text)

    return text
}
```

## Чек-лист

- [ ] Выбор STT провайдера
- [ ] Backend service для транскрипции
- [ ] REST API endpoint
- [ ] Нормализация чисел (слова → цифры)
- [ ] Frontend запись аудио
- [ ] UI: кнопка микрофона, визуализация
- [ ] Оптимизация: сжатие аудио перед отправкой
- [ ] Rate limiting для voice endpoints
- [ ] Тестирование с разными акцентами

---

[← Family Quests](./17-family-quests.md) | [Назад к Roadmap](./roadmap.md) | [Далее: Mini Games →](./19-mini-games.md)
