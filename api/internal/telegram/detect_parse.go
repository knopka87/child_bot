package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"child-bot/api/internal/ocr"
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

func (r *Router) runDetectThenParse(ctx context.Context, chatID int64, merged []byte, mediaGroupID string) {
	mime := util.SniffMimeHTTP(merged)
	llmName := r.EngManager.Get(chatID)

	// DETECT через llmproxy
	var dres ocr.DetectResult
	if dr, err := r.LLM.Detect(ctx, llmName, merged, mime, 0); err == nil {
		dres = dr
	} else {
		// Мягкий фолбэк: продолжаем без детекта (используем значения по умолчанию),
		// просто логируем ошибку и сообщаем пользователю, что попробуем распознать весь снимок.
		log.Printf("detect failed (chat=%d): %v; fallback to parse without detect", chatID, err)
		r.send(chatID, "ℹ️ Не удалось выделить области на фото, попробую распознать задание целиком.")
	}

	// базовая политика
	if dres.FinalState == "inappropriate_image" {
		r.send(chatID, "⚠️ Неподходящее изображение. Пришлите фото учебного задания без личных данных.")
		return
	}
	if dres.NeedsRescan {
		msg := "Пожалуйста, переснимите фото"
		if dres.RescanReason != "" {
			msg += ": " + dres.RescanReason
		}
		r.send(chatID, "📷 "+msg)
		return
	}
	if dres.HasFaces {
		r.send(chatID, "ℹ️ На фото видны лица. Лучше переснять без лиц.")
	}

	// несколько заданий — спросить номер
	if dres.MultipleTasksDetected && !(dres.AutoChoiceSuggested && dres.TopCandidateIndex != nil &&
		*dres.TopCandidateIndex >= 0 && *dres.TopCandidateIndex < len(dres.TasksBrief) &&
		dres.Confidence >= 0.80) {
		if len(dres.TasksBrief) > 0 {
			pendingChoice.Store(chatID, dres.TasksBrief)
			pendingCtx.Store(chatID, &selectionContext{Image: merged, Mime: mime, MediaGroupID: mediaGroupID, Detect: dres})
			var b strings.Builder
			b.WriteString("Нашёл несколько заданий. Выберите номер:\n")
			for i, t := range dres.TasksBrief {
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
	sc := &selectionContext{Image: merged, Mime: mime, MediaGroupID: mediaGroupID, Detect: dres}
	r.runParseAndMaybeConfirm(ctx, chatID, sc, -1, "")
}

func (r *Router) runParseAndMaybeConfirm(ctx context.Context, chatID int64, sc *selectionContext, selectedIdx int, selectedBrief string) {
	imgHash := util.SHA256Hex(sc.Image)
	llmName := r.EngManager.Get(chatID)

	// 1) кэш из БД: принят ли PARSE
	if prRow, err := r.ParseRepo.FindByHash(ctx, imgHash, llmName, 30*24*time.Hour); err == nil && prRow.Accepted {
		r.showTaskAndPrepareHints(chatID, sc, prRow.Parse, llmName)
		return
	}

	// 2) LLM.Parse
	pr, err := r.LLM.Parse(ctx, llmName, sc.Image, ocr.ParseOptions{
		SubjectHint:       sc.Detect.SubjectGuess,
		ChatID:            chatID,
		MediaGroupID:      sc.MediaGroupID,
		ImageHash:         imgHash,
		SelectedTaskIndex: selectedIdx,
		SelectedTaskBrief: selectedBrief,
	})
	if err != nil {
		r.SendError(chatID, fmt.Errorf("parse: %w", err))
		return
	}

	// сохранить черновик
	_ = r.ParseRepo.Upsert(ctx, chatID, sc.MediaGroupID, imgHash, llmName, pr, false, "")

	// 3) подтверждение, если нужно
	if pr.ConfirmationNeeded {
		r.askParseConfirmation(chatID, pr)
		parseWait.Store(chatID, &parsePending{Sc: sc, PR: pr, LLM: llmName})
		return
	}

	// 4) автоподтверждение
	_ = r.ParseRepo.MarkAccepted(ctx, imgHash, llmName, "auto")
	r.showTaskAndPrepareHints(chatID, sc, pr, llmName)
}
