package telegram

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"

	"child-bot/api/internal/store"
	"child-bot/api/internal/v2/types"
)

// tryTransitionWithPersist выполняет переход состояния и сохраняет в БД
// Возвращает (предыдущее состояние, успех перехода)
func (r *Router) tryTransitionWithPersist(chatID int64, newState State) (State, bool) {
	prev, ok := tryTransition(chatID, newState)
	if ok {
		// Сохраняем новое состояние в БД асинхронно
		r.persistCurrentState(chatID)
	}
	return prev, ok
}

// persistCurrentState сохраняет текущее состояние и режим в БД
func (r *Router) persistCurrentState(chatID int64) {
	shutdown := GetShutdownManager()
	if shutdown.IsShutdown() {
		return
	}

	// Копируем данные перед запуском горутины
	stateStr := string(getState(chatID))
	modeStr := getMode(chatID)
	var modePtr *string
	if modeStr != "" {
		modePtr = &modeStr
	}

	done := shutdown.TrackGoroutine()
	go func() {
		defer done()
		if shutdown.IsShutdown() {
			return
		}
		if err := r.Store.UpdateSessionState(context.Background(), chatID, &stateStr, modePtr); err != nil {
			log.Printf("[state_persistence] failed to persist state for chat %d: %v", chatID, err)
		}
	}()
}

// setStateWithPersist устанавливает состояние и сохраняет в БД
func (r *Router) setStateWithPersist(chatID int64, s State) {
	setState(chatID, s)

	// Проверяем shutdown перед запуском горутины
	shutdown := GetShutdownManager()
	if shutdown.IsShutdown() {
		return
	}

	// Регистрируем горутину для отслеживания
	done := shutdown.TrackGoroutine()

	go func() {
		defer done()

		// Повторная проверка после запуска
		if shutdown.IsShutdown() {
			return
		}

		stateStr := string(s)
		modeStr := getMode(chatID)
		var modePtr *string
		if modeStr != "" {
			modePtr = &modeStr
		}
		if err := r.Store.UpdateSessionState(context.Background(), chatID, &stateStr, modePtr); err != nil {
			log.Printf("[state_persistence] failed to persist state for chat %d: %v", chatID, err)
		}
	}()
}

// setModeWithPersist устанавливает режим и сохраняет в БД
func (r *Router) setModeWithPersist(chatID int64, mode string) {
	setMode(chatID, mode)

	shutdown := GetShutdownManager()
	if shutdown.IsShutdown() {
		return
	}
	done := shutdown.TrackGoroutine()

	go func() {
		defer done()
		if shutdown.IsShutdown() {
			return
		}
		stateStr := string(getState(chatID))
		if err := r.Store.UpdateSessionState(context.Background(), chatID, &stateStr, &mode); err != nil {
			log.Printf("[state_persistence] failed to persist mode for chat %d: %v", chatID, err)
		}
	}()
}

// clearModeWithPersist очищает режим и сохраняет в БД
func (r *Router) clearModeWithPersist(chatID int64) {
	clearMode(chatID)

	shutdown := GetShutdownManager()
	if shutdown.IsShutdown() {
		return
	}
	done := shutdown.TrackGoroutine()

	go func() {
		defer done()
		if shutdown.IsShutdown() {
			return
		}
		stateStr := string(getState(chatID))
		if err := r.Store.UpdateSessionState(context.Background(), chatID, &stateStr, nil); err != nil {
			log.Printf("[state_persistence] failed to clear mode for chat %d: %v", chatID, err)
		}
	}()
}

