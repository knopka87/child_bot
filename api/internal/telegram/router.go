package telegram

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/llmclient"
	"child-bot/api/internal/ocr"
	"child-bot/api/internal/store"
)

type Router struct {
	Bot        *tgbotapi.BotAPI
	EngManager *ocr.Manager
	ParseRepo  *store.ParseRepo
	HintRepo   *store.HintRepo
	LLM        *llmclient.Client
	Metrics    *store.MetricsRepo
}

func (r *Router) HandleCommand(upd tgbotapi.Update) {
	cid := upd.Message.Chat.ID
	switch upd.Message.Command() {
	case "start":
		r.send(cid, "Пришли фото задачи — верну распознанный текст и подскажу, с чего начать.\nКоманды: /health")
	case "health":
		r.send(cid, "✅ OK")
	case "engine":
		args := strings.Fields(strings.TrimSpace(strings.TrimPrefix(upd.Message.Text, "/engine")))
		cur := r.EngManager.Get(cid)
		if len(args) == 0 {
			r.send(cid, "Текущий LLM-провайдер: "+cur+
				"\nИспользование:\n/engine gemini\n/engine gpt")
			return
		}
		// применим через общий обработчик ниже
		r.handleEngineCommand(cid, upd.Message.Text)
		return
	default:
		r.send(cid, "Неизвестная команда")
	}
}

func (r *Router) HandleUpdate(upd tgbotapi.Update) {
	// 1) Callback-кнопки
	if upd.CallbackQuery != nil {
		r.handleCallback(*upd.CallbackQuery)
		return
	}

	// 2) Сообщений нет — выходим
	if upd.Message == nil {
		return
	}
	cid := upd.Message.Chat.ID

	// 3) Если ждём текстовую правку после «Нет» — приоритетно принимаем её
	if r.hasPendingCorrection(cid) && upd.Message.Text != "" {
		r.applyTextCorrectionThenShowHints(cid, upd.Message.Text)
		return
	}

	// 4) «Жёсткий» режим ввода: если ждём решение — принимаем и текст, и фото;
	//    если ждём новую задачу — просим фото задачи; в остальных случаях — как раньше.
	if upd.Message.Text != "" && !upd.Message.IsCommand() {
		switch getMode(cid) {
		case "await_solution":
			// Нормализуем текстовый ответ ученика
			r.normalizeText(context.Background(), cid, upd.Message.Contact.UserID, upd.Message.Text)
			return
		case "await_new_task":
			r.send(cid, "Я жду фото новой задачи. Пожалуйста, пришлите фото.")
			return
		}
	}

	// 5) Ветвь выбора пункта при multiple tasks (ожидаем число 1..N)
	if v, ok := pendingChoice.Load(cid); ok && upd.Message.Text != "" {
		briefs := v.([]string)
		if n, err := strconv.Atoi(strings.TrimSpace(upd.Message.Text)); err == nil && n >= 1 && n <= len(briefs) {
			if ctxv, ok2 := pendingCtx.Load(cid); ok2 {
				pendingChoice.Delete(cid)
				pendingCtx.Delete(cid)
				sc := ctxv.(*selectionContext)
				r.send(cid, fmt.Sprintf("Ок, беру задание: %s — обрабатываю.", briefs[n-1]))
				r.runParseAndMaybeConfirm(context.Background(), cid, sc, n-1, briefs[n-1])
				return
			}
			pendingChoice.Delete(cid)
			r.send(cid, "Не нашёл предыдущее изображение. Пришлите фото ещё раз.")
			return
		}
		// иначе ждём корректный номер
	}

	// 6) Команды (в т.ч. /engine)
	if upd.Message.IsCommand() && strings.HasPrefix(upd.Message.Text, "/engine") {
		r.handleEngineCommand(cid, upd.Message.Text)
		return
	}
	if upd.Message.IsCommand() {
		r.HandleCommand(upd)
		return
	}

	// 7) Фото/альбом
	if len(upd.Message.Photo) > 0 {
		if getMode(cid) == "await_solution" {
			// Фото с ответом ученика → нормализация
			r.normalizePhoto(context.Background(), *upd.Message)
			clearMode(cid)
			return
		}
		// Иначе — это фото задачи/страницы
		clearMode(cid)
		r.acceptPhoto(*upd.Message)
		return
	}

	// 8) Остальное — игнорируем
}

