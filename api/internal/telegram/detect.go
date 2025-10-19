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

// TaskChoice хранит отображаемый номер варианта и краткое описание + индекс задачи
type TaskChoice struct {
	Number      string
	Description string
	TaskIndex   int
}

// briefTaskDesc возвращает короткое описание задания для списка выбора
func briefTaskDesc(t types.DetectTask) string {
	brief := strings.TrimSpace(t.TitleRaw)
	if brief == "" && len(t.Blocks) > 0 {
		br := t.Blocks[0].BlockRaw
		br = strings.SplitN(br, "\n", 2)[0]
		brief = strings.TrimSpace(br)
	}
	return brief
}

// fillPendingChoice подготавливает и сохраняет pendingChoice/pendingCtx,
// возвращая строки для пользовательского списка.
func (r *Router) fillPendingChoice(chatID int64, merged []byte, mime, mediaGroupID string, dres types.DetectResult) []string {
	// Соберём занятые номера из оригинальных номеров задач
	used := make(map[string]struct{})
	for _, t := range dres.Tasks {
		if n := strings.TrimSpace(t.OriginalNumber); n != "" {
			used[n] = struct{}{}
		}
	}
	// Генератор свободных числовых номеров, не пересекающихся с оригинальными
	next := 1
	genFree := func() string {
		for {
			n := fmt.Sprintf("%d", next)
			next++
			if _, exists := used[n]; !exists {
				used[n] = struct{}{}
				return n
			}
		}
	}

	var choices []TaskChoice
	var lines []string

	for idx, t := range dres.Tasks {
		brief := briefTaskDesc(t)
		number := strings.TrimSpace(t.OriginalNumber)
		if number == "" {
			number = genFree()
		}
		desc := brief
		if desc == "" && number != "" {
			desc = number
		}
		if desc == "" {
			desc = "Задание"
		}
		choices = append(choices, TaskChoice{
			Number:      number,
			Description: desc,
			TaskIndex:   idx,
		})
		lines = append(lines, fmt.Sprintf("%s — %s", number, desc))
	}

	// Сохраняем в pendingChoice список объектов {номер, описание, индекс}
	pendingChoice.Store(chatID, choices)
	pendingCtx.Store(chatID, &selectionContext{
		Image:        merged,
		Mime:         mime,
		MediaGroupID: mediaGroupID,
		Detect:       dres,
	})
	return lines
}

func (r *Router) hasPendingCorrection(chatID int64) bool { _, ok := parseWait.Load(chatID); return ok }
func (r *Router) clearPendingCorrection(chatID int64)    { parseWait.Delete(chatID) }

func (r *Router) runDetectThenParse(ctx context.Context, chatID int64, userID *int64, merged []byte, mediaGroupID string) {
	setState(chatID, Detect)
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
		setState(chatID, NotATask)
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
		setState(chatID, Inappropriate)
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, "ℹ️ На фото видны лица. Лучше переснять без лиц.", b)
		return
	}
	if piiAny {
		setState(chatID, Inappropriate)
		b := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сообщить об ошибке", "report"))
		r.send(chatID, "ℹ️ На фото обнаружены личные данные. Пожалуйста, замажьте их или переснимите без них.", b)
		return
	}

	// Несколько заданий — попросить выбрать
	if len(dres.Tasks) > 1 {
		setState(chatID, AskChoice)
		lines := r.fillPendingChoice(chatID, merged, mime, mediaGroupID, dres)

		var b strings.Builder
		b.WriteString("Нашёл несколько заданий. Ответьте номером из списка:\n")
		for _, line := range lines {
			fmt.Fprintf(&b, "%s\n", line)
		}
		r.send(chatID, b.String(), nil)
		return
	}

	// без выбора — сразу PARSE
	setState(chatID, DecideTasks)
	r.send(chatID, "Изображение распознано, перехожу к парсингу.", nil)
	sc := &selectionContext{Image: merged, Mime: mime, MediaGroupID: mediaGroupID, Detect: dres}
	r.runParseAndMaybeConfirm(ctx, chatID, userID, sc, -1, "")
	util.PrintInfo("runDetectThenParse", llmName, chatID, fmt.Sprintf("Total time: %d", time.Since(start).Milliseconds()))
}
