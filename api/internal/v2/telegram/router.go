package telegram

import (
	"context"
	"fmt"
	"html"
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
	LLMClient  *llmclient.Client
	Store      *store.Store
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
		r.send(cid, StartMessageText, nil)
	case "health":
		r.send(cid, OkText, nil)
	default:
		r.send(cid, UnderFoundCommandText, nil)
	}
}

func (r *Router) HandleUpdate(upd tgbotapi.Update, llmName string) {
	ctx := context.Background()
	util.PrintInfo("HandleUpdate", llmName, util.GetChatIDByTgUpdate(upd), "Start")
	cid := util.GetChatIDByTgUpdate(upd)

	// r.sendDebug(cid, "telegram_message", upd)
	stopTyping := r.startTyping(cid, upd.Message, tgbotapi.ChatTyping, 4*time.Second)
	defer stopTyping()

	cur := getState(cid)

	if cur != AwaitGrade {
		if _, ok := userState.Load(cid); !ok {
			user, err := r.Store.FindUserByChatID(ctx, cid)
			if err != nil || user.Grade == nil {
				setState(cid, AwaitGrade)
				r.send(cid, GradePreviewText, makeGradeListButtons())
				return
			}
			userState.Store(cid, user)
		}
	}

	// r.sendDebug(cid, "last_state", cur)

	if ns, ok := inferNextState(upd, cur); ok && ns != cur {
		// r.sendDebug(cid, "new_state", ns)

		if !canTransition(cur, ns) {
			// Запрещённый переход — сообщим пользователю
			msg := fmt.Sprintf("Нельзя выполнить действие сейчас: %s → %s.%s",
				friendlyState(cur), friendlyState(ns), allowedStateHints(cur))
			b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
			b = append(b, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(SendReportButton, "report")))
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
		b = append(b, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📝 Сообщить об ошибке", "report")))
		r.send(cid, msg, b)

		if upd.Message != nil && upd.Message.Text != "" {
			if sid, ok := r.getSession(cid); ok {
				_ = r.Store.InsertHistory(ctx, store.TimelineEvent{
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

	// 4) «Жёсткий» режим ввода: если ждём решение — принимаем и текст, и фото;
	//    если ждём новую задачу — просим фото задачи; в остальных случаях — как раньше.
	if upd.Message.Text != "" && !upd.Message.IsCommand() {
		switch getState(cid) {
		case Report:
			resetContext(cid)
			r.send(cid, SendReportText, nil)
			_ = r.SendSessionReport(ctx, cid, upd.Message.Text)
		case AwaitSolution:
			userID := util.GetUserIDFromTgUpdate(upd)
			r.normalizeText(ctx, cid, userID, upd.Message.Text)
			return
		case AwaitingTask:
			sid, _ := r.getSession(cid)
			_ = r.Store.InsertHistory(ctx, store.TimelineEvent{
				ChatID:        cid,
				TaskSessionID: sid,
				Direction:     "in",
				EventType:     string(AwaitingTask),
				Provider:      llmName,
				OK:            true,
				TgMessageID:   &upd.Message.MessageID,
				Text:          upd.Message.Text,
			})
			r.send(cid, NewTaskText, makeErrorButtons())
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
			r.send(cid, CheckAnswerText, nil)
			r.OCR(ctx, *upd.Message)
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
	message := ""
	switch getMode(cid) {
	case "await_solution":
		message = AwaitSolutionText
	case "await_new_task":
		message = AwaitNewTaskText
	}

	r.send(cid, message, nil)
}

func (r *Router) send(chatID int64, text string, buttons [][]tgbotapi.InlineKeyboardButton) {
	r._sendWithError(chatID, text, "", buttons, nil)
}

func (r *Router) sendMarkdown(chatID int64, text string, buttons [][]tgbotapi.InlineKeyboardButton) {
	r._sendWithError(chatID, text, "Markdown", buttons, nil)
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

func (r *Router) sendError(chatID int64, err error) {
	r._sendWithError(chatID, ErrorText, "", makeErrorButtons(), err)
}

func (r *Router) _sendWithError(chatID int64, text, parseMode string, buttons [][]tgbotapi.InlineKeyboardButton, err error) {
	msg := tgbotapi.NewMessage(chatID, text)
	if buttons != nil {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	}
	if parseMode != "" {
		msg.ParseMode = parseMode
	}

	m, _ := r.Bot.Send(msg)

	sid, _ := r.getSession(chatID)

	if textLen := len(text); textLen > 4000 {
		text = text[:4000] + "…"
	}
	_ = r.Store.InsertHistory(context.Background(), store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Direction:     "out",
		EventType:     "tg_text",
		Provider:      r.LlmManager.Get(chatID),
		TgMessageID:   &m.MessageID,
		Text:          text,
		OK:            err == nil,
		Error:         err,
	})
}

// startTyping sends a repeated chat action (e.g. typing) to the chat.
// NOTE: Some versions of tgbotapi don't expose MessageThreadID on Message.
// If you need per-topic typing in forum chats, upgrade the library and set
// cfg.MessageThreadID at the call site where the thread id is available.
func (r *Router) startTyping(chatID int64, _ *tgbotapi.Message, action string, interval time.Duration) (stop func()) {
	done := make(chan struct{})

	// базовый конфиг; без thread id для совместимости со старыми версиями
	cfg := tgbotapi.NewChatAction(chatID, action)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		_, _ = r.Bot.Request(cfg) // первая отсылка сразу
		for {
			select {
			case <-ticker.C:
				_, _ = r.Bot.Request(cfg)
			case <-done:
				return
			}
		}
	}()
	return func() { close(done) }
}
