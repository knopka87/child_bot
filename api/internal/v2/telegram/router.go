package telegram

import (
	"context"
	"fmt"
	"html"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/llmclient"
	"child-bot/api/internal/service"
	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	llmclientv2 "child-bot/api/internal/v2/llmclient"
)

type Router struct {
	Bot        BotSender
	LlmManager *service.LlmManager
	LLMClient  *llmclient.Client
	Store      *store.Store
}

func (r *Router) GetToken() string {
	return r.Bot.GetToken()
}

func (r *Router) GetLLMClient() *llmclientv2.Client {
	return llmclientv2.New(r.LLMClient)
}

func (r *Router) HandleCommand(upd tgbotapi.Update) {
	cid := util.GetChatIDByTgUpdate(upd)
	switch upd.Message.Command() {
	case "start":
		r.resetContextWithPersist(cid)
		r.send(cid, StartMessageText, nil)
	case "health":
		r.send(cid, OkText, nil)
	case "cachestats":
		if IsAdmin(cid) {
			r.sendDebug(cid, "cache_stats", GetCacheStats())
		}
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

	// Восстанавливаем состояние из БД если его нет в кэше (после редеплоя)
	r.restoreStateFromDB(cid)

	// Скрываем кнопки предыдущего сообщения при получении нового сообщения от пользователя
	r.hidePreviousButtons(cid)

	cur := getState(cid)

	if cur != AwaitGrade {
		if _, ok := userInfo.Load(cid); !ok {
			user, err := r.Store.FindUserByChatID(ctx, cid)
			if err != nil || user.Grade == nil {
				r.setStateWithPersist(cid, AwaitGrade)
				r.send(cid, GradePreviewText, makeGradeListButtons())
				return
			}
			userInfo.Store(cid, user)
		}
	}

	chatVal, chatOk := chatInfo.Load(cid)
	chat, chatTypeOk := chatVal.(store.Chat)
	if !chatOk || !chatTypeOk || chat.Username == nil || *chat.Username == "" {
		chat, err := r.Store.FindChatByID(ctx, cid)
		if err != nil || chat.Username == nil || *chat.Username == "" {
			chat = store.Chat{
				ID: cid,
			}
			if upd.Message != nil && upd.Message.Chat != nil {
				chat.Type = &upd.Message.Chat.Type
				chat.Username = &upd.Message.Chat.UserName
				chat.FirstName = &upd.Message.Chat.FirstName
				chat.LastName = &upd.Message.Chat.LastName
			}
			if chat.Username == nil || *chat.Username == "" {
				if upd.Message != nil && upd.Message.From != nil {
					chat.Username = &upd.Message.From.UserName
					chat.FirstName = &upd.Message.From.FirstName
					chat.LastName = &upd.Message.From.LastName
				}
			}
			_ = r.Store.UpsertChat(ctx, chat)
		}
		chatInfo.Store(cid, chat)
	}

	// r.sendDebug(cid, "last_state", cur)

	if ns, inferred := inferNextState(upd, cur); inferred && ns != cur {
		// r.sendDebug(cid, "new_state", ns)

		// Используем атомарный переход состояния
		actualCur, transitioned := tryTransition(cid, ns)
		if !transitioned {
			// Запрещённый переход — сообщим пользователю
			msg := fmt.Sprintf("Нельзя выполнить действие сейчас: %s → %s.%s",
				friendlyState(actualCur), friendlyState(ns), allowedStateHints(actualCur))
			b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
			b = append(b, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(SendReportButton, "report")))
			r.send(cid, msg, b)
			return
		}
		// Переход выполнен успешно
	} else if !inferred {
		// Не удалось определить следующее состояние — сообщим пользователю
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
			r.resetContextWithPersist(cid)
			r.send(cid, SendReportText, nil)
			_ = r.SendSessionReport(ctx, cid, upd.Message.Text)
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
			// Фото с ответом ученика — прогресс показывается в checkSolution
			userID := util.GetUserIDFromTgUpdate(upd)
			r.checkSolution(ctx, cid, userID, *upd.Message)
			r.clearModeWithPersist(cid)
			return
		}
		// Иначе — это фото задачи/страницы
		r.clearModeWithPersist(cid)
		// Session создаётся только если нет активного batch (для альбомов session один на все фото)
		r.ensureSessionForNewTask(cid, upd.Message.MediaGroupID)

		r.acceptPhoto(cid, *upd.Message)
		return
	}

	// 7.1) Документ-изображение (фото отправлено как файл)
	if upd.Message.Document != nil && strings.HasPrefix(upd.Message.Document.MimeType, "image/") {
		if getMode(cid) == "await_solution" {
			// Фото с ответом ученика — прогресс показывается в checkSolutionFromDocument
			userID := util.GetUserIDFromTgUpdate(upd)
			r.checkSolutionFromDocument(ctx, cid, userID, *upd.Message)
			r.clearModeWithPersist(cid)
			return
		}
		r.clearModeWithPersist(cid)
		// Session создаётся только если нет активного batch (для альбомов session один на все фото)
		r.ensureSessionForNewTask(cid, upd.Message.MediaGroupID)

		r.acceptDocument(cid, *upd.Message)
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

// hidePreviousButtons скрывает кнопки у предыдущего сообщения, если они есть
func (r *Router) hidePreviousButtons(chatID int64) {
	if msgID := getLastButtonMsgID(chatID); msgID > 0 {
		edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
		_, _ = r.Bot.Send(edit)
		clearLastButtonMsgID(chatID)
	}
}

func (r *Router) sendMarkdown(chatID int64, text string, buttons [][]tgbotapi.InlineKeyboardButton) {
	r._sendWithError(chatID, text, "Markdown", buttons, nil)
}

func (r *Router) sendAlert(chatID int64, text string, postpone, delay time.Duration) *time.Timer {
	shutdown := GetShutdownManager()

	// Не запускаем если идёт shutdown
	if shutdown.IsShutdown() {
		return time.NewTimer(0) // возвращаем dummy timer
	}

	// Регистрируем горутину ДО создания таймера
	// Используем канал для сигнализации о завершении
	outerDone := shutdown.TrackGoroutine()

	return time.AfterFunc(postpone*time.Second, func() {
		defer outerDone() // освобождаем регистрацию при выходе из callback

		// Проверяем shutdown перед отправкой
		if shutdown.IsShutdown() {
			return
		}

		msg := tgbotapi.NewMessage(chatID, text)
		sent, err := r.Bot.Send(msg)
		if err != nil {
			return
		}

		// Регистрируем внутреннюю горутину ДО создания таймера
		innerDone := shutdown.TrackGoroutine()

		time.AfterFunc(delay*time.Second, func() {
			defer innerDone()

			if shutdown.IsShutdown() {
				return
			}

			del := tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: sent.MessageID}
			_, _ = r.Bot.Request(del)
		})
	})
}