func (r *Router) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, _ = r.Bot.Send(msg)
}

func (r *Router) SendResult(chatID int64, text string) {
	if len(text) > 3900 {
		text = text[:3900] + "…"
	}
	r.send(chatID, "📝 Распознанный текст:\n\n"+text)
}

func (r *Router) SendError(chatID int64, err error) {
	r.send(chatID, fmt.Sprintf("Ошибка OCR: %v", err))
}

// handleEngineCommand парсит команду /engine и переключает провайдера LLM для чата.
// Поддерживаются только gemini и gpt.
func (r *Router) handleEngineCommand(chatID int64, cmd string) {
	args := strings.Fields(strings.TrimSpace(strings.TrimPrefix(cmd, "/engine")))
	if len(args) == 0 {
		r.send(chatID, "Использование: /engine {gemini|gpt}")
		return
	}
	name := strings.ToLower(args[0])
	switch name {
	case "gemini", "google":
		r.EngManager.Set(chatID, "gemini")
		r.send(chatID, "✅ Провайдер LLM: gemini")
	case "gpt", "openai":
		r.EngManager.Set(chatID, "gpt")
		r.send(chatID, "✅ Провайдер LLM: gpt")
	default:
		r.send(chatID, "Неизвестный провайдер. Доступны: gemini | gpt")
	}
}

// Показ запроса подтверждения распознанного текста
func (r *Router) askParseConfirmation(chatID int64, pr ocr.ParseResult) {
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

// PhotoAcceptedText — первый ответ после получения фото/первой страницы альбома.
func (r *Router) PhotoAcceptedText() string {
	return "Фото принято. Если задание на нескольких фото — просто пришлите их подряд, я склею страницы перед обработкой."
}

// normalizeText — отправляет текст ученика на нормализацию в LLM-прокси
func (r *Router) normalizeText(ctx context.Context, chatID, userID int64, text string) {
	llmName := r.EngManager.Get(chatID)
	shape := r.suggestSolutionShape(chatID)
	in := ocr.NormalizeInput{
		UserIDAnon:    fmt.Sprint(userID),
		SolutionShape: shape,
		Provider:      llmName,
		Answer:        ocr.NormalizeAnswer{Source: "text", Text: strings.TrimSpace(text)},
	}
	start := time.Now()
	res, err := r.LLM.Normalize(ctx, llmName, in)
	if err != nil {
		if r.Metrics != nil {
			_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
				Stage:      "normalize",
				Provider:   llmName,
				OK:         false,
				Error:      err.Error(),
				DurationMS: time.Since(start).Milliseconds(),
				ChatID:     &chatID,
				UserIDAnon: &userID,
				Details: map[string]any{
					"source":      "text",
					"input_chars": len(text),
				},
			})
		}
		r.send(chatID, fmt.Sprintf("Не удалось нормализовать ответ: %v", err))
		return
	}
	r.sendNormalizePreview(chatID, res)
	if r.Metrics != nil {
		_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "normalize",
			Provider:   llmName,
			OK:         true,
			DurationMS: time.Since(start).Milliseconds(),
			ChatID:     &chatID,
			UserIDAnon: &userID,
			Details: map[string]any{
				"source":          "text",
				"shape":           res.Shape,
				"needs_clarify":   res.NeedsClarification,
				"uncertain_count": len(res.UncertainReasons),
			},
		})
	}
	// Попробуем сразу проверить решение, если в системе есть ожидаемое решение
	r.maybeCheckSolution(ctx, chatID, userID, res)
	clearMode(chatID)
}

