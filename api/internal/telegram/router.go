package telegram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/llmclient"
	"child-bot/api/internal/ocr"
	"child-bot/api/internal/store"
)

type Router struct {
	Bot        *tgbotapi.BotAPI
	EngManager *ocr.Manager
	ParseRepo  *store.ParseRepo
	HintRepo   *store.HintRepo
	LLM        *llmclient.Client
}

func (r *Router) HandleCommand(upd tgbotapi.Update) {
	cid := upd.Message.Chat.ID
	switch upd.Message.Command() {
	case "start":
		r.send(cid, "Пришли фото задачи — верну распознанный текст и подскажу, с чего начать.\nКоманды: /health, /engine (gemini|gpt)")
	case "health":
		r.send(cid, "✅ OK")
	case "engine":
		args := strings.Fields(strings.TrimSpace(strings.TrimPrefix(upd.Message.Text, "/engine")))
		cur := r.EngManager.Get(cid)
		if len(args) == 0 {
			r.send(cid, "Текущий LLM-провайдер: "+cur+
				"\nИспользование:\n/engine gemini\n/engine gpt")
			return
		}
		// применим через общий обработчик ниже
		r.handleEngineCommand(cid, upd.Message.Text)
		return
	default:
		r.send(cid, "Неизвестная команда")
	}
}

func (r *Router) HandleUpdate(upd tgbotapi.Update) {
	// 1) Callback-кнопки
	if upd.CallbackQuery != nil {
		r.handleCallback(*upd.CallbackQuery)
		return
	}

	// 2) Сообщений нет — выходим
	if upd.Message == nil {
		return
	}
	cid := upd.Message.Chat.ID

	// 3) Если ждём текстовую правку после «Нет» — приоритетно принимаем её
	if r.hasPendingCorrection(cid) && upd.Message.Text != "" {
		r.applyTextCorrectionThenShowHints(cid, upd.Message.Text)
		return
	}

	// 4) «Жёсткий» режим ввода: если ждём решение — принимаем и текст, и фото;
	//    если ждём новую задачу — просим фото задачи; в остальных случаях — как раньше.
	if upd.Message.Text != "" && !upd.Message.IsCommand() {
		switch getMode(cid) {
		case "await_solution":
			// Нормализуем текстовый ответ ученика
			r.normalizeText(context.Background(), cid, upd.Message.Text)
			return
		case "await_new_task":
			r.send(cid, "Я жду фото новой задачи. Пожалуйста, пришлите фото.")
			return
		}
	}

	// 5) Ветвь выбора пункта при multiple tasks (ожидаем число 1..N)
	if v, ok := pendingChoice.Load(cid); ok && upd.Message.Text != "" {
		briefs := v.([]string)
		if n, err := strconv.Atoi(strings.TrimSpace(upd.Message.Text)); err == nil && n >= 1 && n <= len(briefs) {
			if ctxv, ok2 := pendingCtx.Load(cid); ok2 {
				pendingChoice.Delete(cid)
				pendingCtx.Delete(cid)
				sc := ctxv.(*selectionContext)
				r.send(cid, fmt.Sprintf("Ок, беру задание: %s — обрабатываю.", briefs[n-1]))
				r.runParseAndMaybeConfirm(context.Background(), cid, sc, n-1, briefs[n-1])
				return
			}
			pendingChoice.Delete(cid)
			r.send(cid, "Не нашёл предыдущее изображение. Пришлите фото ещё раз.")
			return
		}
		// иначе ждём корректный номер
	}

	// 6) Команды (в т.ч. /engine)
	if upd.Message.IsCommand() && strings.HasPrefix(upd.Message.Text, "/engine") {
		r.handleEngineCommand(cid, upd.Message.Text)
		return
	}
	if upd.Message.IsCommand() {
		r.HandleCommand(upd)
		return
	}

	// 7) Фото/альбом
	if len(upd.Message.Photo) > 0 {
		if getMode(cid) == "await_solution" {
			// Фото с ответом ученика → нормализация
			r.normalizePhoto(context.Background(), *upd.Message)
			clearMode(cid)
			return
		}
		// Иначе — это фото задачи/страницы
		clearMode(cid)
		r.acceptPhoto(*upd.Message)
		return
	}

	// 8) Остальное — игнорируем
}

func (r *Router) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, _ = r.Bot.Send(msg)
}

func (r *Router) SendResult(chatID int64, text string) {
	if len(text) > 3900 {
		text = text[:3900] + "…"
	}
	r.send(chatID, "📝 Распознанный текст:\n\n"+text)
}

func (r *Router) SendError(chatID int64, err error) {
	r.send(chatID, fmt.Sprintf("Ошибка OCR: %v", err))
}

// handleEngineCommand парсит команду /engine и переключает провайдера LLM для чата.
// Поддерживаются только gemini и gpt.
func (r *Router) handleEngineCommand(chatID int64, cmd string) {
	args := strings.Fields(strings.TrimSpace(strings.TrimPrefix(cmd, "/engine")))
	if len(args) == 0 {
		r.send(chatID, "Использование: /engine {gemini|gpt}")
		return
	}
	name := strings.ToLower(args[0])
	switch name {
	case "gemini", "google":
		r.EngManager.Set(chatID, "gemini")
		r.send(chatID, "✅ Провайдер LLM: gemini")
	case "gpt", "openai":
		r.EngManager.Set(chatID, "gpt")
		r.send(chatID, "✅ Провайдер LLM: gpt")
	default:
		r.send(chatID, "Неизвестный провайдер. Доступны: gemini | gpt")
	}
}

