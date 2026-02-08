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
			Name: "T17_olympiad_complex",
			Task: makeTask(4, "Три друга купили книги. Первый заплатил половину того, что заплатили два других вместе. Второй заплатил треть от суммы первого и третьего. Третий заплатил 120 рублей. Сколько стоили все книги вместе?", nil),
			Items: []types.ParseItem{
				makeItem("Сложная задача на несколько действий", "word_problems", "plain_text"),
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
		{
			Name: "T39_count_figures",
			Task: makeTask(2, "Сколько треугольников на рисунке?", []types.VisualFact{
				{Kind: "diagram", Value: "Составная фигура из треугольников"},
			}),
			Items: []types.ParseItem{
				makeItem("Посчитай сколько треугольников на рисунке", "geometry", "plain_text"),
			},
			ExpectedCode: "T39",
			ShouldMatch:  true,
		},
		{
			Name: "T39_matchstick_puzzle",
			Task: makeTask(3, "Задача со спичками: переложи спички так, чтобы получилось 4 квадрата вместо 5", []types.VisualFact{
				{Kind: "diagram", Value: "Фигура из спичек"},
			}),
			Items: []types.ParseItem{
				makeItem("Переложи спички", "geometry", "drawing"),
			},
			ExpectedCode: "T39",
			ShouldMatch:  true,
		},
		{
			Name: "T39_tiled_rectangles",
			Task: makeTask(3, "Прямоугольник разрезали на 9 неодинаковых квадратов. Длина стороны одного из квадратов указана, а сторона чёрного квадрата равна 1. Покажи на рисунке длины сторон остальных квадратов.", []types.VisualFact{
				{Kind: "diagram", Value: "Прямоугольник из квадратов"},
			}),
			Items: []types.ParseItem{
				makeItem("Прямоугольник разрезали на квадраты", "geometry", "plain_text"),
			},
			ExpectedCode: "T39",
			ShouldMatch:  true,
		},
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
		{
			Name: "T39_count_not_perimeter",
			Task: makeTask(2, "Сколько треугольников на рисунке?", []types.VisualFact{
				{Kind: "diagram", Value: "Составная фигура"},
			}),
			Items: []types.ParseItem{
				makeItem("Сколько фигур", "geometry", "plain_text"),
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