// startProgress отправляет сообщение с прогрессом и периодически его обновляет.
// stages — этапы прогресса, interval — интервал между этапами.
// Возвращает функцию stop(), которую нужно вызвать после завершения операции.
func (r *Router) startProgress(chatID int64, stages []string, interval time.Duration) func() {
	shutdown := GetShutdownManager()

	if shutdown.IsShutdown() || len(stages) == 0 {
		return func() {}
	}

	msg := tgbotapi.NewMessage(chatID, stages[0])
	sent, err := r.Bot.Send(msg)
	if err != nil {
		return func() {}
	}

	done := make(chan struct{})
	var once sync.Once

	outerDone := shutdown.TrackGoroutine()

	go func() {
		defer outerDone()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		stage := 1
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if shutdown.IsShutdown() {
					return
				}
				if stage < len(stages) {
					edit := tgbotapi.NewEditMessageText(chatID, sent.MessageID, stages[stage])
					_, _ = r.Bot.Send(edit)
					stage++
				}
			}
		}
	}()

	return func() {
		once.Do(func() {
			close(done)
			if !shutdown.IsShutdown() {
				del := tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: sent.MessageID}
				_, _ = r.Bot.Request(del)
			}
		})
	}
}

// startParseProgress — прогресс для этапа парсинга.
func (r *Router) startParseProgress(chatID int64) func() {
	stages := []string{ParseProgress1, ParseProgress2, ParseProgress3, ParseProgress4}
	return r.startProgress(chatID, stages, 3*time.Second)
}

// startHintProgress — прогресс для генерации подсказки.
func (r *Router) startHintProgress(chatID int64) func() {
	stages := []string{HintProgress1, HintProgress2, HintProgress3, HintProgress4}
	return r.startProgress(chatID, stages, 10*time.Second)
}

// startCheckProgress — прогресс для проверки решения.
func (r *Router) startCheckProgress(chatID int64) func() {
	stages := []string{CheckProgress1, CheckProgress2, CheckProgress3, CheckProgress4}
	return r.startProgress(chatID, stages, 4*time.Second)
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
	_ = r.SendSessionReport(context.Background(), chatID, "Внимание!! Техническая ошибка!")
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

	// Отслеживаем сообщение с кнопками для последующего скрытия
	if buttons != nil && m.MessageID > 0 {
		setLastButtonMsgID(chatID, m.MessageID)
	}

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
	shutdown := GetShutdownManager()

	// Не запускаем новые горутины если идёт shutdown
	if shutdown.IsShutdown() {
		return func() {}
	}

	// базовый конфиг; без thread id для совместимости со старыми версиями
	cfg := tgbotapi.NewChatAction(chatID, action)

	// Регистрируем горутину для отслеживания
	goroutineDone := shutdown.TrackGoroutine()

	go func() {
		defer goroutineDone()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		_, _ = r.Bot.Request(cfg) // первая отсылка сразу
		for {
			select {
			case <-ticker.C:
				_, _ = r.Bot.Request(cfg)
			case <-done:
				return
			case <-shutdown.Done():
				// Graceful shutdown — завершаем горутину
				return
			}
		}
	}()

	var once sync.Once
	return func() {
		once.Do(func() { close(done) })
	}
}
