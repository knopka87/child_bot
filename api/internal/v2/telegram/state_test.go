package telegram

import (
	"sync"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestCanTransition(t *testing.T) {
	tests := []struct {
		name     string
		from     State
		to       State
		expected bool
	}{
		// Valid transitions
		{"AwaitingTask -> CollectingPages", AwaitingTask, CollectingPages, true},
		{"AwaitingTask -> Report", AwaitingTask, Report, true},
		{"CollectingPages -> Detect", CollectingPages, Detect, true},
		{"Detect -> Parse", Detect, Parse, true},
		{"Parse -> Hints", Parse, Hints, true},
		{"Parse -> AwaitSolution", Parse, AwaitSolution, true},
		{"Hints -> AwaitSolution", Hints, AwaitSolution, true},
		{"Hints -> AwaitingTask", Hints, AwaitingTask, true},
		{"Hints -> Analogue", Hints, Analogue, true},
		{"Hints -> Hints (self)", Hints, Hints, true},
		{"AwaitSolution -> Check", AwaitSolution, Check, true},
		{"AwaitSolution -> AwaitingTask", AwaitSolution, AwaitingTask, true},
		{"Check -> Correct", Check, Correct, true},
		{"Check -> Incorrect", Check, Incorrect, true},
		{"Incorrect -> Analogue", Incorrect, Analogue, true},
		{"Incorrect -> AwaitingTask", Incorrect, AwaitingTask, true},
		{"Correct -> AwaitingTask", Correct, AwaitingTask, true},

		// Invalid transitions
		{"AwaitingTask -> Check (invalid)", AwaitingTask, Check, false},
		{"AwaitingTask -> Correct (invalid)", AwaitingTask, Correct, false},
		{"Detect -> Hints (invalid)", Detect, Hints, false},
		{"Parse -> Check (invalid)", Parse, Check, false},
		{"Hints -> Check (valid)", Hints, Check, true},
		{"AwaitSolution -> Hints (invalid)", AwaitSolution, Hints, false},
		{"Correct -> Check (invalid)", Correct, Check, false},

		// Unknown state
		{"Unknown -> AwaitingTask", State("unknown"), AwaitingTask, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := canTransition(tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("canTransition(%s, %s) = %v, want %v", tt.from, tt.to, result, tt.expected)
			}
		})
	}
}

func TestGetSetState(t *testing.T) {
	chatID := int64(12345)

	// Clear any existing state
	chatState.Delete(chatID)

	// Test default state
	state := getState(chatID)
	if state != AwaitingTask {
		t.Errorf("Default state = %s, want %s", state, AwaitingTask)
	}

	// Test setState
	setState(chatID, Hints)
	state = getState(chatID)
	if state != Hints {
		t.Errorf("After setState(Hints), getState() = %s, want %s", state, Hints)
	}

	// Test different states
	testStates := []State{Parse, Check, Correct, Incorrect, AwaitSolution}
	for _, s := range testStates {
		setState(chatID, s)
		got := getState(chatID)
		if got != s {
			t.Errorf("After setState(%s), getState() = %s", s, got)
		}
	}

	// Cleanup
	chatState.Delete(chatID)
}

func TestGetAndSetState(t *testing.T) {
	chatID := int64(12346)

	// Clear any existing state
	chatState.Delete(chatID)

	// Initial state should be AwaitingTask
	prev := getAndSetState(chatID, Hints)
	if prev != AwaitingTask {
		t.Errorf("Initial getAndSetState returned %s, want %s", prev, AwaitingTask)
	}

	// Current state should be Hints
	current := getState(chatID)
	if current != Hints {
		t.Errorf("After getAndSetState, state = %s, want %s", current, Hints)
	}

	// Next call should return Hints and set Parse
	prev = getAndSetState(chatID, Parse)
	if prev != Hints {
		t.Errorf("Second getAndSetState returned %s, want %s", prev, Hints)
	}

	current = getState(chatID)
	if current != Parse {
		t.Errorf("After second getAndSetState, state = %s, want %s", current, Parse)
	}

	// Cleanup
	chatState.Delete(chatID)
}

