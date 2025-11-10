package telegram

import (
	"context"
	"fmt"
	"html"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"child-bot/api/internal/llmclient"
	"child-bot/api/internal/service"
	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	llmclientv2 "child-bot/api/internal/v2/llmclient"
)

type Router struct {
	Bot        *tgbotapi.BotAPI
	LlmManager *service.LlmManager
	ParseRepo  *store.ParseRepo
	HintRepo   *store.HintRepo
	LLMClient  *llmclient.Client
	Metrics    *store.MetricsRepo
	History    *store.HistoryRepo
	Session    *store.SessionRepo
}

func (r *Router) GetToken() string {
	return r.Bot.Token
}

func (r *Router) GetLLMClient() *llmclientv2.Client {
	return llmclientv2.New(r.LLMClient)
}

func (r *Router) HandleCommand(upd tgbotapi.Update) {
	cid := util.GetChatIDByTgUpdate(upd)
	switch upd.Message.Command() {
	case "start":
		resetContext(cid)
		r.send(cid, "👋 Ура, мы начинаем!\n\n\nПогнали! 🎒\nСкидывай своё задание — и разберёмся вместе! 🤓", nil)
	case "health":
		r.send(cid, "✅ OK", nil)
	// case "engine":
	// 	args := strings.Fields(strings.TrimSpace(strings.TrimPrefix(upd.Message.Text, "/engine")))
	// 	cur := r.LlmManager.Get(cid)
	// 	if len(args) == 0 {
	// 		r.send(cid, "Текущий LLMClient-провайдер: "+cur+
	// 			"\nИспользование:\n/engine gemini\n/engine gpt", nil)
	// 		return
	// 	}
	// 	// применим через общий обработчик ниже
	// 	r.handleEngineCommand(cid, upd.Message.Text)
	// 	return
	// case "hintL1":
	// 	// Everything after the subcommand is treated as the prompt text
	// 	rest := strings.TrimSpace(strings.TrimPrefix(upd.Message.Text, "/hintL1"))
	// 	if rest == "" {
	// 		r.send(cid, "Использование: /hintL1  <текст промпта>", nil)
	// 		return
	// 	}
	// 	r.postUpdatePrompt(context.Background(), cid, upd.Message.Command(), rest)
	// 	return
	// case "hintL2":
	// 	rest := strings.TrimSpace(strings.TrimPrefix(upd.Message.Text, "/hintL2"))
	// 	if rest == "" {
	// 		r.send(cid, "Использование: /hintL2  <текст промпта>", nil)
	// 		return
	// 	}
	// 	r.postUpdatePrompt(context.Background(), cid, upd.Message.Command(), rest)
	// 	return
	// case "hintL3":
	// 	rest := strings.TrimSpace(strings.TrimPrefix(upd.Message.Text, "/hintL3"))
	// 	if rest == "" {
	// 		r.send(cid, "Использование: /hintL3  <текст промпта>", nil)
	// 		return
	// 	}
	// 	r.postUpdatePrompt(context.Background(), cid, upd.Message.Command(), rest)
	// 	return
	default:
		r.send(cid, "Неизвестная команда. Я знаю только команду /start", nil)
	}
}

