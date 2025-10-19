package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/ocr/types"
	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
)

type hintSession struct {
	Image        []byte
	Mime         string
	MediaGroupID string
	Parse        types.ParseResult
	Detect       types.DetectResult
	EngineName   string
	NextLevel    int
}

func (r *Router) sendHint(chatID int64, msgID int, hs *hintSession) {
	imgHash := util.SHA256Hex(hs.Image)
	level := hs.NextLevel

	// кэш подсказок
	if hr, err := r.HintRepo.Find(context.Background(), imgHash, hs.EngineName, level, 90*24*time.Hour); err == nil {
		r.send(chatID, formatHint(level, hr), nil)
	} else {
		in := types.HintInput{
			Level:            lvlToConst(level),
			RawText:          hs.Parse.RawText,
			Subject:          hs.Parse.Subject,
			TaskType:         hs.Parse.TaskType,
			Grade:            hs.Parse.Grade,
			SolutionShape:    hs.Parse.SolutionShape,
			TerminologyLevel: levelTerminology(level),
		}
		llmName := r.EngManager.Get(chatID)
		start := time.Now()
		hrNew, err := r.LLM.Hint(context.Background(), llmName, in)
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
			b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
			r.send(chatID, fmt.Sprintf("Не удалось получить подсказку L%d: %s", level, err.Error()), b)
			return
		}
		_ = r.HintRepo.Upsert(context.Background(), imgHash, hs.EngineName, level, hrNew)
		r.send(chatID, formatHint(level, hrNew), nil)
	}
	// После того как отправили подсказку текстом:
	// Отправляем новую клавиатуру с тремя кнопками под НОВЫМ сообщением
	reply := tgbotapi.NewMessage(chatID, "Выберите дальнейшее действие:")
	reply.ReplyMarkup = makeActionsKeyboard(level)
	_, _ = r.Bot.Send(reply)
}

func (r *Router) showTaskAndPrepareHints(chatID int64, sc *selectionContext, pr types.ParseResult, llmName string) {
	var b strings.Builder
	b.WriteString("📄 *Текст задания:*\n```\n")
	if strings.TrimSpace(pr.RawText) != "" {
		b.WriteString(pr.RawText)
	} else {
		b.WriteString("(не удалось чётко переписать текст)")
	}
	b.WriteString("\n```\n")
	if q := strings.TrimSpace(pr.Question); q != "" {
		b.WriteString("\n*Вопрос:* " + q + "\n")
	}

	msg := tgbotapi.NewMessage(chatID, b.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = makeActionsKeyboard(0)
	_, _ = r.Bot.Send(msg)

	// в этом месте бот ждёт дальнейших действий — снимем любые «узкие» режимы
	clearMode(chatID)

	hs := &hintSession{
		Image: sc.Image, Mime: sc.Mime, MediaGroupID: sc.MediaGroupID,
		Parse: pr, Detect: sc.Detect, EngineName: llmName, NextLevel: 1,
	}
	hintState.Store(chatID, hs)
}

func (r *Router) applyTextCorrectionThenShowHints(chatID int64, corrected string) {
	v, ok := parseWait.Load(chatID)
	if !ok {
		return
	}
	p := v.(*parsePending)
	parseWait.Delete(chatID)

	llmName := r.EngManager.Get(chatID)
	imgHash := util.SHA256Hex(p.Sc.Image)

	pr := p.PR
	pr.RawText = corrected
	pr.ConfirmationNeeded = false
	pr.ConfirmationReason = "user_fix"

	_ = r.ParseRepo.Upsert(context.Background(), chatID, p.Sc.MediaGroupID, imgHash, llmName, pr, true, "user_fix")
	r.showTaskAndPrepareHints(chatID, &selectionContext{
		Image: p.Sc.Image, Mime: p.Sc.Mime, MediaGroupID: p.Sc.MediaGroupID, Detect: p.Sc.Detect,
	}, pr, llmName)
}

func formatHint(level int, hr types.HintResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "💡 *Подсказка L%d*: %s\n", level, safe(hr.HintTitle))
	for _, s := range hr.HintSteps {
		if t := strings.TrimSpace(s); t != "" {
			fmt.Fprintf(&b, "• %s\n", safe(t))
		}
	}
	if t := strings.TrimSpace(hr.ControlQuestion); t != "" {
		fmt.Fprintf(&b, "\n*Проверь себя:* %s\n", safe(t))
	}
	// Дополнительные поля (при наличии)
	if hr.RuleHint != "" {
		fmt.Fprintf(&b, "_Подсказка по правилу:_ %s\n", safe(hr.RuleHint))
	}
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
