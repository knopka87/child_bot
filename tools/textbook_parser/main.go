package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Textbook представляет учебник
type Textbook struct {
	Subject   string `json:"subject"`
	Grade     int    `json:"grade"`
	Authors   string `json:"authors"`
	Title     string `json:"title"`
	Part      int    `json:"part"`
	Year      int    `json:"year"`
	Publisher string `json:"publisher"`
	SourceURL string `json:"source_url"`
}

// Task представляет задачу из учебника
type Task struct {
	PageNumber    int         `json:"page_number"`
	TaskNumber    string      `json:"task_number"`
	TaskOrder     int         `json:"task_order"`
	ConditionText string      `json:"condition_text"`
	ConditionHTML string      `json:"condition_html"`
	SolutionText  string      `json:"solution_text"`
	SolutionHTML  string      `json:"solution_html"`
	HintsText     string      `json:"hints_text"`
	HintsHTML     string      `json:"hints_html"`
	HasSubItems   bool        `json:"has_sub_items"`
	SubItems      []SubItem   `json:"sub_items,omitempty"`
	Images        []TaskImage `json:"images"`
	SourceURL     string      `json:"source_url"`
	Index         *TaskIndex  `json:"index,omitempty"`
}

// SubItem представляет подпункт задачи
type SubItem struct {
	Letter    string `json:"letter"`
	Condition string `json:"condition"`
	Solution  string `json:"solution"`
}

// TaskImage представляет картинку к задаче
type TaskImage struct {
	ImageType     string `json:"image_type"` // condition, solution, hint
	ImageOrder    int    `json:"image_order"`
	SubItemLetter string `json:"sub_item_letter,omitempty"`
	OriginalURL   string `json:"original_url"`
	LocalPath     string `json:"local_path"`
	AltText       string `json:"alt_text"`
}

// TaskIndex представляет индексные данные для быстрого поиска
type TaskIndex struct {
	NormalizedHash   string   `json:"normalized_hash"`
	NumbersSignature string   `json:"numbers_signature"`
	Keywords         []string `json:"keywords"`
	NormalizedText   string   `json:"normalized_text"`
}

// Config конфигурация парсера
type Config struct {
	BaseURL        string
	OutputDir      string
	ImagesDir      string
	Part           int
	StartPage      int
	EndPage        int
	DelayMs        int
	DownloadImages bool
}

var config = Config{
	BaseURL:        "https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik",
	OutputDir:      "output",
	ImagesDir:      "images",
	Part:           1,
	StartPage:      3,
	EndPage:        111,
	DelayMs:        500,
	DownloadImages: false,
}

func init() {
	// Парсим аргументы командной строки
	for i, arg := range os.Args[1:] {
		switch arg {
		case "-part":
			if i+2 < len(os.Args) {
				config.Part, _ = strconv.Atoi(os.Args[i+2])
			}
		case "-start":
			if i+2 < len(os.Args) {
				config.StartPage, _ = strconv.Atoi(os.Args[i+2])
			}
		case "-end":
			if i+2 < len(os.Args) {
				config.EndPage, _ = strconv.Atoi(os.Args[i+2])
			}
		case "-delay":
			if i+2 < len(os.Args) {
				config.DelayMs, _ = strconv.Atoi(os.Args[i+2])
			}
		case "-download-images":
			config.DownloadImages = true
		}
	}
}

func main() {
	// Создаём директории
	os.MkdirAll(filepath.Join(config.OutputDir, config.ImagesDir), 0755)

	textbook := Textbook{
		Subject:   "math",
		Grade:     3,
		Authors:   "Петерсон Л.Г.",
		Title:     "Математика 3 класс",
		Part:      config.Part,
		Year:      2022,
		Publisher: "Просвещение",
		SourceURL: config.BaseURL,
	}

	var allTasks []Task

	for page := config.StartPage; page <= config.EndPage; page++ {
		fmt.Printf("Парсинг страницы %d...\n", page)

		tasks, err := parsePage(config.Part, page)
		if err != nil {
			fmt.Printf("Ошибка на странице %d: %v\n", page, err)
			continue
		}

		allTasks = append(allTasks, tasks...)
		fmt.Printf("  Найдено задач: %d\n", len(tasks))

		// Задержка между запросами
		time.Sleep(time.Duration(config.DelayMs) * time.Millisecond)
	}

	fmt.Printf("\nВсего задач: %d\n", len(allTasks))

	// Скачиваем изображения если нужно
	if config.DownloadImages {
		downloadAllImages(allTasks)
	}

	// Сохраняем JSON
	saveJSON(textbook, allTasks)

	// Генерируем SQL
	generateSQL(textbook, allTasks)

	fmt.Println("Готово!")
}