// normalizePhoto — скачивает фото из Telegram и отправляет на нормализацию
func (r *Router) normalizePhoto(ctx context.Context, msg tgbotapi.Message) {
	if len(msg.Photo) == 0 {
		return
	}
	llmName := r.EngManager.Get(msg.Chat.ID)
	ph := msg.Photo[len(msg.Photo)-1] // самое большое
	data, mime, err := r.downloadFileBytes(ph.FileID)
	if err != nil {
		r.send(msg.Chat.ID, fmt.Sprintf("Не удалось получить фото: %v", err))
		return
	}
	shape := r.suggestSolutionShape(msg.Chat.ID)
	in := ocr.NormalizeInput{
		SolutionShape: shape,
		Provider:      llmName,
		Answer: ocr.NormalizeAnswer{
			Source:   "photo",
			PhotoB64: base64.StdEncoding.EncodeToString(data),
			Mime:     mime,
		},
	}
	start := time.Now()
	res, err := r.LLM.Normalize(ctx, llmName, in)
	if err != nil {
		if r.Metrics != nil {
			_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
				Stage:      "normalize",
				Provider:   llmName,
				OK:         false,
				Error:      err.Error(),
				DurationMS: time.Since(start).Milliseconds(),
				ChatID:     &msg.Chat.ID,
				UserIDAnon: &msg.Contact.UserID,
				Details: map[string]any{
					"source": "photo",
					"mime":   mime,
					"bytes":  len(data),
				},
			})
		}
		r.send(msg.Chat.ID, fmt.Sprintf("Не удалось нормализовать ответ (фото): %v", err))
		return
	}
	r.sendNormalizePreview(msg.Chat.ID, res)
	if r.Metrics != nil {
		_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "normalize",
			Provider:   llmName,
			OK:         true,
			DurationMS: time.Since(start).Milliseconds(),
			ChatID:     &msg.Chat.ID,
			UserIDAnon: &msg.Contact.UserID,
			Details: map[string]any{
				"source":          "photo",
				"mime":            mime,
				"bytes":           len(data),
				"shape":           res.Shape,
				"needs_clarify":   res.NeedsClarification,
				"uncertain_count": len(res.UncertainReasons),
			},
		})
	}
	// Попробуем сразу проверить решение, если в системе есть ожидаемое решение
	r.maybeCheckSolution(ctx, msg.Chat.ID, msg.Contact.UserID, res)
}

// suggestSolutionShape — простая эвристика: если по парсингу известна форма — берём её, иначе number
func (r *Router) suggestSolutionShape(chatID int64) string {
	// TODO: можно взять из последнего ParseResult из БД (ParseRepo) subject/task_type → shape
	return "number"
}

// sendNormalizePreview — короткий текст для пользователя по NormalizeResult
func (r *Router) sendNormalizePreview(chatID int64, nr ocr.NormalizeResult) {
	shape := strings.ToLower(strings.TrimSpace(nr.Shape))
	val := ""
	switch v := nr.Value.(type) {
	case string:
		val = v
	case float64:
		val = strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		val = strconv.Itoa(v)
	case []string:
		val = strings.Join(v, "; ")
	default:
		val = "(не удалось отобразить значение)"
	}
	b := &strings.Builder{}
	b.WriteString("✅ Принял ответ.")
	if shape != "" {
		b.WriteString("\nФорма: ")
		b.WriteString(shape)
	}
	if val != "" {
		b.WriteString("\nЗначение: ")
		b.WriteString(val)
	}
	if nr.UncertainReasons != nil && len(nr.UncertainReasons) > 0 {
		b.WriteString("\nПредупреждения: ")
		b.WriteString(strings.Join(nr.UncertainReasons, ", "))
	}
	if nr.NeedsClarification && nr.NeedsUserActionMessage != "" {
		b.WriteString("\nНужно уточнение: ")
		b.WriteString(nr.NeedsUserActionMessage)
	}
	r.send(chatID, b.String())
}

