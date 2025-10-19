package telegram

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/ocr/types"
	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
)

type selectionContext struct {
	Image        []byte
	Mime         string
	MediaGroupID string
	Detect       types.DetectResult
}

type parsePending struct {
	Sc  *selectionContext
	PR  types.ParseResult
	LLM string // "gemini"|"gpt"
}

func (r *Router) runParseAndMaybeConfirm(ctx context.Context, chatID int64, userID *int64, sc *selectionContext, selectedIdx int, selectedBrief string) {
	setState(chatID, Parse)
	imgHash := util.SHA256Hex(sc.Image)
	llmName := r.EngManager.Get(chatID)

	// 1) Проверка кэша: если уже было подтверждено ранее — используем сразу
	if prRow, err := r.ParseRepo.FindByHash(ctx, imgHash, llmName, 30*24*time.Hour); err == nil && prRow.Accepted {
		r.showTaskAndPrepareHints(chatID, sc, prRow.Parse, llmName)
		return
	}

	// 2) Запрос к LLM.Parse по новой схеме (v1.2)
	in := types.ParseInput{
		ImageB64: base64.StdEncoding.EncodeToString(sc.Image),
		Options: types.ParseOptions{
			SubjectHint:       "",
			ChatID:            chatID,
			MediaGroupID:      sc.MediaGroupID,
			ImageHash:         imgHash,
			SelectedTaskIndex: selectedIdx,
			SelectedTaskBrief: selectedBrief,
		},
	}
	start := time.Now()
	pr, err := r.LLM.Parse(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	sid, _ := r.getSession(chatID)
	_ = r.History.Insert(ctx, store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Direction:     "api",
		EventType:     string(Parse),
		Provider:      llmName,
		OK:            err == nil,
		LatencyMS:     &latency,
		InputPayload:  in,
		OutputPayload: pr,
		Error:         err,
	})
	if err != nil {
		_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "parse",
			Provider:   llmName,
			OK:         false,
			Error:      err.Error(),
			DurationMS: latency,
			ChatID:     &chatID,
			UserIDAnon: userID,
		})

		util.PrintError("runParseAndMaybeConfirm", llmName, chatID, "parse", err)
		r.SendError(chatID, fmt.Errorf("parse: %w", err))
		return
	}

	// 3) Метрики строго по новой структуре
	_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
		Stage:      "parse",
		Provider:   llmName,
		OK:         true,
		DurationMS: time.Since(start).Milliseconds(),
		ChatID:     &chatID,
		UserIDAnon: userID,
		Details: map[string]any{
			"final_state":                 pr.FinalState,
			"rescan_reason":               pr.RescanReason,
			"confirmation_reason":         pr.ConfirmationReason,
			"grade_alignment":             pr.GradeAlignment,
			"grade":                       pr.Grade,
			"solution_shape":              pr.SolutionShape,
			"needs_rescan":                pr.NeedsRescan,
			"confidence":                  pr.Confidence,
			"meaning_change_risk":         pr.MeaningChangeRisk,
			"bracketed_spans_count":       pr.BracketedSpansCount,
			"original_number":             pr.OriginalNumber,
			"has_subparts":                pr.HasSubparts,
			"subparts_labels_len":         len(pr.SubpartsLabels),
			"solution_fragments_detected": pr.SolutionFragmentsDetected,
			"stripped_solution_spans_len": len(pr.StrippedSolutionSpans),
			"has_diagrams_or_formulas":    pr.HasDiagramsOrFormulas,
			"routing_hint":                pr.RoutingHint,
			"attachments_used_len":        len(pr.AttachmentsUsed),
			"subject":                     pr.Subject,
			"task_type":                   pr.TaskType,
		},
	})
	util.PrintInfo("runParseAndMaybeConfirm", llmName, chatID, fmt.Sprintf("Received a response from LLM: %d", time.Since(start).Milliseconds()))

	// 4) Сохраняем черновик PARSE в БД
	if errP := r.ParseRepo.Upsert(ctx, chatID, sc.MediaGroupID, imgHash, llmName, pr, false, ""); errP != nil {
		util.PrintError("runParseAndMaybeConfirm", llmName, chatID, "error upsert parsed_tasks", errP)
	}

	// 5) Если нужен рескан — сообщаем и выходим
	if pr.NeedsRescan {
		setState(chatID, NeedsRescan)
		msg := pr.RescanReason
		if strings.TrimSpace(msg) == "" {
			msg = "Нужно переснять фото: постарайтесь сделать его чётким, без бликов и поближе к задаче."
		}
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, "ℹ️ "+msg, b)
		return
	}

	// 6) Если требуется подтверждение — спрашиваем пользователя
	if pr.ConfirmationNeeded {
		setState(chatID, Confirm)
		r.askParseConfirmation(chatID, pr)
		parseWait.Store(chatID, &parsePending{Sc: sc, PR: pr, LLM: llmName})
		return
	}

	// 7) Иначе — автоподтверждение и переход к подсказкам
	setState(chatID, AutoPick)
	_ = r.ParseRepo.MarkAccepted(ctx, imgHash, llmName, "auto")
	r.showTaskAndPrepareHints(chatID, sc, pr, llmName)
	util.PrintInfo("runParseAndMaybeConfirm", llmName, chatID, fmt.Sprintf("total time: %d", time.Since(start).Milliseconds()))
}

// Показ запроса подтверждения распознанного текста
func (r *Router) askParseConfirmation(chatID int64, pr types.ParseResult) {
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
