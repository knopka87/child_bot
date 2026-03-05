package telegram

import (
	"strings"
	"testing"

	"child-bot/api/internal/v2/types"
)

func TestLvlToConst(t *testing.T) {
	tests := []struct {
		level    int
		expected types.HintLevel
	}{
		{1, types.HintL1},
		{2, types.HintL2},
		{3, types.HintL3},
		{4, types.HintL3}, // Defaults to L3 for levels > 3
		{0, types.HintL3}, // Defaults to L3 for level < 1
		{-1, types.HintL3},
	}

	for _, tt := range tests {
		t.Run(string(tt.expected), func(t *testing.T) {
			result := lvlToConst(tt.level)
			if result != tt.expected {
				t.Errorf("lvlToConst(%d) = %s, want %s", tt.level, result, tt.expected)
			}
		})
	}
}

func TestFormatHint_Level1(t *testing.T) {
	hr := types.HintResponse{
		Items: []types.HintItem{
			{
				Hints: []types.Hint{
					{Level: types.HintL1, HintText: "Подсказка первого уровня"},
					{Level: types.HintL2, HintText: "Подсказка второго уровня"},
				},
			},
		},
	}

	result := formatHint(hr, 1)

	if !strings.Contains(result, "Подсказка первого уровня") {
		t.Error("Expected L1 hint text in result")
	}
	if strings.Contains(result, "Подсказка второго уровня") {
		t.Error("L2 hint text should not be in L1 result")
	}
	if !strings.Contains(result, "Подсказка 1") {
		t.Error("Expected 'Подсказка 1' header in result")
	}
}

func TestFormatHint_Level2(t *testing.T) {
	hr := types.HintResponse{
		Items: []types.HintItem{
			{
				Hints: []types.Hint{
					{Level: types.HintL1, HintText: "Подсказка первого уровня"},
					{Level: types.HintL2, HintText: "Подсказка второго уровня"},
				},
			},
		},
	}

	result := formatHint(hr, 2)

	if !strings.Contains(result, "Подсказка второго уровня") {
		t.Error("Expected L2 hint text in result")
	}
	if strings.Contains(result, "Подсказка первого уровня") {
		t.Error("L1 hint text should not be in L2 result")
	}
	if !strings.Contains(result, "Подсказка 2") {
		t.Error("Expected 'Подсказка 2' header in result")
	}
}

func TestFormatHint_Level3(t *testing.T) {
	hr := types.HintResponse{
		Items: []types.HintItem{
			{
				Hints: []types.Hint{
					{Level: types.HintL3, HintText: "Подсказка третьего уровня"},
				},
			},
		},
	}

	result := formatHint(hr, 3)

	if !strings.Contains(result, "Подсказка третьего уровня") {
		t.Error("Expected L3 hint text in result")
	}
	if !strings.Contains(result, "Подсказка 3") {
		t.Error("Expected 'Подсказка 3' header in result")
	}
}

func TestFormatHint_MultipleItems(t *testing.T) {
	hr := types.HintResponse{
		Items: []types.HintItem{
			{
				Hints: []types.Hint{
					{Level: types.HintL1, HintText: "Первая подсказка для первого пункта"},
				},
			},
			{
				Hints: []types.Hint{
					{Level: types.HintL1, HintText: "Первая подсказка для второго пункта"},
				},
			},
		},
	}

	result := formatHint(hr, 1)

	if !strings.Contains(result, "Первая подсказка для первого пункта") {
		t.Error("Expected first item hint text in result")
	}
	if !strings.Contains(result, "Первая подсказка для второго пункта") {
		t.Error("Expected second item hint text in result")
	}
}

func TestFormatHint_NoHintsForLevel(t *testing.T) {
	hr := types.HintResponse{
		Items: []types.HintItem{
			{
				Hints: []types.Hint{
					{Level: types.HintL1, HintText: "Only L1 hint"},
				},
			},
		},
	}

	result := formatHint(hr, 2)

	if result != "Подсказка не найдена" {
		t.Errorf("Expected 'Подсказка не найдена', got %q", result)
	}
}

func TestFormatHint_EmptyItems(t *testing.T) {
	hr := types.HintResponse{
		Items: []types.HintItem{},
	}

	result := formatHint(hr, 1)

	if result != "Подсказка не найдена" {
		t.Errorf("Expected 'Подсказка не найдена', got %q", result)
	}
}

func TestFormatHint_EmptyHints(t *testing.T) {
	hr := types.HintResponse{
		Items: []types.HintItem{
			{
				Hints: []types.Hint{},
			},
		},
	}

	result := formatHint(hr, 1)

	if result != "Подсказка не найдена" {
		t.Errorf("Expected 'Подсказка не найдена', got %q", result)
	}
}

func TestFormatHint_ContainsDivider(t *testing.T) {
	hr := types.HintResponse{
		Items: []types.HintItem{
			{
				Hints: []types.Hint{
					{Level: types.HintL1, HintText: "Test hint"},
				},
			},
		},
	}

	result := formatHint(hr, 1)

	if !strings.Contains(result, Divider) {
		t.Error("Expected result to contain divider")
	}
}

func TestFormatHint_MultipleHintsSameLevel(t *testing.T) {
	hr := types.HintResponse{
		Items: []types.HintItem{
			{
				Hints: []types.Hint{
					{Level: types.HintL1, HintText: "First L1 hint"},
					{Level: types.HintL1, HintText: "Second L1 hint"},
				},
			},
		},
	}

	result := formatHint(hr, 1)

	if !strings.Contains(result, "First L1 hint") {
		t.Error("Expected first L1 hint in result")
	}
	if !strings.Contains(result, "Second L1 hint") {
		t.Error("Expected second L1 hint in result")
	}
}