// restoreStateFromDB восстанавливает состояние из БД если его нет в кэше
// Возвращает true если состояние было восстановлено
func (r *Router) restoreStateFromDB(chatID int64) bool {
	// Проверяем, есть ли состояние в кэше
	stateInCache := false
	if v, ok := chatState.Load(chatID); ok {
		if _, ok2 := v.(State); ok2 {
			stateInCache = true
		}
	}

	// Проверяем, есть ли session_id в кэше
	sessionInCache := false
	if v, ok := sessionByChat.Load(chatID); ok {
		if sid, ok2 := v.(string); ok2 && sid != "" {
			sessionInCache = true
		}
	}

	// Если оба в кэше — ничего не делаем
	if stateInCache && sessionInCache {
		return false
	}

	// Пробуем загрузить из БД
	session, err := r.Store.FindSession(context.Background(), chatID)
	if err != nil || session.SessionID == "" {
		return false
	}

	// Кэшируем session_id для быстрого доступа через getSession
	if !sessionInCache {
		sessionByChat.Store(chatID, session.SessionID)
	}

	// Если состояние уже в кэше — не восстанавливаем остальное
	if stateInCache {
		return false
	}

	restored := false
	var restoredState State

	// Восстанавливаем состояние
	if session.CurrentState != nil && *session.CurrentState != "" {
		state := State(*session.CurrentState)
		// Проверяем что это валидное состояние
		if isValidState(state) {
			restoredState = state
			chatState.Store(chatID, state)
			restored = true
			log.Printf("[state_persistence] restored state '%s' for chat %d", state, chatID)
		}
	}

	// Восстанавливаем режим
	if session.ChatMode != nil && *session.ChatMode != "" {
		chatMode.Store(chatID, *session.ChatMode)
		log.Printf("[state_persistence] restored mode '%s' for chat %d", *session.ChatMode, chatID)
	}

	// Восстанавливаем контекст подсказок
	hintContextRestored := false
	if len(session.HintContext) > 0 {
		if hs, err := r.restoreHintSession(session.HintContext); err == nil && hs != nil {
			hintState.Store(chatID, hs)
			log.Printf("[state_persistence] restored hint context for chat %d (level %d/%d)",
				chatID, hs.NextLevel, hs.MaxHints)
			hintContextRestored = true
			restored = true
		}
	}

	// Проверяем консистентность состояния и контекста
	// Некоторые состояния требуют наличия контекста для продолжения работы
	if restored {
		var newState State
		var newMode string

		switch restoredState {
		case CollectingPages:
			// Сбор фото альбома — pendingCtx не персистится
			// После редеплоя пользователю нужно заново загрузить фото
			log.Printf("[state_persistence] state '%s' requires pending context which is not persisted, resetting to AwaitingTask for chat %d",
				restoredState, chatID)
			newState = AwaitingTask
			newMode = "await_new_task"

		case Detect, Parse:
			// Эти состояния требуют parseWait контекст, который не персистится
			// После редеплоя пользователю нужно заново загрузить фото
			log.Printf("[state_persistence] state '%s' requires parse context which is not persisted, resetting to AwaitingTask for chat %d",
				restoredState, chatID)
			newState = AwaitingTask
			newMode = "await_new_task"

		case Hints:
			// Состояние Hints требует hintContext
			if !hintContextRestored {
				log.Printf("[state_persistence] state 'Hints' requires hint context which is missing, resetting to AwaitingTask for chat %d",
					chatID)
				newState = AwaitingTask
				newMode = "await_new_task"
			}

		case Check:
			// Проверка решения — промежуточное состояние, данные в памяти
			// После редеплоя пользователю нужно заново отправить решение
			log.Printf("[state_persistence] state 'Check' is transient, resetting to AwaitSolution for chat %d", chatID)
			newState = AwaitSolution
			newMode = "await_solution"

		case Analogue:
			// Генерация аналогичной задачи — данные в памяти
			// После редеплоя пользователю нужно запросить аналогию заново или начать новую задачу
			log.Printf("[state_persistence] state 'Analogue' is transient, resetting to AwaitingTask for chat %d", chatID)
			newState = AwaitingTask
			newMode = "await_new_task"
		}

		if newState != "" {
			chatState.Store(chatID, newState)
			chatMode.Store(chatID, newMode)
			log.Printf("[state_persistence] reset to %s due to missing context for chat %d", newState, chatID)
		}
	}

	return restored
}

// isValidState проверяет валидность состояния
func isValidState(s State) bool {
	validStates := []State{
		AwaitingTask, CollectingPages, Detect, Parse, Report,
		Hints, AwaitSolution, Check, Correct, Incorrect, Analogue, AwaitGrade,
	}
	for _, vs := range validStates {
		if s == vs {
			return true
		}
	}
	return false
}

// saveHintContext сохраняет контекст подсказок в БД
func (r *Router) saveHintContext(chatID int64, hs *hintSession) {
	if hs == nil {
		return
	}

	shutdown := GetShutdownManager()
	if shutdown.IsShutdown() {
		return
	}

	// Копируем данные под мьютексом перед запуском горутины
	hs.mu.Lock()
	parseData := hs.Parse
	detectData := hs.Detect
	engineName := hs.EngineName
	nextLevel := hs.NextLevel
	maxHints := hs.MaxHints
	cachedHints := hs.CachedHints
	// Для больших изображений очищаем ссылку чтобы освободить память
	var imageCopy []byte
	if len(hs.Image) > 0 && len(hs.Image) < 100*1024 {
		imageCopy = make([]byte, len(hs.Image))
		copy(imageCopy, hs.Image)
	}
	// Очищаем большие изображения из памяти (>100KB не сохраняются в БД)
	if len(hs.Image) > 100*1024 {
		hs.Image = nil
	}
	hs.mu.Unlock()

	done := shutdown.TrackGoroutine()

	go func() {
		defer done()
		if shutdown.IsShutdown() {
			return
		}

		parseJSON, _ := json.Marshal(parseData)
		detectJSON, _ := json.Marshal(detectData)

		var imageBase64 string
		if len(imageCopy) > 0 {
			imageBase64 = base64.StdEncoding.EncodeToString(imageCopy)
		}

		var cachedHintsJSON []byte
		if cachedHints != nil {
			cachedHintsJSON, _ = json.Marshal(cachedHints)
		}

		data := store.HintContextData{
			ParseJSON:       parseJSON,
			DetectJSON:      detectJSON,
			EngineName:      engineName,
			NextLevel:       nextLevel,
			MaxHints:        maxHints,
			ImageBase64:     imageBase64,
			CachedHintsJSON: cachedHintsJSON,
		}

		hintContextJSON, err := json.Marshal(data)
		if err != nil {
			log.Printf("[state_persistence] failed to marshal hint context: %v", err)
			return
		}

		if err := r.Store.UpdateSessionHintContext(context.Background(), chatID, hintContextJSON); err != nil {
			log.Printf("[state_persistence] failed to save hint context for chat %d: %v", chatID, err)
		}
	}()
}