func parsePage(part, page int) ([]Task, error) {
	url := fmt.Sprintf("%s/%d-chast-stranitsa-%d/", config.BaseURL, part, page)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var tasks []Task

	// Ищем контент страницы
	content := doc.Find(".entry-content, .content, article, main").First()
	if content.Length() == 0 {
		content = doc.Find("body")
	}

	html, _ := content.Html()

	// Разбиваем по задачам (ищем паттерн "Номер N." в разных вариантах тегов)
	taskRegex := regexp.MustCompile(`(?is)<(?:strong|b)[^>]*>\s*Номер\s+(\d+)\s*\.?\s*</(?:strong|b)>`)
	matches := taskRegex.FindAllStringSubmatchIndex(html, -1)

	for i, match := range matches {
		taskNum := html[match[2]:match[3]]

		// Определяем границы задачи
		start := match[0]
		end := len(html)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}

		taskHTML := html[start:end]
		task := parseTask(taskNum, taskHTML, part, page, i, url)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func parseTask(taskNum, taskHTML string, part, page, order int, sourceURL string) Task {
	task := Task{
		PageNumber: page,
		TaskNumber: taskNum,
		TaskOrder:  order,
		SourceURL:  sourceURL,
	}

	// Разделяем условие и ответ по "<b class="otvet">Ответ" или "<strong>Ответ"
	answerRegex := regexp.MustCompile(`(?is)<(?:strong|b)[^>]*>\s*Ответ\s*:?\s*</(?:strong|b)>`)
	parts := answerRegex.Split(taskHTML, 2)

	// Условие
	if len(parts) > 0 {
		condHTML := parts[0]
		// Убираем заголовок "Номер N." (разные варианты тегов)
		condHTML = regexp.MustCompile(`(?is)<(?:p|div)[^>]*>\s*<(?:strong|b)[^>]*>\s*Номер\s+\d+\s*\.?\s*</(?:strong|b)>\s*</(?:p|div)>`).ReplaceAllString(condHTML, "")
		condHTML = regexp.MustCompile(`(?is)<(?:strong|b)[^>]*>\s*Номер\s+\d+\s*\.?\s*</(?:strong|b)>`).ReplaceAllString(condHTML, "")
		task.ConditionHTML = strings.TrimSpace(condHTML)
		task.ConditionText = stripHTML(condHTML)

		// Извлекаем картинки из условия
		task.Images = append(task.Images, extractImages(condHTML, "condition", part, page, taskNum)...)
	}

	// Ответ
	if len(parts) > 1 {
		solHTML := parts[1]

		// Извлекаем рекомендации из <div class="recomended-block">
		recomRegex := regexp.MustCompile(`(?is)<div[^>]*class="[^"]*recomended-block[^"]*"[^>]*>(.*?)</div>`)
		recomMatches := recomRegex.FindAllStringSubmatch(solHTML, -1)
		if len(recomMatches) > 0 {
			var hintsHTML []string
			var hintsText []string
			for _, m := range recomMatches {
				hintsHTML = append(hintsHTML, m[0])
				hintsText = append(hintsText, stripHTML(m[1]))
			}
			task.HintsHTML = strings.Join(hintsHTML, "\n")
			task.HintsText = strings.Join(hintsText, "\n")

			// Убираем рекомендации из ответа
			solHTML = recomRegex.ReplaceAllString(solHTML, "")
		}

		// Также извлекаем из <div class="sdvig"> само решение
		sdvigRegex := regexp.MustCompile(`(?is)<div[^>]*class="[^"]*sdvig[^"]*"[^>]*>(.*?)</div>`)
		sdvigMatches := sdvigRegex.FindAllStringSubmatch(solHTML, -1)
		if len(sdvigMatches) > 0 {
			var solParts []string
			for _, m := range sdvigMatches {
				solParts = append(solParts, m[1])
			}
			solHTML = strings.Join(solParts, "\n")
		}

		task.SolutionHTML = strings.TrimSpace(solHTML)
		task.SolutionText = stripHTML(solHTML)

		// Извлекаем картинки из ответа
		task.Images = append(task.Images, extractImages(solHTML, "solution", part, page, taskNum)...)
	}

	// Проверяем наличие подпунктов (а), б), в) и т.д.)
	subItemRegex := regexp.MustCompile(`(?m)^[а-е]\)`)
	if subItemRegex.MatchString(task.ConditionText) {
		task.HasSubItems = true
		task.SubItems = parseSubItems(task.ConditionText, task.SolutionText)
	}

	// Строим индекс для быстрого поиска
	task.Index = buildTaskIndex(task.ConditionText)

	return task
}

func extractImages(html, imageType string, part, page int, taskNum string) []TaskImage {
	var images []TaskImage

	imgRegex := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["'][^>]*(?:alt=["']([^"']*)["'])?[^>]*>`)
	matches := imgRegex.FindAllStringSubmatch(html, -1)

	for i, match := range matches {
		imgURL := match[1]
		altText := ""
		if len(match) > 2 {
			altText = match[2]
		}

		// Формируем полный URL если нужно
		if !strings.HasPrefix(imgURL, "http") {
			imgURL = "https://gdz-raketa.ru" + imgURL
		}

		// Локальный путь (структура: peterson/{grade}/part{part}/page{page}/...)
		ext := filepath.Ext(imgURL)
		if ext == "" {
			ext = ".jpg"
		}
		localPath := fmt.Sprintf("peterson/3/part%d/page%d/task%s_%s_%d%s",
			part, page, taskNum, imageType, i, ext)

		images = append(images, TaskImage{
			ImageType:   imageType,
			ImageOrder:  i,
			OriginalURL: imgURL,
			LocalPath:   localPath,
			AltText:     altText,
		})
	}

	return images
}

func parseSubItems(condition, solution string) []SubItem {
	var items []SubItem

	letters := []string{"а", "б", "в", "г", "д", "е"}
	for _, letter := range letters {
		pattern := regexp.MustCompile(fmt.Sprintf(`(?s)%s\)\s*(.+?)(?:[б-е]\)|$)`, letter))
		condMatch := pattern.FindStringSubmatch(condition)
		solMatch := pattern.FindStringSubmatch(solution)

		if condMatch != nil || solMatch != nil {
			item := SubItem{Letter: letter}
			if condMatch != nil {
				item.Condition = strings.TrimSpace(condMatch[1])
			}
			if solMatch != nil {
				item.Solution = strings.TrimSpace(solMatch[1])
			}
			items = append(items, item)
		}
	}

	return items
}

func stripHTML(html string) string {
	// Убираем теги
	text := regexp.MustCompile(`<[^>]+>`).ReplaceAllString(html, " ")
	// Декодируем HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	// Убираем лишние пробелы
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

// =============================================================================
// Функции индексирования для быстрого поиска задач
// =============================================================================

// buildTaskIndex создаёт индексные данные для задачи
func buildTaskIndex(conditionText string) *TaskIndex {
	if conditionText == "" {
		return nil
	}

	normalized := normalizeForIndex(conditionText)
	hash := computeHash(normalized)
	numbers := extractNumbers(conditionText)
	keywords := extractKeywords(conditionText)

	return &TaskIndex{
		NormalizedHash:   hash,
		NumbersSignature: numbers,
		Keywords:         keywords,
		NormalizedText:   normalized,
	}
}

// normalizeForIndex нормализует текст для индексирования
func normalizeForIndex(text string) string {
	// Приводим к нижнему регистру
	text = strings.ToLower(text)

	// Нормализуем математические символы
	replacements := map[string]string{
		"×":  "*",
		"·":  "*",
		"÷":  ":",
		"−":  "-",
		"–":  "-",
		"—":  "-",
		"…":  "...",
		"«":  "\"",
		"»":  "\"",
		"\t": " ",
		"\n": " ",
		"\r": "",
	}
	for old, new := range replacements {
		text = strings.ReplaceAll(text, old, new)
	}

	// Убираем лишние пробелы вокруг знаков
	text = regexp.MustCompile(`\s*([+\-*/:=<>])\s*`).ReplaceAllString(text, "$1")

	// Убираем множественные пробелы
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Убираем пунктуацию в конце
	text = strings.TrimRight(text, ".,;:!? ")

	return strings.TrimSpace(text)
}

// computeHash вычисляет SHA256 хэш текста
func computeHash(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// extractNumbers извлекает числа из текста и возвращает сигнатуру
func extractNumbers(text string) string {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(text, -1)

	if len(matches) == 0 {
		return ""
	}

	// Убираем дубликаты и сортируем
	seen := make(map[string]bool)
	var unique []string
	for _, m := range matches {
		if !seen[m] {
			seen[m] = true
			unique = append(unique, m)
		}
	}

	// Сортируем числа
	sort.Slice(unique, func(i, j int) bool {
		ni, _ := strconv.Atoi(unique[i])
		nj, _ := strconv.Atoi(unique[j])
		return ni < nj
	})

	// Ограничиваем длину сигнатуры
	if len(unique) > 10 {
		unique = unique[:10]
	}

	return strings.Join(unique, ",")
}

// extractKeywords извлекает ключевые слова из текста
func extractKeywords(text string) []string {
	text = strings.ToLower(text)

	// Математические ключевые слова
	mathKeywords := []string{
		"сложи", "вычти", "умножь", "раздели", "вычисли",
		"найди", "реши", "сравни", "заполни",
		"сумма", "разность", "произведение", "частное",
		"слагаемое", "уменьшаемое", "вычитаемое",
		"множитель", "делитель", "делимое",
		"периметр", "площадь", "сторона",
		"столбик", "письменно", "уголком",
		"уравнение", "неизвестное",
		"больше", "меньше", "равно",
		"раз", "раза",
		"остаток",
		"чётн", "нечётн",
		"римск",
		"таблица умножения",
		"числовая прямая", "числовой луч",
	}

	var found []string
	for _, kw := range mathKeywords {
		if strings.Contains(text, kw) {
			found = append(found, kw)
		}
	}

	return found
}

func saveJSON(textbook Textbook, tasks []Task) {
	data := map[string]interface{}{
		"textbook": textbook,
		"tasks":    tasks,
	}

	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	filename := filepath.Join(config.OutputDir, fmt.Sprintf("peterson_3_part%d.json", config.Part))
	os.WriteFile(filename, jsonBytes, 0644)
	fmt.Printf("JSON сохранён: %s\n", filename)
}

func generateSQL(textbook Textbook, tasks []Task) {
	var sb strings.Builder

	// SQL для учебника (PostgreSQL)
	sb.WriteString("-- Данные учебника Петерсон 3 класс, часть " + strconv.Itoa(textbook.Part) + "\n")
	sb.WriteString("-- PostgreSQL compatible\n\n")

	// Вставка учебника с ON CONFLICT (PostgreSQL upsert)
	sb.WriteString(fmt.Sprintf(`INSERT INTO textbooks (subject, grade, authors, title, part, year, publisher, source_url)
VALUES ('%s', %d, '%s', '%s', %d, %d, '%s', '%s')
ON CONFLICT (subject, grade, authors, part) DO UPDATE SET year = EXCLUDED.year;

`,
		textbook.Subject, textbook.Grade, escapePG(textbook.Authors), escapePG(textbook.Title),
		textbook.Part, textbook.Year, escapePG(textbook.Publisher), escapePG(textbook.SourceURL)))

	// SQL для задач - используем DO $$ блок для работы с переменными
	sb.WriteString("-- Задачи\n")
	sb.WriteString("DO $$\n")
	sb.WriteString("DECLARE\n")
	sb.WriteString("    v_textbook_id BIGINT;\n")
	sb.WriteString("    v_task_id BIGINT;\n")
	sb.WriteString("BEGIN\n")
	sb.WriteString(fmt.Sprintf("    SELECT id INTO v_textbook_id FROM textbooks WHERE subject = '%s' AND grade = %d AND authors = '%s' AND part = %d;\n\n",
		textbook.Subject, textbook.Grade, escapePG(textbook.Authors), textbook.Part))

	for _, task := range tasks {
		subItemsJSON := "NULL"
		if task.HasSubItems && len(task.SubItems) > 0 {
			jsonBytes, _ := json.Marshal(task.SubItems)
			subItemsJSON = fmt.Sprintf("'%s'::jsonb", escapePG(string(jsonBytes)))
		}

		hasSubItems := "FALSE"
		if task.HasSubItems {
			hasSubItems = "TRUE"
		}

		sb.WriteString(fmt.Sprintf(`    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, %d, '%s', %d, '%s', '%s', '%s', '%s', '%s', '%s', %s, %s, '%s')
    RETURNING id INTO v_task_id;

`,
			task.PageNumber, escapePG(task.TaskNumber), task.TaskOrder,
			escapePG(task.ConditionText), escapePG(task.ConditionHTML),
			escapePG(task.SolutionText), escapePG(task.SolutionHTML),
			escapePG(task.HintsText), escapePG(task.HintsHTML),
			hasSubItems, subItemsJSON, escapePG(task.SourceURL)))

		// SQL для картинок
		for _, img := range task.Images {
			subItemLetter := "NULL"
			if img.SubItemLetter != "" {
				subItemLetter = fmt.Sprintf("'%s'", escapePG(img.SubItemLetter))
			}

			sb.WriteString(fmt.Sprintf(`    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, '%s', %d, %s, '%s', '%s', '%s');

`,
				img.ImageType, img.ImageOrder, subItemLetter,
				escapePG(img.OriginalURL), escapePG(img.LocalPath), escapePG(img.AltText)))
		}

		// SQL для индекса поиска
		if task.Index != nil {
			keywordsJSON := "NULL"
			if len(task.Index.Keywords) > 0 {
				jsonBytes, _ := json.Marshal(task.Index.Keywords)
				keywordsJSON = fmt.Sprintf("'%s'::jsonb", escapePG(string(jsonBytes)))
			}

			numbersSignature := "NULL"
			if task.Index.NumbersSignature != "" {
				numbersSignature = fmt.Sprintf("'%s'", escapePG(task.Index.NumbersSignature))
			}

			sb.WriteString(fmt.Sprintf(`    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, %d, '%s', %s, %s, '%s');

`,
				textbook.Grade,
				escapePG(task.Index.NormalizedHash),
				numbersSignature,
				keywordsJSON,
				escapePG(task.Index.NormalizedText)))
		}
	}

	sb.WriteString("END $$;\n")

	migrationNum := 20 + config.Part // Part 1 = 021, Part 2 = 022, Part 3 = 023
	filename := filepath.Join(config.OutputDir, fmt.Sprintf("%03d_peterson_3_part%d_data.up.sql", migrationNum, config.Part))
	os.WriteFile(filename, []byte(sb.String()), 0644)
	fmt.Printf("SQL сохранён: %s\n", filename)

	// Down миграция (PostgreSQL syntax - использует подзапросы вместо DELETE ... JOIN)
	downSQL := fmt.Sprintf(`-- Удаление данных части %d учебника Петерсон 3 класс
-- Удаление индексов (каскадно удалится при удалении задач, но для явности)
DELETE FROM textbook_task_index
WHERE textbook_id = (SELECT id FROM textbooks WHERE subject = 'matematika' AND grade = 3 AND authors = 'Петерсон Л.Г.' AND part = %d);

-- Удаление картинок (каскадно удалится при удалении задач, но для явности)
DELETE FROM textbook_task_images
WHERE task_id IN (
    SELECT t.id FROM textbook_tasks t
    INNER JOIN textbooks tb ON t.textbook_id = tb.id
    WHERE tb.subject = 'matematika' AND tb.grade = 3 AND tb.authors = 'Петерсон Л.Г.' AND tb.part = %d
);

-- Удаление задач
DELETE FROM textbook_tasks
WHERE textbook_id = (SELECT id FROM textbooks WHERE subject = 'matematika' AND grade = 3 AND authors = 'Петерсон Л.Г.' AND part = %d);

-- Удаление учебника
DELETE FROM textbooks WHERE subject = 'matematika' AND grade = 3 AND authors = 'Петерсон Л.Г.' AND part = %d;
`, config.Part, config.Part, config.Part, config.Part, config.Part)

	downFilename := filepath.Join(config.OutputDir, fmt.Sprintf("%03d_peterson_3_part%d_data.down.sql", migrationNum, config.Part))
	os.WriteFile(downFilename, []byte(downSQL), 0644)
	fmt.Printf("SQL down сохранён: %s\n", downFilename)
}

// escapePG экранирует строку для PostgreSQL (двойные кавычки для апострофа)
func escapePG(s string) string {
	// В PostgreSQL апостроф экранируется удвоением
	s = strings.ReplaceAll(s, "'", "''")
	// Обрабатываем переносы строк
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

func escape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

func downloadAllImages(tasks []Task) {
	var totalImages int
	for _, task := range tasks {
		totalImages += len(task.Images)
	}

	fmt.Printf("\nСкачивание изображений: %d шт.\n", totalImages)

	downloaded := 0
	failed := 0

	for _, task := range tasks {
		for _, img := range task.Images {
			err := downloadImage(img.OriginalURL, img.LocalPath)
			if err != nil {
				fmt.Printf("  Ошибка скачивания %s: %v\n", img.OriginalURL, err)
				failed++
			} else {
				downloaded++
			}

			// Небольшая задержка между запросами
			time.Sleep(100 * time.Millisecond)

			// Выводим прогресс каждые 50 изображений
			if (downloaded+failed)%50 == 0 {
				fmt.Printf("  Прогресс: %d/%d\n", downloaded+failed, totalImages)
			}
		}
	}

	fmt.Printf("Скачано: %d, ошибок: %d\n", downloaded, failed)
}

func downloadImage(url, localPath string) error {
	fullPath := filepath.Join(config.OutputDir, config.ImagesDir, localPath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