func TestTryTransition(t *testing.T) {
	chatID := int64(12347)

	tests := []struct {
		name          string
		initialState  State
		targetState   State
		expectSuccess bool
	}{
		{"Valid: AwaitingTask -> CollectingPages", AwaitingTask, CollectingPages, true},
		{"Valid: CollectingPages -> Detect", CollectingPages, Detect, true},
		{"Valid: Detect -> Parse", Detect, Parse, true},
		{"Valid: Parse -> Hints", Parse, Hints, true},
		{"Invalid: AwaitingTask -> Check", AwaitingTask, Check, false},
		{"Invalid: Parse -> Correct", Parse, Correct, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set initial state
			chatState.Delete(chatID)
			setState(chatID, tt.initialState)

			prev, success := tryTransition(chatID, tt.targetState)

			if success != tt.expectSuccess {
				t.Errorf("tryTransition success = %v, want %v", success, tt.expectSuccess)
			}

			if prev != tt.initialState {
				t.Errorf("tryTransition returned prev = %s, want %s", prev, tt.initialState)
			}

			if tt.expectSuccess {
				current := getState(chatID)
				if current != tt.targetState {
					t.Errorf("After successful transition, state = %s, want %s", current, tt.targetState)
				}
			} else {
				current := getState(chatID)
				if current != tt.initialState {
					t.Errorf("After failed transition, state = %s, want %s (unchanged)", current, tt.initialState)
				}
			}
		})
	}

	// Cleanup
	chatState.Delete(chatID)
}

func TestTryTransition_Concurrent(t *testing.T) {
	chatID := int64(12348)
	chatState.Delete(chatID)
	setState(chatID, AwaitingTask)

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// Try to transition from AwaitingTask -> CollectingPages concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, success := tryTransition(chatID, CollectingPages)
			if success {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// At least one should succeed (likely all since they're transitioning to the same state)
	if successCount == 0 {
		t.Error("Expected at least one successful transition")
	}

	// Final state should be CollectingPages
	finalState := getState(chatID)
	if finalState != CollectingPages {
		t.Errorf("Final state = %s, want %s", finalState, CollectingPages)
	}

	// Cleanup
	chatState.Delete(chatID)
}

func TestFriendlyState(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{AwaitGrade, "Укажите класс"},
		{AwaitingTask, "Жду фото задачи"},
		{CollectingPages, "Сбор фото"},
		{Detect, "Детект"},
		{Parse, "Парсинг"},
		{Hints, "Подсказки"},
		{AwaitSolution, "Жду решение"},
		{Check, "Проверка решения"},
		{Correct, "Верно"},
		{Incorrect, "Есть ошибка"},
		{Analogue, "Похожее задание"},
		{Report, "📝 Сообщить об ошибке"},
		{State("unknown"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			result := friendlyState(tt.state)
			if result != tt.expected {
				t.Errorf("friendlyState(%s) = %q, want %q", tt.state, result, tt.expected)
			}
		})
	}
}

func TestInferNextState_Callbacks(t *testing.T) {
	tests := []struct {
		name         string
		callbackData string
		currentState State
		expectedNext State
		expectChange bool
	}{
		{"hint_next from Hints", "hint_next", Hints, Hints, true},
		{"parse_yes from Parse", "parse_yes", Parse, Hints, true},
		{"dont_like_hint from Hints", "dont_like_hint", Hints, Hints, true},
		{"ready_solution from Hints", "ready_solution", Hints, AwaitSolution, true},
		{"new_task from any", "new_task", Hints, AwaitingTask, true},
		{"report from any", "report", Hints, Report, true},
		{"analogue from Incorrect", "analogue", Incorrect, Analogue, true},
		{"analogue_task from Hints", "analogue_task", Hints, Analogue, true},
		{"grade1 from AwaitGrade", "grade1", AwaitGrade, AwaitingTask, true},
		{"grade2 from AwaitGrade", "grade2", AwaitGrade, AwaitingTask, true},
		{"unknown callback", "unknown_callback", Hints, Hints, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upd := tgbotapi.Update{
				CallbackQuery: &tgbotapi.CallbackQuery{
					Data: tt.callbackData,
				},
			}

			next, changed := inferNextState(upd, tt.currentState)

			if changed != tt.expectChange {
				t.Errorf("inferNextState changed = %v, want %v", changed, tt.expectChange)
			}

			if next != tt.expectedNext {
				t.Errorf("inferNextState next = %s, want %s", next, tt.expectedNext)
			}
		})
	}
}

