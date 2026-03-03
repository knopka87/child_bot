package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// TextbookTaskMatch представляет найденную задачу с решением
type TextbookTaskMatch struct {
	TaskID        int64   `json:"task_id"`
	TextbookID    int64   `json:"textbook_id"`
	PageNumber    int     `json:"page_number"`
	TaskNumber    string  `json:"task_number"`
	ConditionText string  `json:"condition_text"`
	SolutionText  string  `json:"solution_text"`
	HintsText     string  `json:"hints_text"`
	MatchScore    float64 `json:"match_score"`
	MatchMethod   string  `json:"match_method"` // "exact_hash", "numbers_fulltext", "fulltext"
}

// TextbookSearchParams параметры поиска
type TextbookSearchParams struct {
	ConditionText string // Текст условия задачи (обязательный)
	Grade         int    // Класс (опционально, 0 = любой)
	TextbookID    int64  // ID учебника (опционально, 0 = любой)
	MaxResults    int    // Максимум результатов (по умолчанию 5)
}

// FindMatchingTask ищет задачу в базе по условию с использованием каскадного поиска
// Порядок поиска: точный хэш → числа + полнотекстовый → только полнотекстовый
func (s *Store) FindMatchingTask(ctx context.Context, params TextbookSearchParams) ([]TextbookTaskMatch, error) {
	if params.MaxResults == 0 {
		params.MaxResults = 5
	}

	// Нормализуем входной текст
	normalized := normalizeForSearch(params.ConditionText)
	hash := computeSearchHash(normalized)
	numbers := extractSearchNumbers(params.ConditionText)

	// 1. Точный поиск по хэшу
	results, err := s.searchByHash(ctx, hash, params)
	if err != nil {
		return nil, err
	}
	if len(results) > 0 {
		return results, nil
	}

	// 2. Поиск по числам + полнотекстовый
	if numbers != "" {
		results, err = s.searchByNumbersAndFulltext(ctx, numbers, normalized, params)
		if err != nil {
			return nil, err
		}
		if len(results) > 0 {
			return results, nil
		}
	}

	// 3. Только полнотекстовый поиск
	results, err = s.searchFulltext(ctx, normalized, params)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// searchByHash ищет точное совпадение по нормализованному хэшу
func (s *Store) searchByHash(ctx context.Context, hash string, params TextbookSearchParams) ([]TextbookTaskMatch, error) {
	query := `
		SELECT
			t.id, t.textbook_id, t.page_number, t.task_number,
			t.condition_text, t.solution_text, t.hints_text
		FROM textbook_task_index idx
		JOIN textbook_tasks t ON t.id = idx.task_id
		WHERE idx.normalized_hash = $1
	`
	args := []interface{}{hash}
	argNum := 2

	if params.Grade > 0 {
		query += " AND idx.grade = $" + strconv.Itoa(argNum)
		args = append(args, params.Grade)
		argNum++
	}
	if params.TextbookID > 0 {
		query += " AND idx.textbook_id = $" + strconv.Itoa(argNum)
		args = append(args, params.TextbookID)
	}
	query += " LIMIT $" + strconv.Itoa(argNum)
	args = append(args, params.MaxResults)

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TextbookTaskMatch
	for rows.Next() {
		var m TextbookTaskMatch
		var conditionText, solutionText, hintsText sql.NullString
		if err := rows.Scan(&m.TaskID, &m.TextbookID, &m.PageNumber, &m.TaskNumber,
			&conditionText, &solutionText, &hintsText); err != nil {
			return nil, err
		}
		m.ConditionText = conditionText.String
		m.SolutionText = solutionText.String
		m.HintsText = hintsText.String
		m.MatchScore = 1.0 // Точное совпадение
		m.MatchMethod = "exact_hash"
		results = append(results, m)
	}
	return results, rows.Err()
}

// searchByNumbersAndFulltext ищет по сигнатуре чисел + полнотекстовый поиск
func (s *Store) searchByNumbersAndFulltext(ctx context.Context, numbers, normalizedText string, params TextbookSearchParams) ([]TextbookTaskMatch, error) {
	// Подготавливаем tsquery - разбиваем на слова и соединяем через &
	words := strings.Fields(normalizedText)
	if len(words) > 10 {
		words = words[:10] // Ограничиваем количество слов
	}
	tsQuery := strings.Join(words, " & ")

	query := `
		SELECT
			t.id, t.textbook_id, t.page_number, t.task_number,
			t.condition_text, t.solution_text, t.hints_text,
			ts_rank(idx.search_vector, plainto_tsquery('russian', $1)) AS rank
		FROM textbook_task_index idx
		JOIN textbook_tasks t ON t.id = idx.task_id
		WHERE idx.numbers_signature = $2
		  AND idx.search_vector @@ plainto_tsquery('russian', $1)
	`
	args := []interface{}{tsQuery, numbers}
	argNum := 3

	if params.Grade > 0 {
		query += " AND idx.grade = $" + strconv.Itoa(argNum)
		args = append(args, params.Grade)
		argNum++
	}
	if params.TextbookID > 0 {
		query += " AND idx.textbook_id = $" + strconv.Itoa(argNum)
		args = append(args, params.TextbookID)
		argNum++
	}

	query += " ORDER BY rank DESC LIMIT $" + strconv.Itoa(argNum)
	args = append(args, params.MaxResults)

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TextbookTaskMatch
	for rows.Next() {
		var m TextbookTaskMatch
		var conditionText, solutionText, hintsText sql.NullString
		var rank float64
		if err := rows.Scan(&m.TaskID, &m.TextbookID, &m.PageNumber, &m.TaskNumber,
			&conditionText, &solutionText, &hintsText, &rank); err != nil {
			return nil, err
		}
		m.ConditionText = conditionText.String
		m.SolutionText = solutionText.String
		m.HintsText = hintsText.String
		m.MatchScore = rank
		m.MatchMethod = "numbers_fulltext"
		results = append(results, m)
	}
	return results, rows.Err()
}

// searchFulltext ищет только по полнотекстовому индексу
func (s *Store) searchFulltext(ctx context.Context, normalizedText string, params TextbookSearchParams) ([]TextbookTaskMatch, error) {
	// Подготавливаем tsquery
	words := strings.Fields(normalizedText)
	if len(words) > 10 {
		words = words[:10]
	}
	tsQuery := strings.Join(words, " & ")

	query := `
		SELECT
			t.id, t.textbook_id, t.page_number, t.task_number,
			t.condition_text, t.solution_text, t.hints_text,
			ts_rank(idx.search_vector, plainto_tsquery('russian', $1)) AS rank
		FROM textbook_task_index idx
		JOIN textbook_tasks t ON t.id = idx.task_id
		WHERE idx.search_vector @@ plainto_tsquery('russian', $1)
	`
	args := []interface{}{tsQuery}
	argNum := 2

	if params.Grade > 0 {
		query += " AND idx.grade = $" + strconv.Itoa(argNum)
		args = append(args, params.Grade)
		argNum++
	}
	if params.TextbookID > 0 {
		query += " AND idx.textbook_id = $" + strconv.Itoa(argNum)
		args = append(args, params.TextbookID)
		argNum++
	}

	query += " ORDER BY rank DESC LIMIT $" + strconv.Itoa(argNum)
	args = append(args, params.MaxResults)

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TextbookTaskMatch
	for rows.Next() {
		var m TextbookTaskMatch
		var conditionText, solutionText, hintsText sql.NullString
		var rank float64
		if err := rows.Scan(&m.TaskID, &m.TextbookID, &m.PageNumber, &m.TaskNumber,
			&conditionText, &solutionText, &hintsText, &rank); err != nil {
			return nil, err
		}
		m.ConditionText = conditionText.String
		m.SolutionText = solutionText.String
		m.HintsText = hintsText.String
		m.MatchScore = rank
		m.MatchMethod = "fulltext"
		results = append(results, m)
	}
	return results, rows.Err()
}

// =============================================================================
// Вспомогательные функции нормализации (дублируют логику парсера)
// =============================================================================

// normalizeForSearch нормализует текст для поиска
func normalizeForSearch(text string) string {
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
	re := regexp.MustCompile(`\s*([+\-*/:=<>])\s*`)
	text = re.ReplaceAllString(text, "$1")

	// Убираем множественные пробелы
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	// Убираем пунктуацию в конце
	text = strings.TrimRight(text, ".,;:!? ")

	return strings.TrimSpace(text)
}

// computeSearchHash вычисляет SHA256 хэш для поиска
func computeSearchHash(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// extractSearchNumbers извлекает сигнатуру чисел из текста
func extractSearchNumbers(text string) string {
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

	// Ограничиваем длину
	if len(unique) > 10 {
		unique = unique[:10]
	}

	return strings.Join(unique, ",")
}
