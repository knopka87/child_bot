package telegram

import (
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"child-bot/api/internal/util"
)

const (
	debounce  = 1200 * time.Millisecond
	maxPixels = 18_000_000
)

var (
	pendingChoice sync.Map // chatID -> []string (tasks brief)
	pendingCtx    sync.Map // chatID -> *selectionContext
	parseWait     sync.Map // chatID -> *parsePending
	hintState     sync.Map // chatID -> *hintSession
	chatMode      sync.Map // chatID -> string: "", "await_solution", "await_new_task"
	chatState     sync.Map // chatID ->
)

// хелперы
func setMode(chatID int64, mode string) { chatMode.Store(chatID, mode) }
func getMode(chatID int64) string {
	if v, ok := chatMode.Load(chatID); ok {
		if s, _ := v.(string); s != "" {
			return s
		}
	}
	return ""
}
func clearMode(chatID int64) { chatMode.Delete(chatID) }

type State string

var (
	AwaitingTask    State = "awaiting_task"
	CollectingPages State = "collecting_pages"
	Detect          State = "detect"
	NeedsRescan     State = "need_rescan"
	NotATask        State = "not_a_task"
	Inappropriate   State = "inappropriate"
	DecideTasks     State = "decide_task"
	Parse           State = "parse"
	AutoPick        State = "auto_pick"
	AskChoice       State = "ask_choice"
	Report          State = "report"
	AnalyzeChoice   State = "analyze_choice"
	Hints           State = "hint"
	Confirm         State = "confirm"
	AnalogueTask    State = "analogue_task"
	AwaitSolution   State = "await_solution"
	OCR             State = "ocr"
	Normalize       State = "normalize"
	Check           State = "check"
	Correct         State = "correct"
	Incorrect       State = "incorrect"
	Uncertain       State = "uncertain"
	Analogue        State = "analogue"
)

var States = map[State][]State{
	AwaitingTask:    {CollectingPages, AwaitingTask, Report},
	CollectingPages: {Detect, Report, AwaitingTask},
	Detect:          {NeedsRescan, NotATask, Inappropriate, DecideTasks, Report},
	NeedsRescan:     {AwaitingTask, CollectingPages, Report},
	NotATask:        {AwaitingTask, CollectingPages, Report},
	Inappropriate:   {AwaitingTask, CollectingPages, Report},
	DecideTasks:     {Parse, AskChoice},
	AskChoice:       {Report, AnalyzeChoice},
	AnalyzeChoice:   {Parse, AwaitingTask, AnalyzeChoice, Report},
	Parse:           {Hints, AwaitSolution, Confirm, NeedsRescan, Report},
	Confirm:         {Hints, AwaitSolution, AwaitingTask, CollectingPages, Report},
	AutoPick:        {Hints, AwaitSolution, AwaitingTask, CollectingPages, Report},
	Hints:           {AwaitSolution, AwaitingTask, Analogue, Hints, Report},
	AwaitSolution:   {Normalize, OCR, Report},
	OCR:             {Normalize, Report},
	Normalize:       {Check, Report, AwaitingTask},
	Check:           {Correct, Incorrect, Uncertain, Report, AwaitingTask, Analogue},
	Correct:         {AwaitingTask, CollectingPages, Report},
	Incorrect:       {Analogue, AwaitingTask, CollectingPages, Report},
	Uncertain:       {Analogue, AwaitingTask, Report},
	Analogue:        {AwaitingTask, CollectingPages, AwaitSolution, Report},
}

// canTransition проверяет, можно ли перейти из from в to.
func canTransition(from, to State) bool {
	nexts, ok := States[from]
	if !ok {
		return false
	}
	for _, n := range nexts {
		if n == to {
			return true
		}
	}
	return false
}

func getState(chatID int64) State {
	if v, ok := chatState.Load(chatID); ok {
		if s, ok2 := v.(State); ok2 {
			return s
		}
	}

	chatState.Store(chatID, AwaitingTask)
	return AwaitingTask
}

func setState(chatID int64, s State) {
	chatState.Store(chatID, s)
}

// В схемe Mermaid помечено, что текст явно допустим в L0/L1/L2/L3 и AnalogueTask.
// У нас этих под-состояний нет, поэтому используем ближайшие «узлы», где мы реально ждём текст:
func isCanUserText(s State) bool {
	switch s {
	case Hints, AskChoice, AwaitSolution, Analogue: // упрощённое соответствие
		return true
	default:
		return false
	}
}