// Показ запроса подтверждения распознанного текста
func (r *Router) askParseConfirmation(chatID int64, pr ocr.ParseResult) {
	var b strings.Builder
	b.WriteString("Я так прочитал задание. Всё верно?\n")
	if s := strings.TrimSpace(pr.RawText); s != "" {
		b.WriteString("```\n")
		b.WriteString(s)
		b.WriteString("\n```\n")
	}
	if q := strings.TrimSpace(pr.Question); q != "" {
		b.WriteString("\nВопрос: ")
		b.WriteString(esc(q))
		b.WriteString("\n")
	}

	msg := tgbotapi.NewMessage(chatID, b.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = makeParseConfirmKeyboard()
	_, _ = r.Bot.Send(msg)
}

// PhotoAcceptedText — первый ответ после получения фото/первой страницы альбома.
func (r *Router) PhotoAcceptedText() string {
	return "Фото принято. Если задание на нескольких фото — просто пришлите их подряд, я склею страницы перед обработкой."
}

// normalizeText — отправляет текст ученика на нормализацию в LLM-прокси
func (r *Router) normalizeText(ctx context.Context, chatID int64, text string) {
	llmName := r.EngManager.Get(chatID)
	shape := r.suggestSolutionShape(chatID)
	in := ocr.NormalizeInput{
		SolutionShape: shape,
		Provider:      llmName,
		Answer:        ocr.NormalizeAnswer{Source: "text", Text: strings.TrimSpace(text)},
	}
	res, err := r.LLM.Normalize(ctx, llmName, in)
	if err != nil {
		r.send(chatID, fmt.Sprintf("Не удалось нормализовать ответ: %v", err))
		return
	}
	r.sendNormalizePreview(chatID, res)
	clearMode(chatID)
}

// normalizePhoto — скачивает фото из Telegram и отправляет на нормализацию
func (r *Router) normalizePhoto(ctx context.Context, msg tgbotapi.Message) {
	if len(msg.Photo) == 0 {
		return
	}
	llmName := r.EngManager.Get(msg.Chat.ID)
	ph := msg.Photo[len(msg.Photo)-1] // самое большое
	data, mime, err := r.downloadFileBytes(ph.FileID)
	if err != nil {
		r.send(msg.Chat.ID, fmt.Sprintf("Не удалось получить фото: %v", err))
		return
	}
	shape := r.suggestSolutionShape(msg.Chat.ID)
	in := ocr.NormalizeInput{
		SolutionShape: shape,
		Provider:      llmName,
		Answer:        ocr.NormalizeAnswer{Source: "photo", PhotoB64: string(data), Mime: mime},
	}
	res, err := r.LLM.Normalize(ctx, llmName, in)
	if err != nil {
		r.send(msg.Chat.ID, fmt.Sprintf("Не удалось нормализовать ответ (фото): %v", err))
		return
	}
	r.sendNormalizePreview(msg.Chat.ID, res)
}

// suggestSolutionShape — простая эвристика: если по парсингу известна форма — берём её, иначе number
func (r *Router) suggestSolutionShape(chatID int64) string {
	// TODO: можно взять из последнего ParseResult из БД (ParseRepo) subject/task_type → shape
	return "number"
}

// sendNormalizePreview — короткий текст для пользователя по NormalizeResult
func (r *Router) sendNormalizePreview(chatID int64, nr ocr.NormalizeResult) {
	shape := strings.ToLower(strings.TrimSpace(nr.Shape))
	val := ""
	switch v := nr.Value.(type) {
	case string:
		val = v
	case float64:
		val = strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		val = strconv.Itoa(v)
	case []string:
		val = strings.Join(v, "; ")
	default:
		val = "(не удалось отобразить значение)"
	}
	b := &strings.Builder{}
	b.WriteString("✅ Принял ответ. Форма: ")
	b.WriteString(shape)
	if val != "" {
		b.WriteString("\nЗначение: ")
		b.WriteString(val)
	}
	if nr.UncertainReasons != nil && len(nr.UncertainReasons) > 0 {
		b.WriteString("\nПредупреждения: ")
		b.WriteString(strings.Join(nr.UncertainReasons, ", "))
	}
	if nr.NeedsClarification && nr.NeedsUserActionMessage != "" {
		b.WriteString("\nНужно уточнение: ")
		b.WriteString(nr.NeedsUserActionMessage)
	}
	r.send(chatID, b.String())
}

// downloadFileBytes — скачивает файл Telegram по fileID и возвращает bytes и mime
func (r *Router) downloadFileBytes(fileID string) ([]byte, string, error) {
	url, err := r.Bot.GetFileDirectURL(fileID)
	if err != nil {
		return nil, "", err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	mime := resp.Header.Get("Content-Type")
	if mime == "" {
		mime = "image/jpeg"
	}
	return b, mime, nil
}
