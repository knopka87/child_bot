package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"child-bot/api/internal/ocr"
	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
)

type selectionContext struct {
	Image        []byte
	Mime         string
	MediaGroupID string
	Detect       ocr.DetectResult
}

type parsePending struct {
	Sc  *selectionContext
	PR  ocr.ParseResult
	LLM string // "gemini"|"gpt"
}

func (r *Router) hasPendingCorrection(chatID int64) bool { _, ok := parseWait.Load(chatID); return ok }
func (r *Router) clearPendingCorrection(chatID int64)    { parseWait.Delete(chatID) }

func (r *Router) runDetectThenParse(ctx context.Context, chatID int64, userID *int64, merged []byte, mediaGroupID string) {
	mime := util.SniffMimeHTTP(merged)
	llmName := r.EngManager.Get(chatID)

	r.sendDebug(chatID, fmt.Sprintf("Detect: merged: %s, mime: %s"))
	// DETECT через llmproxy
	var dres ocr.DetectResult
	start := time.Now()
	if dr, err := r.LLM.Detect(ctx, llmName, merged, mime, 0); err == nil {
		r.sendDebug(chatID, fmt.Sprintf("Detect Res: ```%+v```", dres))
		dres = dr
		errM := r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "detect",
			Provider:   llmName,
			OK:         true,
			DurationMS: time.Since(start).Milliseconds(),
			ChatID:     &chatID,
			UserIDAnon: userID,
			Details: map[string]any{
				"needs_rescan":             dr.NeedsRescan,
				"rescan_reason":            dr.RescanReason,
				"multi_task":               dr.IsMultipleTasks(),
				"final_state":              dr.FinalState,
				"has_faces":                dr.HasFaces,
				"has_diagrams_or_formulas": dr.HasDiagramsOrFormulas,
				"auto_choice_suggested":    dr.AutoChoiceSuggested,
				"pii_detected":             dr.PIIDetected,
			},
		})
		if errM != nil {
			util.PrintError("runDetectThenParse", llmName, chatID, "error insert metrics", errM)
		}
	} else {
		// Мягкий фолбэк: продолжаем без детекта (используем значения по умолчанию),
		// просто логируем ошибку и сообщаем пользователю, что попробуем распознать весь снимок.
		if r.Metrics != nil {
			_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
				Stage:      "detect",
				Provider:   llmName,
				OK:         false,
				DurationMS: time.Since(start).Milliseconds(),
				ChatID:     &chatID,
				UserIDAnon: userID,
				Error:      err.Error(),
			})
		}
		log.Printf("detect failed (chat=%d): %v; fallback to parse without detect", chatID, err)
		r.send(chatID, "ℹ️ Не удалось выделить области на фото, попробую распознать задание целиком.")
	}
	util.PrintInfo("runDetectThenParse", llmName, chatID, fmt.Sprintf("Received a response from LLM: %d", time.Since(start).Milliseconds()))

	// базовая политика
	if dres.FinalState == "inappropriate_image" {
		r.send(chatID, "⚠️ Неподходящее изображение. Пришлите фото учебного задания без личных данных.")
		return
	}
	if dres.FinalState == "not_a_task" {
		r.send(chatID, "ℹ️ Похоже, на фото нет учебного задания. Пришлите фото условия задачи (1–4 класс).")
		return
	}
	if dres.FinalState == "needs_rescan" {
		msg := "Пожалуйста, переснимите фото"
		if dres.RescanReason != "" {
			msg += ": " + dres.RescanReason
		}
		if dres.RescanCode != "" {
			msg += " (код: " + dres.RescanCode + ")"
		}
		r.send(chatID, "📷 "+msg)
		return
	}
	if dres.HasFaces {
		r.send(chatID, "ℹ️ На фото видны лица. Лучше переснять без лиц.")
	}
	if dres.PIIDetected {
		r.send(chatID, "ℹ️ На фото обнаружены личные данные. Пожалуйста, замажьте их или переснимите без них.")
	}

	// несколько заданий — авто-выбор или запрос номера
	if dres.IsMultipleTasks() {
		// собрать список для показа: prefer tasks_brief, иначе из candidates
		tasks := make([]string, 0)
		if len(dres.TasksBrief) > 0 {
			tasks = append(tasks, dres.TasksBrief...)
		} else if len(dres.TasksCandidates) > 0 {
			for _, c := range dres.TasksCandidates {
				tasks = append(tasks, c.Title)
			}
		}

		// можно ли авто-выбрать?
		canAuto := false
		pickedIdx := -1
		if dres.AutoChoiceSuggested != nil && *dres.AutoChoiceSuggested && dres.TopCandidateIndex != nil {
			if *dres.TopCandidateIndex >= 0 && *dres.TopCandidateIndex < len(tasks) && dres.Confidence >= 0.80 {
				canAuto = true
				pickedIdx = *dres.TopCandidateIndex
			}
		}

		if canAuto && pickedIdx >= 0 {
			brief := ""
			if pickedIdx < len(tasks) {
				brief = tasks[pickedIdx]
			}
			sc := &selectionContext{Image: merged, Mime: mime, MediaGroupID: mediaGroupID, Detect: dres}
			r.runParseAndMaybeConfirm(ctx, chatID, userID, sc, pickedIdx, brief)
			return
		}

		// иначе — спросить у пользователя
		if len(tasks) > 0 {
			pendingChoice.Store(chatID, tasks)
			pendingCtx.Store(chatID, &selectionContext{Image: merged, Mime: mime, MediaGroupID: mediaGroupID, Detect: dres})
			var b strings.Builder
			b.WriteString("Нашёл несколько заданий. Выберите номер:\n")
			for i, t := range tasks {
				fmt.Fprintf(&b, "%d) %s\n", i+1, t)
			}
			if dres.DisambiguationQuestion != "" {
				b.WriteString("\n" + dres.DisambiguationQuestion)
			}
			r.send(chatID, b.String())
			return
		}
	}

	// без выбора — сразу PARSE
	r.send(chatID, "Изображение распознано, перехожу к парсингу.")
	sc := &selectionContext{Image: merged, Mime: mime, MediaGroupID: mediaGroupID, Detect: dres}
	r.runParseAndMaybeConfirm(ctx, chatID, userID, sc, -1, "")
	util.PrintInfo("runDetectThenParse", llmName, chatID, fmt.Sprintf("Total time: %d", time.Since(start).Milliseconds()))
}

