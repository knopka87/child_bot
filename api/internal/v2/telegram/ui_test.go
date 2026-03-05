package telegram

import (
	"strings"
	"testing"
)

func TestMakeGradeListButtons(t *testing.T) {
	buttons := makeGradeListButtons()

	if len(buttons) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(buttons))
	}

	// First row: grades 1 and 2
	if len(buttons[0]) != 2 {
		t.Errorf("Expected 2 buttons in first row, got %d", len(buttons[0]))
	}
	if buttons[0][0].Text != Grade1Button {
		t.Errorf("First button text = %q, want %q", buttons[0][0].Text, Grade1Button)
	}
	if buttons[0][1].Text != Grade2Button {
		t.Errorf("Second button text = %q, want %q", buttons[0][1].Text, Grade2Button)
	}

	// Second row: grades 3 and 4
	if len(buttons[1]) != 2 {
		t.Errorf("Expected 2 buttons in second row, got %d", len(buttons[1]))
	}
	if buttons[1][0].Text != Grade3Button {
		t.Errorf("Third button text = %q, want %q", buttons[1][0].Text, Grade3Button)
	}
	if buttons[1][1].Text != Grade4Button {
		t.Errorf("Fourth button text = %q, want %q", buttons[1][1].Text, Grade4Button)
	}

	// Check callback data
	if *buttons[0][0].CallbackData != "grade1" {
		t.Errorf("Grade1 callback = %q, want %q", *buttons[0][0].CallbackData, "grade1")
	}
}

func TestMakeErrorButtons(t *testing.T) {
	buttons := makeErrorButtons()

	if len(buttons) != 1 {
		t.Errorf("Expected 1 row, got %d", len(buttons))
	}

	if len(buttons[0]) != 1 {
		t.Errorf("Expected 1 button in row, got %d", len(buttons[0]))
	}

	if buttons[0][0].Text != NewTaskButton {
		t.Errorf("Button text = %q, want %q", buttons[0][0].Text, NewTaskButton)
	}

	if *buttons[0][0].CallbackData != "new_task" {
		t.Errorf("Callback data = %q, want %q", *buttons[0][0].CallbackData, "new_task")
	}
}

func TestMakeParseConfirmButtons(t *testing.T) {
	buttons := makeParseConfirmButtons()

	if len(buttons) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(buttons))
	}

	// First row: Yes button
	if buttons[0][0].Text != YesButton {
		t.Errorf("First button text = %q, want %q", buttons[0][0].Text, YesButton)
	}
	if *buttons[0][0].CallbackData != "parse_yes" {
		t.Errorf("First button callback = %q, want %q", *buttons[0][0].CallbackData, "parse_yes")
	}

	// Second row: Check answer button
	if buttons[1][0].Text != CheckAnswerButton {
		t.Errorf("Second button text = %q, want %q", buttons[1][0].Text, CheckAnswerButton)
	}

	// Third row: Report button
	if buttons[2][0].Text != SendReportButton {
		t.Errorf("Third button text = %q, want %q", buttons[2][0].Text, SendReportButton)
	}
}

func TestMakeHintButtons(t *testing.T) {
	tests := []struct {
		name         string
		level        int
		maxHints     int
		showAnalogue bool
		expectRows   int
		firstButton  string
	}{
		{
			name:         "Level 1 of 3 - shows next hint",
			level:        1,
			maxHints:     3,
			showAnalogue: true,
			expectRows:   4,
			firstButton:  NextHintButton,
		},
		{
			name:         "Level 2 of 3 - shows next hint",
			level:        2,
			maxHints:     3,
			showAnalogue: true,
			expectRows:   4,
			firstButton:  NextHintButton,
		},
		{
			name:         "Level 3 of 3 with analogue - shows analogue",
			level:        3,
			maxHints:     3,
			showAnalogue: true,
			expectRows:   4,
			firstButton:  AnalogueTaskButton,
		},
		{
			name:         "Level 3 of 3 without analogue - no first button",
			level:        3,
			maxHints:     3,
			showAnalogue: false,
			expectRows:   3,
			firstButton:  CheckAnswerButton,
		},
		{
			name:         "Level 1 of 1 with analogue - shows analogue",
			level:        1,
			maxHints:     1,
			showAnalogue: true,
			expectRows:   4,
			firstButton:  AnalogueTaskButton,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buttons := makeHintButtons(tt.level, tt.maxHints, tt.showAnalogue)

			if len(buttons) != tt.expectRows {
				t.Errorf("Expected %d rows, got %d", tt.expectRows, len(buttons))
			}

			if len(buttons) > 0 && len(buttons[0]) > 0 {
				if buttons[0][0].Text != tt.firstButton {
					t.Errorf("First button = %q, want %q", buttons[0][0].Text, tt.firstButton)
				}
			}

			// Check that common buttons are present
			hasReadySolution := false
			hasDontLike := false
			hasNewTask := false

			for _, row := range buttons {
				for _, btn := range row {
					switch btn.Text {
					case CheckAnswerButton:
						hasReadySolution = true
					case DontLikeHintButton:
						hasDontLike = true
					case NewTaskButton:
						hasNewTask = true
					}
				}
			}

			if !hasReadySolution {
				t.Error("Missing ready_solution button")
			}
			if !hasDontLike {
				t.Error("Missing dont_like_hint button")
			}
			if !hasNewTask {
				t.Error("Missing new_task button")
			}
		})
	}
}

