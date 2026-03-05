package telegram

import (
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	debounce  = 1200 * time.Millisecond
	maxPixels = 18_000_000
)

// TTL-кэши для разных типов данных
var (
	// Временные данные обработки (короткий TTL)
	pendingCtx = NewTTLCache("pendingCtx", PendingTTL)
	parseWait  = NewTTLCache("parseWait", PendingTTL)

	// Сессионные данные (средний TTL)
	hintState = NewTTLCache("hintState", SessionTTL)
	chatMode  = NewTTLCache("chatMode", SessionTTL)
	chatState = NewTTLCache("chatState", SessionTTL)

	// Данные пользователя (длинный TTL)
	userInfo = NewTTLCache("userInfo", UserDataTTL)
	chatInfo = NewTTLCache("chatInfo", UserDataTTL)

	// Флаг инициализации очистки
	cleanupInitOnce sync.Once

	// Per-chat locks для синхронизации state transitions
	chatLocks sync.Map // chatID -> *sync.Mutex
)

// InitCacheCleanup запускает фоновую очистку кэшей
// Вызывать один раз при старте приложения
func InitCacheCleanup() {
	cleanupInitOnce.Do(func() {
		manager := GetCacheManager()
		manager.Register(pendingCtx)
		manager.Register(parseWait)
		manager.Register(hintState)
		manager.Register(chatMode)
		manager.Register(chatState)
		manager.Register(userInfo)
		manager.Register(chatInfo)
		manager.Register(sessionByChat)
		manager.Register(batchSessionKeys) // для отслеживания session ID альбомов
		manager.Start()
	})
}

// StopCacheCleanup останавливает фоновую очистку (для graceful shutdown)
func StopCacheCleanup() {
	GetCacheManager().Stop()
}

// GetCacheStats возвращает статистику всех кэшей
func GetCacheStats() map[string]map[string]int64 {
	return GetCacheManager().GetAllStats()
}

// хелперы для режима чата
func setMode(chatID int64, mode string) {
	chatMode.Store(chatID, mode)
}

func getMode(chatID int64) string {
	if v, ok := chatMode.Load(chatID); ok {
		if s, _ := v.(string); s != "" {
			return s
		}
	}
	return ""
}

func clearMode(chatID int64) {
	chatMode.Delete(chatID)
}

type State string

var (
	AwaitingTask    State = "awaiting_task"
	CollectingPages State = "collecting_pages"
	Detect          State = "detect"
	Parse           State = "parse"
	Report          State = "report"
	Hints           State = "hint"
	AwaitSolution   State = "await_solution"
	Check           State = "check"
	Correct         State = "correct"
	Incorrect       State = "incorrect"
	Analogue        State = "analogue"
	AwaitGrade      State = "await_grade"
)

var States = map[State][]State{
	AwaitGrade:      {AwaitingTask, Report},
	AwaitingTask:    {CollectingPages, AwaitingTask, Report},
	CollectingPages: {Detect, Report, AwaitingTask},
	Detect:          {Parse, Report},
	Parse:           {Hints, AwaitSolution, Report},
	Hints:           {AwaitSolution, AwaitingTask, Analogue, Hints, Report},
	AwaitSolution:   {Check, Report, AwaitingTask},
	Check:           {Correct, Incorrect, Report, AwaitingTask, CollectingPages, Analogue},
	Correct:         {AwaitingTask, CollectingPages, Report},
	Incorrect:       {Analogue, AwaitingTask, CollectingPages, Report},
	Analogue:        {AwaitingTask, CollectingPages, AwaitSolution, Report},
	Report:          {AwaitingTask, CollectingPages, Report},
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

// getChatLock возвращает mutex для конкретного chatID
func getChatLock(chatID int64) *sync.Mutex {
	mu, _ := chatLocks.LoadOrStore(chatID, &sync.Mutex{})
	if m, ok := mu.(*sync.Mutex); ok {
		return m
	}
	// Fallback: создаём новый мьютекс если тип неверный (не должно происходить)
	newMu := &sync.Mutex{}
	chatLocks.Store(chatID, newMu)
	return newMu
}

func getState(chatID int64) State {
	if v, ok := chatState.Load(chatID); ok {
		if s, ok2 := v.(State); ok2 {
			// Touch безопасен даже если запись была удалена между Load и Touch
			// (он просто ничего не сделает если запись не найдена)
			chatState.Touch(chatID)
			return s
		}
	}

	// Используем Store вместо LoadOrStore для простоты
	// Если конкурентный Store произошёл - ничего страшного, состояние идемпотентно
	chatState.Store(chatID, AwaitingTask)
	return AwaitingTask
}

func setState(chatID int64, s State) {
	chatState.Store(chatID, s)
}

// getAndSetState атомарно получает текущее состояние и устанавливает новое
// Возвращает предыдущее состояние
func getAndSetState(chatID int64, newState State) State {
	mu := getChatLock(chatID)
	mu.Lock()
	defer mu.Unlock()

	current := getState(chatID)
	setState(chatID, newState)
	return current
}

// tryTransition пытается выполнить переход в новое состояние
// Возвращает (текущее состояние, успех перехода)
func tryTransition(chatID int64, newState State) (State, bool) {
	mu := getChatLock(chatID)
	mu.Lock()
	defer mu.Unlock()

	current := getState(chatID)
	if !canTransition(current, newState) {
		return current, false
	}
	setState(chatID, newState)
	return current, true
}

func friendlyState(s State) string {
	switch s {
	case AwaitGrade:
		return "Укажите класс"
	case AwaitingTask:
		return "Жду фото задачи"
	case CollectingPages:
		return "Сбор фото"
	case Detect:
		return "Детект"
	case Parse:
		return "Парсинг"
	case Hints:
		return "Подсказки"
	case AwaitSolution:
		return "Жду решение"
	case Check:
		return "Проверка решения"
	case Correct:
		return "Верно"
	case Incorrect:
		return "Есть ошибка"
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
	case Hints:
		return "\nДоступно: «Получить подсказку», «Готов дать решение», «Перейти к новой задаче»."
	case AwaitSolution:
		return "\nПришлите ваш ответ текстом или фото. Либо «Перейти к новой задаче»."
	case Incorrect:
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
		case "analogue", "analogue_task":
			return Analogue, true // после подсказок
		case "hint_next":
			return Hints, true
		case "parse_yes":
			return Hints, true
		case "dont_like_hint":
			return Hints, true
		case "ready_solution":
			return AwaitSolution, true
		case "new_task":
			return AwaitingTask, true
		case "report":
			return Report, true
		case "grade1", "grade2", "grade3", "grade4":
			return AwaitingTask, true
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
			return Check, true // прислано решение фото → сразу check
		}
		return CollectingPages, true // прислано фото задания/страницы
	}

	// 4.1) Документ-изображение (фото отправлено как файл)
	if upd.Message.Document != nil && strings.HasPrefix(upd.Message.Document.MimeType, "image/") {
		if cur == AwaitSolution {
			return Check, true
		}
		return CollectingPages, true
	}

	// 5) Текст
	if s := strings.TrimSpace(upd.Message.Text); s != "" {
		if cur == AwaitSolution {
			return Check, true // текстовое решение → сразу check
		}
		if cur == Report {
			return Report, true
		}
		if cur == AwaitGrade {
			return AwaitingTask, true
		}
		// Иначе текст вне контекста: останемся, где были
		return cur, false
	}

	return cur, false
}