func (r *Router) runParseAndMaybeConfirm(ctx context.Context, chatID int64, userID *int64, sc *selectionContext, selectedIdx int, selectedBrief string) {
	imgHash := util.SHA256Hex(sc.Image)
	llmName := r.EngManager.Get(chatID)

	// 1) кэш из БД: принят ли PARSE
	if prRow, err := r.ParseRepo.FindByHash(ctx, imgHash, llmName, 30*24*time.Hour); err == nil && prRow.Accepted {
		r.showTaskAndPrepareHints(chatID, sc, prRow.Parse, llmName)
		return
	}

	// 2) LLM.Parse
	start := time.Now()
	pr, err := r.LLM.Parse(ctx, llmName, sc.Image, ocr.ParseOptions{
		SubjectHint: func() string {
			if sc.Detect.FinalState == "recognized_ready_to_parse" {
				return sc.Detect.SubjectGuess
			}
			return ""
		}(),
		ChatID:            chatID,
		MediaGroupID:      sc.MediaGroupID,
		ImageHash:         imgHash,
		SelectedTaskIndex: selectedIdx,
		SelectedTaskBrief: selectedBrief,
	})
	if err != nil {
		_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "parse",
			Provider:   llmName,
			OK:         false,
			Error:      err.Error(),
			DurationMS: time.Since(start).Milliseconds(),
			ChatID:     &chatID,
			UserIDAnon: userID,
		})
		util.PrintError("runParseAndMaybeConfirm", llmName, chatID, "parse", err)
		r.SendError(chatID, fmt.Errorf("parse: %w", err))
		return
	}
	_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
		Stage:      "parse",
		Provider:   llmName,
		OK:         false,
		DurationMS: time.Since(start).Milliseconds(),
		ChatID:     &chatID,
		UserIDAnon: userID,
		Details: map[string]any{
			"final_state":    pr.FinalState,
			"rescan_reason":  pr.RescanReason,
			"confirm_reason": pr.ConfirmationReason,
			"grade_aligment": pr.GradeAlignment,
			"grade":          pr.Grade,
			"solution_shape": pr.SolutionShape,
			"need_rescan":    pr.NeedsRescan,
			"confidence":     pr.Confidence,
		},
	})
	util.PrintInfo("runParseAndMaybeConfirm", llmName, chatID, fmt.Sprintf("Received a response from LLM: %d", time.Since(start).Milliseconds()))

	// сохранить черновик
	errP := r.ParseRepo.Upsert(ctx, chatID, sc.MediaGroupID, imgHash, llmName, pr, false, "")
	if errP != nil {
		util.PrintError("runParseAndMaybeConfirm", llmName, chatID, "error upsert parsed_tasks", errP)
	}

	// 3) подтверждение, если нужно
	if pr.ConfirmationNeeded {
		r.askParseConfirmation(chatID, pr)
		parseWait.Store(chatID, &parsePending{Sc: sc, PR: pr, LLM: llmName})
		return
	}

	// 4) автоподтверждение
	_ = r.ParseRepo.MarkAccepted(ctx, imgHash, llmName, "auto")
	r.showTaskAndPrepareHints(chatID, sc, pr, llmName)
	util.PrintInfo("runParseAndMaybeConfirm", llmName, chatID, fmt.Sprintf("total time: %d", time.Since(start).Milliseconds()))
}
