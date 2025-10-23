package telegram

import (
	"context"
	"fmt"
	"html"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"child-bot/api/internal/llmclient"
	"child-bot/api/internal/ocr"
	"child-bot/api/internal/ocr/types"
	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
)

type Router struct {
	Bot        *tgbotapi.BotAPI
	EngManager *ocr.Manager
	ParseRepo  *store.ParseRepo
	HintRepo   *store.HintRepo
	LLM        *llmclient.Client
	Metrics    *store.MetricsRepo
	History    *store.HistoryRepo
	Session    *store.SessionRepo
}

func (r *Router) HandleCommand(upd tgbotapi.Update) {
	cid := util.GetChatIDByTgUpdate(upd)
	switch upd.Message.Command() {
	case "start":
		r.send(cid, "Пришли фото задачи — верну распознанный текст и подскажу, с чего начать.\nКоманды: /health", nil)
	case "health":
		r.send(cid, "✅ OK", nil)
	case "engine":
		args := strings.Fields(strings.TrimSpace(strings.TrimPrefix(upd.Message.Text, "/engine")))
		cur := r.EngManager.Get(cid)
		if len(args) == 0 {
			r.send(cid, "Текущий LLM-провайдер: "+cur+
				"\nИспользование:\n/engine gemini\n/engine gpt", nil)
			return
		}
		// применим через общий обработчик ниже
		r.handleEngineCommand(cid, upd.Message.Text)
		return
	case "hintL1":
		// Everything after the subcommand is treated as the prompt text
		rest := strings.TrimSpace(strings.TrimPrefix(upd.Message.Text, "/hintL1"))
		if rest == "" {
			r.send(cid, "Использование: /hintL1  <текст промпта>", nil)
			return
		}
		r.postUpdatePrompt(context.Background(), cid, upd.Message.Command(), rest)
		return
	case "hintL2":
		rest := strings.TrimSpace(strings.TrimPrefix(upd.Message.Text, "/hintL2"))
		if rest == "" {
			r.send(cid, "Использование: /hintL2  <текст промпта>", nil)
			return
		}
		r.postUpdatePrompt(context.Background(), cid, upd.Message.Command(), rest)
		return
	case "hintL3":
		rest := strings.TrimSpace(strings.TrimPrefix(upd.Message.Text, "/hintL3"))
		if rest == "" {
			r.send(cid, "Использование: /hintL3  <текст промпта>", nil)
			return
		}
		r.postUpdatePrompt(context.Background(), cid, upd.Message.Command(), rest)
		return
	default:
		r.send(cid, "Неизвестная команда", nil)
	}
}