func TestInferNextState_Photo(t *testing.T) {
	photoSizes := []tgbotapi.PhotoSize{{FileID: "test", Width: 100, Height: 100}}

	tests := []struct {
		name         string
		currentState State
		expectedNext State
	}{
		{"Photo from AwaitingTask -> CollectingPages", AwaitingTask, CollectingPages},
		{"Photo from AwaitSolution -> Check", AwaitSolution, Check},
		{"Photo from Hints -> Check", Hints, Check},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upd := tgbotapi.Update{
				Message: &tgbotapi.Message{
					Photo: photoSizes,
				},
			}

			next, changed := inferNextState(upd, tt.currentState)

			if !changed {
				t.Error("Expected state change for photo")
			}

			if next != tt.expectedNext {
				t.Errorf("inferNextState next = %s, want %s", next, tt.expectedNext)
			}
		})
	}
}

func TestInferNextState_Document(t *testing.T) {
	tests := []struct {
		name         string
		mimeType     string
		currentState State
		expectedNext State
		expectChange bool
	}{
		{"Image document from AwaitingTask", "image/jpeg", AwaitingTask, CollectingPages, true},
		{"Image document from AwaitSolution", "image/png", AwaitSolution, Check, true},
		{"PDF document", "application/pdf", AwaitingTask, AwaitingTask, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upd := tgbotapi.Update{
				Message: &tgbotapi.Message{
					Document: &tgbotapi.Document{
						MimeType: tt.mimeType,
					},
				},
			}

			next, changed := inferNextState(upd, tt.currentState)

			if changed != tt.expectChange {
				t.Errorf("inferNextState changed = %v, want %v", changed, tt.expectChange)
			}

			if next != tt.expectedNext {
				t.Errorf("inferNextState next = %s, want %s", next, tt.expectedNext)
			}
		})
	}
}

func TestInferNextState_Text(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		currentState State
		expectedNext State
		expectChange bool
	}{
		{"Text from AwaitSolution -> Check", "42", AwaitSolution, Check, true},
		{"Text from Report -> Report", "error description", Report, Report, true},
		{"Text from AwaitGrade -> AwaitingTask", "2", AwaitGrade, AwaitingTask, true},
		{"Text from AwaitingTask (no change)", "hello", AwaitingTask, AwaitingTask, false},
		{"Empty text", "   ", AwaitingTask, AwaitingTask, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upd := tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text: tt.text,
				},
			}

			next, changed := inferNextState(upd, tt.currentState)

			if changed != tt.expectChange {
				t.Errorf("inferNextState changed = %v, want %v", changed, tt.expectChange)
			}

			if next != tt.expectedNext {
				t.Errorf("inferNextState next = %s, want %s", next, tt.expectedNext)
			}
		})
	}
}

func TestInferNextState_Commands(t *testing.T) {
	tests := []struct {
		name         string
		command      string
		currentState State
		expectChange bool
	}{
		{"/start command", "/start", AwaitingTask, true},
		{"/health command", "/health", Hints, true},
		{"/engine command", "/engine gemini", Parse, true},
		{"/unknown command", "/unknown", AwaitingTask, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upd := tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text: tt.command,
					Entities: []tgbotapi.MessageEntity{
						{Type: "bot_command", Offset: 0, Length: len(tt.command)},
					},
				},
			}

			_, changed := inferNextState(upd, tt.currentState)

			if changed != tt.expectChange {
				t.Errorf("inferNextState changed = %v, want %v", changed, tt.expectChange)
			}
		})
	}
}

func TestModeHelpers(t *testing.T) {
	chatID := int64(12349)

	// Clear any existing mode
	chatMode.Delete(chatID)

	// Test default mode (empty)
	mode := getMode(chatID)
	if mode != "" {
		t.Errorf("Default mode = %q, want empty", mode)
	}

	// Test setMode
	setMode(chatID, "await_solution")
	mode = getMode(chatID)
	if mode != "await_solution" {
		t.Errorf("After setMode, getMode() = %q, want %q", mode, "await_solution")
	}

	// Test clearMode
	clearMode(chatID)
	mode = getMode(chatID)
	if mode != "" {
		t.Errorf("After clearMode, getMode() = %q, want empty", mode)
	}

	// Cleanup
	chatMode.Delete(chatID)
}