func (r *Router) HandleUpdate(upd tgbotapi.Update, llmName string) {
	util.PrintInfo("HandleUpdate", llmName, util.GetChatIDByTgUpdate(upd), "Start")
	cid := util.GetChatIDByTgUpdate(upd)

	// r.sendDebug(cid, "telegram_message", upd)
	done := make(chan struct{})
	defer close(done)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		// сразу показать индикатор
		_, _ = r.Bot.Send(tgbotapi.NewChatAction(cid, tgbotapi.ChatTyping))
		for {
			select {
			case <-ticker.C:
				_, _ = r.Bot.Send(tgbotapi.NewChatAction(cid, tgbotapi.ChatTyping))
			case <-done:
				return
			}
		}
	}()

	cur := getState(cid)
	// r.sendDebug(cid, "last_state", cur)

	if ns, ok := inferNextState(upd, cur); ok && ns != cur {
		// r.sendDebug(cid, "new_state", ns)

		if !canTransition(cur, ns) {
			// Запрещённый переход — сообщим пользователю
			msg := fmt.Sprintf("Нельзя выполнить действие сейчас: %s → %s.%s",
				friendlyState(cur), friendlyState(ns), allowedStateHints(cur))
			b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
			b = append(b, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report")))
			r.send(cid, msg, b)

			return
		}
		// Переход допустим — фиксируем новое состояние
		setState(cid, ns)
	} else if !ok {
		// Запрещённый переход — сообщим пользователю
		msg := fmt.Sprintf("Нельзя выполнить действие сейчас: %s → %s.%s",
			friendlyState(cur), friendlyState(ns), allowedStateHints(cur))
		b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
		b = append(b, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report")))
		r.send(cid, msg, b)

		if upd.Message != nil && upd.Message.Text != "" {
			if sid, ok := r.getSession(cid); ok {
				_ = r.History.Insert(context.Background(), store.TimelineEvent{
					ChatID:        cid,
					TaskSessionID: sid,
					Direction:     "in",
					EventType:     string(cur),
					Text:          upd.Message.Text,
					TgMessageID:   &upd.Message.MessageID,
				})
			}
		}
		return
	}

	// 1) Callback-кнопки
	if upd.CallbackQuery != nil {
		r.handleCallback(*upd.CallbackQuery, llmName)
		return
	}

	// 2) Сообщений нет — выходим
	if upd.Message == nil {
		util.PrintInfo("HandleUpdate", llmName, 0, "not found telegram message")
		return
	}

	// 3) Если ждём текстовую правку после «Нет» — приоритетно принимаем её
	if r.hasPendingCorrection(cid) && upd.Message.Text != "" {
		sid, _ := r.getSession(cid)
		_ = r.History.Insert(context.Background(), store.TimelineEvent{
			ChatID:        cid,
			TaskSessionID: sid,
			Direction:     "in",
			EventType:     "pending_correction",
			Provider:      llmName,
			OK:            true,
			TgMessageID:   &upd.Message.MessageID,
			Text:          upd.Message.Text,
		})

		r.applyTextCorrectionThenShowHints(context.Background(), cid, upd.Message.Text)
		return
	}

	// 4) «Жёсткий» режим ввода: если ждём решение — принимаем и текст, и фото;
	//    если ждём новую задачу — просим фото задачи; в остальных случаях — как раньше.
	if upd.Message.Text != "" && !upd.Message.IsCommand() {
		switch getState(cid) {
		case AwaitSolution:
			// Нормализуем текстовый ответ ученика
			r.send(cid, "Начинаю нормализацию твоего ответа.", nil)
			userID := util.GetUserIDFromTgUpdate(upd)
			r.normalizeText(context.Background(), cid, userID, upd.Message.Text)
			return
		case AwaitingTask:
			sid, _ := r.getSession(cid)
			_ = r.History.Insert(context.Background(), store.TimelineEvent{
				ChatID:        cid,
				TaskSessionID: sid,
				Direction:     "in",
				EventType:     string(AwaitingTask),
				Provider:      llmName,
				OK:            true,
				TgMessageID:   &upd.Message.MessageID,
				Text:          upd.Message.Text,
			})
			r.send(cid, "Я жду фото новой задачи. Пожалуйста, пришлите фото.", nil)
			return
		}
	}

	// 6) Команды (в т.ч. /engine)
	// if upd.Message.IsCommand() && strings.HasPrefix(upd.Message.Text, "/engine") {
	// 	r.handleEngineCommand(cid, upd.Message.Text)
	// 	return
	// }
	if upd.Message.IsCommand() {
		r.HandleCommand(upd)
		return
	}

	// 7) Фото/альбом
	if len(upd.Message.Photo) > 0 {
		if getMode(cid) == "await_solution" {
			// Фото с ответом ученика → OCR
			r.send(cid, "Начинаю парсинг твоего ответа.", nil)
			r.OCR(context.Background(), *upd.Message)
			clearMode(cid)
			return
		}
		// Иначе — это фото задачи/страницы
		clearMode(cid)
		sid := uuid.NewString()
		r.setSession(cid, sid)

		r.acceptPhoto(cid, *upd.Message)
		return
	}

	// 8) Остальное — игнорируем
	message := "Не смог понять, что Вы от меня хотите."
	switch getMode(cid) {
	case "await_solution":
		message += " Я жду от вас фото с решением."
	case "await_new_task":
		message += " Я жду от тебя фото с задачей."
	}

	r.send(cid, message, nil)
}

func (r *Router) send(chatID int64, text string, buttons [][]tgbotapi.InlineKeyboardButton) {
	msg := tgbotapi.NewMessage(chatID, text)
	if buttons != nil {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	}

	m, _ := r.Bot.Send(msg)

	sid, _ := r.getSession(chatID)

	if textLen := len(text); textLen > 4000 {
		text = text[:4000] + "…"
	}
	_ = r.History.Insert(context.Background(), store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Direction:     "out",
		EventType:     "tg_text",
		Provider:      r.LlmManager.Get(chatID),
		TgMessageID:   &m.MessageID,
		Text:          text,
		OK:            true,
	})
}

