package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	"child-bot/api/internal/v2/types"
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
	_ = r.Store.InsertHistory(context.Background(), store.TimelineEvent{
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
		r.onParseYes(cid, cb.Message.MessageID)
	case "dont_like_hint":
		r.onDontLikeHint(cid, cb.Message.MessageID)
	case "ready_solution":
		sid, _ := r.getSession(cid)
		_ = r.Store.MarkAcceptedParseBySID(context.Background(), sid, "user_yes")
		// Скрыть старые кнопки у сообщения с колбэком
		_ = hideKeyboard(cid, cb.Message.MessageID, r)
		setMode(cid, "await_solution")
		r.send(cid, CheckAnswerClick, makeCheckAnswerClickButtons())
	case "analogue_task":
		_ = hideKeyboard(cid, cb.Message.MessageID, r)
		r.send(cid, AnalogueTaskWaitingText, nil)

		timer1 := r.sendAlert(cid, AnalogueAlert1, 4, 4)
		timer2 := r.sendAlert(cid, AnalogueAlert2, 8, 4)
		timer3 := r.sendAlert(cid, AnalogueAlert3, 12, 4)

		userID := util.GetUserIDFromTgCB(cb)
		if getState(cid) == Incorrect {
			r.HandleAnalogueCallback(cid, userID, types.ReasonAfterIncorrect)
		} else {
			r.HandleAnalogueCallback(cid, userID, types.ReasonAfter3Hints)
		}
		timer3.Stop()
		timer2.Stop()
		timer1.Stop()
	case "new_task":
		_ = hideKeyboard(cid, cb.Message.MessageID, r)
		resetContext(cid)
		r.send(cid, NewTaskText, nil)
	case "report":
		setState(cid, Report)
		r.send(cid, ReportText, nil)
	case "grade1":
		r.updateGradeUser(cid, 1)
	case "grade2":
		r.updateGradeUser(cid, 2)
	case "grade3":
		r.updateGradeUser(cid, 3)
	case "grade4":
		r.updateGradeUser(cid, 4)
	}
}

func (r *Router) onParseYes(chatID int64, msgID int) {
	v, ok := parseWait.Load(chatID)
	if !ok {
		r.sendError(chatID, fmt.Errorf("not found Parse"))
		return
	}
	parseWait.Delete(chatID)
	p := v.(*parsePending)

	sid, _ := r.getSession(chatID)
	_ = r.Store.MarkAcceptedParseBySID(context.Background(), sid, "user_yes")

	llmName := r.LlmManager.Get(chatID)
	hs := &hintSession{
		Image: p.Sc.Image, Mime: p.Sc.Mime, MediaGroupID: p.Sc.MediaGroupID,
		Parse: p.PR, Detect: p.Sc.Detect, EngineName: llmName, NextLevel: 1,
	}
	hintState.Store(chatID, hs)

	r.onHintNext(chatID, msgID)

}

func (r *Router) onDontLikeHint(chatID int64, msgID int) {
	r.send(chatID, DontLikeHint, nil)
	r.onHintNext(chatID, msgID)
}

func (r *Router) onHintNext(chatID int64, msgID int) {
	v, ok := hintState.Load(chatID)
	if !ok {
		r.send(chatID, HintNotFoundText, makeErrorButtons())
		return
	}
	hs := v.(*hintSession)
	if hs.NextLevel > 3 {
		edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
		_, _ = r.Bot.Send(edit)
		r.send(chatID, HintFinishText, makeFinishHintButtons())
		return
	}

	_ = hideKeyboard(chatID, msgID, r)

	r.sendHint(context.Background(), chatID, msgID, hs)

	hs.NextLevel++
	if hs.NextLevel > 3 {
		edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
		_, _ = r.Bot.Send(edit)
	}
	hintState.Store(chatID, hs)
}

func (r *Router) updateGradeUser(cid, grade int64) {
	user := store.User{
		ID:    cid,
		Grade: &grade,
	}
	_ = r.Store.UpsertUser(context.Background(), user)
	userInfo.Store(cid, user)
	setState(cid, AwaitingTask)
	r.send(cid, StartMessageText, nil)
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
	pendingCtx.Delete(cid)
	parseWait.Delete(cid)
	setMode(cid, "await_new_task")
	setState(cid, AwaitingTask)
}