// clearHintContext очищает контекст подсказок в БД
func (r *Router) clearHintContext(chatID int64) {
	shutdown := GetShutdownManager()
	if shutdown.IsShutdown() {
		return
	}
	done := shutdown.TrackGoroutine()

	go func() {
		defer done()
		if shutdown.IsShutdown() {
			return
		}
		if err := r.Store.ClearSessionHintContext(context.Background(), chatID); err != nil {
			log.Printf("[state_persistence] failed to clear hint context for chat %d: %v", chatID, err)
		}
	}()
}

// restoreHintSession восстанавливает hintSession из JSON
func (r *Router) restoreHintSession(hintContextJSON []byte) (*hintSession, error) {
	var data store.HintContextData
	if err := json.Unmarshal(hintContextJSON, &data); err != nil {
		return nil, err
	}

	hs := &hintSession{
		EngineName: data.EngineName,
		NextLevel:  data.NextLevel,
		MaxHints:   data.MaxHints,
	}

	// Восстанавливаем Parse
	if len(data.ParseJSON) > 0 {
		if err := json.Unmarshal(data.ParseJSON, &hs.Parse); err != nil {
			log.Printf("[state_persistence] failed to unmarshal parse: %v", err)
		}
	}

	// Восстанавливаем Detect
	if len(data.DetectJSON) > 0 {
		if err := json.Unmarshal(data.DetectJSON, &hs.Detect); err != nil {
			log.Printf("[state_persistence] failed to unmarshal detect: %v", err)
		}
	}

	// Восстанавливаем изображение
	if data.ImageBase64 != "" {
		if imageBytes, err := base64.StdEncoding.DecodeString(data.ImageBase64); err == nil {
			hs.Image = imageBytes
		}
	}

	// Восстанавливаем кэш подсказок
	if len(data.CachedHintsJSON) > 0 {
		var cachedHints types.HintResponse
		if err := json.Unmarshal(data.CachedHintsJSON, &cachedHints); err == nil {
			hs.CachedHints = &cachedHints
		} else {
			log.Printf("[state_persistence] failed to unmarshal cached hints: %v", err)
		}
	}

	return hs, nil
}

// resetContextWithPersist сбрасывает контекст, очищает в БД и создаёт новую сессию
func (r *Router) resetContextWithPersist(cid int64) {
	// Очищаем кэши
	hintState.Delete(cid)
	pendingCtx.Delete(cid)
	parseWait.Delete(cid)
	batchSessionKeys.Delete(cid)
	setMode(cid, "await_new_task")
	setState(cid, AwaitingTask)

	// Очищаем старую сессию и создаём новую, чтобы последующие действия были привязаны к ней
	r.clearSession(cid)
	r.ensureSession(cid)

	shutdown := GetShutdownManager()
	if shutdown.IsShutdown() {
		return
	}
	done := shutdown.TrackGoroutine()

	// Очищаем в БД
	go func() {
		defer done()
		if shutdown.IsShutdown() {
			return
		}
		stateStr := string(AwaitingTask)
		mode := "await_new_task"
		_ = r.Store.UpdateSessionState(context.Background(), cid, &stateStr, &mode)
		_ = r.Store.ClearSessionHintContext(context.Background(), cid)
	}()
}

// HintContextForPersist — облегчённая версия hintSession для сохранения
// (без mutex и с минимумом данных)
type HintContextForPersist struct {
	Parse       types.ParseResponse  `json:"parse"`
	Detect      types.DetectResponse `json:"detect"`
	EngineName  string               `json:"engine_name"`
	NextLevel   int                  `json:"next_level"`
	MaxHints    int                  `json:"max_hints"`
	CachedHints *types.HintResponse  `json:"cached_hints,omitempty"`
}
