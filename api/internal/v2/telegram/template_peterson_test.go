package telegram

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"child-bot/api/internal/v2/types"
)

// PetersonTask структура задачи из peterson JSON
type PetersonTask struct {
	PageNumber    int    `json:"page_number"`
	TaskNumber    string `json:"task_number"`
	ConditionText string `json:"condition_text"`
}

type PetersonFile struct {
	Tasks []PetersonTask `json:"tasks"`
}

func TestAnalyzeUnmatchedTasks(t *testing.T) {
	SetTemplatesDir("../../../internal/v2/templates")
	ResetTemplatesCache()

	files := []string{
		"../../../../tools/textbook_parser/output/peterson_3_part1.json",
		"../../../../tools/textbook_parser/output/peterson_3_part2.json",
		"../../../../tools/textbook_parser/output/peterson_3_part3.json",
	}

	// Категоризация задач без шаблона
	categories := map[string][]string{
		"множеств": {},
		"диаграмм": {},
		"найди":    {},
		"выражени": {},
		"вычисли":  {},
		"реши":     {},
		"сравни":   {},
		"other":    {},
	}

	for _, filePath := range files {
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var pf PetersonFile
		if err := json.Unmarshal(data, &pf); err != nil {
			continue
		}

		for _, task := range pf.Tasks {
			parseTask := types.ParseTask{
				Subject:       types.SubjectMath,
				Grade:         3,
				TaskTextClean: task.ConditionText,
			}
			parseItem := types.ParseItem{
				ItemTextClean: task.ConditionText,
				PedKeys:       types.PedKeys{Format: "plain_text"},
			}

			if getTemplateID(parseTask, []types.ParseItem{parseItem}) == "" {
				text := task.ConditionText
				lowText := strings.ToLower(text)

				categorized := false
				for cat := range categories {
					if cat != "other" && strings.Contains(lowText, cat) {
						if len(text) > 100 {
							text = text[:100] + "..."
						}
						categories[cat] = append(categories[cat], text)
						categorized = true
						break
					}
				}
				if !categorized {
					if len(text) > 100 {
						text = text[:100] + "..."
					}
					categories["other"] = append(categories["other"], text)
				}
			}
		}
	}

	fmt.Println("\n=== ДЕТАЛЬНЫЙ АНАЛИЗ ЗАДАЧ БЕЗ ШАБЛОНА ===")
	for cat, tasks := range categories {
		if len(tasks) > 0 {
			fmt.Printf("\n--- %s (%d задач) ---\n", strings.ToUpper(cat), len(tasks))
			for i, t := range tasks {
				if i >= 5 { // показываем только первые 5
					fmt.Printf("  ... и ещё %d\n", len(tasks)-5)
					break
				}
				fmt.Printf("  %d. %s\n", i+1, t)
			}
		}
	}
}

func TestTaskTypeMapping(t *testing.T) {
	SetTemplatesDir("../../../internal/v2/templates")
	ResetTemplatesCache()

	// Тест с task_type: "expressions" (как в примере пользователя)
	parseTask := types.ParseTask{
		TaskId:        "test",
		Subject:       types.SubjectMath,
		Grade:         3,
		TaskTextClean: "Вычисли: 70 * 4 = 280, 42 : 3 = 14",
	}

	parseItem := types.ParseItem{
		ItemId:        "i1",
		ItemTextClean: "Вычисли: 70 * 4 = 280, 42 : 3 = 14",
		PedKeys: types.PedKeys{
			TaskType: "expressions", // LLM вернула этот task_type
			Format:   "plain_text",
		},
	}

	templateID := getTemplateID(parseTask, []types.ParseItem{parseItem})
	fmt.Printf("task_type='expressions' → шаблон: %s\n", templateID)

	if templateID == "" {
		t.Error("Шаблон не найден для task_type='expressions'")
	}
}

