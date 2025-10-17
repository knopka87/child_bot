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

	r.sendDebug(chatID, "merged", merged)
	r.sendDebug(chatID, "mime", mime)
	// DETECT через llmproxy
	var dres ocr.DetectResult
	start := time.Now()
	dr, err := r.LLM.Detect(ctx, llmName, merged, mime, 0)
	if err == nil {
		dres = dr
		r.sendDebug(chatID, "detect_res", dres)

		// агрегируем флаги по задачам
		tasksCount := len(dres.Tasks)
		hasFacesAny := false
		piiAny := false
		multipleDetected := false
		for _, t := range dres.Tasks {
			if t.HasFaces {
				hasFacesAny = true
			}
			if t.PIIDetected {
				piiAny = true
			}
			if t.MultipleTasksDetected {
				multipleDetected = true
			}
		}
		if tasksCount > 1 {
			multipleDetected = true
		}

		errM := r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "detect",
			Provider:   llmName,
			OK:         true,
			DurationMS: time.Since(start).Milliseconds(),
			ChatID:     &chatID,
			UserIDAnon: userID,
			Details: map[string]any{
				"tasks_count":       tasksCount,
				"verbatim_mode":     dres.VerbatimMode,
				"operators_strict":  dres.OperatorsStrict,
				"whitespace_policy": dres.WhitespacePolicy,
				"page_number":       dres.PageMeta.PageNumber,
				"multiple_tasks":    multipleDetected,
				"has_faces_any":     hasFacesAny,
				"pii_detected_any":  piiAny,
			},
		})
		if errM != nil {
			util.PrintError("runDetectThenParse", llmName, chatID, "error insert metrics", errM)
		}
	} else {
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

	// Базовая политика по результату
	if len(dres.Tasks) == 0 {
		r.send(chatID, "ℹ️ Похоже, на фото не распознано учебное задание. Пришлите фото условия задачи (1–4 класс).")
		return
	}
	// предупредим о лицах/PII, если встречаются в любой задаче
	hasFacesAny := false
	piiAny := false
	for _, t := range dres.Tasks {
		if t.HasFaces {
			hasFacesAny = true
		}
		if t.PIIDetected {
			piiAny = true
		}
	}
	if hasFacesAny {
		r.send(chatID, "ℹ️ На фото видны лица. Лучше переснять без лиц.")
	}
	if piiAny {
		r.send(chatID, "ℹ️ На фото обнаружены личные данные. Пожалуйста, замажьте их или переснимите без них.")
	}

	// Несколько заданий — попросить выбрать
	if len(dres.Tasks) > 1 {
		tasks := make([]string, 0, len(dres.Tasks))
		for _, t := range dres.Tasks {
			// попытка краткого описания: номер + первая строка из первого блока либо TitleRaw
			brief := strings.TrimSpace(t.TitleRaw)
			if brief == "" && len(t.Blocks) > 0 {
				// берём первую строку из block_raw
				br := t.Blocks[0].BlockRaw
				br = strings.SplitN(br, "\n", 2)[0]
				brief = strings.TrimSpace(br)
			}
			title := strings.TrimSpace(t.OriginalNumber)
			if title != "" && brief != "" {
				tasks = append(tasks, title+" — "+brief)
			} else if brief != "" {
				tasks = append(tasks, brief)
			} else if title != "" {
				tasks = append(tasks, title)
			} else {
				tasks = append(tasks, "Задание")
			}
		}
		pendingChoice.Store(chatID, tasks)
		pendingCtx.Store(chatID, &selectionContext{Image: merged, Mime: mime, MediaGroupID: mediaGroupID, Detect: dres})

		var b strings.Builder
		b.WriteString("Нашёл несколько заданий. Выберите номер:\n")
		for i, t := range tasks {
			fmt.Fprintf(&b, "%d) %s\n", i+1, t)
		}
		b.WriteString("\nЕсли номер не виден на фото — укажите позицию из списка.")
		r.send(chatID, b.String())
		return
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
		SubjectHint:       "",
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
		OK:         true,
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