func (r *Router) sendDebug(chatID int64, name string, v any) {
	find := false
	for _, adminID := range adminsChatID {
		if chatID == adminID {
			find = true
			break
		}
	}
	if !find {
		return
	}

	const limit = 4096 // лимит длины сообщения в Telegram
	raw := util.PrettyJSON(v)
	// экранируем HTML-символы и оборачиваем в pre/code
	body := name + ":\n<pre><code class=\"language-json\">" + html.EscapeString(raw) + "</code></pre>"

	// если не помещается — отправим как файл
	if len(body) > limit {
		r.sendJSONAsDocument(chatID, []byte(raw), name+".json")
		return
	}

	msg := tgbotapi.NewMessage(chatID, body)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	_, _ = r.Bot.Send(msg)

}

func (r *Router) sendJSONAsDocument(chatID int64, data []byte, filename string) {
	doc := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{
		Name:  filename,
		Bytes: data,
	})
	_, _ = r.Bot.Send(doc)
	return
}

// func (r *Router) SendResult(chatID int64, text string) {
// 	if len(text) > 3900 {
// 		text = text[:3900] + "…"
// 	}
// 	r.send(chatID, "📝 Распознанный текст:\n\n"+text)
// }

func (r *Router) SendError(chatID int64, err error) {
	r.send(chatID, fmt.Sprintf("Ошибка OCR: %v", err), nil)
}

// handleEngineCommand парсит команду /engine и переключает провайдера LLMClient для чата.
// Поддерживаются только gemini и gpt.
func (r *Router) handleEngineCommand(chatID int64, cmd string) {
	args := strings.Fields(strings.TrimSpace(strings.TrimPrefix(cmd, "/engine")))
	if len(args) == 0 {
		r.send(chatID, "Использование: /engine {gemini|gpt}", nil)
		return
	}
	name := strings.ToLower(args[0])
	switch name {
	case "gemini", "google":
		r.LlmManager.Set(chatID, "gemini")
		r.send(chatID, "✅ Провайдер LLMClient: gemini", nil)
	case "gpt", "openai":
		r.LlmManager.Set(chatID, "gpt")
		r.send(chatID, "✅ Провайдер LLMClient: gpt", nil)
	default:
		r.send(chatID, "Неизвестный провайдер. Доступны: gemini | gpt", nil)
	}
}
