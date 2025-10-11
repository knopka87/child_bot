package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/llmclient"
	"child-bot/api/internal/ocr"
	"child-bot/api/internal/store"
)

// per-chat предпочтение провайдера LLM: "gemini" или "gpt"
var providerByChat sync.Map

func getProvider(cid int64) string {
	if v, ok := providerByChat.Load(cid); ok {
		if s, ok2 := v.(string); ok2 && s != "" {
			return s
		}
	}
	return "gemini" // значение по умолчанию
}
func setProvider(cid int64, p string) {
	providerByChat.Store(cid, strings.ToLower(strings.TrimSpace(p)))
}

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
		cur := getProvider(cid)
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

	// 4) «Жёсткий» режим: если ждём фото (решение/новая задача) и пришёл произвольный ТЕКСТ — мягко игнорируем.
	// Команды разрешаем, чтобы можно было переключать движки/проверять health.
	if upd.Message.Text != "" && !upd.Message.IsCommand() {
		switch getMode(cid) {
		case "await_solution":
			r.send(cid, "Я жду фото с вашим решением. Пожалуйста, пришлите фото.")
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

	// 7) Фото/альбом — это снимает «режим ожидания фото»
	if len(upd.Message.Photo) > 0 {
		clearMode(cid) // получили фото — разблокируем пайплайн
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
		setProvider(chatID, "gemini")
		r.send(chatID, "✅ Провайдер LLM: gemini")
	case "gpt", "openai":
		setProvider(chatID, "gpt")
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

// CurrentProvider returns per-chat preferred LLM provider ("gemini"|"gpt").
func (r *Router) CurrentProvider(chatID int64) string { return getProvider(chatID) }