func TestPetersonTemplateRouting(t *testing.T) {
	// Устанавливаем путь к шаблонам
	SetTemplatesDir("../../../internal/v2/templates")
	ResetTemplatesCache()

	files := []string{
		"../../../../tools/textbook_parser/output/peterson_3_part1.json",
		"../../../../tools/textbook_parser/output/peterson_3_part2.json",
		"../../../../tools/textbook_parser/output/peterson_3_part3.json",
	}

	var totalTasks int
	var matchedTasks int
	var unmatchedTasks []string
	var unmatchedTexts []string // полные тексты для анализа
	templateStats := make(map[string]int)

	for _, filePath := range files {
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Logf("Не удалось прочитать файл %s: %v", filePath, err)
			continue
		}

		var pf PetersonFile
		if err := json.Unmarshal(data, &pf); err != nil {
			t.Logf("Не удалось распарсить %s: %v", filePath, err)
			continue
		}

		for _, task := range pf.Tasks {
			totalTasks++

			// Создаём минимальный ParseTask для маршрутизации
			parseTask := types.ParseTask{
				TaskId:        fmt.Sprintf("page%d_task%s", task.PageNumber, task.TaskNumber),
				Subject:       types.SubjectMath,
				Grade:         3,
				TaskTextClean: task.ConditionText,
				VisualFacts:   []types.VisualFact{},
				Quality:       types.ParseTaskQuality{Flags: []string{}},
			}

			// Создаём минимальный ParseItem (без task_type - будет fallback)
			parseItem := types.ParseItem{
				ItemId:        "i1",
				ItemTextClean: task.ConditionText,
				PedKeys: types.PedKeys{
					TaskType:       "", // пустой - fallback режим
					Format:         "plain_text",
					TemplateParams: map[string]interface{}{},
				},
			}

			templateID := getTemplateID(parseTask, []types.ParseItem{parseItem})

			if templateID != "" {
				matchedTasks++
				templateStats[templateID]++
			} else {
				unmatchedTexts = append(unmatchedTexts, task.ConditionText)
				if len(unmatchedTasks) < 20 { // ограничиваем список
					text := task.ConditionText
					if len(text) > 80 {
						text = text[:80] + "..."
					}
					unmatchedTasks = append(unmatchedTasks, fmt.Sprintf("стр.%d №%s: %s", task.PageNumber, task.TaskNumber, text))
				}
			}
		}
	}

	// Выводим результаты
	fmt.Printf("\n=== РЕЗУЛЬТАТЫ АНАЛИЗА PETERSON 3 КЛАСС ===\n")
	fmt.Printf("Всего задач: %d\n", totalTasks)
	fmt.Printf("Шаблон найден: %d (%.1f%%)\n", matchedTasks, float64(matchedTasks)*100/float64(totalTasks))
	fmt.Printf("Шаблон НЕ найден: %d (%.1f%%)\n", totalTasks-matchedTasks, float64(totalTasks-matchedTasks)*100/float64(totalTasks))

	// Статистика по шаблонам
	fmt.Printf("\nРаспределение по шаблонам:\n")
	type kv struct {
		k string
		v int
	}
	var sorted []kv
	for k, v := range templateStats {
		sorted = append(sorted, kv{k, v})
	}
	// Сортируем по убыванию
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].v > sorted[i].v {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	for _, item := range sorted {
		fmt.Printf("  %s: %d (%.1f%%)\n", item.k, item.v, float64(item.v)*100/float64(totalTasks))
	}

	if len(unmatchedTasks) > 0 {
		fmt.Printf("\nПримеры задач без шаблона (первые 20):\n")
		for i, task := range unmatchedTasks {
			fmt.Printf("  %d. %s\n", i+1, task)
		}
	}

	// Анализ ключевых слов в задачах без шаблона
	keywords := map[string]int{
		"множеств":   0,
		"элемент":    0,
		"перечисли":  0,
		"назови":     0,
		"опиши":      0,
		"верн":       0,
		"обоснуй":    0,
		"докажи":     0,
		"объясни":    0,
		"почему":     0,
		"диаграмм":   0,
		"график":     0,
		"таблиц":     0,
		"схем":       0,
		"рисун":      0,
		"вычисли":    0,
		"реши":       0,
		"найди":      0,
		"сравни":     0,
		"раздели":    0,
		"умнож":      0,
		"сложи":      0,
		"вычти":      0,
		"уравнен":    0,
		"неравенств": 0,
		"выражени":   0,
		"пример":     0,
		"задач":      0,
	}

	for _, text := range unmatchedTexts {
		lowText := strings.ToLower(text)
		for kw := range keywords {
			if strings.Contains(lowText, kw) {
				keywords[kw]++
			}
		}
	}

	fmt.Printf("\n=== АНАЛИЗ ЗАДАЧ БЕЗ ШАБЛОНА ===\n")
	fmt.Printf("Ключевые слова (из %d задач):\n", len(unmatchedTexts))

	// Сортируем по частоте
	type kwStat struct {
		word  string
		count int
	}
	var kwList []kwStat
	for k, v := range keywords {
		if v > 0 {
			kwList = append(kwList, kwStat{k, v})
		}
	}
	for i := 0; i < len(kwList); i++ {
		for j := i + 1; j < len(kwList); j++ {
			if kwList[j].count > kwList[i].count {
				kwList[i], kwList[j] = kwList[j], kwList[i]
			}
		}
	}
	for _, kw := range kwList {
		fmt.Printf("  '%s': %d (%.1f%%)\n", kw.word, kw.count, float64(kw.count)*100/float64(len(unmatchedTexts)))
	}
}
