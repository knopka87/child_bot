package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	"child-bot/api/internal/v1/types"
)

type hintSession struct {
	Image        []byte
	Mime         string
	MediaGroupID string
	Parse        types.ParseResponse
	Detect       types.DetectResponse
	EngineName   string
	NextLevel    int
}

func (r *Router) sendHint(ctx context.Context, chatID int64, msgID int, hs *hintSession) {
	imgHash := util.SHA256Hex(hs.Image)
	level := hs.NextLevel

	// кэш подсказок
	hc, err := r.HintRepo.Find(ctx, imgHash, hs.EngineName, level)
	if err == nil && time.Since(hc.CreatedAt) <= 90*24*time.Hour {
		var hr types.HintResponse
		_ = json.Unmarshal(hc.HintJson, &hr)
		r.send(chatID, formatHint(level, hr), nil)
	} else {
		in := types.HintRequest{
			RawTaskText: hs.Parse.RawTaskText,
			Level:       lvlToConst(level),
			Grade:       hs.Detect.GradeHint,
			TaskStruct:  hs.Parse.TaskStruct,
			Locale:      "ru_RU",
		}
		hintLevel := level - 1
		for hintLevel > 0 {
			h, err := r.HintRepo.Find(ctx, imgHash, hs.EngineName, hintLevel)
			if err == nil {
				var hr types.HintResponse
				_ = json.Unmarshal(h.HintJson, &hr)
				in.PreviousHints = append(in.PreviousHints, hr.HintText)
			}
			hintLevel--
		}
		llmName := r.LlmManager.Get(chatID)
		start := time.Now()
		hrNew, err := r.GetLLMClient().Hint(context.Background(), llmName, in)
		latency := time.Since(start).Milliseconds()
		sid, _ := r.getSession(chatID)
		_ = r.History.Insert(context.Background(), store.TimelineEvent{
			ChatID:        chatID,
			TaskSessionID: sid,
			Direction:     "api",
			EventType:     string(Hints),
			Provider:      llmName,
			OK:            err == nil,
			LatencyMS:     &latency,
			TgMessageID:   &msgID,
			InputPayload:  in,
			OutputPayload: hrNew,
			Error:         err,
		})
		if err != nil {
			b := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
			b = append(b, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report")))
			r.send(chatID, fmt.Sprintf("Не удалось получить подсказку L%d: %s", level, err.Error()), b)
			return
		}
		js, _ := json.Marshal(hrNew)
		data := store.HintCache{
			CreatedAt: time.Now(),
			Engine:    llmName,
			HintJson:  js,
			Level:     string(lvlToConst(level)),
			ImageHash: imgHash,
		}
		_ = r.HintRepo.Upsert(context.Background(), data)
		r.send(chatID, formatHint(level, hrNew), nil)
	}
	// После того как отправили подсказку текстом:
	// Отправляем новую клавиатуру с тремя кнопками под НОВЫМ сообщением
	reply := makeActionsKeyboardRow(level, true)
	r.send(chatID, "Выберите дальнейшее действие:", reply)
}

func (r *Router) showTaskAndPrepareHints(chatID int64, sc *selectionContext, pr types.ParseResponse, llmName string) {
	var b strings.Builder
	b.WriteString("📄 *Текст задания:*\n```\n")
	if strings.TrimSpace(pr.RawTaskText) != "" {
		b.WriteString(pr.RawTaskText)
	} else {
		b.WriteString("(не удалось чётко переписать текст)")
	}

	buttons := makeActionsKeyboardRow(0, true)
	r.send(chatID, b.String(), buttons)

	// в этом месте бот ждёт дальнейших действий — снимем любые «узкие» режимы
	clearMode(chatID)

	hs := &hintSession{
		Image: sc.Image, Mime: sc.Mime, MediaGroupID: sc.MediaGroupID,
		Parse: pr, Detect: sc.Detect, EngineName: llmName, NextLevel: 1,
	}
	hintState.Store(chatID, hs)
}

func (r *Router) applyTextCorrectionThenShowHints(ctx context.Context, chatID int64, corrected string) {
	v, ok := parseWait.Load(chatID)
	if !ok {
		return
	}
	p := v.(*parsePending)
	parseWait.Delete(chatID)

	llmName := r.LlmManager.Get(chatID)
	imgHash := util.SHA256Hex(p.Sc.Image)
	sid, _ := r.getSession(chatID)

	pr, ok := r.ParseRepo.FindLastConfirmed(ctx, sid)
	if !ok {
		pr = &store.ParsedTasks{
			CreatedAt:         time.Now(),
			ChatID:            chatID,
			SessionID:         sid,
			ImageHash:         imgHash,
			Engine:            llmName,
			RawTaskText:       corrected,
			CombinedSubpoints: false,
			ResultJSON:        make(json.RawMessage, 0),
		}
	}
	pr.NeedsUserConfirmation = false
	pr.Accepted = true
	pr.AcceptReason = "user_fix"

	_ = r.ParseRepo.Upsert(ctx, *pr)

	r.showTaskAndPrepareHints(chatID, &selectionContext{
		Image: p.Sc.Image, Mime: p.Sc.Mime, MediaGroupID: p.Sc.MediaGroupID, Detect: p.Sc.Detect,
	}, p.PR, llmName)
}

func formatHint(level int, hr types.HintResponse) string {
	var b strings.Builder

	// Человеко-понятная подпись уровня в соответствии с промптом:
	// L1 — наводящий вопрос, L2 — практический совет, L3 — общий алгоритм.
	var ruTitle string
	switch hr.Level {
	case types.HintL1:
		ruTitle = "наводящий вопрос"
	case types.HintL2:
		ruTitle = "практический совет"
	case types.HintL3:
		ruTitle = "общий алгоритм"
	default:
		ruTitle = ""
	}

	if ruTitle != "" {
		_, _ = fmt.Fprintf(&b, "💡 *Подсказка L%d* — %s\n", level, ruTitle)
	} else {
		_, _ = fmt.Fprintf(&b, "💡 *Подсказка L%d*\n", level)
	}

	_, _ = fmt.Fprintf(&b, "• %s\n", safe(hr.HintText))

	msg := tgbotapi.NewMessage(0, "") // заглушка для ParseMode
	_ = msg                           // просто, чтобы напомнить: используйте Markdown, поэтому экранируем
	return markdown(b.String())
}

func safe(s string) string {
	// лёгкая защита от Markdown-вставок
	s = strings.ReplaceAll(s, "`", "'")
	s = strings.ReplaceAll(s, "_", "\\_")
	s = strings.ReplaceAll(s, "*", "\\*")
	s = strings.ReplaceAll(s, "[", "\\[")
	return s
}

func markdown(s string) string {
	// Возвращаем как есть — в месте отправки задаём ParseMode=Markdown при необходимости
	return s
}
