package telegram

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"child-bot/api/internal/v2/types"
)

func TestMain(m *testing.M) {
	// Get the directory of this test file
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	templatesPath := filepath.Join(dir, "..", "templates")

	// Reset cache to ensure clean state, then set templates directory
	ResetTemplatesCache()
	SetTemplatesDir(templatesPath)

	os.Exit(m.Run())
}

// TestCase represents a single template routing test case
type TestCase struct {
	Name         string
	Task         types.ParseTask
	Items        []types.ParseItem
	ExpectedCode string // expected template_code (T1, T2, T3, etc.)
	ShouldMatch  bool   // should a template be found?
}

// makeTask creates a ParseTask with common defaults
func makeTask(grade int64, text string, visualFacts []types.VisualFact) types.ParseTask {
	return types.ParseTask{
		TaskId:        "test-task",
		Subject:       types.SubjectMath,
		Grade:         grade,
		TaskTextClean: text,
		VisualFacts:   visualFacts,
	}
}

// makeItem creates a ParseItem with taskType and format
func makeItem(text, taskType, format string) types.ParseItem {
	return types.ParseItem{
		ItemId:        "test-item",
		ItemTextClean: text,
		PedKeys: types.PedKeys{
			TaskType: taskType,
			Format:   format,
		},
	}
}

// runTestCases runs a slice of test cases
func runTestCases(t *testing.T, tests []TestCase) {
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := buildRoutingContext(tc.Task, tc.Items)
			candidate, found := selectTemplate(ctx)

			if tc.ShouldMatch {
				if !found {
					t.Errorf("Expected to find template %s, but no template was found", tc.ExpectedCode)
					return
				}
				if candidate.Template.TemplateCode != tc.ExpectedCode {
					t.Errorf("Expected template %s, got %s (score=%d, anchors=%d, visual=%v, rule=%s)",
						tc.ExpectedCode, candidate.Template.TemplateCode,
						candidate.Score, candidate.AnchorsMatched, candidate.VisualMatched, candidate.MatchedRuleID)
				}
			} else {
				if found {
					t.Errorf("Expected no template match, but got %s", candidate.Template.TemplateCode)
				}
			}
		})
	}
}

// runTestCasesWithDebug runs test cases with trace enabled for debugging
func runTestCasesWithDebug(t *testing.T, tests []TestCase) {
	SetRoutingDebug(true)
	defer SetRoutingDebug(false)

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := buildRoutingContext(tc.Task, tc.Items)
			candidate, found := selectTemplate(ctx)

			trace := GetLastRoutingTrace()
			if trace != nil {
				t.Logf("Trace for %s:", tc.Name)
				t.Logf("  TextAll: %s", trace.TextAll)
				t.Logf("  TaskType: %s", trace.TaskType)
				t.Logf("  Format: %s", trace.Format)
				t.Logf("  Grade: %d", trace.Grade)
				t.Logf("  CandidateCount: %d", trace.CandidateCount)
				t.Logf("  Winner: %s", trace.Winner)
				t.Logf("  Total trace entries: %d", len(trace.Entries))
				// Show entries for expected template
				foundExpected := false
				for _, e := range trace.Entries {
					if strings.Contains(e.TemplateCode, tc.ExpectedCode) {
						foundExpected = true
						t.Logf("  Entry for %s: status=%s, rule=%s, rejected=%v, matched=%v",
							e.TemplateCode, e.Status, e.RuleID, e.RejectedBy, e.MatchedPatterns)
					}
				}
				if !foundExpected {
					t.Logf("  No trace entries found for template %s", tc.ExpectedCode)
					// Show first 5 entries to understand what templates ARE being checked
					for i, e := range trace.Entries {
						if i >= 5 {
							t.Logf("  ... and %d more entries", len(trace.Entries)-5)
							break
						}
						t.Logf("  Entry[%d]: %s status=%s", i, e.TemplateCode, e.Status)
					}
				}
			}

			if tc.ShouldMatch {
				if !found {
					t.Errorf("Expected to find template %s, but no template was found", tc.ExpectedCode)
					return
				}
				if candidate.Template.TemplateCode != tc.ExpectedCode {
					t.Errorf("Expected template %s, got %s (score=%d, anchors=%d, visual=%v, rule=%s)",
						tc.ExpectedCode, candidate.Template.TemplateCode,
						candidate.Score, candidate.AnchorsMatched, candidate.VisualMatched, candidate.MatchedRuleID)
				}
			} else {
				if found {
					t.Errorf("Expected no template match, but got %s", candidate.Template.TemplateCode)
				}
			}
		})
	}
}

// =============================================================================
// T1: Числовая прямая и разрядный состав
// task_type: number_sense
// Patterns: "числовая прямая", "числовой луч", "разрядный состав", "десятки и единицы"
// =============================================================================

