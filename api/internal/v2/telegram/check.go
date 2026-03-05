package telegram

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/store"
	"child-bot/api/internal/util"
	"child-bot/api/internal/v2/types"
)

// checkSolution — если есть ожидаемое решение для текущей задачи, проверяем ответ
func (r *Router) checkSolution(ctx context.Context, chatID int64, userID *int64, msg tgbotapi.Message) {
	llmName := r.LlmManager.Get(util.GetChatIDFromTgMessage(msg))

	if len(msg.Photo) == 0 {
		util.PrintInfo("Check", llmName, chatID, "not found photo")
		return
	}

	ph := msg.Photo[len(msg.Photo)-1] // последнее
	r.checkSolutionWithFileID(ctx, chatID, userID, ph.FileID, llmName)
}

// checkSolutionFromDocument — проверка решения из документа-изображения
func (r *Router) checkSolutionFromDocument(ctx context.Context, chatID int64, userID *int64, msg tgbotapi.Message) {
	llmName := r.LlmManager.Get(util.GetChatIDFromTgMessage(msg))

	if msg.Document == nil {
		util.PrintInfo("Check", llmName, chatID, "not found document")
		return
	}

	r.checkSolutionWithFileID(ctx, chatID, userID, msg.Document.FileID, llmName)
}

// checkSolutionWithFileID — общая логика проверки решения по fileID
func (r *Router) checkSolutionWithFileID(ctx context.Context, chatID int64, userID *int64, fileID string, llmName string) {
	r.setStateWithPersist(chatID, Check)
	sid, _ := r.getSession(chatID)

	stopProgress := r.startCheckProgress(chatID)
	defer stopProgress()

	data, mime, err := r.downloadFileBytes(fileID)
	if err != nil {
		r.sendError(chatID, fmt.Errorf("не удалось получить фото: %v", err))
		return
	}
	if mime == "application/octet-stream" {
		// Попробуем руками распознать распространённые форматы и HEIC/AVIF
		if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xD8 {
			mime = "image/jpeg"
		}
		if len(data) >= 8 &&
			data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
			data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
			mime = "image/png"
		}
		if heicAvif := util.SniffHEICorAVIF(data); heicAvif != "" {
			mime = heicAvif
		}
	}

	// 0) Подтянем метаданные предмета/класса из последнего подтверждённого парсинга
	subj := "generic"
	grade := int64(0)
	rawTaskText := ""
	var taskStruct types.TaskStructCheck
	taskStructLoaded := false

	if pr, ok := r.Store.FindLastConfirmedParse(ctx, sid); ok {
		if s := strings.TrimSpace(pr.Subject); s != "" {
			subj = s
			grade = pr.Grade
			rawTaskText = pr.RawTaskText

			// Восстанавливаем ParseResponse из сохранённого JSON
			var parseResp types.ParseResponse
			if err := json.Unmarshal(pr.ResultJSON, &parseResp); err == nil {
				taskStruct = types.TaskStructCheck{
					TaskTextClean:   parseResp.Task.TaskTextClean,
					VisualReasoning: parseResp.Task.VisualReasoning,
					VisualFacts:     parseResp.Task.VisualFacts,
					QualityFlags:    parseResp.Task.Quality,
					Items:           parseResp.Items,
				}
				taskStructLoaded = true
			} else {
				util.PrintError("checkSolutionWithFileID", llmName, chatID,
					"failed to unmarshal ParseResponse from DB, check will proceed without task context", err)
			}
		}
	}

	// Логируем предупреждение если проверка идёт без контекста задачи
	if !taskStructLoaded {
		util.PrintInfo("checkSolutionWithFileID", llmName, chatID,
			"proceeding without task structure context - check accuracy may be reduced")
	}

	in := types.CheckRequest{
		Image:            base64.StdEncoding.EncodeToString(data),
		TaskStruct:       taskStruct,
		RawTaskText:      rawTaskText,
		PhotoQualityHint: "auto",
		Student: types.StudentCheck{
			Grade:   grade,
			Subject: subj,
			Locale:  "ru_RU",
		},
	}

	start := time.Now()
	res, err := r.GetLLMClient().CheckSolution(ctx, llmName, in)
	latency := time.Since(start).Milliseconds()
	_ = r.Store.InsertHistory(ctx, store.TimelineEvent{
		ChatID:        chatID,
		TaskSessionID: sid,
		Direction:     "api",
		EventType:     string(Check),
		Provider:      llmName,
		OK:            err == nil,
		LatencyMS:     &latency,
		InputPayload:  in,
		OutputPayload: res,
		Error:         err,
	})
	if err != nil {
		_ = r.Store.InsertEvent(ctx, store.MetricEvent{
			Stage:      "check",
			Provider:   llmName,
			OK:         false,
			Error:      err.Error(),
			DurationMS: latency,
			ChatID:     &chatID,
			UserIDAnon: userID,
			Details: map[string]any{
				"subject": subj,
			},
		})

		r.sendError(chatID, fmt.Errorf("не удалось проверить решение: %v", err))
		return
	}

	r.sendDebug(chatID, "check input", in)
	r.sendDebug(chatID, "check res", res)

	_ = r.Store.InsertEvent(ctx, store.MetricEvent{
		Stage:      "check",
		Provider:   llmName,
		OK:         true,
		DurationMS: latency,
		ChatID:     &chatID,
		UserIDAnon: userID,
		Details: map[string]any{
			"subject":    subj,
			"confidence": res.Confidence,
			"is_correct": res.IsCorrect,
		},
	})

	r.sendCheckResponse(chatID, res)
}

