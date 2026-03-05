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

	// Безопасное получение MessageID (cb.Message может быть nil для inline callbacks)
	var msgID int
	var tgMsgID *int
	if cb.Message != nil {
		msgID = cb.Message.MessageID
		tgMsgID = &msgID
	}

	sid, _ := r.getSession(cid)
	fromState := string(getState(cid))

	_ = r.Store.InsertHistory(context.Background(), store.TimelineEvent{
		ChatID:        cid,
		TaskSessionID: sid,
		Direction:     "button",
		EventType:     "callback_" + data,
		Provider:      llmName,
		OK:            true,
		TgMessageID:   tgMsgID,
	})

	// Получаем grade пользователя для метрик
	var userGrade int64
	if user, err := r.Store.FindUserByChatID(context.Background(), cid); err == nil && user.Grade != nil {
		userGrade = *user.Grade
	}

	// Метрика действия пользователя
	_ = r.Store.InsertEvent(context.Background(), store.MetricEvent{
		Stage:    "user_action",
		Provider: llmName,
		OK:       true,
		ChatID:   &cid,
		Details: map[string]any{
			"action":     data,
			"from_state": fromState,
			"session_id": sid,
			"grade":      userGrade,
		},
	})

	// Для большинства callback'ов требуется Message
	// Если его нет — игнорируем (кроме grade callbacks, которые не требуют MessageID)
	if cb.Message == nil {
		switch data {
		case "grade1", "grade2", "grade3", "grade4", "report":
			// Эти callbacks не требуют MessageID
		default:
			util.PrintError("handleCallback", llmName, cid, "cb.Message is nil", nil)
			return
		}
	}

	switch data {
	case "hint_next":
		r.onHintNext(cid, msgID)
	case "parse_yes":
		r.onParseYes(cid, msgID)
	case "dont_like_hint":
		r.onDontLikeHint(cid, msgID)
	case "ready_solution":
		sid, _ := r.getSession(cid)
		_ = r.Store.MarkAcceptedParseBySID(context.Background(), sid, "user_yes")
		// Скрыть старые кнопки у сообщения с колбэком
		_ = hideKeyboard(cid, msgID, r)
		r.setModeWithPersist(cid, "await_solution")
		r.send(cid, CheckAnswerClick, makeCheckAnswerClickButtons())
	case "analogue_task":
		_ = hideKeyboard(cid, msgID, r)
		r.send(cid, AnalogueTaskWaitingText, nil)

		// Используем анонимную функцию с defer для гарантированной остановки таймеров
		func() {
			timer1 := r.sendAlert(cid, AnalogueAlert1, 5, 5)
			timer2 := r.sendAlert(cid, AnalogueAlert2, 10, 5)
			timer3 := r.sendAlert(cid, AnalogueAlert3, 15, 5)
			defer timer1.Stop()
			defer timer2.Stop()
			defer timer3.Stop()

			userID := util.GetUserIDFromTgCB(cb)
			if getState(cid) == Incorrect {
				r.HandleAnalogueCallback(cid, userID, types.ReasonAfterIncorrect)
			} else {
				r.HandleAnalogueCallback(cid, userID, types.ReasonAfter3Hints)
			}
		}()
	case "new_task":
		_ = hideKeyboard(cid, msgID, r)
		r.resetContextWithPersist(cid)
		r.send(cid, NewTaskText, nil)
	case "report":
		_ = hideKeyboard(cid, msgID, r)
		r.setStateWithPersist(cid, Report)
		r.send(cid, ReportText, nil)
	case "grade1":
		_ = hideKeyboard(cid, msgID, r)
		r.updateGradeUser(cid, 1)
	case "grade2":
		_ = hideKeyboard(cid, msgID, r)
		r.updateGradeUser(cid, 2)
	case "grade3":
		_ = hideKeyboard(cid, msgID, r)
		r.updateGradeUser(cid, 3)
	case "grade4":
		_ = hideKeyboard(cid, msgID, r)
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
	p, ok := v.(*parsePending)
	if !ok {
		r.sendError(chatID, fmt.Errorf("invalid parse context type"))
		return
	}

	sid, _ := r.getSession(chatID)
	_ = r.Store.MarkAcceptedParseBySID(context.Background(), sid, "user_yes")

	llmName := r.LlmManager.Get(chatID)
	maxHints := 3 // default
	if len(p.PR.Items) > 0 {
		maxHints = p.PR.Items[0].HintPolicy.MaxHints
	}
	hs := &hintSession{
		Image: p.Sc.Image, Mime: p.Sc.Mime, MediaGroupID: p.Sc.MediaGroupID,
		Parse: p.PR, Detect: p.Sc.Detect, EngineName: llmName, NextLevel: 1,
		MaxHints: maxHints,
	}
	hintState.Store(chatID, hs)

	// Сохраняем контекст подсказок в БД для восстановления после редеплоя
	r.saveHintContext(chatID, hs)

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
	hs, ok := v.(*hintSession)
	if !ok {
		r.send(chatID, HintNotFoundText, makeErrorButtons())
		return
	}

	// Защита от concurrent access
	hs.mu.Lock()
	currentLevel := hs.NextLevel
	maxHints := hs.MaxHints
	hs.mu.Unlock()

	if currentLevel > maxHints {
		edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
		_, _ = r.Bot.Send(edit)
		r.send(chatID, HintFinishText, makeFinishHintButtons())
		return
	}

	_ = hideKeyboard(chatID, msgID, r)

	r.sendHint(context.Background(), chatID, msgID, hs)

	// Инкремент уровня под mutex
	hs.mu.Lock()
	hs.NextLevel++
	nextLevel := hs.NextLevel
	hs.mu.Unlock()

	if nextLevel > maxHints {
		edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
		_, _ = r.Bot.Send(edit)
	}
	hintState.Store(chatID, hs)

	// Сохраняем обновлённый контекст в БД
	r.saveHintContext(chatID, hs)
}

func (r *Router) updateGradeUser(cid, grade int64) {
	// Получаем предыдущий grade для метрики
	var prevGrade int64
	if user, err := r.Store.FindUserByChatID(context.Background(), cid); err == nil && user.Grade != nil {
		prevGrade = *user.Grade
	}

	user := store.User{
		ID:    cid,
		Grade: &grade,
	}
	_ = r.Store.UpsertUser(context.Background(), user)
	userInfo.Store(cid, user)

	// Метрика изменения класса
	_ = r.Store.InsertEvent(context.Background(), store.MetricEvent{
		Stage:  "grade_change",
		OK:     true,
		ChatID: &cid,
		Details: map[string]any{
			"prev_grade": prevGrade,
			"new_grade":  grade,
		},
	})

	r.setStateWithPersist(cid, AwaitingTask)
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
	// Очищаем кэш, чтобы не пытаться скрыть те же кнопки повторно
	clearLastButtonMsgID(chatID)
	return err
}