// maybeCheckSolution — если есть ожидаемое решение для текущей задачи, проверяем ответ
func (r *Router) maybeCheckSolution(ctx context.Context, chatID, userID int64, nr ocr.NormalizeResult) {
	// 0) Подтянем метаданные предмета/класса из последнего подтверждённого парсинга
	subj := "math"
	grade := 0
	if r.ParseRepo != nil {
		if pr, ok := r.ParseRepo.FindLastConfirmed(ctx, chatID); ok {
			subj = strings.TrimSpace(pr.Subject)
			grade = pr.Grade
		}
	}

	// 1) Пытаемся взять ожидаемое решение из БД
	exp, ok := r.getExpectedForChat(ctx, chatID)
	if !ok {
		// 2) Фолбэк: строим «policy-only» ожидание по данным нормализации ученика
		shape := strings.TrimSpace(nr.Shape)
		if shape == "" {
			shape = strings.TrimSpace(nr.ShapeDetected)
		}
		if shape == "" {
			shape = "number"
		}

		var units *ocr.UnitsExpectedSpec
		if nr.Units != nil {
			policy := "optional"
			if nr.Units.Kept {
				policy = "required"
			}
			primary := ""
			if nr.Units.Canonical != nil {
				primary = strings.TrimSpace(*nr.Units.Canonical)
			}
			alts := []string{}
			if nr.Units.Detected != nil {
				det := strings.TrimSpace(*nr.Units.Detected)
				if det != "" && det != primary {
					alts = append(alts, det)
				}
			}
			units = &ocr.UnitsExpectedSpec{
				Policy:          policy,  // требуем/не требуем единицы
				ExpectedPrimary: primary, // если нормализация вывела канон. единицу
				Alternatives:    alts,    // допустимые альтернативы
			}
		}

		exp = ocr.ExpectedSolution{
			Shape: shape,
			Units: units,
			// Number/String/List/Steps — не задаём без эталона, чтобы не «подгонять» под ответ
		}
	}

	llmName := r.EngManager.Get(chatID)
	in := ocr.CheckSolutionInput{
		UserIDAnon: fmt.Sprint(chatID),
		Subject:    subj,
		Grade:      grade,
		Student:    nr,
		Expected:   exp,
	}
	start := time.Now()
	res, err := r.LLM.CheckSolution(ctx, llmName, in)
	if err != nil {
		if r.Metrics != nil {
			_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
				Stage:      "check",
				Provider:   llmName,
				OK:         false,
				Error:      err.Error(),
				DurationMS: time.Since(start).Milliseconds(),
				ChatID:     &chatID,
				UserIDAnon: &userID,
				Details: map[string]any{
					"subject": subj,
					"grade":   grade,
				},
			})
		}
		r.send(chatID, fmt.Sprintf("Не удалось проверить решение: %v", err))
		r.offerAnalogueButton(chatID)
		return
	}
	r.sendCheckResult(chatID, res)
	if r.Metrics != nil {
		_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "check",
			Provider:   llmName,
			OK:         true,
			DurationMS: time.Since(start).Milliseconds(),
			ChatID:     &chatID,
			UserIDAnon: &userID,
			Details: map[string]any{
				"subject": subj,
				"grade":   grade,
				"verdict": res.Verdict,
			},
		})
	}
}

// getExpectedForChat — извлекает ожидаемое решение из вашей БД для текущей задачи чата
func (r *Router) getExpectedForChat(ctx context.Context, chatID int64) (ocr.ExpectedSolution, bool) {
	// if r.ParseRepo != nil {
	// 	if pr, ok := r.ParseRepo.FindLastConfirmed(ctx, chatID); ok {
	// 		return pr.Expected, true
	// 	}
	// }
	var exp ocr.ExpectedSolution
	return exp, false
}

// sendCheckResult — вывод краткого результата проверки
func (r *Router) sendCheckResult(chatID int64, cr ocr.CheckSolutionResult) {
	var b strings.Builder
	switch cr.Verdict {
	case "correct":
		b.WriteString("✅ Задача решена верно\n")
	case "incorrect":
		b.WriteString("⚠️ Похоже, есть ошибка\n")
	case "uncertain":
		b.WriteString("🤔 Я не уверен в оценке\n")
	default:
		b.WriteString("Результат проверки получен\n")
	}
	if s := strings.TrimSpace(cr.ShortHint); s != "" {
		b.WriteString("Подсказка: ")
		b.WriteString(s)
		b.WriteString("\n")
	}
	// Доп. диагностическая сводка без раскрытия ответа
	if cr.Comparison.Units != nil && cr.Comparison.Units.Policy != "" {
		b.WriteString("Единицы: ")
		if cr.Comparison.Units.Detected == "" {
			b.WriteString("(не указаны)")
		} else {
			b.WriteString(cr.Comparison.Units.Detected)
		}
		if cr.Comparison.Units.Applied != "" {
			b.WriteString("; конверсия: ")
			b.WriteString(cr.Comparison.Units.Applied)
		}
		b.WriteString("\n")
	}
	if nd := cr.Comparison.NumberDiff; nd != nil {
		if nd.WithinTolerance {
			b.WriteString("Число в допуске\n")
		} else if nd.EquivalentByRule {
			b.WriteString("Число эквивалентно по правилу\n")
		}
	}
	if sm := cr.Comparison.StringMatch; sm != nil && sm.Method != "" {
		b.WriteString("Проверка слова: ")
		b.WriteString(sm.Method)
		if sm.Passed {
			b.WriteString(" — ок\n")
		} else {
			b.WriteString(" — есть расхождение\n")
		}
	}
	if lm := cr.Comparison.ListMatch; lm != nil && lm.Total > 0 {
		b.WriteString(fmt.Sprintf("Элементов совпало: %d/%d\n", lm.Matched, lm.Total))
	}
	if st := cr.Comparison.StepsMatch; st != nil && st.Total > 0 {
		b.WriteString(fmt.Sprintf("Шагов покрыто: %d/%d\n", st.Covered, st.Total))
	}
	// Сообщение для озвучивания (не более 140 симв.)
	if s := strings.TrimSpace(cr.SpeakableMessage); s != "" {
		b.WriteString("\n")
		b.WriteString(s)
	}
	r.send(chatID, b.String())
	// Предложить аналогичное задание при ошибке/неуверенности
	if cr.Verdict == "incorrect" || cr.Verdict == "uncertain" {
		r.offerAnalogueButton(chatID)
	}
}

