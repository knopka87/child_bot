package telegram

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/ocr/types"
	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
)

func (r *Router) hasPendingCorrection(chatID int64) bool { _, ok := parseWait.Load(chatID); return ok }
func (r *Router) clearPendingCorrection(chatID int64)    { parseWait.Delete(chatID) }

func (r *Router) runDetectThenParse(ctx context.Context, chatID int64, userID *int64, merged []byte, mediaGroupID string) {
	mime := util.SniffMimeHTTP(merged)
	llmName := r.EngManager.Get(chatID)

	r.sendDebug(chatID, "mime", mime)
	// DETECT через llmproxy
	var dres types.DetectResult
	in := types.DetectInput{
		ImageB64:  base64.StdEncoding.EncodeToString(merged),
		Mime:      mime,
		GradeHint: 0,
	}
	start := time.Now()
	dr, err := r.LLM.Detect(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	if err == nil {
		dres = dr
		r.sendDebug(chatID, "detect_req", in)
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
			DurationMS: latency,
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
				DurationMS: latency,
				ChatID:     &chatID,
				UserIDAnon: userID,
				Error:      err.Error(),
			})
		}
		log.Printf("detect failed (chat=%d): %v; fallback to parse without detect", chatID, err)
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, "ℹ️ Не удалось выделить области на фото, попробую распознать задание целиком.", b)
	}
	util.PrintInfo("runDetectThenParse", llmName, chatID, fmt.Sprintf("Received a response from LLM: %d", time.Since(start).Milliseconds()))

	sid, _ := r.getSession(chatID)
	_ = r.History.Insert(ctx, store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Provider:      llmName,
		Direction:     "api",
		EventType:     string(Detect),
		InputPayload:  in,
		OutputPayload: dres,
		Error:         err,
		OK:            err == nil,
		LatencyMS:     &latency,
	})

	// Базовая политика по результату
	if len(dres.Tasks) == 0 {
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, "ℹ️ Похоже, на фото не распознано учебное задание. Пришлите фото условия задачи (1–4 класс).", b)
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
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, "ℹ️ На фото видны лица. Лучше переснять без лиц.", b)
	}
	if piiAny {
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, "ℹ️ На фото обнаружены личные данные. Пожалуйста, замажьте их или переснимите без них.", b)
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
		r.send(chatID, b.String(), nil)
		return
	}

	// без выбора — сразу PARSE
	r.send(chatID, "Изображение распознано, перехожу к парсингу.", nil)
	sc := &selectionContext{Image: merged, Mime: mime, MediaGroupID: mediaGroupID, Detect: dres}
	r.runParseAndMaybeConfirm(ctx, chatID, userID, sc, -1, "")
	util.PrintInfo("runDetectThenParse", llmName, chatID, fmt.Sprintf("Total time: %d", time.Since(start).Milliseconds()))
}