// sendCheckResponse — вывод краткого результата проверки (P0.3: использует Decision enum)
func (r *Router) sendCheckResponse(chatID int64, cr types.CheckResponse) {
	// P0.3: Используем Decision вместо IsCorrect
	switch cr.Decision {
	case types.CheckDecisionCorrect:
		r.setStateWithPersist(chatID, Correct)
		r.clearModeWithPersist(chatID)
		r.clearSession(chatID)
		r.send(chatID, AnswerCorrectText, makeCorrectAnswerButtons())
		return

	case types.CheckDecisionIncorrect:
		r.setStateWithPersist(chatID, Incorrect)
		text := fmt.Sprintf(AnswerIncorrectText, cr.Feedback)
		r.send(chatID, text, makeIncorrectAnswerButtons())
		return

	case types.CheckDecisionNeedAnnotation:
		// Нужна аннотация рисунка
		r.setStateWithPersist(chatID, AwaitSolution)
		r.send(chatID, "Подпиши вершины или обведи искомую область на рисунке, чтобы я мог точно проверить.", nil)
		return

	case types.CheckDecisionInvalidExpected:
		// P0.1: Противоречие в эталоне - внутренняя ошибка, не показываем техническое сообщение
		r.setStateWithPersist(chatID, AwaitSolution)
		r.send(chatID, "Не удалось проверить ответ. Попробуй отправить фото ещё раз.", nil)
		return

	case types.CheckDecisionCannotEvaluate:
		// Не удалось проверить — показываем feedback если есть (не технический)
		r.setStateWithPersist(chatID, AwaitSolution)
		if cr.Feedback != "" && !strings.Contains(cr.Feedback, "эталон") {
			r.send(chatID, cr.Feedback, nil)
		} else {
			r.send(chatID, "Не удалось проверить ответ. Попробуй переснять фото.", nil)
		}
		return

	default:
		// Fallback для неизвестного Decision (обратная совместимость)
		if cr.IsCorrect != nil && *cr.IsCorrect {
			r.setStateWithPersist(chatID, Correct)
			r.clearModeWithPersist(chatID)
			r.clearSession(chatID)
			r.send(chatID, AnswerCorrectText, makeCorrectAnswerButtons())
		} else {
			r.setStateWithPersist(chatID, Incorrect)
			text := fmt.Sprintf(AnswerIncorrectText, cr.Feedback)
			r.send(chatID, text, makeIncorrectAnswerButtons())
		}
	}
}

// downloadFileBytes — скачивает файл Telegram по fileID и возвращает bytes и mime
// Использует BotSender.DownloadFile для совместимости с моками в тестах
func (r *Router) downloadFileBytes(fileID string) ([]byte, string, error) {
	b, err := r.Bot.DownloadFile(fileID)
	if err != nil {
		return nil, "", err
	}

	// Определяем MIME-тип по сигнатуре файла
	mime := "image/jpeg" // default
	if len(b) >= 8 {
		if b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4E && b[3] == 0x47 {
			mime = "image/png"
		} else if len(b) >= 12 {
			// Check for HEIC/AVIF (ftyp box)
			if string(b[4:8]) == "ftyp" {
				brand := string(b[8:12])
				if brand == "heic" || brand == "heix" || brand == "mif1" {
					mime = "image/heic"
				} else if brand == "avif" {
					mime = "image/avif"
				}
			}
		}
	}

	return b, mime, nil
}