// downloadFileBytes — скачивает файл Telegram по fileID и возвращает bytes и mime
func (r *Router) downloadFileBytes(fileID string) ([]byte, string, error) {
	url, err := r.Bot.GetFileDirectURL(fileID)
	if err != nil {
		return nil, "", err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	mime := resp.Header.Get("Content-Type")
	if mime == "" {
		mime = "image/jpeg"
	}
	return b, mime, nil
}

// --- ANALOGUE SOLUTION (v1.1) ----------------------------------------------
// По кнопке «Похожее задание» генерируем аналог по тем же приёмам, но с другими данными.
// Основано на инструкции ANALOGUE_SOLUTION v1.1.

// offerAnalogueButton — показывает кнопку для вызова аналога
func (r *Router) offerAnalogueButton(chatID int64) {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Похожее задание", "analogue_solution"),
		},
	)
	msg := tgbotapi.NewMessage(chatID, "Если нужно, покажу похожее задание тем же приёмом (без ответа исходной задачи).")
	msg.ReplyMarkup = kb
	_, _ = r.Bot.Send(msg)
}

// HandleAnalogueCallback — публичный помощник для существующего handleCallback
// Вызовите его из вашего обработчика, когда callback.Data == "ANALOGUE".
func (r *Router) HandleAnalogueCallback(chatID, userID int64) {
	ctx := context.Background()
	if err := r.runAnalogue(ctx, chatID, userID); err != nil {
		r.send(chatID, "Не удалось подготовить аналогичное задание: "+err.Error())
	}
}

// runAnalogue — собирает вход из последнего подтверждённого парсинга и вызывает LLM-прокси
func (r *Router) runAnalogue(ctx context.Context, chatID, userID int64) error {
	in, err := r.buildAnalogueInput(ctx, chatID)
	if err != nil {
		return err
	}
	llmName := r.EngManager.Get(chatID)
	start := time.Now()
	ar, err := r.LLM.AnalogueSolution(ctx, llmName, in)
	if err != nil {
		if r.Metrics != nil {
			_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
				Stage:      "analogue",
				Provider:   llmName,
				OK:         false,
				Error:      err.Error(),
				DurationMS: time.Since(start).Milliseconds(),
				ChatID:     &chatID,
				UserIDAnon: &userID,
			})
		}
		return err
	}
	r.sendAnalogueResult(chatID, ar)
	if r.Metrics != nil {
		_ = r.Metrics.InsertEvent(ctx, store.MetricEvent{
			Stage:      "analogue",
			Provider:   llmName,
			OK:         true,
			DurationMS: time.Since(start).Milliseconds(),
			ChatID:     &chatID,
			UserIDAnon: &userID,
			Details: map[string]any{
				"has_minichecks": len(ar.MiniChecks) > 0,
			},
		})
	}
	return nil
}