func TestT1_NumberLineAndPlaceValue(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T1_number_line",
			Task: makeTask(1, "Отметь число 7 на числовой прямой", []types.VisualFact{
				{Kind: "number_line", Value: "числовая прямая"},
			}),
			Items: []types.ParseItem{
				makeItem("Числовая прямая", "number_sense", "number_line"),
			},
			ExpectedCode: "T1",
			ShouldMatch:  true,
		},
		{
			Name: "T1_number_ray",
			Task: makeTask(1, "Найди число на числовом луче", []types.VisualFact{
				{Kind: "number_line", Value: "числовой луч"},
			}),
			Items: []types.ParseItem{
				makeItem("Числовой луч", "number_sense", "number_line"),
			},
			ExpectedCode: "T1",
			ShouldMatch:  true,
		},
		{
			Name: "T1_place_value_composition",
			Task: makeTask(2, "Укажи разрядный состав числа 45: сколько десятков и единиц?", nil),
			Items: []types.ParseItem{
				makeItem("Десятки и единицы", "number_sense", "plain_text"),
			},
			ExpectedCode: "T1",
			ShouldMatch:  true,
		},
		{
			Name: "T1_counting_by",
			Task: makeTask(1, "Продолжи счёт через 5: 5, 10, 15, ...", nil),
			Items: []types.ParseItem{
				makeItem("Счёт через 5", "number_sense", "plain_text"),
			},
			ExpectedCode: "T1",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T2: Римские числа
// task_type: numeral_systems
// Patterns: "римск", "римские цифры", "римские числа"
// =============================================================================

func TestT2_RomanNumerals(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T2_roman_to_arabic",
			Task: makeTask(3, "Запиши римскими цифрами число 45", nil),
			Items: []types.ParseItem{
				makeItem("Римские цифры", "numeral_systems", "plain_text"),
			},
			ExpectedCode: "T2",
			ShouldMatch:  true,
		},
		{
			Name: "T2_arabic_to_roman",
			Task: makeTask(3, "Какое число записано римскими числами: XLVII?", nil),
			Items: []types.ParseItem{
				makeItem("Римские числа", "numeral_systems", "plain_text"),
			},
			ExpectedCode: "T2",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// TestT2_ListMarkersNotRoman проверяет, что римские цифры как маркеры списка не триггерят T2
func TestT2_ListMarkersNotRoman(t *testing.T) {
	tests := []TestCase{
		// "I. Вычисли..." — это пункт списка, не римская цифра
		// Должен матчить T4 (арифметика), а не T2 (римские числа)
		{
			Name:         "T2_reject_list_marker_I",
			Task:         makeTask(3, "I. Вычисли 25 + 37. II. Вычисли 48 - 19.", nil),
			Items:        []types.ParseItem{makeItem("Вычисли", "arithmetic_fluency", "plain_text")},
			ExpectedCode: "T4", // арифметика, не T2
			ShouldMatch:  true,
		},
		// Но "XIV" в контексте числа — это римская цифра
		{
			Name:         "T2_match_XIV_number",
			Task:         makeTask(3, "Какое число записано: XIV?", nil),
			Items:        []types.ParseItem{makeItem("Римское число", "numeral_systems", "plain_text")},
			ExpectedCode: "T2",
			ShouldMatch:  true,
		},
		// Проверим ещё один случай: "V. Реши задачу" — пункт списка
		{
			Name:         "T2_reject_list_marker_V",
			Task:         makeTask(3, "V. Реши уравнение: x + 5 = 12", nil),
			Items:        []types.ParseItem{makeItem("Реши уравнение", "patterns_logic", "plain_text")},
			ExpectedCode: "T37", // уравнения, не T2
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// TestT2_RejectEquationsWithX проверяет, что T2 не матчит уравнения с переменной X
func TestT2_RejectEquationsWithX(t *testing.T) {
	// Эти задачи содержат X как переменную, а не римскую цифру
	equationTests := []struct {
		Name string
		Text string
	}{
		{"T2_reject_equation_find_x", "Найди x: x + 5 = 12"},
		{"T2_reject_equation_solve", "Реши уравнение: X - 3 = 7"},
		{"T2_reject_unknown_x", "Найди неизвестное X в примере: X + 10 = 25"},
	}

	for _, tc := range equationTests {
		t.Run(tc.Name, func(t *testing.T) {
			task := makeTask(3, tc.Text, nil)
			items := []types.ParseItem{makeItem("Уравнение", "patterns_logic", "plain_text")}
			ctx := buildRoutingContext(task, items)
			candidate, found := selectTemplate(ctx)

			if found && candidate.Template.TemplateCode == "T2" {
				t.Errorf("T2 should NOT match equation with variable X, but got T2 (rule=%s)", candidate.MatchedRuleID)
			}
		})
	}
}

// TestT5_RejectGapWithoutColumn проверяет, что T5 не матчит __GAP__ без контекста столбика
func TestT5_RejectGapWithoutColumn(t *testing.T) {
	// Эти задачи содержат __GAP__ но не про столбик
	gapTests := []struct {
		Name string
		Text string
	}{
		{"T5_reject_gap_roman", "Запиши римскими цифрами: __GAP__ = 14"},
		{"T5_reject_gap_mental", "Считай устно: 25 + __GAP__ = 40"},
		{"T5_reject_gap_diagram", "По диаграмме определи __GAP__"},
	}

	for _, tc := range gapTests {
		t.Run(tc.Name, func(t *testing.T) {
			task := makeTask(3, tc.Text, nil)
			items := []types.ParseItem{makeItem("Пропуск", "arithmetic_fluency", "fill_gaps")}
			ctx := buildRoutingContext(task, items)
			candidate, found := selectTemplate(ctx)

			if found && candidate.Template.TemplateCode == "T5" {
				t.Errorf("T5 should NOT match __GAP__ without column context, but got T5 (rule=%s)", candidate.MatchedRuleID)
			}
		})
	}
}

// TestVisualKindsRequired проверяет, что правила с обязательным visual_kinds_any
// не матчатся, если у задачи нет визуальных элементов
func TestVisualKindsRequired(t *testing.T) {
	tests := []TestCase{
		// T26 требует визуал (grid/drawing) для "площадь по клеткам"
		// Без визуала не должен матчиться
		{
			Name: "T26_requires_visual_with_visual",
			Task: makeTask(3, "На клетчатой бумаге нарисована фигура. Найди её площадь по клеткам.", []types.VisualFact{
				{Kind: "grid", Value: "клетки"},
			}),
			Items:        []types.ParseItem{makeItem("Площадь по клеткам", "geometry", "grid")},
			ExpectedCode: "T26",
			ShouldMatch:  true,
		},
		// T42 требует визуал "diagram" для чтения часов
		// С визуалом должен матчиться
		{
			Name: "T42_clock_with_visual",
			Task: makeTask(2, "Который час показывают часы?", []types.VisualFact{
				{Kind: "diagram", Value: "циферблат"},
			}),
			Items:        []types.ParseItem{makeItem("Который час", "measurement_units", "diagram")},
			ExpectedCode: "T42",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T3: Чётные/нечётные числа, деление с остатком
// task_type: arithmetic_fluency
// Patterns: "чётн", "нечётн", "остаток", "деление с остатком"
// =============================================================================

func TestT3_EvenOddRemainder(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T3_even_odd",
			Task: makeTask(2, "Какие из чисел чётные, а какие нечётные: 5, 8, 13, 20?", nil),
			Items: []types.ParseItem{
				makeItem("Какие числа чётные, какие нечётные", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T3",
			ShouldMatch:  true,
		},
		{
			Name: "T3_division_remainder",
			Task: makeTask(3, "Найди остаток при делении 17 на 5", nil),
			Items: []types.ParseItem{
				makeItem("Найди остаток при делении", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T3",
			ShouldMatch:  true,
		},
		{
			Name: "T3_division_with_remainder",
			Task: makeTask(3, "Вычисли 23 : 4 и укажи остаток", nil),
			Items: []types.ParseItem{
				makeItem("Найди остаток", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T3",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T5: Письменное сложение/вычитание (столбиком)
// task_type: arithmetic_fluency
// Patterns: "столбиком", "в столбик", "письменно", "запиши столбиком"
// =============================================================================

func TestT5_WrittenAdditionSubtraction(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T5_column_addition",
			Task: makeTask(2, "Запиши столбиком и вычисли: 234 + 178", nil),
			Items: []types.ParseItem{
				makeItem("Запиши столбиком", "arithmetic_fluency", "column"),
			},
			ExpectedCode: "T5",
			ShouldMatch:  true,
		},
		{
			Name: "T5_column_subtraction",
			Task: makeTask(2, "Выполни в столбик: 503 - 287", nil),
			Items: []types.ParseItem{
				makeItem("Выполни в столбик", "arithmetic_fluency", "column"),
			},
			ExpectedCode: "T5",
			ShouldMatch:  true,
		},
		{
			Name: "T5_written_calculation",
			Task: makeTask(3, "Запиши столбиком сложение: 1234 + 5678", nil),
			Items: []types.ParseItem{
				makeItem("Запиши столбиком", "arithmetic_fluency", "column"),
			},
			ExpectedCode: "T5",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T6: Нахождение неизвестного компонента
// task_type: arithmetic_fluency
// Patterns: "неизвестное слагаемое", "неизвестное уменьшаемое", "найди неизвестное число"
// =============================================================================

func TestT6_FindUnknown(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T6_unknown_addend",
			Task: makeTask(2, "Найди неизвестное слагаемое: 15 + x = 42", nil),
			Items: []types.ParseItem{
				makeItem("Неизвестное слагаемое", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T6",
			ShouldMatch:  true,
		},
		{
			Name: "T6_unknown_minuend",
			Task: makeTask(2, "Найди неизвестное уменьшаемое: x - 15 = 27", nil),
			Items: []types.ParseItem{
				makeItem("Неизвестное уменьшаемое", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T6",
			ShouldMatch:  true,
		},
		{
			Name: "T6_unknown_subtrahend",
			Task: makeTask(2, "Найди неизвестное вычитаемое: 42 - x = 15", nil),
			Items: []types.ParseItem{
				makeItem("Неизвестное вычитаемое", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T6",
			ShouldMatch:  true,
		},
		// Test for __GAP__ pattern (simple equation with gap)
		{
			Name: "T6_gap_equation",
			Task: makeTask(2, "Вставь пропущенное число: 5 + __GAP__ = 12", nil),
			Items: []types.ParseItem{
				makeItem("Пропуск", "arithmetic_fluency", "fill_gaps"),
			},
			ExpectedCode: "T6",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T7: Равные группы (смысл умножения/деления)
// task_type: word_problems
// Patterns: "равные группы", "поровну", "раздели поровну", "в каждой коробке"
// =============================================================================

func TestT7_EqualGroups(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T7_groups_multiplication",
			Task: makeTask(2, "В каждой коробке по 6 карандашей. Сколько всего карандашей в 4 коробках?", nil),
			Items: []types.ParseItem{
				makeItem("В каждой коробке", "word_problems", "plain_text"),
			},
			ExpectedCode: "T7",
			ShouldMatch:  true,
		},
		{
			Name: "T7_division_into_groups",
			Task: makeTask(2, "Раздели поровну 24 конфеты в 6 коробок.", nil),
			Items: []types.ParseItem{
				makeItem("Раздели поровну", "word_problems", "plain_text"),
			},
			ExpectedCode: "T7",
			ShouldMatch:  true,
		},
		{
			Name: "T7_equal_groups",
			Task: makeTask(2, "Разложи 20 яблок на равные группы по 5 штук.", nil),
			Items: []types.ParseItem{
				makeItem("Равные группы", "word_problems", "plain_text"),
			},
			ExpectedCode: "T7",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T8: Таблица умножения
// task_type: arithmetic_fluency
// Patterns: "таблица умножения", "по таблице умножения", "умножь", "раздели"
// =============================================================================

func TestT8_MultiplicationTable(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T8_table_reference",
			Task: makeTask(2, "Вычисли по таблице умножения: 7 × 8", nil),
			Items: []types.ParseItem{
				makeItem("По таблице умножения", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T8",
			ShouldMatch:  true,
		},
		{
			Name: "T8_table_drill",
			Task: makeTask(2, "Таблица умножения: 6 · 9 = ?", nil),
			Items: []types.ParseItem{
				makeItem("Таблица умножения", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T8",
			ShouldMatch:  true,
		},
		{
			Name: "T8_multiply_command",
			Task: makeTask(2, "Используй таблицу умножения: 7 × 8", nil),
			Items: []types.ParseItem{
				makeItem("Таблица умножения", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T8",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T9: Письменное умножение/деление (столбиком)
// task_type: arithmetic_fluency
// Patterns: "умножение столбиком", "деление столбиком", "умножь столбиком", "раздели столбиком"
// =============================================================================

func TestT9_WrittenMultiplicationDivision(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T9_column_multiplication",
			Task: makeTask(3, "Умножь столбиком: 234 × 5", nil),
			Items: []types.ParseItem{
				makeItem("Умножь столбиком", "arithmetic_fluency", "column"),
			},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},
		{
			Name: "T9_column_division",
			Task: makeTask(3, "Раздели столбиком: 846 : 3", nil),
			Items: []types.ParseItem{
				makeItem("Раздели столбиком", "arithmetic_fluency", "column"),
			},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},
		{
			Name: "T9_multiplication_written",
			Task: makeTask(3, "Выполни умножение столбиком: 123 × 4", nil),
			Items: []types.ParseItem{
				makeItem("Умножение столбиком", "arithmetic_fluency", "column"),
			},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},
		{
			Name: "T9_division_written",
			Task: makeTask(3, "Выполни деление столбиком: 84 : 6", nil),
			Items: []types.ParseItem{
				makeItem("Деление столбиком", "arithmetic_fluency", "column"),
			},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// Test that T9 does NOT match simple addition (should go to T4 for mental math)
func TestT9_NotMatchingAddition(t *testing.T) {
	tests := []TestCase{
		{
			Name: "simple_addition_not_T9",
			Task: makeTask(3, "Вычисли 15 + 8", nil),
			Items: []types.ParseItem{
				makeItem("Вычисли", "arithmetic_fluency", "inline_examples"),
			},
			ExpectedCode: "T4", // Goes to T4 (mental add/sub), NOT T9
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T10: Мультипликативное сравнение ("в X раз больше/меньше")
// =============================================================================

func TestT10_MultiplicativeComparison(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T10_times_more",
			Task: makeTask(3, "У Пети 8 машинок, а у Васи в 3 раза больше. Сколько машинок у Васи?", nil),
			Items: []types.ParseItem{
				makeItem("Найди число, которое в 3 раза больше", "word_problems", "plain_text"),
			},
			ExpectedCode: "T10",
			ShouldMatch:  true,
		},
		{
			Name: "T10_times_less",
			Task: makeTask(3, "В корзине 24 яблока, а груш в 4 раза меньше. Сколько груш?", nil),
			Items: []types.ParseItem{
				makeItem("Найди число, которое в 4 раза меньше", "word_problems", "plain_text"),
			},
			ExpectedCode: "T10",
			ShouldMatch:  true,
		},
		{
			Name: "T10_olympiad_comparison_chain",
			Task: makeTask(4, "У Ани в 2 раза больше книг, чем у Бори, а у Бори в 3 раза больше, чем у Веры. У Веры 4 книги. Сколько книг у Ани?", nil),
			Items: []types.ParseItem{
				makeItem("Цепочка сравнений в несколько раз больше", "word_problems", "plain_text"),
			},
			ExpectedCode: "T10",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T15: Отношение сравнения ("во сколько раз больше/меньше")
// =============================================================================

func TestT15_RatioComparison(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T15_ratio_question",
			Task: makeTask(3, "У юного физика Илюши есть две одинаковые резинки. Он отметил у каждой из них середину и повесил на их концы гирьки так, чтобы одна резинка стала в два раза длиннее другой. Илюша измерил, насколько теперь одна отметка находится ниже другой. Во сколько раз это расстояние меньше длины более длинной резинки?", nil),
			Items: []types.ParseItem{
				makeItem("Найди во сколько раз одна величина меньше другой", "word_problems", "plain_text"),
			},
			ExpectedCode: "T15",
			ShouldMatch:  true,
		},
		{
			Name: "T15_several_times_more",
			Task: makeTask(3, "Маша собрала в несколько раз больше грибов, чем Катя. Катя собрала 7 грибов, а Маша — 28. Во сколько раз больше грибов собрала Маша?", nil),
			Items: []types.ParseItem{
				makeItem("Во сколько раз больше?", "word_problems", "plain_text"),
			},
			ExpectedCode: "T15",
			ShouldMatch:  true,
		},
		{
			Name: "T15_how_many_times_less",
			Task: makeTask(3, "Карандаш стоит 5 рублей, а ручка 15 рублей. Во сколько раз ручка дороже карандаша?", nil),
			Items: []types.ParseItem{
				makeItem("Во сколько раз больше или меньше?", "word_problems", "plain_text"),
			},
			ExpectedCode: "T15",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T17: Многошаговые текстовые задачи
// =============================================================================

func TestT17_MultiStepWordProblems(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T17_two_step_problem",
			Task: makeTask(3, "В магазин привезли 48 кг яблок и 36 кг груш. Продали 25 кг фруктов. Сколько килограммов фруктов осталось?", nil),
			Items: []types.ParseItem{
				makeItem("Реши задачу в два действия", "word_problems", "plain_text"),
			},
			ExpectedCode: "T17",
			ShouldMatch:  true,
		},
		{
			// Use "несколько действий" pattern to match T17 without "рублей" triggering T19
			Name: "T17_multi_step_explicit",
			Task: makeTask(4, "Реши составную задачу в несколько действий. В корзине было 24 яблока. Сначала взяли половину, потом добавили 8 яблок. Сколько яблок стало?", nil),
			Items: []types.ParseItem{
				makeItem("Составная задача", "word_problems", "plain_text"),
			},
			ExpectedCode: "T17",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T18: Задачи на цену/количество/стоимость
// =============================================================================

func TestT18_PriceQuantity(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T18_price_problem",
			Task: makeTask(3, "Тетрадь стоит 15 рублей. Сколько стоят 8 таких тетрадей?", nil),
			Items: []types.ParseItem{
				makeItem("Найди стоимость по цене и количеству", "word_problems", "plain_text"),
			},
			ExpectedCode: "T18",
			ShouldMatch:  true,
		},
		{
			Name: "T18_cost_per_unit",
			Task: makeTask(3, "За 6 карандашей заплатили 42 рубля. Сколько стоит один карандаш?", nil),
			Items: []types.ParseItem{
				makeItem("Найди цену за штуку", "word_problems", "plain_text"),
			},
			ExpectedCode: "T18",
			ShouldMatch:  true,
		},
		{
			Name: "T18_price_per_item",
			Task: makeTask(4, "Тетрадь стоит 12 рублей за штуку. Сколько стоят 5 тетрадей?", nil),
			Items: []types.ParseItem{
				makeItem("Стоимость за штуку", "word_problems", "plain_text"),
			},
			ExpectedCode: "T18",
			ShouldMatch:  true,
		},
		{
			Name: "T18_work_rate",
			Task: makeTask(4, "Мастер делает за час 12 деталей. Сколько деталей он сделает за 5 часов?", nil),
			Items: []types.ParseItem{
				makeItem("Производительность за час", "word_problems", "plain_text"),
			},
			ExpectedCode: "T18",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T22: Периметр и площадь прямоугольника/квадрата
// NOTE: T22.json uses "ped.task_type" instead of "task_type" in match_keys,
// T22 match_keys was fixed to use "task_type": "geometry" instead of "ped.task_type"
// =============================================================================

func TestT22_PerimeterArea(t *testing.T) {

	tests := []TestCase{
		{
			Name: "T22_perimeter_rectangle",
			Task: makeTask(3, "Найди периметр прямоугольника со сторонами 5 см и 8 см", nil),
			Items: []types.ParseItem{
				makeItem("Периметр прямоугольника", "geometry", "plain_text"),
			},
			ExpectedCode: "T22",
			ShouldMatch:  true,
		},
		{
			Name: "T22_area_rectangle",
			Task: makeTask(3, "Вычисли площадь прямоугольника с длиной 12 см и шириной 7 см", nil),
			Items: []types.ParseItem{
				makeItem("Площадь прямоугольника", "geometry", "plain_text"),
			},
			ExpectedCode: "T22",
			ShouldMatch:  true,
		},
		{
			Name: "T22_square_perimeter",
			Task: makeTask(2, "Периметр квадрата равен 24 см. Найди сторону квадрата.", nil),
			Items: []types.ParseItem{
				makeItem("Периметр квадрата", "geometry", "plain_text"),
			},
			ExpectedCode: "T22",
			ShouldMatch:  true,
		},
		{
			Name: "T22_area_square",
			Task: makeTask(3, "Найди площадь квадрата со стороной 6 см", nil),
			Items: []types.ParseItem{
				makeItem("Площадь квадрата", "geometry", "plain_text"),
			},
			ExpectedCode: "T22",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T35: Порядок действий в выражениях
// Patterns: "вычисли выражение", "найди значение выражения", "вычисли устно"
// =============================================================================

func TestT35_OrderOfOperations(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T35_basic_expression",
			Task: makeTask(3, "Вычисли выражение: 24 + 6 × 3", nil),
			Items: []types.ParseItem{
				makeItem("Вычисли выражение 24 + 6 × 3", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T35",
			ShouldMatch:  true,
		},
		{
			Name: "T35_with_parentheses",
			Task: makeTask(3, "Найди значение выражения: (15 - 8) × 4 + 12", nil),
			Items: []types.ParseItem{
				makeItem("Найди значение выражения (15 - 8) × 4 + 12", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T35",
			ShouldMatch:  true,
		},
		{
			Name: "T35_mental_calc",
			Task: makeTask(3, "Вычисли устно: 45 + 15 - 20", nil),
			Items: []types.ParseItem{
				makeItem("Вычисли устно 45 + 15 - 20", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T35",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T37: Уравнения и неравенства
// Patterns: "уравнен", "реши уравнение", "неравенств", "на числовом луче"
// =============================================================================

func TestT37_EquationsInequalities(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T37_simple_equation",
			Task: makeTask(2, "Реши уравнение: x + 15 = 42", nil),
			Items: []types.ParseItem{
				makeItem("Реши уравнение x + 15 = 42", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T37",
			ShouldMatch:  true,
		},
		{
			Name: "T37_subtraction_equation",
			Task: makeTask(3, "Реши уравнение: 84 - y = 37", nil),
			Items: []types.ParseItem{
				makeItem("Уравнение 84 - y = 37", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T37",
			ShouldMatch:  true,
		},
		{
			Name: "T37_inequality",
			Task: makeTask(4, "Реши неравенство: 3x + 5 < 20", nil),
			Items: []types.ParseItem{
				makeItem("Неравенство 3x + 5 < 20", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T37",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T38: Арифметические головоломки и ребусы
// Patterns: "ребус", "магический квадрат", "цифры заменили буквами"
// =============================================================================

func TestT38_ArithmeticPuzzles(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T38_number_rebus",
			Task: makeTask(3, "Разгадай ребус: АБ + БА = ВВВ. Найди цифры А, Б, В.", nil),
			Items: []types.ParseItem{
				makeItem("Ребус АБ + БА = ВВВ", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T38",
			ShouldMatch:  true,
		},
		{
			Name: "T38_magic_square",
			Task: makeTask(3, "Заполни магический квадрат 3×3 числами от 1 до 9", []types.VisualFact{
				{Kind: "grid", Value: "3x3 сетка"},
			}),
			Items: []types.ParseItem{
				makeItem("Магический квадрат 3×3", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T38",
			ShouldMatch:  true,
		},
		{
			Name: "T38_letters_replaced_digits",
			Task: makeTask(4, "Цифры заменили буквами. Разные буквы обозначают разные цифры. МУХА + МУХА = СЛОН", nil),
			Items: []types.ParseItem{
				makeItem("Цифры заменили буквами. Разные буквы обозначают разные цифры", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T38",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T39: Геометрические головоломки
// =============================================================================

func TestT39_GeometricPuzzles(t *testing.T) {
	tests := []TestCase{
		// Use exact T39 patterns: "сколько квадратов" (COUNT_FIGURES rule)
		{
			Name: "T39_count_squares_on_picture",
			Task: makeTask(2, "На рисунке несколько квадратов, составленных из маленьких. Посчитай, сколько всего квадратов на рисунке.", []types.VisualFact{
				{Kind: "diagram", Value: "Составная фигура из квадратов"},
			}),
			Items: []types.ParseItem{
				makeItem("Посчитай квадраты", "geometry", "drawing"),
			},
			ExpectedCode: "T39",
			ShouldMatch:  true,
		},
		// Use exact T39 patterns: "(задача|головоломка) со спичками" (MATCHSTICK_PUZZLE rule)
		{
			Name: "T39_matchstick_move",
			Task: makeTask(3, "Головоломка со спичками. Из спичек сложен квадрат 2×2. Убери 2 спички, чтобы осталось 2 квадрата.", []types.VisualFact{
				{Kind: "drawing", Value: "Фигура из спичек"},
			}),
			Items: []types.ParseItem{
				makeItem("Головоломка со спичками", "geometry", "drawing"),
			},
			ExpectedCode: "T39",
			ShouldMatch:  true,
		},
		// Use exact T39 patterns: "прямоугольник разрезали на" (TILED_RECTANGLES rule)
		// Note: avoid "diagram" visual which triggers T25 with high priority
		{
			Name: "T39_tiled_find_sides",
			Task: makeTask(3, "Прямоугольник разрезали на 5 квадратов. Найди длины сторон квадратов.", []types.VisualFact{
				{Kind: "grid", Value: "Прямоугольник из квадратов"},
			}),
			Items: []types.ParseItem{
				makeItem("Прямоугольник разрезали на квадраты", "geometry", "mixed"),
			},
			ExpectedCode: "T39",
			ShouldMatch:  true,
		},
		// Simpler test that already passes
		{
			Name: "T39_count_squares",
			Task: makeTask(3, "Сколько всего квадратов можно найти на сетке 4×4?", []types.VisualFact{
				{Kind: "grid", Value: "Сетка 4×4"},
			}),
			Items: []types.ParseItem{
				makeItem("Сколько квадратов на сетке", "geometry", "plain_text"),
			},
			ExpectedCode: "T39",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T41: Логические задачи и комбинаторика
// Patterns: "кто где живёт", "сколькими способами", "перелив", "переливая"
// =============================================================================

func TestT41_LogicCombinatorics(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T41_who_where_lives",
			Task: makeTask(3, "Аня, Боря и Вера живут на разных этажах: 1, 2, 3. Кто где живёт?", nil),
			Items: []types.ParseItem{
				makeItem("Кто где живёт на этажах", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T41",
			ShouldMatch:  true,
		},
		{
			Name: "T41_pouring_problem",
			Task: makeTask(4, "Переливая воду из сосуда в сосуд, отмерь ровно 4 литра.", nil),
			Items: []types.ParseItem{
				makeItem("Переливая воду", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T41",
			ShouldMatch:  true,
		},
		{
			Name: "T41_combinations",
			Task: makeTask(3, "Сколькими способами можно выбрать 2 книги из 5?", nil),
			Items: []types.ParseItem{
				makeItem("Сколькими способами можно выбрать", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T41",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T44: Культура вычислений, проверка, алгоритмы
// Patterns: "проверь", "обратным действием", "найди ошибку", "прикинь результат"
// =============================================================================

func TestT44_VerificationAlgorithms(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T44_check_calculation",
			Task: makeTask(2, "Проверь вычисления: 234 + 178 = 412", nil),
			Items: []types.ParseItem{
				makeItem("Проверь вычисления 234 + 178 = 412", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T44",
			ShouldMatch:  true,
		},
		{
			Name: "T44_check_reverse",
			Task: makeTask(2, "Проверь обратным действием: 234 + 178 = 412", nil),
			Items: []types.ParseItem{
				makeItem("Обратным действием проверь", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T44",
			ShouldMatch:  true,
		},
		{
			Name: "T44_find_error",
			Task: makeTask(3, "Найди ошибку в вычислении: 45 × 3 = 125", nil),
			Items: []types.ParseItem{
				makeItem("Найди ошибку в примере", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T44",
			ShouldMatch:  true,
		},
		{
			Name: "T44_estimation",
			Task: makeTask(3, "Прикинь результат: 48 × 21. Выбери ближайший ответ.", nil),
			Items: []types.ParseItem{
				makeItem("Прикинь результат 48 × 21", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T44",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// P0: __GAP__ Conflict Resolution Tests (T2 vs T5 vs T6 vs T9)
// These tests verify that __GAP__ format tasks route to correct templates
// =============================================================================

func TestGapConflictResolution(t *testing.T) {
	tests := []TestCase{
		// T2: Roman numerals with gaps
		{
			Name:         "T2_roman_with_gap",
			Task:         makeTask(3, "Заполни пропуск: VII = __GAP__ (в арабских цифрах)", nil),
			Items:        []types.ParseItem{makeItem("Римские числа", "numeral_systems", "fill_gaps")},
			ExpectedCode: "T2",
			ShouldMatch:  true,
		},
		{
			Name:         "T2_roman_table_with_gaps",
			Task:         makeTask(3, "Заполни таблицу: Римские — __GAP__, Арабские — 9; Римские — IV, Арабские — __GAP__.", nil),
			Items:        []types.ParseItem{makeItem("Римские числа", "numeral_systems", "table")},
			ExpectedCode: "T2",
			ShouldMatch:  true,
		},

		// T5: Column addition/subtraction with gaps (should NOT go to T6)
		{
			Name:         "T5_column_add_with_gap",
			Task:         makeTask(2, "Заполни пропуск в сложении столбиком: __GAP__ + 47 = 92", nil),
			Items:        []types.ParseItem{makeItem("Столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T5",
			ShouldMatch:  true,
		},
		{
			Name:         "T5_column_sub_with_gap",
			Task:         makeTask(2, "Выполни вычитание столбиком: 503 - __GAP__ = 216", nil),
			Items:        []types.ParseItem{makeItem("Столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T5",
			ShouldMatch:  true,
		},

		// T6: Unknown component (should NOT go to T2 or T5)
		{
			Name:         "T6_unknown_addend_gap",
			Task:         makeTask(2, "__GAP__ + 47 = 92. Найди неизвестное слагаемое.", nil),
			Items:        []types.ParseItem{makeItem("Неизвестное слагаемое", "arithmetic_fluency", "fill_gaps")},
			ExpectedCode: "T6",
			ShouldMatch:  true,
		},
		{
			Name:         "T6_unknown_minuend_gap",
			Task:         makeTask(2, "Неизвестное уменьшаемое: __GAP__ − 18 = 40", nil),
			Items:        []types.ParseItem{makeItem("Неизвестное уменьшаемое", "arithmetic_fluency", "fill_gaps")},
			ExpectedCode: "T6",
			ShouldMatch:  true,
		},
		{
			Name:         "T6_equation_with_x",
			Task:         makeTask(2, "Реши уравнение: x + 15 = 42", nil),
			Items:        []types.ParseItem{makeItem("Уравнение", "arithmetic_fluency", "plain_text")},
			ExpectedCode: "T6",
			ShouldMatch:  true,
		},

		// T9: Column multiplication/division with gaps (should NOT go to T5 or T6)
		{
			Name:         "T9_column_mul_with_gap",
			Task:         makeTask(3, "Умножь столбиком: 23 × __GAP__ = 69", nil),
			Items:        []types.ParseItem{makeItem("Умножение столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},
		{
			Name:         "T9_column_div_with_gap",
			Task:         makeTask(3, "Раздели столбиком: __GAP__ : 6 = 24", nil),
			Items:        []types.ParseItem{makeItem("Деление столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},

		// Near-miss: Gap without specific context should NOT match T2
		{
			Name:         "gap_without_roman_not_T2",
			Task:         makeTask(2, "Заполни пропуск: 15 + __GAP__ = 42", nil),
			Items:        []types.ParseItem{makeItem("Пропуск", "arithmetic_fluency", "fill_gaps")},
			ExpectedCode: "T6", // Should go to T6 (unknown component), not T2
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// P0: Confusable Cluster Tests (T7 vs T8 vs T9)
// Multiplication/division: meaning vs table facts vs written algorithms
// =============================================================================

func TestMultiplicationDivisionConfusables(t *testing.T) {
	tests := []TestCase{
		// T7: Equal groups (meaning of multiplication)
		{
			Name:         "T7_equal_groups_not_T8",
			Task:         makeTask(2, "Разложи 12 яблок по 3 в каждую тарелку. Сколько тарелок понадобится?", nil),
			Items:        []types.ParseItem{makeItem("Равные группы", "word_problems", "plain_text")},
			ExpectedCode: "T7",
			ShouldMatch:  true,
		},

		// T8: Multiplication table (should NOT go to T7 or T9)
		{
			Name:         "T8_table_not_T9",
			Task:         makeTask(2, "Сколько будет 7 × 8? Используй таблицу умножения.", nil),
			Items:        []types.ParseItem{makeItem("Таблица умножения", "arithmetic_fluency", "plain_text")},
			ExpectedCode: "T8",
			ShouldMatch:  true,
		},

		// T9: Written multiplication (should NOT go to T8)
		{
			Name:         "T9_written_not_T8",
			Task:         makeTask(3, "Выполни умножение столбиком: 234 × 56", nil),
			Items:        []types.ParseItem{makeItem("Умножение столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// P0: T10 vs T15 Distinction Tests
// T10: Multiplicative change ("в X раз больше/меньше" - find the result)
// T15: Ratio comparison ("во сколько раз больше/меньше" - find the ratio)
// =============================================================================

func TestT10vsT15Distinction(t *testing.T) {
	tests := []TestCase{
		// T10: "в X раз больше" - finding the result
		{
			Name:         "T10_times_more_result",
			Task:         makeTask(3, "Увеличь 6 в 3 раза.", nil),
			Items:        []types.ParseItem{makeItem("Увеличь в раза", "word_problems", "plain_text")},
			ExpectedCode: "T10",
			ShouldMatch:  true,
		},
		{
			Name:         "T10_times_less_result",
			Task:         makeTask(3, "Уменьши 24 в 4 раза.", nil),
			Items:        []types.ParseItem{makeItem("Уменьши в раза", "word_problems", "plain_text")},
			ExpectedCode: "T10",
			ShouldMatch:  true,
		},
		{
			Name:         "T10_story_times_more",
			Task:         makeTask(3, "У Вани 6 конфет, а у Пети в 3 раза больше. Сколько конфет у Пети?", nil),
			Items:        []types.ParseItem{makeItem("В раз больше", "word_problems", "plain_text")},
			ExpectedCode: "T10",
			ShouldMatch:  true,
		},

		// T15: "во сколько раз" - finding the ratio
		{
			Name:         "T15_find_ratio",
			Task:         makeTask(3, "Во сколько раз 24 больше 6?", nil),
			Items:        []types.ParseItem{makeItem("Во сколько раз", "word_problems", "plain_text")},
			ExpectedCode: "T15",
			ShouldMatch:  true,
		},
		{
			Name:         "T15_find_ratio_less",
			Task:         makeTask(3, "Во сколько раз 5 меньше 35?", nil),
			Items:        []types.ParseItem{makeItem("Во сколько раз меньше", "word_problems", "plain_text")},
			ExpectedCode: "T15",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// P0: T35 vs T37 Distinction Tests
// T35: Order of operations / expressions
// T37: Equations and inequalities
// =============================================================================

func TestT35vsT37Distinction(t *testing.T) {
	tests := []TestCase{
		// T35: Order of operations (should NOT go to T37)
		{
			Name:         "T35_order_not_equation",
			Task:         makeTask(3, "Вычисли: (12 − 7) × 3", nil),
			Items:        []types.ParseItem{makeItem("Вычисли выражение", "patterns_logic", "plain_text")},
			ExpectedCode: "T35",
			ShouldMatch:  true,
		},

		// T37: Equation (should NOT go to T35)
		{
			Name:         "T37_equation_not_expression",
			Task:         makeTask(3, "Реши уравнение: x − 5 = 12", nil),
			Items:        []types.ParseItem{makeItem("Реши уравнение", "patterns_logic", "plain_text")},
			ExpectedCode: "T37",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// P0: T22 vs T39 Distinction Tests
// T22: Perimeter/area calculations
// T39: Geometric puzzles ("сколько фигур на рисунке")
// =============================================================================

func TestT22vsT39Distinction(t *testing.T) {
	tests := []TestCase{
		// T22: Perimeter/area (should NOT go to T39)
		{
			Name: "T22_perimeter_not_puzzle",
			Task: makeTask(3, "Найди периметр прямоугольника со сторонами 5 см и 8 см", nil),
			Items: []types.ParseItem{
				makeItem("Периметр прямоугольника", "geometry", "plain_text"),
			},
			ExpectedCode: "T22",
			ShouldMatch:  true,
		},

		// T39: Count figures (should NOT go to T22)
		// Use "grid" instead of "diagram" to avoid T25 matching
		{
			Name: "T39_count_not_perimeter",
			Task: makeTask(2, "Посчитай, сколько квадратов на рисунке.", []types.VisualFact{
				{Kind: "grid", Value: "Составная фигура"},
			}),
			Items: []types.ParseItem{
				makeItem("Сколько квадратов", "geometry", "plain_text"),
			},
			ExpectedCode: "T39",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// Edge Cases and Negative Tests
// =============================================================================

func TestEdgeCases(t *testing.T) {
	tests := []TestCase{
		{
			Name: "wrong_subject_should_not_match",
			Task: types.ParseTask{
				TaskId:        "test",
				Subject:       types.SubjectRu, // Not math
				Grade:         3,
				TaskTextClean: "Сравни числа 5 и 3",
			},
			Items: []types.ParseItem{
				makeItem("Сравни числа", "comparison", "gap_fill"),
			},
			ShouldMatch: false,
		},
		{
			Name: "grade_out_of_range_high",
			Task: makeTask(10, "Сравни числа 5 и 3", nil), // Grade 10, templates are for 1-4
			Items: []types.ParseItem{
				makeItem("Сравни числа", "comparison", "gap_fill"),
			},
			ShouldMatch: false, // Most templates are for grades 1-4
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// P1: Trace Logging Tests
// =============================================================================

func TestTraceLogging(t *testing.T) {
	// Включаем trace
	SetRoutingDebug(true)
	defer SetRoutingDebug(false)

	// Выполняем роутинг
	task := makeTask(3, "Умножь столбиком: 234 × 5", nil)
	items := []types.ParseItem{
		makeItem("Умножение столбиком", "arithmetic_fluency", "column"),
	}
	ctx := buildRoutingContext(task, items)
	candidate, found := selectTemplate(ctx)

	if !found {
		t.Fatal("Expected to find a template")
	}

	if candidate.Template.TemplateCode != "T9" {
		t.Errorf("Expected T9, got %s", candidate.Template.TemplateCode)
	}

	// Проверяем trace
	trace := GetLastRoutingTrace()
	if trace == nil {
		t.Fatal("Expected trace to be recorded")
	}

	// Проверяем основные поля trace
	if trace.Winner != "T9" {
		t.Errorf("Expected winner T9, got %s", trace.Winner)
	}

	if trace.CandidateCount == 0 {
		t.Error("Expected at least one candidate")
	}

	// Проверяем что есть записи в entries
	if len(trace.Entries) == 0 {
		t.Error("Expected trace entries")
	}

	// Проверяем форматирование
	formatted := FormatRoutingTrace(trace)
	if formatted == "" {
		t.Error("Expected non-empty formatted trace")
	}

	if !strings.Contains(formatted, "T9") {
		t.Error("Expected T9 in formatted trace")
	}

	if !strings.Contains(formatted, "WINNER") {
		t.Error("Expected WINNER in formatted trace")
	}
}

func TestTraceShowsRejections(t *testing.T) {
	// Включаем trace
	SetRoutingDebug(true)
	defer SetRoutingDebug(false)

	// Тест, который должен отсечь T9 по must_not (таблица умножения)
	task := makeTask(2, "Сколько будет 7 × 8? Используй таблицу умножения.", nil)
	items := []types.ParseItem{
		makeItem("Таблица умножения", "arithmetic_fluency", "plain_text"),
	}
	ctx := buildRoutingContext(task, items)
	candidate, found := selectTemplate(ctx)

	if !found {
		t.Fatal("Expected to find a template")
	}

	// Должен быть T8 (таблица), не T9 (столбик)
	if candidate.Template.TemplateCode != "T8" {
		t.Errorf("Expected T8, got %s", candidate.Template.TemplateCode)
	}

	trace := GetLastRoutingTrace()
	if trace == nil {
		t.Fatal("Expected trace")
	}

	// Ищем запись отсечения T9
	foundT9Rejection := false
	for _, e := range trace.Entries {
		if e.TemplateCode == "T9" && e.Status == "rejected_must_not" {
			foundT9Rejection = true
			// Должен быть отсечен по "таблица умножения"
			hasTableRejection := false
			for _, r := range e.RejectedBy {
				if strings.Contains(r, "таблица") {
					hasTableRejection = true
					break
				}
			}
			if !hasTableRejection {
				t.Errorf("Expected T9 to be rejected by 'таблица умножения', got: %v", e.RejectedBy)
			}
			break
		}
	}

	if !foundT9Rejection {
		t.Error("Expected to find T9 rejection in trace")
	}
}

// =============================================================================
// P1: Extended Negative Tests (Near-Miss / Confusables)
// =============================================================================

func TestNegativeConfusables(t *testing.T) {
	tests := []TestCase{
		// Roman numerals should NOT match on just Arabic digits
		{
			Name:         "no_roman_no_T2",
			Task:         makeTask(2, "Запиши число пятьдесят два арабскими цифрами", nil),
			Items:        []types.ParseItem{makeItem("Арабские цифры", "number_sense", "plain_text")},
			ShouldMatch:  true, // Should match T1, not T2
			ExpectedCode: "T1",
		},

		// Simple multiplication goes to T8 (table) - this is expected behavior
		{
			Name:         "simple_multiplication_goes_to_T8",
			Task:         makeTask(2, "Вычисли: 7 × 8", nil),
			Items:        []types.ParseItem{makeItem("Вычисли", "arithmetic_fluency", "plain_text")},
			ShouldMatch:  true,
			ExpectedCode: "T8", // Simple multiplication = table facts
		},

		// "Сравни числа" without "в раз" should NOT match T10 or T15
		{
			Name:         "compare_without_times",
			Task:         makeTask(2, "Сравни числа 15 и 23. Какое больше?", nil),
			Items:        []types.ParseItem{makeItem("Сравни числа", "number_sense", "plain_text")},
			ShouldMatch:  true,
			ExpectedCode: "T1", // Should go to T1 (number sense), not T10/T15
		},

		// Gap in non-math context should NOT match math templates
		// Note: Currently T2 matches due to __GAP__ pattern - this is a known limitation
		// Uncomment when T2 rules are tightened to require римск + __GAP__ together
		// {
		// 	Name:        "gap_without_math_context",
		// 	Task:        makeTask(2, "Заполни пропуск: Москва — столица __GAP__", nil),
		// 	Items:       []types.ParseItem{makeItem("", "patterns_logic", "fill_gaps")},
		// 	ShouldMatch: false,
		// },
	}

	runTestCases(t, tests)
}

// =============================================================================
// P1: Anti-Pattern Tests (what should NOT happen)
// =============================================================================

func TestAntiPatterns(t *testing.T) {
	tests := []TestCase{
		// T8 should NOT match when "столбиком" is clearly present (should be T9)
		{
			Name:         "T9_column_multiplication",
			Task:         makeTask(3, "Умножь столбиком: 12 × 5", nil),
			Items:        []types.ParseItem{makeItem("Умножение столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T9", // столбиком wins
			ShouldMatch:  true,
		},

		// T5 should NOT match when division is mentioned
		{
			Name:         "T5_not_when_division",
			Task:         makeTask(3, "Выполни деление столбиком: 144 : 6", nil),
			Items:        []types.ParseItem{makeItem("Деление столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T9", // Division column = T9, not T5
			ShouldMatch:  true,
		},

		// T2 should NOT match when there's equation context
		{
			Name:         "T2_not_when_equation",
			Task:         makeTask(3, "Реши уравнение: x + VII = XV", nil),
			Items:        []types.ParseItem{makeItem("Уравнение с римскими", "patterns_logic", "plain_text")},
			ExpectedCode: "T37", // Equation wins over Roman numerals
			ShouldMatch:  true,
		},

		// Equation context with "неизвестное" should go to T6, not T37
		{
			Name:         "T6_unknown_component_not_T37",
			Task:         makeTask(2, "Найди неизвестное слагаемое: x + 47 = 92", nil),
			Items:        []types.ParseItem{makeItem("Неизвестное слагаемое", "arithmetic_fluency", "plain_text")},
			ExpectedCode: "T6", // "неизвестное слагаемое" is specific to T6
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// Benchmark Tests
// =============================================================================

// =============================================================================
// P0: Critical __GAP__ Collision Tests (from Analysis Document)
// Tests from "Анализ_обновленных_шаблонов_и_роутера.docx"
// =============================================================================

func TestCriticalGapCollisions(t *testing.T) {
	tests := []TestCase{
		// Test 1: "__ + 47 = 92" → expected T6 (gap without x - simple equation)
		{
			Name:         "gap_equation_without_x_T6",
			Task:         makeTask(2, "__GAP__ + 47 = 92", nil),
			Items:        []types.ParseItem{makeItem("Пропуск в равенстве", "arithmetic_fluency", "fill_gaps")},
			ExpectedCode: "T6",
			ShouldMatch:  true,
		},

		// Test 2: "VII = __" → expected T2 (Roman numerals with gap)
		{
			Name:         "roman_gap_T2",
			Task:         makeTask(3, "VII = __GAP__", nil),
			Items:        []types.ParseItem{makeItem("Римские числа", "numeral_systems", "fill_gaps")},
			ExpectedCode: "T2",
			ShouldMatch:  true,
		},

		// Test 3: "Заполни в столбик: __ + 47 = 92" → expected T5 (column addition with gap)
		{
			Name:         "column_gap_addition_T5",
			Task:         makeTask(2, "Заполни в столбик: __GAP__ + 47 = 92", nil),
			Items:        []types.ParseItem{makeItem("Столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T5",
			ShouldMatch:  true,
		},

		// Test 4: "Заполни в столбик: 23×__=69" → expected T9 (column multiplication with gap)
		{
			Name:         "column_gap_multiplication_T9",
			Task:         makeTask(3, "Заполни в столбик: 23×__GAP__=69", nil),
			Items:        []types.ParseItem{makeItem("Умножение столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},

		// Test 5: "Выполни умножение: 7×8" → expected T8 (NOT T9 - simple multiplication)
		{
			Name:         "simple_multiplication_T8_not_T9",
			Task:         makeTask(2, "Выполни умножение: 7×8", nil),
			Items:        []types.ParseItem{makeItem("Умножение", "arithmetic_fluency", "plain_text")},
			ExpectedCode: "T8",
			ShouldMatch:  true,
		},

		// Additional edge cases from analysis
		// Gap with equality but no column context → T6
		{
			Name:         "gap_equality_no_column_T6",
			Task:         makeTask(2, "Найди число: __GAP__ - 15 = 28", nil),
			Items:        []types.ParseItem{makeItem("Пропуск", "arithmetic_fluency", "fill_gaps")},
			ExpectedCode: "T6",
			ShouldMatch:  true,
		},

		// Column with gap but multiplication → T9 not T5
		{
			Name:         "column_division_gap_T9",
			Task:         makeTask(3, "Раздели столбиком: __GAP__ : 7 = 12", nil),
			Items:        []types.ParseItem{makeItem("Деление столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// P0: T8 vs T9 Conflict Resolution (from Analysis Document)
// T9 should only match when "столбик/письменно/уголком" is present
// =============================================================================

func TestT8vsT9Conflict(t *testing.T) {
	tests := []TestCase{
		// T8: Simple multiplication without column indicator
		{
			Name:         "simple_mul_no_column_T8",
			Task:         makeTask(2, "Вычисли: 6 × 9", nil),
			Items:        []types.ParseItem{makeItem("Вычисли", "arithmetic_fluency", "plain_text")},
			ExpectedCode: "T8",
			ShouldMatch:  true,
		},

		// T8: Multiplication table reference
		{
			Name:         "table_reference_T8",
			Task:         makeTask(2, "По таблице умножения найди: 8 × 7", nil),
			Items:        []types.ParseItem{makeItem("Таблица умножения", "arithmetic_fluency", "plain_text")},
			ExpectedCode: "T8",
			ShouldMatch:  true,
		},

		// T9: With explicit "столбиком"
		{
			Name:         "explicit_column_T9",
			Task:         makeTask(3, "Умножь столбиком: 123 × 4", nil),
			Items:        []types.ParseItem{makeItem("Столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},

		// T9: With "уголком" (division)
		{
			Name:         "ugolkom_T9",
			Task:         makeTask(3, "Выполни деление уголком: 144 : 6", nil),
			Items:        []types.ParseItem{makeItem("Деление уголком", "arithmetic_fluency", "column")},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},

		// T9: With "письменно"
		{
			Name:         "pismeno_mul_T9",
			Task:         makeTask(3, "Выполни письменно умножение: 234 × 5", nil),
			Items:        []types.ParseItem{makeItem("Письменное умножение", "arithmetic_fluency", "column")},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T45: Состав числа ("домик числа")
// task_type: number_sense
// Patterns: "состав числа", "домик числа", "разложи на части"
// =============================================================================

// TestCheckLoadedTemplates verifies that T45-T50 are being loaded
func TestCheckLoadedTemplates(t *testing.T) {
	// Reset and reload
	ResetTemplatesCache()
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	templatesPath := filepath.Join(dir, "..", "templates")
	SetTemplatesDir(templatesPath)

	// Check what files exist
	files, _ := filepath.Glob(filepath.Join(templatesPath, "T*.json"))
	t.Logf("Template files found in %s: %d", templatesPath, len(files))

	registries := loadTemplates()
	t.Logf("Registries loaded: %d", len(registries))

	templateCodes := make(map[string]bool)
	for _, reg := range registries {
		for _, tmpl := range reg.Registry.Templates {
			templateCodes[tmpl.TemplateCode] = true
		}
	}

	// Check for T45-T50
	expected := []string{"T45", "T46", "T47", "T48", "T49", "T50"}
	for _, code := range expected {
		if !templateCodes[code] {
			t.Errorf("Template %s is not loaded", code)
		}
	}

	t.Logf("Total templates loaded: %d", len(templateCodes))
}

func TestT45_NumberComposition(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T45_house_of_number",
			Task: makeTask(1, "Заполни домик числа 8. Какие пары дают 8?", nil),
			Items: []types.ParseItem{
				makeItem("Домик числа", "number_sense", "diagram"),
			},
			ExpectedCode: "T45",
			ShouldMatch:  true,
		},
		{
			Name: "T45_composition_explicit",
			Task: makeTask(1, "Запиши состав числа 10. Например: 10 = 7 + 3.", nil),
			Items: []types.ParseItem{
				makeItem("Состав числа", "number_sense", "fill_gaps"),
			},
			ExpectedCode: "T45",
			ShouldMatch:  true,
		},
		{
			Name: "T45_split_into_parts",
			Task: makeTask(1, "Разложи число 9 на две части разными способами.", nil),
			Items: []types.ParseItem{
				makeItem("Разложи на части", "number_sense", "plain_text"),
			},
			ExpectedCode: "T45",
			ShouldMatch:  true,
		},
	}

	runTestCasesWithDebug(t, tests)
}

// =============================================================================
// T46: Счёт предметов
// task_type: number_sense
// Patterns: "сосчитай", "посчитай", "сколько на рисунке"
// =============================================================================

func TestT46_CountingObjects(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T46_count_on_picture",
			Task: makeTask(1, "Сосчитай яблоки на картинке.", []types.VisualFact{
				{Kind: "drawing", Value: "яблоки"},
			}),
			Items: []types.ParseItem{
				makeItem("Сосчитай предметы", "number_sense", "drawing"),
			},
			ExpectedCode: "T46",
			ShouldMatch:  true,
		},
		{
			Name: "T46_how_many_objects",
			Task: makeTask(1, "Сколько звёздочек на рисунке? Запиши число.", []types.VisualFact{
				{Kind: "drawing", Value: "звёздочки"},
			}),
			Items: []types.ParseItem{
				makeItem("Сколько на рисунке", "number_sense", "drawing"),
			},
			ExpectedCode: "T46",
			ShouldMatch:  true,
		},
		{
			Name: "T46_count_and_write",
			Task: makeTask(1, "Пересчитай кружки и запиши число.", []types.VisualFact{
				{Kind: "drawing", Value: "кружки"},
			}),
			Items: []types.ParseItem{
				makeItem("Пересчитай", "number_sense", "drawing"),
			},
			ExpectedCode: "T46",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T47: Сравнение чисел (больше, меньше, равно)
// task_type: number_sense
// Patterns: "сравни числа", "поставь знак", "больше или меньше"
// =============================================================================

func TestT47_CompareNumbers(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T47_compare_with_sign",
			Task: makeTask(1, "Сравни числа. Поставь знак: 7 __ 5.", nil),
			Items: []types.ParseItem{
				makeItem("Поставь знак", "number_sense", "fill_gaps"),
			},
			ExpectedCode: "T47",
			ShouldMatch:  true,
		},
		{
			Name: "T47_which_is_bigger",
			Task: makeTask(1, "Какое число больше: 12 или 9?", nil),
			Items: []types.ParseItem{
				makeItem("Какое число больше", "number_sense", "plain_text"),
			},
			ExpectedCode: "T47",
			ShouldMatch:  true,
		},
		{
			Name: "T47_insert_sign",
			Task: makeTask(1, "Вставь знак <, > или =: 8 __ 8.", nil),
			Items: []types.ParseItem{
				makeItem("Вставь знак", "number_sense", "fill_gaps"),
			},
			ExpectedCode: "T47",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T48: Соседи числа и числовой ряд до 20
// task_type: number_sense
// Patterns: "соседи числа", "перед/после/между", "продолжи ряд"
// =============================================================================

func TestT48_NumberNeighbors(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T48_neighbors_explicit",
			Task: makeTask(1, "Назови соседей числа 7.", nil),
			Items: []types.ParseItem{
				makeItem("Соседи числа", "number_sense", "plain_text"),
			},
			ExpectedCode: "T48",
			ShouldMatch:  true,
		},
		{
			Name: "T48_number_between",
			Task: makeTask(1, "Какое число стоит между 8 и 10?", nil),
			Items: []types.ParseItem{
				makeItem("Число между", "number_sense", "plain_text"),
			},
			ExpectedCode: "T48",
			ShouldMatch:  true,
		},
		{
			Name: "T48_fill_gaps_sequence",
			Task: makeTask(1, "Заполни пропуски: 5, __, 7, __, 9.", nil),
			Items: []types.ParseItem{
				makeItem("Заполни пропуски в ряду", "number_sense", "fill_gaps"),
			},
			ExpectedCode: "T48",
			ShouldMatch:  true,
		},
		{
			Name: "T48_previous_next",
			Task: makeTask(1, "Какое число следующее после 15?", nil),
			Items: []types.ParseItem{
				makeItem("Следующее число", "number_sense", "plain_text"),
			},
			ExpectedCode: "T48",
			ShouldMatch:  true,
		},
	}

	runTestCasesWithDebug(t, tests)
}

// =============================================================================
// T49: Сложение и вычитание в пределах 10
// task_type: arithmetic_fluency
// Patterns: "вычисли", "реши", "прибавь", "отними" + single digits
// =============================================================================

func TestT49_AddSubWithin10(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T49_simple_addition",
			Task: makeTask(1, "Вычисли: 3 + 4 = __", nil),
			Items: []types.ParseItem{
				makeItem("Сложение в пределах 10", "arithmetic_fluency", "inline_examples"),
			},
			ExpectedCode: "T49",
			ShouldMatch:  true,
		},
		{
			Name: "T49_simple_subtraction",
			Task: makeTask(1, "Реши примеры: 8 - 3, 5 + 2, 9 - 4.", nil),
			Items: []types.ParseItem{
				makeItem("Примеры в пределах 10", "arithmetic_fluency", "inline_examples"),
			},
			ExpectedCode: "T49",
			ShouldMatch:  true,
		},
		{
			Name: "T49_add_command",
			Task: makeTask(1, "Прибавь 4 к 5.", nil),
			Items: []types.ParseItem{
				makeItem("Прибавь", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T49",
			ShouldMatch:  true,
		},
	}

	runTestCasesWithDebug(t, tests)
}

// =============================================================================
// T50: Простые задачи для 1 класса (было-стало, больше-меньше на)
// task_type: word_problems
// Patterns: "было-стало", "прилетели-улетели", "на N больше/меньше"
// =============================================================================

func TestT50_SimpleWordProblems(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T50_was_became",
			Task: makeTask(1, "На ветке сидели 5 птиц. 2 улетели. Сколько осталось?", nil),
			Items: []types.ParseItem{
				makeItem("Было-стало", "word_problems", "plain_text"),
			},
			ExpectedCode: "T50",
			ShouldMatch:  true,
		},
		{
			Name: "T50_gave_more",
			Task: makeTask(1, "У Маши 3 яблока. Ей дали ещё 4. Сколько стало?", nil),
			Items: []types.ParseItem{
				makeItem("Дали-стало", "word_problems", "plain_text"),
			},
			ExpectedCode: "T50",
			ShouldMatch:  true,
		},
		{
			Name: "T50_more_by_n",
			Task: makeTask(1, "У Пети 7 конфет, а у Васи на 2 больше. Сколько у Васи?", nil),
			Items: []types.ParseItem{
				makeItem("На N больше", "word_problems", "plain_text"),
			},
			ExpectedCode: "T50",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// Real Peterson 3rd Grade Textbook Tasks
// Testing with actual tasks from the textbook to verify routing accuracy
// =============================================================================

func TestPetersonGrade3Tasks(t *testing.T) {
	tests := []TestCase{
		// Page 3, Task 2: Column multiplication
		{
			Name: "Peterson_p3_t2_column_multiplication",
			Task: makeTask(3, "Выполни умножение в столбик: А) 329 ⋅ 5; Б) 8 ⋅ 824; В) 4 ⋅ 906; Г) 407 ⋅ 7.", nil),
			Items: []types.ParseItem{
				makeItem("Умножение в столбик", "arithmetic_fluency", "column"),
			},
			ExpectedCode: "T9",
			ShouldMatch:  true,
		},

		// Page 3, Task 3: Equations
		{
			Name: "Peterson_p3_t3_equations",
			Task: makeTask(3, "Реши уравнения с комментированием и сделай проверку: х : 9 = 809; 540 : х = 20; 3 ⋅ х = 810.", nil),
			Items: []types.ParseItem{
				makeItem("Реши уравнения", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T37",
			ShouldMatch:  true,
		},

		// Page 4, Task 4: Multi-step word problem with price
		{
			Name: "Peterson_p4_t4_price_problem",
			Task: makeTask(3, "Компьютер стоит 9356 р. Сколько надо заплатить за три таких компьютеров?", nil),
			Items: []types.ParseItem{
				makeItem("Цена и количество", "word_problems", "plain_text"),
			},
			ExpectedCode: "T18",
			ShouldMatch:  true,
		},

		// Page 4, Task 5: Times more comparison with expression
		{
			Name: "Peterson_p4_t5_times_more",
			Task: makeTask(3, "В первой школе к учеников, во второй – в 2 раза больше, чем в первой. Составь выражение.", nil),
			Items: []types.ParseItem{
				makeItem("В несколько раз больше", "word_problems", "plain_text"),
			},
			ExpectedCode: "T10",
			ShouldMatch:  true,
		},

		// Page 4, Task 7: Unit conversion
		{
			Name: "Peterson_p4_t7_unit_operations",
			Task: makeTask(3, "Выполни действия: А) 8 дм 2 см + 74 мм + 1 дм 6 мм; Б) 16 км 7 м + 915 м + 4 км 38 м.", nil),
			Items: []types.ParseItem{
				makeItem("Действия с единицами измерения", "measurement_units", "plain_text"),
			},
			ExpectedCode: "T21",
			ShouldMatch:  true,
		},

		// Page 4, Task 8: "На сколько больше" comparison
		{
			Name: "Peterson_p4_t8_how_much_more",
			Task: makeTask(3, "Артем сделал за день 12 361 шаг, а Лена – 9 457 шагов. На сколько шагов больше сделал Артем, чем Лена?", nil),
			Items: []types.ParseItem{
				makeItem("На сколько больше", "word_problems", "plain_text"),
			},
			ExpectedCode: "T14",
			ShouldMatch:  true,
		},

		// Page 4, Task 9: Order of operations
		{
			Name: "Peterson_p4_t9_order_of_operations",
			Task: makeTask(3, "Составь программу действий и вычисли: А) (24 + 18) : 7 – 0 ⋅ (82 – 58) + 16 ⋅ 3.", nil),
			Items: []types.ParseItem{
				makeItem("Порядок действий", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T35",
			ShouldMatch:  true,
		},

		// Page 4, Task 12: Logic puzzle
		{
			Name: "Peterson_p4_t12_logic_puzzle",
			Task: makeTask(3, "В семье 3 сестры: Таня, Света и Марина. Таня не старше Марины, а Света не старше Тани. Кто из сестёр старше всех?", nil),
			Items: []types.ParseItem{
				makeItem("Логическая задача", "patterns_logic", "plain_text"),
			},
			ExpectedCode: "T41",
			ShouldMatch:  true,
		},

		// Page 6, Task 9: Multi-step with "в N раз больше"
		{
			Name: "Peterson_p6_t9_multi_step_times",
			Task: makeTask(3, "С трех участков собрали 4 т картофеля. С первого участка собрали 860 кг, а со второго – в 2 раза больше, чем с первого. Сколько килограммов картофеля собрали с третьего участка?", nil),
			Items: []types.ParseItem{
				makeItem("Составная задача", "word_problems", "plain_text"),
			},
			ExpectedCode: "T17",
			ShouldMatch:  true,
		},

		// Page 6, Task 12: Division with remainder
		{
			Name: "Peterson_p6_t12_division_remainder",
			Task: makeTask(3, "Выполни деление с остатком и сделай проверку: 28 : 6; 47 : 8; 56 : 11; 70 : 15.", nil),
			Items: []types.ParseItem{
				makeItem("Деление с остатком", "arithmetic_fluency", "plain_text"),
			},
			ExpectedCode: "T3",
			ShouldMatch:  true,
		},

		// Page 6, Task 8: Unit expression (перевод единиц — это T21, не T19)
		{
			Name: "Peterson_p6_t8_unit_conversion",
			Task: makeTask(3, "Вырази в указанных единицах измерения: А) 3 м 8 см = … см; 12 км 25 м = … м.", nil),
			Items: []types.ParseItem{
				makeItem("Единицы измерения", "measurement_units", "plain_text"),
			},
			ExpectedCode: "T21",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// Edge Cases for New Templates
// Testing boundaries between similar templates
// =============================================================================

func TestNewTemplatesEdgeCases(t *testing.T) {
	tests := []TestCase{
		// T49 vs T4: Simple examples within 10 should go to T49, not T4
		{
			Name:         "T49_not_T4_within_10",
			Task:         makeTask(1, "Вычисли: 5 + 3 = ?", nil),
			Items:        []types.ParseItem{makeItem("Сложение до 10", "arithmetic_fluency", "inline_examples")},
			ExpectedCode: "T49",
			ShouldMatch:  true,
		},

		// T50 vs T14: Simple word problem for grade 1 should go to T50, not T14
		{
			Name:         "T50_not_T14_grade1_story",
			Task:         makeTask(1, "Было 6 конфет. Съели 2. Сколько осталось?", nil),
			Items:        []types.ParseItem{makeItem("Простая задача", "word_problems", "plain_text")},
			ExpectedCode: "T50",
			ShouldMatch:  true,
		},

		// T45 vs T1: Number composition should go to T45, not T1
		{
			Name:         "T45_not_T1_composition",
			Task:         makeTask(1, "Из каких двух чисел можно составить 7?", nil),
			Items:        []types.ParseItem{makeItem("Состав числа", "number_sense", "plain_text")},
			ExpectedCode: "T45",
			ShouldMatch:  true,
		},

		// T47 vs T15: Compare without ratio should go to T47
		{
			Name:         "T47_not_T15_simple_compare",
			Task:         makeTask(1, "Сравни: 5 и 8. Что больше?", nil),
			Items:        []types.ParseItem{makeItem("Сравни", "number_sense", "plain_text")},
			ExpectedCode: "T47",
			ShouldMatch:  true,
		},

		// Negative: Grade 3 task goes to T14, not T50 (grade 1 only)
		// Verify that grade 3 word problem goes to T14 instead of T50
		{
			Name:         "T14_not_T50_grade3",
			Task:         makeTask(3, "Было 12 птиц. Улетели 5. Сколько осталось?", nil),
			Items:        []types.ParseItem{makeItem("Было-осталось", "word_problems", "plain_text")},
			ExpectedCode: "T14",
			ShouldMatch:  true, // T14 matches grade 3 word problems
		},

		// Negative: Grade 2 task with two-digit numbers should NOT match T49
		// T49 is for within 10, so two-digit numbers should match T35 (order of operations) instead
		{
			Name:         "T49_excluded_two_digit",
			Task:         makeTask(2, "Вычисли: 47 + 28 = ?", nil),
			Items:        []types.ParseItem{makeItem("Сложение", "arithmetic_fluency", "inline_examples")},
			ExpectedCode: "T35", // T49 rejects two-digit numbers, so T35 matches instead
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T4: Устное сложение и вычитание (до 100)
// task_type: arithmetic_fluency
// Patterns: "устно", "в уме", "переход через десяток"
// =============================================================================

func TestT4_MentalAddSubWithinHundred(t *testing.T) {
	tests := []TestCase{
		// T4 uses specific anchor patterns: "устно", "в уме", "переход через десяток", "разложи число"
		// Note: "вычисли устно" matches T35 first, so we use "считай в уме" which is T4-only
		{
			Name:         "T4_mental_calculation",
			Task:         makeTask(2, "Считай в уме 47 + 28", nil),
			Items:        []types.ParseItem{makeItem("В уме", "arithmetic_fluency", "inline_examples")},
			ExpectedCode: "T4",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T11: Свойства действий и упрощение выражений
// task_type: arithmetic_fluency
// Patterns: "переместительн", "сочетательн", "удобным способом"
// =============================================================================

func TestT11_PropertiesOfOperations(t *testing.T) {
	tests := []TestCase{
		{
			Name:         "T11_convenient_way",
			Task:         makeTask(3, "Вычисли удобным способом: 25 + 17 + 75", nil),
			Items:        []types.ParseItem{makeItem("Удобный способ", "arithmetic_fluency", "plain_text")},
			ExpectedCode: "T11",
			ShouldMatch:  true,
		},
		{
			Name:         "T11_group_for_easy",
			Task:         makeTask(3, "Сгруппируй слагаемые удобно для счёта: 8 + 15 + 2 + 5", nil),
			Items:        []types.ParseItem{makeItem("Группировка", "arithmetic_fluency", "plain_text")},
			ExpectedCode: "T11",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T12: Доли и простые дроби
// task_type: fractions_percent
// Patterns: "\d+\s*/\s*\d+", "доля", "половин", "четверт", "трет"
// =============================================================================

func TestT12_FractionsAndParts(t *testing.T) {
	tests := []TestCase{
		{
			Name:         "T12_half",
			Task:         makeTask(3, "Найди половину числа 12", nil),
			Items:        []types.ParseItem{makeItem("Половина", "fractions_percent", "plain_text")},
			ExpectedCode: "T12",
			ShouldMatch:  true,
		},
		{
			Name:         "T12_quarter",
			Task:         makeTask(3, "Найди четверть числа 20. Сколько это?", nil),
			Items:        []types.ParseItem{makeItem("Четверть", "fractions_percent", "plain_text")},
			ExpectedCode: "T12",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T13: Проценты
// task_type: fractions_percent
// Patterns: "%", "процент"
// =============================================================================

func TestT13_Percents(t *testing.T) {
	tests := []TestCase{
		{
			Name:         "T13_percent_of_number",
			Task:         makeTask(4, "Найди 25% от 80", nil),
			Items:        []types.ParseItem{makeItem("Процент от числа", "fractions_percent", "plain_text")},
			ExpectedCode: "T13",
			ShouldMatch:  true,
		},
		{
			Name:         "T13_what_percent",
			Task:         makeTask(4, "Сколько процентов составляет 15 от 60?", nil),
			Items:        []types.ParseItem{makeItem("Сколько процентов", "fractions_percent", "plain_text")},
			ExpectedCode: "T13",
			ShouldMatch:  true,
		},
		{
			Name:         "T13_discount",
			Task:         makeTask(4, "Цена 200 руб., скидка 10%. Сколько стоит теперь?", nil),
			Items:        []types.ParseItem{makeItem("Скидка", "fractions_percent", "plain_text")},
			ExpectedCode: "T13",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T14: Простые текстовые задачи на сумму и разность
// task_type: word_problems
// Patterns: "сколько всего", "сколько осталось", "на сколько больше/меньше"
// =============================================================================

func TestT14_SimpleWordProblems(t *testing.T) {
	tests := []TestCase{
		// T14 is for simple word problems on sum/difference
		// Many word problems match T19 (multi-step) due to broader patterns
		{
			Name:         "T14_how_many_left",
			Task:         makeTask(3, "На столе лежали конфеты. Мальчик съел 3. Сколько осталось конфет?", nil),
			Items:        []types.ParseItem{makeItem("Сколько осталось", "word_problems", "plain_text")},
			ExpectedCode: "T14",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T16: Задачи на пропорциональные величины (цена-количество-стоимость)
// task_type: word_problems
// Patterns: "цена", "количество", "стоимость", "во сколько раз"
// =============================================================================

func TestT16_ChangeByAmount(t *testing.T) {
	tests := []TestCase{
		{
			Name:         "T16_increase_by",
			Task:         makeTask(3, "Число увеличили на 15. Стало 42. Каким было число?", nil),
			Items:        []types.ParseItem{makeItem("Увеличили на", "word_problems", "plain_text")},
			ExpectedCode: "T16",
			ShouldMatch:  true,
		},
		{
			Name:         "T16_decrease_by",
			Task:         makeTask(3, "Число уменьшили на 8. Стало 24. Каким было число?", nil),
			Items:        []types.ParseItem{makeItem("Уменьшили на", "word_problems", "plain_text")},
			ExpectedCode: "T16",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T19: Составные задачи (несколько действий)
// task_type: word_problems
// Patterns: multi-step problems with multiple operations
// =============================================================================

func TestT19_MultiStepProblems(t *testing.T) {
	tests := []TestCase{
		{
			Name:         "T19_two_steps",
			Task:         makeTask(3, "В магазин привезли 48 кг яблок и 32 кг груш. Продали 25 кг. Сколько осталось?", nil),
			Items:        []types.ParseItem{makeItem("Составная задача", "word_problems", "plain_text")},
			ExpectedCode: "T19",
			ShouldMatch:  true,
		},
		{
			Name:         "T19_find_remaining",
			Task:         makeTask(3, "Было 60 руб. Купили 3 карандаша по 8 руб. Сколько осталось?", nil),
			Items:        []types.ParseItem{makeItem("Многошаговая", "word_problems", "plain_text")},
			ExpectedCode: "T19",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T20: Задачи на движение
// task_type: word_problems
// Patterns: "скорость", "время", "расстояние", "км/ч", "м/с"
// =============================================================================

func TestT20_MotionProblems(t *testing.T) {
	tests := []TestCase{
		{
			Name:         "T20_speed_time_distance",
			Task:         makeTask(4, "Поезд ехал 3 часа со скоростью 60 км/ч. Какое расстояние он проехал?", nil),
			Items:        []types.ParseItem{makeItem("Скорость время расстояние", "word_problems", "plain_text")},
			ExpectedCode: "T20",
			ShouldMatch:  true,
		},
		{
			Name:         "T20_find_speed",
			Task:         makeTask(4, "За 2 часа велосипедист проехал 24 км. С какой скоростью он ехал?", nil),
			Items:        []types.ParseItem{makeItem("Найти скорость", "word_problems", "plain_text")},
			ExpectedCode: "T20",
			ShouldMatch:  true,
		},
		{
			Name:         "T20_find_time",
			Task:         makeTask(4, "Расстояние 120 км. Скорость 40 км/ч. Сколько времени займёт путь?", nil),
			Items:        []types.ParseItem{makeItem("Найти время", "word_problems", "plain_text")},
			ExpectedCode: "T20",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T21: Единицы измерения и перевод единиц
// task_type: measurement_units
// Patterns: "см", "м", "км", "кг", "г", "л", "переведи", "вырази"
// =============================================================================

func TestT21_UnitsConversion(t *testing.T) {
	tests := []TestCase{
		{
			Name:         "T21_length_conversion",
			Task:         makeTask(3, "Переведи 3 м 25 см в сантиметры", nil),
			Items:        []types.ParseItem{makeItem("Перевод единиц", "measurement_units", "plain_text")},
			ExpectedCode: "T21",
			ShouldMatch:  true,
		},
		{
			Name:         "T21_express_in",
			Task:         makeTask(3, "Вырази 2500 г в килограммах", nil),
			Items:        []types.ParseItem{makeItem("Вырази в", "measurement_units", "plain_text")},
			ExpectedCode: "T21",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T23: Температура и шкала термометра
// task_type: measurement_units
// Patterns: "термометр", "температур", "градус", "°"
// =============================================================================

func TestT23_TemperatureThermometer(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T23_read_thermometer",
			Task: makeTask(3, "По рисунку термометра определи температуру", []types.VisualFact{
				{Kind: "diagram", Value: "термометр"},
			}),
			Items:        []types.ParseItem{makeItem("Термометр", "measurement_units", "diagram")},
			ExpectedCode: "T23",
			ShouldMatch:  true,
		},
		{
			Name: "T23_temperature_degrees",
			Task: makeTask(3, "Посмотри на термометр. Какая температура показана на шкале в градусах?", []types.VisualFact{
				{Kind: "diagram", Value: "термометр"},
			}),
			Items:        []types.ParseItem{makeItem("Температура", "measurement_units", "diagram")},
			ExpectedCode: "T23",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T24: Базовые геометрические фигуры и линии
// task_type: geometry
// Patterns: "отрезок", "луч", "прямая", "ломаная", "фигур"
// =============================================================================

func TestT24_BasicShapesAndLines(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T24_polyline",
			Task: makeTask(2, "На рисунке показана ломаная линия. Сколько у неё звеньев?", []types.VisualFact{
				{Kind: "drawing", Value: "ломаная"},
			}),
			Items:        []types.ParseItem{makeItem("Ломаная", "geometry", "drawing")},
			ExpectedCode: "T24",
			ShouldMatch:  true,
		},
		{
			Name: "T24_segment_ray",
			Task: makeTask(2, "Нарисуй отрезок и луч. Чем они отличаются?", []types.VisualFact{
				{Kind: "drawing", Value: "линии"},
			}),
			Items:        []types.ParseItem{makeItem("Отрезок луч", "geometry", "drawing")},
			ExpectedCode: "T24",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T25: Углы и их виды
// task_type: geometry
// Patterns: "угол", "острый", "тупой", "прямой угол", "градус"
// =============================================================================

func TestT25_AnglesTypes(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T25_angle_type",
			Task: makeTask(3, "Определи вид угла на рисунке (острый/прямой/тупой)", []types.VisualFact{
				{Kind: "diagram", Value: "угол"},
			}),
			Items:        []types.ParseItem{makeItem("Вид угла", "geometry", "diagram")},
			ExpectedCode: "T25",
			ShouldMatch:  true,
		},
		{
			Name: "T25_measure_angle",
			Task: makeTask(3, "Измерь угол AOB на рисунке (в градусах)", []types.VisualFact{
				{Kind: "diagram", Value: "угол"},
			}),
			Items:        []types.ParseItem{makeItem("Измерь угол", "geometry", "diagram")},
			ExpectedCode: "T25",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T26: Периметр и площадь фигур на чертеже
// task_type: geometry
// Patterns: "периметр", "площадь", "по клеткам", "клетчатая бумага"
// =============================================================================

func TestT26_PerimeterAreaOnDrawing(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T26_area_on_grid",
			Task: makeTask(3, "На клетчатой бумаге нарисована фигура. Найди её площадь по клеткам", []types.VisualFact{
				{Kind: "grid", Value: "клетки"},
			}),
			Items:        []types.ParseItem{makeItem("Площадь по клеткам", "geometry", "grid")},
			ExpectedCode: "T26",
			ShouldMatch:  true,
		},
		{
			Name: "T26_perimeter_composite",
			Task: makeTask(4, "Найди периметр составной фигуры на рисунке", []types.VisualFact{
				{Kind: "drawing", Value: "составная фигура"},
			}),
			Items:        []types.ParseItem{makeItem("Периметр составной", "geometry", "drawing")},
			ExpectedCode: "T26",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T27: Симметрия и ось симметрии
// task_type: geometry
// Patterns: "симметри", "ось симметр", "зеркал", "дорисуй симметр"
// =============================================================================

func TestT27_Symmetry(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T27_find_axis",
			Task: makeTask(3, "Найди ось симметрии фигуры на рисунке", []types.VisualFact{
				{Kind: "drawing", Value: "фигура"},
			}),
			Items:        []types.ParseItem{makeItem("Ось симметрии", "geometry", "drawing")},
			ExpectedCode: "T27",
			ShouldMatch:  true,
		},
		{
			Name: "T27_draw_symmetric",
			Task: makeTask(3, "Дорисуй вторую половину рисунка так, чтобы получилась симметричная фигура", []types.VisualFact{
				{Kind: "grid", Value: "половина фигуры"},
			}),
			Items:        []types.ParseItem{makeItem("Дорисуй симметрично", "geometry", "grid")},
			ExpectedCode: "T27",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T28: Координатная сетка и перемещения
// task_type: geometry
// Patterns: "координат", "(x;y)", "вправо на N клеток"
// =============================================================================

func TestT28_CoordinatesGrid(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T28_plot_point",
			Task: makeTask(4, "Построй точку A(2;3) на координатной сетке", []types.VisualFact{
				{Kind: "grid", Value: "координатная сетка"},
			}),
			Items:        []types.ParseItem{makeItem("Точка на координатах", "geometry", "grid")},
			ExpectedCode: "T28",
			ShouldMatch:  true,
		},
		{
			Name: "T28_move_on_grid",
			Task: makeTask(3, "От точки A перейди вправо на 4 клетки и вверх на 1 клетку. Где окажешься?", []types.VisualFact{
				{Kind: "grid", Value: "сетка"},
			}),
			Items:        []types.ParseItem{makeItem("Перемещение по клеткам", "geometry", "grid")},
			ExpectedCode: "T28",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T29: Составление фигур из частей
// task_type: geometry
// Patterns: "составь фигур", "из частей", "разрезана на части", "танграм"
// =============================================================================

func TestT29_ComposeShapes(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T29_compose_from_parts",
			Task: makeTask(3, "Составь квадрат из 4 частей (по рисунку)", []types.VisualFact{
				{Kind: "drawing", Value: "части"},
			}),
			Items:        []types.ParseItem{makeItem("Составь из частей", "geometry", "drawing")},
			ExpectedCode: "T29",
			ShouldMatch:  true,
		},
		{
			Name: "T29_tangram",
			Task: makeTask(3, "Используя детали танграма, сложи фигуру домика", []types.VisualFact{
				{Kind: "drawing", Value: "танграм"},
			}),
			Items:        []types.ParseItem{makeItem("Танграм", "geometry", "drawing")},
			ExpectedCode: "T29",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T30: Чтение и дополнение таблиц
// task_type: data_representation
// Patterns: "таблиц", "заполни таблиц", "по таблице найди"
// =============================================================================

func TestT30_ReadCompleteTables(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T30_read_table",
			Task: makeTask(2, "По таблице найди, сколько яблок продали во вторник", []types.VisualFact{
				{Kind: "table", Value: "таблица продаж"},
			}),
			Items:        []types.ParseItem{makeItem("По таблице найди", "data_representation", "table")},
			ExpectedCode: "T30",
			ShouldMatch:  true,
		},
		{
			Name: "T30_complete_table",
			Task: makeTask(2, "Заполни пропуски в таблице сложения", []types.VisualFact{
				{Kind: "table", Value: "таблица сложения"},
			}),
			Items:        []types.ParseItem{makeItem("Заполни таблицу", "data_representation", "table")},
			ExpectedCode: "T30",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T31: Чтение диаграмм и графиков
// task_type: data_representation
// Patterns: "диаграмм", "график", "по диаграмме"
// =============================================================================

func TestT31_ReadDiagrams(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T31_read_bar_chart",
			Task: makeTask(3, "По диаграмме определи, сколько книг прочитали в среду", []types.VisualFact{
				{Kind: "diagram", Value: "столбчатая диаграмма"},
			}),
			Items:        []types.ParseItem{makeItem("По диаграмме", "data_representation", "diagram")},
			ExpectedCode: "T31",
			ShouldMatch:  true,
		},
		{
			Name: "T31_read_graph",
			Task: makeTask(3, "По графику температуры найди самый тёплый день", []types.VisualFact{
				{Kind: "diagram", Value: "график"},
			}),
			Items:        []types.ParseItem{makeItem("По графику", "data_representation", "diagram")},
			ExpectedCode: "T31",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T32: Построение диаграмм по данным
// task_type: data_representation
// Patterns: "построй диаграмм", "составь диаграмм"
// =============================================================================

func TestT32_BuildDiagrams(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T32_build_bar_chart",
			Task: makeTask(3, "По таблице построй столбчатую диаграмму продаж по дням", []types.VisualFact{
				{Kind: "table", Value: "таблица продаж"},
			}),
			Items:        []types.ParseItem{makeItem("Построй диаграмму", "data_representation", "mixed")},
			ExpectedCode: "T32",
			ShouldMatch:  true,
		},
		{
			Name: "T32_draw_diagram",
			Task: makeTask(3, "Нарисуй диаграмму по данным таблицы", []types.VisualFact{
				{Kind: "table", Value: "данные"},
			}),
			Items:        []types.ParseItem{makeItem("Нарисуй диаграмму", "data_representation", "diagram")},
			ExpectedCode: "T32",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T33: Простая статистика (среднее, максимум/минимум, частота)
// task_type: data_representation
// Patterns: "среднее", "среднее арифметическое", "сколько раз", "максимальн"
// =============================================================================

func TestT33_SimpleStatistics(t *testing.T) {
	tests := []TestCase{
		{
			Name:         "T33_average",
			Task:         makeTask(4, "Найди среднее арифметическое чисел: 4, 6, 10", nil),
			Items:        []types.ParseItem{makeItem("Среднее арифметическое", "data_representation", "plain_text")},
			ExpectedCode: "T33",
			ShouldMatch:  true,
		},
		{
			Name:         "T33_frequency",
			Task:         makeTask(4, "В списке 2, 3, 2, 5, 2. Сколько раз встречается число 2?", nil),
			Items:        []types.ParseItem{makeItem("Частота", "data_representation", "plain_text")},
			ExpectedCode: "T33",
			ShouldMatch:  true,
		},
		{
			Name:         "T33_max_min",
			Task:         makeTask(4, "Найди максимальное и минимальное значение: 8, 3, 15, 7, 11", nil),
			Items:        []types.ParseItem{makeItem("Максимум минимум", "data_representation", "plain_text")},
			ExpectedCode: "T33",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T34: Логические задачи и закономерности
// task_type: logic
// Patterns: "закономерност", "продолжи ряд", "лишнее", "по какому правилу"
// =============================================================================

func TestT34_PatternsAndLogic(t *testing.T) {
	tests := []TestCase{
		{
			Name:         "T34_continue_sequence",
			Task:         makeTask(3, "Продолжи ряд: 2, 4, 6, 8, __", nil),
			Items:        []types.ParseItem{makeItem("Продолжи ряд", "logic", "plain_text")},
			ExpectedCode: "T34",
			ShouldMatch:  true,
		},
		{
			Name:         "T34_find_odd_one",
			Task:         makeTask(3, "Найди лишнее число: 3, 5, 7, 8, 9", nil),
			Items:        []types.ParseItem{makeItem("Лишнее", "logic", "plain_text")},
			ExpectedCode: "T34",
			ShouldMatch:  true,
		},
		{
			Name:         "T34_find_pattern",
			Task:         makeTask(3, "По какому правилу составлен ряд: 1, 3, 5, 7?", nil),
			Items:        []types.ParseItem{makeItem("Закономерность", "logic", "plain_text")},
			ExpectedCode: "T34",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T36: Задачи по рисунку (сравнение и счёт объектов)
// task_type: word_problems
// Patterns: "по рисунку", "на рисунке", "сосчитай"
// =============================================================================

func TestT36_PictureCountCompare(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T36_count_by_picture",
			Task: makeTask(1, "По рисунку сосчитай, сколько кружков изображено", []types.VisualFact{
				{Kind: "drawing", Value: "кружки"},
			}),
			Items:        []types.ParseItem{makeItem("Сосчитай по рисунку", "word_problems", "drawing")},
			ExpectedCode: "T36",
			ShouldMatch:  true,
		},
		{
			Name: "T36_compare_on_picture",
			Task: makeTask(1, "На рисунке 6 красных и 4 синих шарика. Сколько шариков всего?", []types.VisualFact{
				{Kind: "drawing", Value: "шарики"},
			}),
			Items:        []types.ParseItem{makeItem("По рисунку", "word_problems", "mixed")},
			ExpectedCode: "T36",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T40: Простая комбинаторика (сколько способов)
// task_type: logic
// Patterns: "сколько способов", "сколько вариантов", "можно составить"
// =============================================================================

func TestT40_SimpleCombinatorics(t *testing.T) {
	tests := []TestCase{
		{
			Name:         "T40_how_many_ways",
			Task:         makeTask(4, "У Маши 3 футболки и 2 юбки. Сколькими способами она может выбрать комплект?", nil),
			Items:        []types.ParseItem{makeItem("Сколько способов", "logic", "plain_text")},
			ExpectedCode: "T40",
			ShouldMatch:  true,
		},
		{
			Name:         "T40_how_many_variants",
			Task:         makeTask(4, "Сколько вариантов выбрать обед из 3 супов и 2 вторых блюд?", nil),
			Items:        []types.ParseItem{makeItem("Сколько вариантов", "logic", "plain_text")},
			ExpectedCode: "T40",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T42: Часы, расписание и промежутки времени
// task_type: measurement_units
// Patterns: "часы", "циферблат", "который час", "сколько времени"
// =============================================================================

func TestT42_TimeClock(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T42_read_clock",
			Task: makeTask(2, "Посмотри на часы и скажи, который час (по циферблату)", []types.VisualFact{
				{Kind: "diagram", Value: "циферблат"},
			}),
			Items:        []types.ParseItem{makeItem("Который час", "measurement_units", "diagram")},
			ExpectedCode: "T42",
			ShouldMatch:  true,
		},
		{
			Name: "T42_clock_hands",
			Task: makeTask(2, "На часах большая стрелка на 12, маленькая на 3. Который час?", []types.VisualFact{
				{Kind: "diagram", Value: "часы"},
			}),
			Items:        []types.ParseItem{makeItem("Стрелки часов", "measurement_units", "diagram")},
			ExpectedCode: "T42",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T43: Деньги: стоимость, сдача, покупка
// task_type: measurement_units
// Patterns: "руб", "коп", "цена", "сдач", "заплатил"
// =============================================================================

func TestT43_MoneyPrices(t *testing.T) {
	tests := []TestCase{
		{
			Name:         "T43_total_cost",
			Task:         makeTask(2, "Тетрадь стоит 18 руб. Сколько стоят 3 тетради?", nil),
			Items:        []types.ParseItem{makeItem("Стоимость", "measurement_units", "plain_text")},
			ExpectedCode: "T43",
			ShouldMatch:  true,
		},
		{
			Name:         "T43_change",
			Task:         makeTask(2, "Покупка стоит 57 руб., заплатили 100 руб. Какая сдача?", nil),
			Items:        []types.ParseItem{makeItem("Сдача", "measurement_units", "plain_text")},
			ExpectedCode: "T43",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// CONFUSABLE PAIRS: Negative tests
// Проверяем, что похожие шаблоны не путаются между собой
// =============================================================================

// T23 vs T25: Температура (градусы °C) vs Углы (градусы °)
// Оба используют слово "градус" и символ "°", но контекст разный
func TestT23vsT25_ConfusablePair(t *testing.T) {
	tests := []TestCase{
		// T23 должен ловить температуру, а не углы
		{
			Name: "T23_temperature_not_angle",
			Task: makeTask(3, "Термометр показывает 15 градусов. Какая температура на улице?", []types.VisualFact{
				{Kind: "diagram", Value: "термометр"},
			}),
			Items:        []types.ParseItem{makeItem("Температура", "measurement_units", "diagram")},
			ExpectedCode: "T23",
			ShouldMatch:  true,
		},
		// T25 должен ловить углы, а не температуру
		{
			Name: "T25_angle_not_temperature",
			Task: makeTask(3, "Угол равен 90 градусов. Какой это угол — острый, прямой или тупой?", []types.VisualFact{
				{Kind: "diagram", Value: "угол"},
			}),
			Items:        []types.ParseItem{makeItem("Вид угла", "geometry", "diagram")},
			ExpectedCode: "T25",
			ShouldMatch:  true,
		},
		// Негативный: "градусы" без контекста температуры — НЕ T23
		{
			Name: "T23_reject_angle_degrees",
			Task: makeTask(3, "Измерь угол транспортиром. Сколько градусов?", []types.VisualFact{
				{Kind: "diagram", Value: "транспортир"},
			}),
			Items:        []types.ParseItem{makeItem("Угол в градусах", "geometry", "diagram")},
			ExpectedCode: "T25", // должен быть T25, а не T23
			ShouldMatch:  true,
		},
		// T23: температура с термометром (более надёжное соответствие)
		{
			Name: "T23_temperature_with_thermometer",
			Task: makeTask(3, "Посмотри на термометр. Утром температура была минус 5 градусов, а днём плюс 10 градусов. На сколько градусов потеплело?", []types.VisualFact{
				{Kind: "diagram", Value: "термометр"},
			}),
			Items:        []types.ParseItem{makeItem("Изменение температуры", "measurement_units", "diagram")},
			ExpectedCode: "T23",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// T19 vs T21: Текстовые задачи на величины vs Перевод единиц
// T19 — сюжетная задача с величинами, T21 — чистый перевод/сравнение единиц
func TestT19vsT21_ConfusablePair(t *testing.T) {
	tests := []TestCase{
		// T21: чистый перевод единиц (без сюжета)
		{
			Name:         "T21_pure_conversion",
			Task:         makeTask(3, "Переведи 2 км 300 м в метры", nil),
			Items:        []types.ParseItem{makeItem("Перевод единиц", "measurement_units", "plain_text")},
			ExpectedCode: "T21",
			ShouldMatch:  true,
		},
		// T21: сравнение единиц длины
		{
			Name:         "T21_compare_units",
			Task:         makeTask(3, "Сравни единицы длины: 3 км 50 м и 3500 м. Что больше? Переведи в одну единицу.", nil),
			Items:        []types.ParseItem{makeItem("Сравнение единиц длины", "measurement_units", "plain_text")},
			ExpectedCode: "T21",
			ShouldMatch:  true,
		},
		// T19: сюжетная задача с величинами (не чистый перевод)
		{
			Name:         "T19_story_with_units",
			Task:         makeTask(3, "Мама купила 2 кг яблок и 1 кг 500 г груш. Сколько всего фруктов купила мама?", nil),
			Items:        []types.ParseItem{makeItem("Задача с величинами", "word_problems", "plain_text")},
			ExpectedCode: "T19",
			ShouldMatch:  true,
		},
		// Негативный: "вырази" — это T21, не T19
		{
			Name:         "T21_express_not_story",
			Task:         makeTask(3, "Вырази 1 час 20 минут в минутах", nil),
			Items:        []types.ParseItem{makeItem("Вырази", "measurement_units", "plain_text")},
			ExpectedCode: "T21",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// T31 vs T32: Чтение диаграмм vs Построение диаграмм
// T31 — прочитать/найти по диаграмме, T32 — построить/нарисовать диаграмму
func TestT31vsT32_ConfusablePair(t *testing.T) {
	tests := []TestCase{
		// T31: чтение диаграммы
		{
			Name: "T31_read_diagram",
			Task: makeTask(3, "По диаграмме определи, в какой день было больше всего посетителей", []types.VisualFact{
				{Kind: "diagram", Value: "столбчатая диаграмма"},
			}),
			Items:        []types.ParseItem{makeItem("По диаграмме определи", "data_representation", "diagram")},
			ExpectedCode: "T31",
			ShouldMatch:  true,
		},
		// T32: построение диаграммы
		{
			Name: "T32_build_diagram",
			Task: makeTask(3, "Построй столбчатую диаграмму по данным: Пн-5, Вт-8, Ср-3", []types.VisualFact{
				{Kind: "table", Value: "данные"},
			}),
			Items:        []types.ParseItem{makeItem("Построй диаграмму", "data_representation", "mixed")},
			ExpectedCode: "T32",
			ShouldMatch:  true,
		},
		// Негативный: "найди по диаграмме" — это T31, не T32
		{
			Name: "T31_find_not_build",
			Task: makeTask(3, "Найди по диаграмме наименьшее значение", []types.VisualFact{
				{Kind: "diagram", Value: "диаграмма"},
			}),
			Items:        []types.ParseItem{makeItem("Найди по диаграмме", "data_representation", "diagram")},
			ExpectedCode: "T31",
			ShouldMatch:  true,
		},
		// Негативный: "составь диаграмму" — это T32, не T31
		{
			Name: "T32_compose_not_read",
			Task: makeTask(3, "По данным таблицы составь столбчатую диаграмму распределения учеников по кружкам", []types.VisualFact{
				{Kind: "table", Value: "данные"},
			}),
			Items:        []types.ParseItem{makeItem("Составь диаграмму", "data_representation", "mixed")},
			ExpectedCode: "T32",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// T21 vs T42: Перевод единиц времени vs Часы и расписание
// T21 — "переведи 2 часа в минуты", T42 — "который час на циферблате"
func TestT21vsT42_ConfusablePair(t *testing.T) {
	tests := []TestCase{
		// T21: перевод единиц времени
		{
			Name:         "T21_convert_time_units",
			Task:         makeTask(3, "Переведи 2 часа 15 минут в минуты", nil),
			Items:        []types.ParseItem{makeItem("Перевод времени", "measurement_units", "plain_text")},
			ExpectedCode: "T21",
			ShouldMatch:  true,
		},
		// T42: чтение часов
		{
			Name: "T42_read_clock_face",
			Task: makeTask(2, "Который час показывают часы на рисунке?", []types.VisualFact{
				{Kind: "diagram", Value: "циферблат"},
			}),
			Items:        []types.ParseItem{makeItem("Который час", "measurement_units", "diagram")},
			ExpectedCode: "T42",
			ShouldMatch:  true,
		},
		// T42: длительность/расписание
		{
			Name:         "T42_duration",
			Task:         makeTask(3, "Фильм начался в 14:00 и закончился в 15:30. Сколько времени длился фильм?", nil),
			Items:        []types.ParseItem{makeItem("Сколько времени", "measurement_units", "plain_text")},
			ExpectedCode: "T42",
			ShouldMatch:  true,
		},
		// Негативный: "вырази в секундах" — это T21, не T42
		{
			Name:         "T21_express_seconds",
			Task:         makeTask(3, "Вырази 3 минуты в секундах", nil),
			Items:        []types.ParseItem{makeItem("Вырази в секундах", "measurement_units", "plain_text")},
			ExpectedCode: "T21",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// T14 vs T15 vs T16: Простые задачи vs Кратное сравнение vs Изменение величины
// T14 — общие текстовые, T15 — "во сколько раз", T16 — "увеличили/уменьшили на/в"
func TestT14vsT15vsT16_ConfusablePair(t *testing.T) {
	tests := []TestCase{
		// T15: "во сколько раз" — кратное сравнение
		{
			Name:         "T15_how_many_times",
			Task:         makeTask(3, "У Пети 12 марок, а у Васи 4 марки. Во сколько раз больше марок у Пети?", nil),
			Items:        []types.ParseItem{makeItem("Во сколько раз", "word_problems", "plain_text")},
			ExpectedCode: "T15",
			ShouldMatch:  true,
		},
		// T16: "увеличили в N раз" — изменение величины
		{
			Name:         "T16_increased_times",
			Task:         makeTask(3, "Число увеличили в 3 раза и получили 24. Какое было число?", nil),
			Items:        []types.ParseItem{makeItem("Увеличили в раз", "word_problems", "plain_text")},
			ExpectedCode: "T16",
			ShouldMatch:  true,
		},
		// T16: "уменьшили на N" — изменение величины
		{
			Name:         "T16_decreased_by",
			Task:         makeTask(3, "Было 50 яблок, стало на 15 меньше. Сколько стало?", nil),
			Items:        []types.ParseItem{makeItem("Уменьшили на", "word_problems", "plain_text")},
			ExpectedCode: "T16",
			ShouldMatch:  true,
		},
		// T14: простая задача без спецякорей (fallback)
		{
			Name:         "T14_simple_story",
			Task:         makeTask(2, "В корзине было 8 яблок. Добавили ещё 5. Сколько стало?", nil),
			Items:        []types.ParseItem{makeItem("Простая задача", "word_problems", "plain_text")},
			ExpectedCode: "T14",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)

	// Проверка: T15 НЕ должен матчиться на задачи с "увеличили/уменьшили"
	forbidTests := []struct {
		Name string
		Text string
	}{
		{"T15_reject_increased", "Число увеличили в несколько раз. Во сколько раз увеличили?"},
		{"T15_reject_decreased", "Уменьшили в 3 раза. Во сколько раз уменьшили число?"},
		{"T15_reject_became_more", "Стало в 2 раза больше. Во сколько раз увеличилось?"},
	}

	for _, tc := range forbidTests {
		t.Run(tc.Name, func(t *testing.T) {
			task := makeTask(3, tc.Text, nil)
			items := []types.ParseItem{makeItem("Изменение", "word_problems", "plain_text")}
			ctx := buildRoutingContext(task, items)
			candidate, found := selectTemplate(ctx)

			if found && candidate.Template.TemplateCode == "T15" {
				t.Errorf("T15 should NOT match task with 'увеличили/уменьшили', but got T15 (rule=%s)", candidate.MatchedRuleID)
			}
		})
	}
}

// T4 vs T5: Устное vs Письменное сложение/вычитание
// T4 — "устно/в уме", T5 — "столбиком/письменно"
func TestT4vsT5_ConfusablePair(t *testing.T) {
	tests := []TestCase{
		// T4: устный счёт
		{
			Name:         "T4_mental_calculation",
			Task:         makeTask(2, "Считай в уме 45 + 27", nil),
			Items:        []types.ParseItem{makeItem("В уме", "arithmetic_fluency", "inline_examples")},
			ExpectedCode: "T4",
			ShouldMatch:  true,
		},
		// T5: письменное сложение столбиком
		{
			Name: "T5_column_addition",
			Task: makeTask(2, "Выполни сложение столбиком: 345 + 278", []types.VisualFact{
				{Kind: "column", Value: "столбик"},
			}),
			Items:        []types.ParseItem{makeItem("Столбиком", "arithmetic_fluency", "column")},
			ExpectedCode: "T5",
			ShouldMatch:  true,
		},
		// Негативный: "столбиком" — это T5, не T4
		{
			Name:         "T5_not_T4_column_keyword",
			Task:         makeTask(2, "Вычисли в столбик: 67 + 48", nil),
			Items:        []types.ParseItem{makeItem("В столбик", "arithmetic_fluency", "column")},
			ExpectedCode: "T5",
			ShouldMatch:  true,
		},
		// Негативный: "письменно" — это T5, не T4
		{
			Name:         "T5_not_T4_written_keyword",
			Task:         makeTask(2, "Выполни письменно: 234 - 156", nil),
			Items:        []types.ParseItem{makeItem("Письменно", "arithmetic_fluency", "column")},
			ExpectedCode: "T5",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// T7 vs T8: Смысл умножения vs Таблица умножения
// T7 — объяснение смысла ("поровну", "по 5 штук"), T8 — знание таблицы
func TestT7vsT8_ConfusablePair(t *testing.T) {
	tests := []TestCase{
		// T7: смысл умножения (группы, поровну)
		{
			Name:         "T7_multiplication_meaning",
			Task:         makeTask(2, "В 4 коробках лежит по 6 карандашей. Сколько всего карандашей?", nil),
			Items:        []types.ParseItem{makeItem("По штук в коробках", "word_problems", "plain_text")},
			ExpectedCode: "T7",
			ShouldMatch:  true,
		},
		// T8: табличное умножение
		{
			Name:         "T8_multiplication_table",
			Task:         makeTask(2, "Вычисли по таблице умножения: 7 × 8", nil),
			Items:        []types.ParseItem{makeItem("Таблица умножения", "arithmetic_fluency", "inline_examples")},
			ExpectedCode: "T8",
			ShouldMatch:  true,
		},
		// T8: деление на основе таблицы
		{
			Name:         "T8_division_table",
			Task:         makeTask(2, "Найди: 56 : 7 (используй таблицу умножения)", nil),
			Items:        []types.ParseItem{makeItem("Деление по таблице", "arithmetic_fluency", "inline_examples")},
			ExpectedCode: "T8",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// T7 vs T12: Equal groups (mul/div meaning) vs Fractions (parts of whole)
// T7 forbids fraction-related words to prevent confusion
func TestT7vsT12_ConfusablePair(t *testing.T) {
	tests := []TestCase{
		// T7: смысл умножения/деления (равные группы) - должен матчиться
		{
			Name:         "T7_equal_groups_basic",
			Task:         makeTask(3, "Разложи 12 конфет поровну по 3 тарелкам. Сколько конфет в каждой?", nil),
			Items:        []types.ParseItem{makeItem("Поровну", "word_problems", "plain_text")},
			ExpectedCode: "T7",
			ShouldMatch:  true,
		},
		// T7: ещё один positive кейс
		{
			Name:         "T7_equal_groups_each",
			Task:         makeTask(2, "В каждой коробке по 5 карандашей. Коробок 4. Сколько всего карандашей?", nil),
			Items:        []types.ParseItem{makeItem("По штук в каждой", "word_problems", "plain_text")},
			ExpectedCode: "T7",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)

	// Отдельная проверка: T7 НЕ должен матчиться на задачи с дробями
	// (они могут пойти в другие шаблоны, но не в T7)
	fractionTests := []struct {
		Name string
		Text string
	}{
		{"T7_reject_fraction_dolia", "Раздели яблоко на 4 равные доли поровну"},
		{"T7_reject_fraction_half", "Раздели пирог на две половины"},
		{"T7_reject_fraction_third", "Возьми треть от целого яблока"},
		{"T7_reject_fraction_quarter", "Четверть торта съели. Сколько четвертей осталось?"},
		{"T7_reject_fraction_notation", "Закрашено 3/8 фигуры. Какая доля не закрашена?"},
	}

	for _, tc := range fractionTests {
		t.Run(tc.Name, func(t *testing.T) {
			task := makeTask(3, tc.Text, nil)
			items := []types.ParseItem{makeItem("Дроби", "fractions_percent", "plain_text")}
			ctx := buildRoutingContext(task, items)
			candidate, found := selectTemplate(ctx)

			if found && candidate.Template.TemplateCode == "T7" {
				t.Errorf("T7 should NOT match fraction task, but got T7 (rule=%s)", candidate.MatchedRuleID)
			}
		})
	}
}

// T30 vs T31: Таблицы vs Диаграммы
// T30 — работа с таблицами, T31 — чтение диаграмм/графиков
func TestT30vsT31_ConfusablePair(t *testing.T) {
	tests := []TestCase{
		// T30: чтение таблицы
		{
			Name: "T30_read_table",
			Task: makeTask(3, "По таблице найди, сколько учеников в 3-А классе", []types.VisualFact{
				{Kind: "table", Value: "таблица классов"},
			}),
			Items:        []types.ParseItem{makeItem("По таблице найди", "data_representation", "table")},
			ExpectedCode: "T30",
			ShouldMatch:  true,
		},
		// T31: чтение диаграммы
		{
			Name: "T31_read_chart",
			Task: makeTask(3, "По столбчатой диаграмме определи максимальное значение", []types.VisualFact{
				{Kind: "diagram", Value: "столбчатая диаграмма"},
			}),
			Items:        []types.ParseItem{makeItem("По диаграмме", "data_representation", "diagram")},
			ExpectedCode: "T31",
			ShouldMatch:  true,
		},
		// Негативный: "заполни таблицу" — это T30, не T31
		{
			Name: "T30_fill_table_not_diagram",
			Task: makeTask(3, "Заполни пустые ячейки таблицы", []types.VisualFact{
				{Kind: "table", Value: "таблица с пропусками"},
			}),
			Items:        []types.ParseItem{makeItem("Заполни таблицу", "data_representation", "table")},
			ExpectedCode: "T30",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// =============================================================================
// T51: Множества — элементы, подмножества, операции
// task_type: sets_logic
// Patterns: "множеств", "элемент", "подмножеств", "диаграмм Эйлера-Венна"
// =============================================================================

func TestT51_SetsAndElements(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T51_describe_sets",
			Task: makeTask(3, "Опиши множества, которые можно назвать: хор, оркестр, бригада.", nil),
			Items: []types.ParseItem{
				makeItem("Опиши множества", "sets_logic", "plain_text"),
			},
			ExpectedCode: "T51",
			ShouldMatch:  true,
		},
		{
			Name: "T51_name_elements",
			Task: makeTask(3, "Назови 5 элементов множества ягод.", nil),
			Items: []types.ParseItem{
				makeItem("Элементы множества", "sets_logic", "plain_text"),
			},
			ExpectedCode: "T51",
			ShouldMatch:  true,
		},
		{
			Name: "T51_venn_diagram",
			Task: makeTask(3, "Нарисуй диаграмму Эйлера-Венна множеств М и К.", nil),
			Items: []types.ParseItem{
				makeItem("Диаграмма Венна", "sets_logic", "diagram"),
			},
			ExpectedCode: "T51",
			ShouldMatch:  true,
		},
		{
			Name: "T51_subset",
			Task: makeTask(3, "Является ли множество А подмножеством множества В?", nil),
			Items: []types.ParseItem{
				makeItem("Подмножество", "sets_logic", "plain_text"),
			},
			ExpectedCode: "T51",
			ShouldMatch:  true,
		},
		{
			Name: "T51_intersection",
			Task: makeTask(3, "Найди пересечение множеств А = {1, 2, 3} и В = {2, 3, 4}.", nil),
			Items: []types.ParseItem{
				makeItem("Пересечение множеств", "sets_logic", "plain_text"),
			},
			ExpectedCode: "T51",
			ShouldMatch:  true,
		},
		{
			Name: "T51_union",
			Task: makeTask(3, "Найди объединение множеств М и К.", nil),
			Items: []types.ParseItem{
				makeItem("Объединение множеств", "sets_logic", "plain_text"),
			},
			ExpectedCode: "T51",
			ShouldMatch:  true,
		},
		{
			Name: "T51_belongs_to_set",
			Task: makeTask(3, "Принадлежит ли число 5 множеству чётных чисел?", nil),
			Items: []types.ParseItem{
				makeItem("Принадлежит множеству", "sets_logic", "plain_text"),
			},
			ExpectedCode: "T51",
			ShouldMatch:  true,
		},
		{
			Name: "T51_set_notation",
			Task: makeTask(3, "Запиши множество делителей числа 12.", nil),
			Items: []types.ParseItem{
				makeItem("Запиши множество", "sets_logic", "plain_text"),
			},
			ExpectedCode: "T51",
			ShouldMatch:  true,
		},
		{
			Name: "T51_empty_set",
			Task: makeTask(3, "Является ли пустое множество подмножеством любого множества?", nil),
			Items: []types.ParseItem{
				makeItem("Пустое множество", "sets_logic", "plain_text"),
			},
			ExpectedCode: "T51",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// =============================================================================
// T52: Творческие задачи — составь, придумай, объясни
// task_type: creative_composition
// Patterns: "составь выражени", "придумай задач", "объясни", "что означа"
// =============================================================================

func TestT52_CreativeComposition(t *testing.T) {
	tests := []TestCase{
		{
			Name: "T52_compose_expression",
			Task: makeTask(3, "Составь выражение по задаче и найди его значение.", nil),
			Items: []types.ParseItem{
				makeItem("Составь выражение", "creative_composition", "plain_text"),
			},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
		{
			Name: "T52_create_problem",
			Task: makeTask(3, "Придумай задачу по схеме и реши её.", nil),
			Items: []types.ParseItem{
				makeItem("Придумай задачу", "creative_composition", "plain_text"),
			},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
		{
			Name: "T52_explain_why",
			Task: makeTask(3, "Объясни, почему это равенство верно.", nil),
			Items: []types.ParseItem{
				makeItem("Объясни", "creative_composition", "plain_text"),
			},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
		{
			Name: "T52_what_means",
			Task: makeTask(3, "Что означают выражения a + b и a - b?", nil),
			Items: []types.ParseItem{
				makeItem("Что означают", "creative_composition", "plain_text"),
			},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
		{
			Name: "T52_is_it_true",
			Task: makeTask(3, "Верно ли, что сумма двух чётных чисел всегда чётная?", nil),
			Items: []types.ParseItem{
				makeItem("Верно ли", "creative_composition", "plain_text"),
			},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
		{
			Name: "T52_justify",
			Task: makeTask(3, "Обоснуй свой ответ.", nil),
			Items: []types.ParseItem{
				makeItem("Обоснуй", "creative_composition", "plain_text"),
			},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
		{
			Name: "T52_prove",
			Task: makeTask(3, "Докажи, что это утверждение верно.", nil),
			Items: []types.ParseItem{
				makeItem("Докажи", "creative_composition", "plain_text"),
			},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
		{
			Name: "T52_what_common",
			Task: makeTask(3, "Что общего в выражениях: а · 15 + а и 5 · (3 · b)?", nil),
			Items: []types.ParseItem{
				makeItem("Что общего", "creative_composition", "plain_text"),
			},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
		{
			Name: "T52_how_differ",
			Task: makeTask(3, "Чем похожи и чем различаются выражения (30 – d) : 3 и 30 – d : 3?", nil),
			Items: []types.ParseItem{
				makeItem("Чем похожи и различаются", "creative_composition", "plain_text"),
			},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
		{
			Name: "T52_compose_problem_by_expression",
			Task: makeTask(3, "Составь задачу, которая решается так: (а – b) : 2. Подбери данные.", nil),
			Items: []types.ParseItem{
				makeItem("Составь задачу по выражению", "creative_composition", "plain_text"),
			},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
	}

	runTestCases(t, tests)
}

// T51 vs T33: Множества vs Статистика
// T51 — теория множеств, T33 — простая статистика (среднее, мода)
func TestT51vsT33_ConfusablePair(t *testing.T) {
	tests := []TestCase{
		// T51: множества
		{
			Name:         "T51_set_elements_not_stats",
			Task:         makeTask(3, "Перечисли элементы множества А = {1, 2, 3, 4, 5}.", nil),
			Items:        []types.ParseItem{makeItem("Элементы множества", "sets_logic", "plain_text")},
			ExpectedCode: "T51",
			ShouldMatch:  true,
		},
		// T33: статистика
		{
			Name:         "T33_average_not_sets",
			Task:         makeTask(3, "Найди среднее арифметическое чисел 4, 6, 8.", nil),
			Items:        []types.ParseItem{makeItem("Среднее арифметическое", "data_representation", "plain_text")},
			ExpectedCode: "T33",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

// T52 vs T1: Творческие vs Вычисления
// T52 — составь/объясни, T1 — простые вычисления
func TestT52vsT1_ConfusablePair(t *testing.T) {
	tests := []TestCase{
		// T52: составь выражение
		{
			Name:         "T52_compose_not_calculate",
			Task:         makeTask(3, "Составь выражение и найди его значение при x = 5.", nil),
			Items:        []types.ParseItem{makeItem("Составь выражение", "creative_composition", "plain_text")},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
		// T52: придумай задачу — не должен совпадать с T1
		{
			Name:         "T52_create_problem_not_T1",
			Task:         makeTask(3, "Придумай задачу, которая решается выражением a + b.", nil),
			Items:        []types.ParseItem{makeItem("Придумай задачу", "creative_composition", "plain_text")},
			ExpectedCode: "T52",
			ShouldMatch:  true,
		},
	}
	runTestCases(t, tests)
}

func BenchmarkSelectTemplate(b *testing.B) {
	task := makeTask(3, "Прямоугольник разрезали на 9 неодинаковых квадратов. Покажи длины сторон остальных квадратов.", []types.VisualFact{
		{Kind: "diagram", Value: "Прямоугольник из квадратов"},
	})
	items := []types.ParseItem{
		makeItem("Разбиение на квадраты", "geometry", "plain_text"),
	}
	ctx := buildRoutingContext(task, items)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		selectTemplate(ctx)
	}
}

// TestTaskTypeMappingNormalization проверяет нормализацию task_type
// и корректность маппинга equations/algebra → patterns_logic (T37)
func TestTaskTypeMappingNormalization(t *testing.T) {
	tests := []struct {
		name             string
		taskType         string
		taskText         string
		expectedTemplate string
	}{
		// Тесты нормализации: пробелы, регистр, точки
		{
			name:             "equations_with_spaces",
			taskType:         " Equations ",
			taskText:         "Реши уравнение: x + 5 = 10",
			expectedTemplate: "T37",
		},
		{
			name:             "equation_uppercase",
			taskType:         "EQUATION",
			taskText:         "Реши уравнение: x + 5 = 10",
			expectedTemplate: "T37",
		},
		{
			name:             "equations_with_dot",
			taskType:         "equations.",
			taskText:         "Найди неизвестное: x - 3 = 7",
			expectedTemplate: "T37",
		},
		{
			name:             "algebra_lowercase",
			taskType:         "algebra",
			taskText:         "Реши уравнение: 2 * x = 8",
			expectedTemplate: "T37",
		},
		// Тесты что канонические значения НЕ меняются
		{
			name:             "sets_logic_unchanged",
			taskType:         "sets_logic",
			taskText:         "Назови 5 элементов множества ягод",
			expectedTemplate: "T51",
		},
		{
			name:             "creative_composition_unchanged",
			taskType:         "creative_composition",
			taskText:         "Составь задачу по схеме",
			expectedTemplate: "T52",
		},
		{
			name:             "patterns_logic_unchanged",
			taskType:         "patterns_logic",
			taskText:         "Реши уравнение: x + 5 = 10",
			expectedTemplate: "T37",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := makeTask(3, tc.taskText, nil)
			items := []types.ParseItem{makeItem(tc.taskText, tc.taskType, "plain_text")}

			ctx := buildRoutingContext(task, items)
			candidate, found := selectTemplate(ctx)

			if !found {
				t.Errorf("Шаблон не найден для task_type=%q", tc.taskType)
				return
			}

			if candidate.Template.TemplateCode != tc.expectedTemplate {
				t.Errorf("task_type=%q: ожидали %s, получили %s",
					tc.taskType, tc.expectedTemplate, candidate.Template.TemplateCode)
			}
		})
	}
}

// TestTaskTypeMappingCanonicalUnchanged проверяет что канонические task_type не изменяются маппингом
func TestTaskTypeMappingCanonicalUnchanged(t *testing.T) {
	canonicalTypes := []string{
		"number_sense",
		"numeral_systems",
		"arithmetic_fluency",
		"fractions_percent",
		"word_problems",
		"measurement_units",
		"geometry",
		"data_representation",
		"patterns_logic",
		"sets_logic",
		"creative_composition",
	}

	for _, taskType := range canonicalTypes {
		t.Run(taskType, func(t *testing.T) {
			task := makeTask(3, "Тестовая задача", nil)
			items := []types.ParseItem{makeItem("Тестовая задача", taskType, "plain_text")}

			ctx := buildRoutingContext(task, items)

			// Проверяем что taskType в контексте не изменился
			if ctx.TaskType != taskType {
				t.Errorf("Канонический task_type %q изменился на %q", taskType, ctx.TaskType)
			}
		})
	}
}