func TestMakeFinishHintButtons(t *testing.T) {
	buttons := makeFinishHintButtons()

	if len(buttons) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(buttons))
	}

	expectedButtons := []string{CheckAnswerButton, NewTaskButton, SendReportButton}
	for i, expected := range expectedButtons {
		if buttons[i][0].Text != expected {
			t.Errorf("Row %d button = %q, want %q", i, buttons[i][0].Text, expected)
		}
	}
}

func TestMakeAnalogueButtons(t *testing.T) {
	buttons := makeAnalogueButtons()

	if len(buttons) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(buttons))
	}

	expectedButtons := []string{CheckAnswerButton, NewTaskButton, SendReportButton}
	for i, expected := range expectedButtons {
		if buttons[i][0].Text != expected {
			t.Errorf("Row %d button = %q, want %q", i, buttons[i][0].Text, expected)
		}
	}
}

func TestMakeCheckAnswerClickButtons(t *testing.T) {
	buttons := makeCheckAnswerClickButtons()

	if len(buttons) != 1 {
		t.Errorf("Expected 1 row, got %d", len(buttons))
	}

	if buttons[0][0].Text != NewTaskButton {
		t.Errorf("Button text = %q, want %q", buttons[0][0].Text, NewTaskButton)
	}
}

func TestMakeCorrectAnswerButtons(t *testing.T) {
	buttons := makeCorrectAnswerButtons()

	if len(buttons) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(buttons))
	}

	if buttons[0][0].Text != NewTaskButton {
		t.Errorf("First button = %q, want %q", buttons[0][0].Text, NewTaskButton)
	}
	if buttons[1][0].Text != SendReportButton {
		t.Errorf("Second button = %q, want %q", buttons[1][0].Text, SendReportButton)
	}
}

func TestMakeIncorrectAnswerButtons(t *testing.T) {
	buttons := makeIncorrectAnswerButtons()

	if len(buttons) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(buttons))
	}

	if buttons[0][0].Text != AnalogueTaskButton {
		t.Errorf("First button = %q, want %q", buttons[0][0].Text, AnalogueTaskButton)
	}
	if buttons[1][0].Text != NewTaskButton {
		t.Errorf("Second button = %q, want %q", buttons[1][0].Text, NewTaskButton)
	}
	if buttons[2][0].Text != SendReportButton {
		t.Errorf("Third button = %q, want %q", buttons[2][0].Text, SendReportButton)
	}
}

func TestUIConstants(t *testing.T) {
	// Check that divider is present
	if Divider == "" {
		t.Error("Divider constant is empty")
	}

	// Check that message templates contain placeholders where expected
	templatesWithPlaceholder := []struct {
		name     string
		template string
	}{
		{"TaskViewText", TaskViewText},
		{"HINT1Text", HINT1Text},
		{"HINT2Text", HINT2Text},
		{"HINT3Text", HINT3Text},
		{"AnalogueTaskText", AnalogueTaskText},
		{"AnswerIncorrectText", AnswerIncorrectText},
	}

	for _, tt := range templatesWithPlaceholder {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.template, "%s") {
				t.Errorf("%s should contain %%s placeholder", tt.name)
			}
		})
	}

	// Check that templates contain dividers
	templatesWithDivider := []struct {
		name     string
		template string
	}{
		{"TaskViewText", TaskViewText},
		{"HINT1Text", HINT1Text},
		{"HINT2Text", HINT2Text},
		{"HINT3Text", HINT3Text},
		{"AnalogueTaskText", AnalogueTaskText},
		{"AnswerCorrectText", AnswerCorrectText},
		{"AnswerIncorrectText", AnswerIncorrectText},
	}

	for _, tt := range templatesWithDivider {
		t.Run(tt.name+"_has_divider", func(t *testing.T) {
			if !strings.Contains(tt.template, Divider) {
				t.Errorf("%s should contain Divider", tt.name)
			}
		})
	}
}

func TestProgressConstants(t *testing.T) {
	// Check that progress messages are not empty
	progressMessages := []struct {
		name    string
		message string
	}{
		{"ParseProgress1", ParseProgress1},
		{"ParseProgress2", ParseProgress2},
		{"ParseProgress3", ParseProgress3},
		{"ParseProgress4", ParseProgress4},
		{"HintProgress1", HintProgress1},
		{"HintProgress2", HintProgress2},
		{"HintProgress3", HintProgress3},
		{"CheckProgress1", CheckProgress1},
		{"CheckProgress2", CheckProgress2},
		{"CheckProgress3", CheckProgress3},
		{"CheckProgress4", CheckProgress4},
	}

	for _, tt := range progressMessages {
		t.Run(tt.name, func(t *testing.T) {
			if tt.message == "" {
				t.Errorf("%s is empty", tt.name)
			}
		})
	}
}