// buildAnalogueInput — конструирует вход для ANALOGUE из данных последнего парсинга
func (r *Router) buildAnalogueInput(ctx context.Context, chatID int64) (ocr.AnalogueSolutionInput, error) {
	if r.ParseRepo == nil {
		return ocr.AnalogueSolutionInput{}, errors.New("ParseRepo is not configured")
	}
	pr, ok := r.ParseRepo.FindLastConfirmed(ctx, chatID)
	if !ok {
		return ocr.AnalogueSolutionInput{}, errors.New("нет подтверждённого задания — пришлите фото и подтвердите распознавание")
	}

	// Берём краткую суть, либо строим её из вопроса/сырого текста, удаляя числа/единицы
	essence := strings.TrimSpace(pr.ShortEssence)
	if essence == "" {
		base := strings.TrimSpace(pr.Question)
		if base == "" {
			base = strings.TrimSpace(pr.RawText)
		}
		norm := stripNumbersUnits(base)
		if norm == "" {
			return ocr.AnalogueSolutionInput{}, errors.New("не удалось получить краткую суть задания")
		}
		essence = norm
	}

	in := ocr.AnalogueSolutionInput{
		TaskID:              pr.TaskID,
		UserIDAnon:          fmt.Sprint(chatID),
		Grade:               pr.Grade,
		Subject:             pr.Subject,   // "math"|"russian"|...
		TaskType:            pr.TaskType,  // если классификатор есть
		MethodTag:           pr.MethodTag, // ключевой приём (если определён)
		DifficultyHint:      pr.DifficultyHint,
		OriginalTaskEssence: essence, // без чисел/единиц исходника
		Locale:              "ru",
	}
	return in, nil
}

var reNums = regexp.MustCompile(`(?i)(\d+[\d\s./,:-]*\d*|\d+)`)
var reUnits = regexp.MustCompile(`(?i)(см|мм|м|кг|г|л|мл|ч|мин|сек|%|грн|руб|р\.|км)\.?`)

// stripNumbersUnits — удаляет из текста числа и типичные единицы/знаки, чтобы
// получить краткую суть без утечки исходных данных (см. anti‑leak в v1.1)
func stripNumbersUnits(s string) string {
	out := reNums.ReplaceAllString(s, "N")
	out = reUnits.ReplaceAllString(out, "U")
	out = strings.TrimSpace(strings.Join(strings.Fields(out), " "))
	return out
}

// sendAnalogueResult — формирует человекочитаемый вывод без раскрытия ответа исходника
func (r *Router) sendAnalogueResult(chatID int64, ar ocr.AnalogueSolutionResult) {
	var b strings.Builder
	if t := strings.TrimSpace(ar.AnalogyTitle); t != "" {
		b.WriteString("📘 ")
		b.WriteString(t)
		b.WriteString("\n\n")
	}
	if t := strings.TrimSpace(ar.AnalogyTask); t != "" {
		b.WriteString("Похожее задание:\n")
		b.WriteString(t)
		b.WriteString("\n\n")
	}
	if len(ar.SolutionSteps) > 0 {
		b.WriteString("Как решать (тот же приём):\n")
		for i, s := range ar.SolutionSteps {
			b.WriteString(strconv.Itoa(i + 1))
			b.WriteString(". ")
			b.WriteString(strings.TrimSpace(s))
			b.WriteString("\n")
		}
	}
	if len(ar.TransferBridge) > 0 {
		b.WriteString("\nМостик переноса:\n")
		for i, s := range ar.TransferBridge {
			b.WriteString("• ")
			b.WriteString(strings.TrimSpace(s))
			if i < len(ar.TransferBridge)-1 {
				b.WriteString("\n")
			}
		}
	}
	if s := strings.TrimSpace(ar.TransferCheck); s != "" {
		b.WriteString("\n\nПроверь себя: ")
		b.WriteString(s)
	}
	// Мини‑проверки: показываем без ответов
	if len(ar.MiniChecks) > 0 {
		b.WriteString("\n\nМини‑проверки:\n")
		for _, mc := range ar.MiniChecks {
			p := strings.TrimSpace(mc.Prompt)
			if p == "" && mc.Raw != "" {
				p = mc.Raw
			}
			if p != "" {
				b.WriteString("— ")
				b.WriteString(p)
				b.WriteString("\n")
			}
		}
	}
	// Короткая подсказка по безопасности/антилику
	if !ar.LeakGuardPassed || !ar.Safety.NoOriginalAnswerLeak {
		b.WriteString("\n(Замечание: аналог без ссылок на исходные данные, ответы не раскрываются.)")
	}
	r.send(chatID, b.String())
}