func (r *Router) HandleUpdate(upd tgbotapi.Update, llmName string) {
	util.PrintInfo("HandleUpdate", llmName, util.GetChatIDByTgUpdate(upd), "Start")
	cid := util.GetChatIDByTgUpdate(upd)

	// r.sendDebug(cid, "telegram_message", upd)
	message := fmt.Sprintf("telegram message: %+v", upd)
	// util.PrintInfo("HandleUpdate", llmName, cid, message)

	cur := getState(cid)
	// r.sendDebug(cid, "last_state", cur)

	if ns, ok := inferNextState(upd, cur); ok && ns != cur {
		// r.sendDebug(cid, "new_state", ns)

		if !canTransition(cur, ns) {
			// Запрещённый переход — сообщим пользователю
			msg := fmt.Sprintf("Нельзя выполнить действие сейчас: %s → %s.%s",
				friendlyState(cur), friendlyState(ns), allowedStateHints(cur))
			b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
			r.send(cid, msg, b)

			return
		}
		// Переход допустим — фиксируем новое состояние
		setState(cid, ns)
	} else if !ok {
		// Запрещённый переход — сообщим пользователю
		msg := fmt.Sprintf("Нельзя выполнить действие сейчас: %s → %s.%s",
			friendlyState(cur), friendlyState(ns), allowedStateHints(cur))
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(cid, msg, b)

		if upd.Message != nil && upd.Message.Text != "" {
			sid := r.ensureSession(cid)
			_ = r.History.Insert(context.Background(), store.TimelineEvent{
				ChatID:        cid,
				TaskSessionID: sid,
				Direction:     "in",
				EventType:     string(cur),
				Text:          upd.Message.Text,
				TgMessageID:   &upd.Message.MessageID,
			})
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

		r.applyTextCorrectionThenShowHints(cid, upd.Message.Text)
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

	// 5) Ветвь выбора пункта при multiple tasks (ожидаем номер из списка)
	if v, ok := pendingChoice.Load(cid); ok && upd.Message.Text != "" {
		setState(cid, AnalyzeChoice)
		sid, _ := r.getSession(cid)
		_ = r.History.Insert(context.Background(), store.TimelineEvent{
			ChatID:        cid,
			TaskSessionID: sid,
			Direction:     "in",
			EventType:     string(AnalyzeChoice),
			Provider:      llmName,
			OK:            true,
			TgMessageID:   &upd.Message.MessageID,
			Text:          upd.Message.Text,
		})

		choices, ok := v.([]TaskChoice)
		if !ok || len(choices) == 0 {
			// Нечего выбирать — очистим и попросим фото снова
			pendingChoice.Delete(cid)
			b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
			r.send(cid, "Не нашёл варианты задач. Пришлите фото ещё раз.", b)
			return
		}

		input := strings.TrimSpace(upd.Message.Text)

		// 1) Пытаемся сопоставить по явному номеру (оригинальному или сгенерированному)
		var chosen *TaskChoice
		for i := range choices {
			if choices[i].Number == input {
				chosen = &choices[i]
				break
			}
		}
		// 2) Fallback: если пользователь прислал порядковый номер 1..N
		if chosen == nil {
			if n, err := strconv.Atoi(input); err == nil && n >= 1 && n <= len(choices) {
				chosen = &choices[n-1]
			}
		}

		if chosen != nil {
			if ctxv, ok2 := pendingCtx.Load(cid); ok2 {
				pendingChoice.Delete(cid)
				pendingCtx.Delete(cid)

				sc := ctxv.(*selectionContext)
				display := fmt.Sprintf("%s — %s", chosen.Number, chosen.Description)
				r.send(cid, "Ок, беру задание: "+display+" — обрабатываю.", nil)

				userID := util.GetUserIDFromTgUpdate(upd)
				r.runParseAndMaybeConfirm(context.Background(), cid, userID, sc, chosen.TaskIndex, display)
				return
			}
			// Нет контекста — сбросим и попросим фото
			pendingChoice.Delete(cid)
			b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
			r.send(cid, "Не нашёл предыдущее изображение. Пришлите фото ещё раз.", b)
			return
		}

		// Неверный ввод — покажем варианты снова
		var lines []string
		for _, c := range choices {
			lines = append(lines, fmt.Sprintf("%s — %s", c.Number, c.Description))
		}
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(cid, "Неверный номер. Выберите один из:\n"+strings.Join(lines, "\n"), b)
		setState(cid, AskChoice)
		return
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
			r.send(cid, "Начинаю нормализацию твоего ответа.", nil)
			r.normalizePhoto(context.Background(), *upd.Message)
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
	message = "Не смог понять, что Вы от меня хотите."
	switch getMode(cid) {
	case "await_solution":
		message += " Я жду от вас фото с решением."
	case "await_new_task":
		message += " Я жду от тебя фото с задачей."
	}

	r.send(cid, message, nil)
}

func (r *Router) send(chatID int64, text string, buttons []tgbotapi.InlineKeyboardButton) {
	msg := tgbotapi.NewMessage(chatID, text)
	if buttons != nil {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons)
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
		Provider:      r.EngManager.Get(chatID),
		TgMessageID:   &m.MessageID,
		Text:          text,
		OK:            true,
	})
}

func (r *Router) sendDebug(chatID int64, name string, v any) {
	if chatID != int64(255509524) {
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

// handleEngineCommand парсит команду /engine и переключает провайдера LLM для чата.
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
		r.EngManager.Set(chatID, "gemini")
		r.send(chatID, "✅ Провайдер LLM: gemini", nil)
	case "gpt", "openai":
		r.EngManager.Set(chatID, "gpt")
		r.send(chatID, "✅ Провайдер LLM: gpt", nil)
	default:
		r.send(chatID, "Неизвестный провайдер. Доступны: gemini | gpt", nil)
	}
}

// PhotoAcceptedText — первый ответ после получения фото/первой страницы альбома.
func (r *Router) PhotoAcceptedText() string {
	return "Фото принято. Если задание на нескольких фото — просто пришлите их подряд, я склею страницы перед обработкой."
}

// postUpdatePrompt sends UpdatePromptRequest to llm-proxy /api/prompt and reports the result back to the chat.
func (r *Router) postUpdatePrompt(ctx context.Context, chatID int64, name, text string) {
	provider := r.EngManager.Get(chatID)

	// Build request payload
	reqBody := types.UpdatePromptRequest{
		Provider: provider,
		Name:     name,
		Text:     text,
	}

	out, err := r.LLM.UpdatePrompt(ctx, reqBody)
	if err != nil {
		r.sendDebug(chatID, "update prompt", err)
	}

	if err != nil {
		// Ответ пришёл с ошибкой
		r.send(chatID, fmt.Sprintf("Не удалось обновить промпт '%s' для провайдера '%s': %v", reqBody.Name, reqBody.Provider, err), nil)
		return
	}
	if !out.OK {
		// Ответ пришёл, но ок == false — покажем пользователю
		r.send(chatID, fmt.Sprintf("Не удалось обновить промпт '%s' для провайдера '%s' (путь: %s)", out.Name, out.Provider, out.Path), nil)
		return
	}

	// Успех
	msg := fmt.Sprintf("✅ Промпт обновлён.\nПровайдер: %s\nИмя: %s\nФайл: %s\nРазмер: %d байт\nОбновлён: %s", out.Provider, out.Name, out.Path, out.Size, out.Updated)
	r.send(chatID, msg, nil)
}
