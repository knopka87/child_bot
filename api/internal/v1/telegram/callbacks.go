package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	"child-bot/api/internal/v1/types"
)

func (r *Router) handleCallback(cb tgbotapi.CallbackQuery, llmName string) {
	cid := util.GetChatIDFromTgCB(cb)
	data := cb.Data
	_, _ = r.Bot.Request(tgbotapi.NewCallback(cb.ID, "")) // ack
	// log
	message := fmt.Sprintf("llmName: %s, chatID: %d, data: %s, message: %+v", llmName, cid, data, cb.Message)
	util.PrintInfo("handleCallback", llmName, cid, message)
	// r.sendDebug(cid, "message", cb.Message)

	sid, _ := r.getSession(cid)
	_ = r.History.Insert(context.Background(), store.TimelineEvent{
		ChatID:        cid,
		TaskSessionID: sid,
		Direction:     "button",
		EventType:     "callback_" + data,
		Provider:      llmName,
		OK:            true,
		TgMessageID:   &cb.Message.MessageID,
	})

	switch data {
	case "hint_next":
		r.onHintNext(cid, cb.Message.MessageID)
	case "parse_yes":
		r.send(cid, "кнопка нажата", nil)
		r.onParseYes(cid, cb.Message.MessageID)
	case "parse_no":
		r.onParseNo(cid, cb.Message.MessageID)
	case "ready_solution":
		// Скрыть старые кнопки у сообщения с колбэком
		_ = hideKeyboard(cid, cb.Message.MessageID, r)
		setMode(cid, "await_solution")
		r.send(cid, "Отлично! Жду фото с вашим решением. Пришлите, пожалуйста, снимок решения — я проверю без раскрытия ответа.", nil)
	case "analogue_solution":
		_ = hideKeyboard(cid, cb.Message.MessageID, r)
		r.send(cid, "Подбираю похожую задачу. Ожидайте.", nil)
		userID := util.GetUserIDFromTgCB(cb)
		if getState(cid) == Incorrect || getState(cid) == Uncertain {
			r.HandleAnalogueCallback(cid, userID, types.ReasonAfterIncorrect)
		} else {
			r.HandleAnalogueCallback(cid, userID, types.ReasonAfter3Hints)
		}
	case "new_task":
		_ = hideKeyboard(cid, cb.Message.MessageID, r)
		resetContext(cid)
		r.send(cid, "Хорошо! Жду фото новой задачи.", nil)
	case "report":
		resetContext(cid)
		_ = r.SendSessionReport(context.Background(), cid)
	}
}

func (r *Router) onParseYes(chatID int64, msgID int) {
	v, ok := parseWait.Load(chatID)
	if !ok {
		b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
		b = append(b, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report")))
		r.send(chatID, "Контекст подтверждения не найден.", b)
		return
	}
	parseWait.Delete(chatID)
	p := v.(*parsePending)

	sid, _ := r.getSession(chatID)
	_ = r.ParseRepo.MarkAcceptedBySession(context.Background(), sid, "user_yes")
	// убрать клавиатуру
	edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
	_, _ = r.Bot.Send(edit)
	// продолжить
	llmName := r.LlmManager.Get(chatID)
	r.showTaskAndPrepareHints(chatID, p.Sc, p.PR, llmName)
}

func (r *Router) onParseNo(chatID int64, msgID int) {
	edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
	_, _ = r.Bot.Send(edit)
	b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
	b = append(b, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report")))
	r.send(chatID, "Напишите, пожалуйста, текст задания так, как он должен быть прочитан (без ответа). Это поможет дать корректные подсказки.", b)
	// остаёмся в состоянии parseWait — следующий текст примем как корректировку
}

func (r *Router) onHintNext(chatID int64, msgID int) {
	v, ok := hintState.Load(chatID)
	if !ok {
		b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
		b = append(b, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report")))
		r.send(chatID, "Подсказки недоступны: сначала пришлите фото задания.", b)
		return
	}
	hs := v.(*hintSession)
	if hs.NextLevel > 3 {
		edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
		_, _ = r.Bot.Send(edit)
		b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
		b = append(b, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Похожее задание", "analogue_solution"),
			tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"),
		))
		r.send(chatID, "Все подсказки уже показаны. Могу показать аналогичную задачу", b)
		return
	}

	_ = hideKeyboard(chatID, msgID, r)

	r.sendHint(context.Background(), chatID, msgID, hs)

	hs.NextLevel++
	if hs.NextLevel > 3 {
		edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
		_, _ = r.Bot.Send(edit)
	}
}

func (r *Router) GetHintLevel(chatID int64) int {
	v, ok := hintState.Load(chatID)
	if !ok {
		b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
		b = append(b, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report")))
		r.send(chatID, "Подсказки недоступны: сначала пришлите фото задания.", b)
		return 0
	}
	hs := v.(*hintSession)
	return hs.NextLevel - 1
}

func lvlToConst(n int) types.HintLevel {
	switch n {
	case 1:
		return types.HintL1
	case 2:
		return types.HintL2
	default:
		return types.HintL3
	}
}

func hideKeyboard(chatID int64, msgID int, r *Router) error {
	edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
	_, err := r.Bot.Send(edit)
	return err
}

func resetContext(cid int64) {
	// Сброс контекстов
	hintState.Delete(cid)
	pendingChoice.Delete(cid)
	pendingCtx.Delete(cid)
	parseWait.Delete(cid)
	setMode(cid, "await_new_task")
	setState(cid, AwaitingTask)
}