func friendlyState(s State) string {
	switch s {
	case AwaitingTask:
		return "Жду фото задачи"
	case CollectingPages:
		return "Сбор фото"
	case Detect:
		return "Детект"
	case NeedsRescan:
		return "Нужно перефотографировать"
	case NotATask:
		return "Это не задание"
	case Inappropriate:
		return "Неподходящее изображение"
	case DecideTasks:
		return "Выбор задачи"
	case Parse:
		return "Парсинг"
	case AutoPick:
		return "Автовыбор задачи"
	case AskChoice:
		return "Ожидаю номер задачи"
	case AnalyzeChoice:
		return "Анализ выбора"
	case Hints:
		return "Подсказки"
	case Confirm:
		return "Подтверждение"
	case AwaitSolution:
		return "Жду решение"
	case OCR:
		return "Парсинг ответа"
	case Normalize:
		return "Нормализация ответа"
	case Check:
		return "Проверка решения"
	case Correct:
		return "Верно"
	case Incorrect:
		return "Есть ошибка"
	case Uncertain:
		return "Не уверен"
	case Analogue:
		return "Похожее задание"
	case Report:
		return "📝 Сообщить об ошибке"
	default:
		return string(s)
	}
}

// Короткие подсказки по доступным действиям в текущем состоянии для пользователя
func allowedStateHints(cur State) string {
	switch cur {
	case AwaitingTask:
		return "\nМожно прислать фото задания (1–2 фото)."
	case AskChoice:
		return "\nПришлите номер задачи из списка (целое число 1..N) или нажмите «Сообщить об ошибке»."
	case Hints:
		return "\nДоступно: «Получить подсказку», «Готов дать решение», «Перейти к новой задаче»."
	case AwaitSolution:
		return "\nПришлите ваш ответ текстом или фото. Либо «Перейти к новой задаче»."
	case Incorrect, Uncertain:
		return "\nМожно запросить «Похожее задание» или «Перейти к новой задаче»."
	default:
		// По умолчанию — перечислим разрешённые состояния по карте переходов
		nexts := States[cur]
		if len(nexts) == 0 {
			return ""
		}
		var names []string
		for _, n := range nexts {
			names = append(names, friendlyState(n))
		}
		return "\nДоступные действия: " + strings.Join(names, ", ")
	}
}

// Пытаемся вывести желаемое следующее состояние по входящему апдейту.
// Второй флаг = true, если вообще есть предложение смены состояния.
func inferNextState(upd tgbotapi.Update, cur State) (State, bool) {
	// 1) Callback-и
	if upd.CallbackQuery != nil {
		switch strings.ToLower(strings.TrimSpace(upd.CallbackQuery.Data)) {
		case "analogue_task", "analogue":
			return Analogue, true
		case "hint_next":
			return Hints, true
		case "parse_yes":
			return Hints, true
		case "parse_no":
			return AwaitingTask, true
		case "ready_solution":
			return AwaitSolution, true
		case "new_task":
			return AwaitingTask, true
		case "report":
			return Report, true
		default:
			return cur, false
		}
	}

	// 2) Без сообщения — не меняем состояние
	if upd.Message == nil {
		return cur, false
	}

	// 3) Команды
	if upd.Message.IsCommand() {
		// /start, /health, /engine — считаем «сервисными», не меняющими логику ветки.
		cmd := strings.Fields(strings.TrimPrefix(upd.Message.Text, "/"))
		if len(cmd) > 0 {
			switch cmd[0] {
			case "start", "health":
				return cur, true
			case "engine":
				// провайдер переключится в другом месте; состояние оставим прежним либо AwaitingTask
				return cur, true
			case "hintL1", "hintL2", "hintL3":
				return cur, true
			}
		}
		// прочие команды — без смены
		return cur, false
	}

	// 4) Фото
	if upd.Message.Photo != nil && len(upd.Message.Photo) > 0 {
		if cur == AwaitSolution {
			return Normalize, true // прислано решение фото → нормализация
		}
		return CollectingPages, true // прислано фото задания/страницы
	}

	// 5) Текст
	if s := strings.TrimSpace(upd.Message.Text); s != "" {
		if v, ok := pendingChoice.Load(util.GetChatIDByTgUpdate(upd)); ok && v != nil {
			return AnalyzeChoice, true // ввод номера задачи 1..N
		}
		if cur == AwaitSolution {
			return Normalize, true // текстовое решение → нормализация
		}
		// Иначе текст вне контекста: останемся, где были
		return cur, false
	}

	return cur, false
}
