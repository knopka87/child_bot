package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"child-bot/api/internal/v2/types"
)

// TemplateRegistry — корневая структура JSON-файла шаблона
type TemplateRegistry struct {
	Registry RegistryMeta               `json:"template_registry"`
	Profiles map[string]TemplateProfile `json:"template_profiles"`
}

type RegistryMeta struct {
	RegistryVersion string     `json:"registry_version"`
	Scope           string     `json:"scope"`
	Templates       []Template `json:"templates"`
}

type Template struct {
	TemplateCode       string            `json:"template_code"`
	TemplateID         string            `json:"template_id"`
	Title              string            `json:"title"`
	GradeMin           int64             `json:"grade_min"`
	GradeMax           int64             `json:"grade_max"`
	FormatsAllowed     []string          `json:"formats_allowed"`
	PedKeysDefaults    PedKeysDefaults   `json:"ped_keys_defaults"`
	Routing            RoutingConfig     `json:"routing"`
	HintPolicyDefaults HintPolicyDefault `json:"hint_policy_defaults"`
}

type PedKeysDefaults struct {
	TaskType    string   `json:"task_type"`
	Format      string   `json:"format"`
	Topic       string   `json:"topic"`
	Constraints []string `json:"constraints"`
}

type RoutingConfig struct {
	MatchKeys    MatchKeys     `json:"match_keys"`
	RoutingRules []RoutingRule `json:"routing_rules"`
	Confusables  []string      `json:"confusables"`
}

type MatchKeys struct {
	TaskType       string   `json:"task_type"`
	Topic          string   `json:"topic"`
	FormatsAllowed []string `json:"formats_allowed"`
}

type RoutingRule struct {
	RuleID          string       `json:"rule_id"`
	MustHave        RulePatterns `json:"must_have"`
	MustNot         RulePatterns `json:"must_not"`
	RoutingPriority int          `json:"routing_priority"`
}

type RulePatterns struct {
	TextPatternsAny        []string               `json:"text_patterns_any"`
	VisualKindsAny         []string               `json:"visual_kinds_any"`
	TemplateParamsRequired map[string]interface{} `json:"template_params_required"`
}

type HintPolicyDefault struct {
	DefaultVisible int    `json:"default_visible"`
	MaxHints       int    `json:"max_hints"`
	H3Reason       string `json:"h3_reason"`
}

// TemplateProfile — профиль для передачи в HINT (template_profile_core)
type TemplateProfile struct {
	HintStyleProfile   interface{}            `json:"hint_style_profile"` // can be string or object
	MaxHintsDefault    int                    `json:"max_hints_default"`
	AgeLanguage        AgeLanguage            `json:"age_language"`
	TeachingPattern    TeachingPattern        `json:"teaching_pattern"`
	CommonMistakes     []string               `json:"common_mistakes"`
	TerminologyRules   interface{}            `json:"terminology_rules"` // can be string array or map
	DisclosureDefaults map[string]interface{} `json:"disclosure_defaults"`
}

type AgeLanguage struct {
	GradeMin        int      `json:"grade_min"`
	GradeMax        int      `json:"grade_max"`
	Tone            string   `json:"tone"`
	ComplexityRules []string `json:"complexity_rules"`
}

type TeachingPattern struct {
	Goal string    `json:"goal"`
	L1   HintLevel `json:"l1"`
	L2   HintLevel `json:"l2"`
	L3   HintLevel `json:"l3"`
}

type HintLevel struct {
	Rules       []string `json:"rules"`
	Format      string   `json:"format"`
	Forbidden   []string `json:"forbidden"`
	WhenAllowed string   `json:"when_allowed,omitempty"`
}

// RoutingContext — контекст для роутинга
type RoutingContext struct {
	TextAll     string
	VisualKinds map[string]bool
	HasGap      bool
	TaskType    string
	Format      string
	Grade       int64
	Subject     types.Subject
}

// TemplateCandidate — кандидат с оценкой
type TemplateCandidate struct {
	Template       *Template
	Profile        *TemplateProfile
	Score          int
	MatchedRuleID  string
	AnchorsMatched int
	VisualMatched  bool
}

// RoutingTraceEntry — запись трассировки для одного шаблона/правила
type RoutingTraceEntry struct {
	TemplateCode    string   `json:"template_code"`
	RuleID          string   `json:"rule_id"`
	Status          string   `json:"status"` // "matched", "rejected_must_not", "rejected_must_have", "rejected_grade", "rejected_task_type"
	Score           int      `json:"score,omitempty"`
	AnchorsMatched  int      `json:"anchors_matched,omitempty"`
	VisualMatched   bool     `json:"visual_matched,omitempty"`
	RejectedBy      []string `json:"rejected_by,omitempty"` // какие паттерны отсекли
	MatchedPatterns []string `json:"matched_patterns,omitempty"`
}

// RoutingTrace — полная трассировка выбора шаблона
type RoutingTrace struct {
	TextAll        string              `json:"text_all"`
	VisualKinds    []string            `json:"visual_kinds"`
	TaskType       string              `json:"task_type"`
	Format         string              `json:"format"`
	Grade          int64               `json:"grade"`
	Entries        []RoutingTraceEntry `json:"entries"`
	Winner         string              `json:"winner,omitempty"`
	WinnerScore    int                 `json:"winner_score,omitempty"`
	WinnerRuleID   string              `json:"winner_rule_id,omitempty"`
	CandidateCount int                 `json:"candidate_count"`
}

// RoutingDebugEnabled — флаг для включения trace-логирования
var RoutingDebugEnabled = false

// lastRoutingTrace — последняя трассировка (для тестирования и отладки)
var lastRoutingTrace *RoutingTrace
var traceMutex sync.Mutex

// GetLastRoutingTrace возвращает последнюю трассировку (для тестирования)
func GetLastRoutingTrace() *RoutingTrace {
	traceMutex.Lock()
	defer traceMutex.Unlock()
	return lastRoutingTrace
}

// SetRoutingDebug включает/выключает trace-логирование
func SetRoutingDebug(enabled bool) {
	RoutingDebugEnabled = enabled
}

// FormatRoutingTrace форматирует трассировку в читаемый вид
func FormatRoutingTrace(trace *RoutingTrace) string {
	if trace == nil {
		return "No trace available"
	}

	var sb strings.Builder
	sb.WriteString("=== ROUTING TRACE ===\n")
	sb.WriteString(fmt.Sprintf("Text: %.100s...\n", trace.TextAll))
	sb.WriteString(fmt.Sprintf("VisualKinds: %v\n", trace.VisualKinds))
	sb.WriteString(fmt.Sprintf("TaskType: %s, Format: %s, Grade: %d\n", trace.TaskType, trace.Format, trace.Grade))
	sb.WriteString(fmt.Sprintf("Candidates found: %d\n", trace.CandidateCount))

	if trace.Winner != "" {
		sb.WriteString(fmt.Sprintf("WINNER: %s (score=%d, rule=%s)\n", trace.Winner, trace.WinnerScore, trace.WinnerRuleID))
	} else {
		sb.WriteString("WINNER: none\n")
	}

	sb.WriteString("\n--- Evaluation details ---\n")

	// Группируем по статусу для лучшей читаемости
	matched := []RoutingTraceEntry{}
	rejected := []RoutingTraceEntry{}

	for _, e := range trace.Entries {
		if e.Status == "matched" {
			matched = append(matched, e)
		} else {
			rejected = append(rejected, e)
		}
	}

	if len(matched) > 0 {
		sb.WriteString("\nMATCHED:\n")
		for _, e := range matched {
			sb.WriteString(fmt.Sprintf("  [%s] rule=%s score=%d anchors=%d visual=%v\n",
				e.TemplateCode, e.RuleID, e.Score, e.AnchorsMatched, e.VisualMatched))
			if len(e.MatchedPatterns) > 0 {
				sb.WriteString(fmt.Sprintf("    patterns: %v\n", e.MatchedPatterns))
			}
		}
	}

	if len(rejected) > 0 {
		sb.WriteString("\nREJECTED:\n")
		for _, e := range rejected {
			sb.WriteString(fmt.Sprintf("  [%s] rule=%s status=%s\n", e.TemplateCode, e.RuleID, e.Status))
			if len(e.RejectedBy) > 0 {
				sb.WriteString(fmt.Sprintf("    rejected_by: %v\n", e.RejectedBy))
			}
		}
	}

	sb.WriteString("=== END TRACE ===\n")
	return sb.String()
}

var (
	templatesCache     []TemplateRegistry
	templatesCacheOnce sync.Once
	templatesDir       = "api/internal/v2/templates"
)

// SetTemplatesDir sets the templates directory (for testing)
func SetTemplatesDir(dir string) {
	templatesDir = dir
}

// ResetTemplatesCache resets the templates cache (for testing)
func ResetTemplatesCache() {
	templatesCache = nil
	templatesCacheOnce = sync.Once{}
}

// loadTemplates загружает все шаблоны из папки templates
func loadTemplates() []TemplateRegistry {
	templatesCacheOnce.Do(func() {
		files, err := filepath.Glob(filepath.Join(templatesDir, "T*.json"))
		if err != nil {
			return
		}
		for _, f := range files {
			data, err := os.ReadFile(f)
			if err != nil {
				log.Printf("[template] failed to read %s: %v", f, err)
				continue
			}
			var reg TemplateRegistry
			if err := json.Unmarshal(data, &reg); err != nil {
				log.Printf("[template] failed to parse %s: %v", f, err)
				continue
			}
			templatesCache = append(templatesCache, reg)
		}
	})
	return templatesCache
}

// normalizeText нормализует текст для сравнения:
// - lower-case
// - ё → е (простая замена, без NFD чтобы сохранить й)
// - унификация математических символов (×, ·, * → *, ÷, : → :)
// - длинные/средние тире → -
// - убрать множественные пробелы
func normalizeText(s string) string {
	s = strings.ToLower(s)
	// ё → е (простая замена, NFD удаляет й поэтому не используем)
	result := strings.ReplaceAll(s, "ё", "е")

	// унификация математических символов умножения → *
	result = strings.ReplaceAll(result, "×", "*")
	result = strings.ReplaceAll(result, "·", "*")

	// унификация символов деления → :
	result = strings.ReplaceAll(result, "÷", ":")

	// длинные/средние тире → обычный минус
	result = strings.ReplaceAll(result, "—", "-") // длинное тире (em dash)
	result = strings.ReplaceAll(result, "–", "-") // среднее тире (en dash)
	result = strings.ReplaceAll(result, "−", "-") // математический минус

	// убрать множественные пробелы
	re := regexp.MustCompile(`\s+`)
	result = re.ReplaceAllString(result, " ")
	return strings.TrimSpace(result)
}

// buildRoutingContext строит контекст роутинга из ParseResponse
func buildRoutingContext(task types.ParseTask, items []types.ParseItem) RoutingContext {
	// Собираем весь текст
	var textParts []string
	textParts = append(textParts, task.TaskTextClean)
	for _, item := range items {
		textParts = append(textParts, item.ItemTextClean)
	}
	textAll := normalizeText(strings.Join(textParts, " "))

	// Собираем visual_kinds
	visualKinds := make(map[string]bool)
	for _, vf := range task.VisualFacts {
		if vf.Kind != "" {
			visualKinds[strings.ToLower(vf.Kind)] = true
		}
	}

	// Определяем task_type и format из items (mode — наиболее частый)
	taskTypeCount := make(map[string]int)
	formatCount := make(map[string]int)
	for _, item := range items {
		if item.PedKeys.TaskType != "" {
			taskTypeCount[item.PedKeys.TaskType]++
		}
		if item.PedKeys.Format != "" {
			formatCount[item.PedKeys.Format]++
		}
	}

	taskType := getMostFrequent(taskTypeCount)
	format := getMostFrequent(formatCount)

	// Нормализация taskType перед маппингом
	rawTaskType := taskType
	taskType = strings.ToLower(strings.TrimSpace(taskType))
	taskType = strings.Trim(taskType, ".:;,")

	// Маппинг нестандартных task_type от LLM к нашим шаблонам
	taskTypeMapping := map[string]string{
		"expressions": "arithmetic_fluency",
		"calculation": "arithmetic_fluency",
		"compute":     "arithmetic_fluency",
		"algebra":     "patterns_logic",
		"equation":    "patterns_logic",
		"equations":   "patterns_logic",
	}
	if mapped, ok := taskTypeMapping[taskType]; ok {
		if rawTaskType != taskType {
			log.Printf("[template] taskType normalized: %q -> %q -> %q", rawTaskType, taskType, mapped)
		}
		taskType = mapped
	}

	return RoutingContext{
		TextAll:     textAll,
		VisualKinds: visualKinds,
		HasGap:      strings.Contains(textAll, "__gap__"),
		TaskType:    taskType,
		Format:      format,
		Grade:       task.Grade,
		Subject:     task.Subject,
	}
}

func getMostFrequent(m map[string]int) string {
	maxCount := 0
	result := ""
	for k, v := range m {
		if v > maxCount {
			maxCount = v
			result = k
		}
	}
	return result
}

// cyrillicWordBoundary — паттерн для русских границ слова
// Используется вместо \b, который не работает корректно с кириллицей в Go/RE2
const cyrillicWordBoundaryStart = `(?:^|[^А-Яа-яЁёA-Za-z0-9])`
const cyrillicWordBoundaryEnd = `(?:[^А-Яа-яЁёA-Za-z0-9]|$)`

// wrapWithCyrillicBoundaries оборачивает слово в русские границы слова
func wrapWithCyrillicBoundaries(word string) string {
	return cyrillicWordBoundaryStart + regexp.QuoteMeta(word) + cyrillicWordBoundaryEnd
}

// matchCyrillicWord проверяет наличие слова с учётом русских границ
func matchCyrillicWord(text, word string) bool {
	pattern := wrapWithCyrillicBoundaries(word)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return strings.Contains(text, word)
	}
	return re.MatchString(text)
}

// isRegexPattern проверяет, содержит ли паттерн regex-метасимволы
func isRegexPattern(pattern string) bool {
	// Проверяем наличие типичных regex-конструкций
	return strings.Contains(pattern, ".*") ||
		strings.Contains(pattern, ".+") ||
		strings.Contains(pattern, "\\d") ||
		strings.Contains(pattern, "\\w") ||
		strings.Contains(pattern, "[") ||
		strings.Contains(pattern, "(") ||
		strings.Contains(pattern, "?") ||
		strings.Contains(pattern, "+") ||
		strings.Contains(pattern, "|")
}

// matchPatternFlexible проверяет паттерн с допуском расстояния между частями.
// Поддерживает три режима:
// 1. Regex-паттерны (если содержат .*, .+, \d и т.д.)
// 2. Точное совпадение подстроки
// 3. Proximity search для паттернов с 3+ словами
func matchPatternFlexible(text, pattern string, maxDistance int) bool {
	normalizedPattern := normalizeText(pattern)

	// 1. Проверяем, является ли паттерн regex
	if isRegexPattern(pattern) {
		// Нормализуем паттерн, но сохраняем regex-конструкции
		regexPattern := strings.ToLower(pattern)
		regexPattern = strings.ReplaceAll(regexPattern, "ё", "е")
		re, err := regexp.Compile(regexPattern)
		if err != nil {
			// Если regex невалиден, пробуем как обычную строку
			return strings.Contains(text, normalizedPattern)
		}
		return re.MatchString(text)
	}

	// 2. Пробуем точное совпадение подстроки
	if strings.Contains(text, normalizedPattern) {
		return true
	}

	// 3. Разбиваем паттерн на части и ищем с proximity
	words := strings.Fields(normalizedPattern)

	// Proximity search только для паттернов с 3+ словами (избегаем false positives для коротких)
	if len(words) < 3 {
		return false
	}

	// Ищем первые два слова как устойчивый якорь
	anchor := words[0] + " " + words[1]
	idx := strings.Index(text, anchor)
	if idx == -1 {
		return false
	}

	// Проверяем, что остальные слова находятся в пределах maxDistance от якоря
	windowEnd := idx + len(anchor) + maxDistance
	if windowEnd > len(text) {
		windowEnd = len(text)
	}
	window := text[idx:windowEnd]

	for _, word := range words[2:] {
		if !strings.Contains(window, word) {
			return false
		}
	}

	return true
}

// checkMustHave проверяет, выполняются ли must_have условия
func checkMustHave(ctx RoutingContext, patterns RulePatterns) (bool, int, bool) {
	matched, anchorsMatched, visualMatched, _ := checkMustHaveWithTrace(ctx, patterns)
	return matched, anchorsMatched, visualMatched
}

// checkMustNot проверяет, нарушены ли must_not условия
func checkMustNot(ctx RoutingContext, patterns RulePatterns) bool {
	rejected, _ := checkMustNotWithTrace(ctx, patterns)
	return rejected
}

// scoreCandidate вычисляет score для кандидата
func scoreCandidate(ctx RoutingContext, tmpl *Template, rule *RoutingRule, anchorsMatched int, visualMatched bool) int {
	score := 0

	// +50 за совпадение visual_kinds
	if visualMatched {
		score += 50
	}

	// +30 за каждый текстовый якорь (max 3)
	if anchorsMatched > 3 {
		anchorsMatched = 3
	}
	score += anchorsMatched * 30

	// +10 за совпадение format
	for _, f := range tmpl.FormatsAllowed {
		if f == ctx.Format {
			score += 10
			break
		}
	}

	// +routing_priority (нормализованный)
	score += rule.RoutingPriority / 10

	// Бонус за специфичность шаблона (узкий диапазон классов)
	// Чем уже диапазон, тем выше бонус
	// Это важно для выбора специализированных шаблонов (T45-T50 для 1 класса)
	// перед общими шаблонами (T1-T14 для классов 1-4)
	// Формула: 80 / gradeRange, так что:
	// - range=1 (один класс): +80
	// - range=2: +40
	// - range=4: +20
	gradeRange := tmpl.GradeMax - tmpl.GradeMin + 1
	if gradeRange > 0 {
		specificityBonus := 80 / int(gradeRange)
		score += specificityBonus
	}

	return score
}

// selectTemplate выбирает лучший шаблон по алгоритму из ТЗ
// Двухпроходный поиск: сначала с точным совпадением task_type, затем fallback по паттернам
func selectTemplate(ctx RoutingContext) (*TemplateCandidate, bool) {
	// Инициализируем trace если включен debug
	var trace *RoutingTrace
	if RoutingDebugEnabled {
		visualKindsList := make([]string, 0, len(ctx.VisualKinds))
		for k := range ctx.VisualKinds {
			visualKindsList = append(visualKindsList, k)
		}
		trace = &RoutingTrace{
			TextAll:     ctx.TextAll,
			VisualKinds: visualKindsList,
			TaskType:    ctx.TaskType,
			Format:      ctx.Format,
			Grade:       ctx.Grade,
			Entries:     []RoutingTraceEntry{},
		}
	}

	// Только для math
	if ctx.Subject != types.SubjectMath {
		if trace != nil {
			traceMutex.Lock()
			lastRoutingTrace = trace
			traceMutex.Unlock()
		}
		return nil, false
	}

	registries := loadTemplates()

	// Первый проход: точное совпадение task_type
	candidates := findCandidatesWithTrace(ctx, registries, true, trace)

	// Fallback: если кандидатов нет, ищем без учёта task_type
	if len(candidates) == 0 {
		candidates = findCandidatesWithTrace(ctx, registries, false, trace)
	}

	// Fallback 2: если всё ещё нет кандидатов, проверяем общие арифметические паттерны
	if len(candidates) == 0 {
		if fallback := tryArithmeticFallback(ctx, registries); fallback != nil {
			candidates = append(candidates, *fallback)
		}
	}

	if len(candidates) == 0 {
		if trace != nil {
			trace.CandidateCount = 0
			traceMutex.Lock()
			lastRoutingTrace = trace
			traceMutex.Unlock()
		}
		return nil, false
	}

	// Сортировка по tie-break правилам
	best := candidates[0]
	for _, c := range candidates[1:] {
		if compareCandidates(c, best) > 0 {
			best = c
		}
	}

	// Сохраняем trace
	if trace != nil {
		trace.CandidateCount = len(candidates)
		trace.Winner = best.Template.TemplateCode
		trace.WinnerScore = best.Score
		trace.WinnerRuleID = best.MatchedRuleID
		traceMutex.Lock()
		lastRoutingTrace = trace
		traceMutex.Unlock()
	}

	return &best, true
}

// checkMustNotWithTrace проверяет must_not и возвращает список отсекающих паттернов
func checkMustNotWithTrace(ctx RoutingContext, patterns RulePatterns) (bool, []string) {
	var rejectedBy []string

	// Проверяем text_patterns_any
	for _, pattern := range patterns.TextPatternsAny {
		normalizedPattern := normalizeText(pattern)

		// Если паттерн содержит regex-конструкции, используем regex
		if isRegexPattern(pattern) {
			re, err := regexp.Compile(strings.ToLower(pattern))
			if err == nil && re.MatchString(ctx.TextAll) {
				rejectedBy = append(rejectedBy, "text:"+pattern)
			}
			continue
		}

		// Для коротких паттернов (1-2 слова) используем русские границы слова
		words := strings.Fields(normalizedPattern)
		if len(words) <= 2 && len(normalizedPattern) >= 3 {
			if matchCyrillicWord(ctx.TextAll, normalizedPattern) {
				rejectedBy = append(rejectedBy, "text:"+pattern)
			}
		} else if strings.Contains(ctx.TextAll, normalizedPattern) {
			rejectedBy = append(rejectedBy, "text:"+pattern)
		}
	}

	// Проверяем visual_kinds_any
	for _, kind := range patterns.VisualKindsAny {
		if ctx.VisualKinds[strings.ToLower(kind)] {
			rejectedBy = append(rejectedBy, "visual:"+kind)
		}
	}

	return len(rejectedBy) > 0, rejectedBy
}

// checkMustHaveWithTrace проверяет must_have и возвращает список совпавших паттернов
func checkMustHaveWithTrace(ctx RoutingContext, patterns RulePatterns) (bool, int, bool, []string) {
	anchorsMatched := 0
	visualMatched := false
	var matchedPatterns []string

	// Проверяем text_patterns_any (OR) с гибким поиском
	textMatched := len(patterns.TextPatternsAny) == 0
	for _, pattern := range patterns.TextPatternsAny {
		// Используем гибкий поиск с окном 100 символов
		if matchPatternFlexible(ctx.TextAll, pattern, 100) {
			textMatched = true
			anchorsMatched++
			matchedPatterns = append(matchedPatterns, "text:"+pattern)
		}
	}

	// Проверяем visual_kinds_any (OR)
	visualKindsMatched := len(patterns.VisualKindsAny) == 0
	for _, kind := range patterns.VisualKindsAny {
		if ctx.VisualKinds[strings.ToLower(kind)] {
			visualKindsMatched = true
			visualMatched = true
			matchedPatterns = append(matchedPatterns, "visual:"+kind)
			break
		}
	}

	return textMatched && visualKindsMatched, anchorsMatched, visualMatched, matchedPatterns
}

// findCandidatesWithTrace ищет кандидатов с опциональной трассировкой
func findCandidatesWithTrace(ctx RoutingContext, registries []TemplateRegistry, strictTaskType bool, trace *RoutingTrace) []TemplateCandidate {
	var candidates []TemplateCandidate

	for i := range registries {
		reg := &registries[i]
		for j := range reg.Registry.Templates {
			tmpl := &reg.Registry.Templates[j]

			// Проверяем grade
			if ctx.Grade > 0 && (ctx.Grade < tmpl.GradeMin || ctx.Grade > tmpl.GradeMax) {
				if trace != nil {
					trace.Entries = append(trace.Entries, RoutingTraceEntry{
						TemplateCode: tmpl.TemplateCode,
						RuleID:       "",
						Status:       "rejected_grade",
						RejectedBy:   []string{fmt.Sprintf("grade %d not in [%d, %d]", ctx.Grade, tmpl.GradeMin, tmpl.GradeMax)},
					})
				}
				continue
			}

			// Проверяем match_keys (task_type) — только в strict режиме
			taskTypeMatched := true
			if strictTaskType && ctx.TaskType != "" && tmpl.Routing.MatchKeys.TaskType != "" {
				if ctx.TaskType != tmpl.Routing.MatchKeys.TaskType {
					if trace != nil {
						trace.Entries = append(trace.Entries, RoutingTraceEntry{
							TemplateCode: tmpl.TemplateCode,
							RuleID:       "",
							Status:       "rejected_task_type",
							RejectedBy:   []string{fmt.Sprintf("task_type '%s' != '%s'", ctx.TaskType, tmpl.Routing.MatchKeys.TaskType)},
						})
					}
					continue
				}
			} else if ctx.TaskType != "" && tmpl.Routing.MatchKeys.TaskType != "" {
				taskTypeMatched = ctx.TaskType == tmpl.Routing.MatchKeys.TaskType
			}

			// Проверяем routing_rules (OR)
			ruleMatched := false
			for k := range tmpl.Routing.RoutingRules {
				rule := &tmpl.Routing.RoutingRules[k]

				// Проверяем must_not
				rejected, rejectedBy := checkMustNotWithTrace(ctx, rule.MustNot)
				if rejected {
					if trace != nil {
						trace.Entries = append(trace.Entries, RoutingTraceEntry{
							TemplateCode: tmpl.TemplateCode,
							RuleID:       rule.RuleID,
							Status:       "rejected_must_not",
							RejectedBy:   rejectedBy,
						})
					}
					continue
				}

				// Проверяем must_have
				matched, anchorsMatched, visualMatched, matchedPatterns := checkMustHaveWithTrace(ctx, rule.MustHave)
				if !matched {
					if trace != nil {
						trace.Entries = append(trace.Entries, RoutingTraceEntry{
							TemplateCode: tmpl.TemplateCode,
							RuleID:       rule.RuleID,
							Status:       "rejected_must_have",
						})
					}
					continue
				}

				// Вычисляем score
				score := scoreCandidate(ctx, tmpl, rule, anchorsMatched, visualMatched)

				// Бонус за совпадение task_type
				if taskTypeMatched {
					score += 20
				}

				// Получаем profile
				var profile *TemplateProfile
				if p, ok := reg.Profiles[tmpl.TemplateID]; ok {
					profile = &p
				}

				candidates = append(candidates, TemplateCandidate{
					Template:       tmpl,
					Profile:        profile,
					Score:          score,
					MatchedRuleID:  rule.RuleID,
					AnchorsMatched: anchorsMatched,
					VisualMatched:  visualMatched,
				})

				if trace != nil {
					trace.Entries = append(trace.Entries, RoutingTraceEntry{
						TemplateCode:    tmpl.TemplateCode,
						RuleID:          rule.RuleID,
						Status:          "matched",
						Score:           score,
						AnchorsMatched:  anchorsMatched,
						VisualMatched:   visualMatched,
						MatchedPatterns: matchedPatterns,
					})
				}

				ruleMatched = true
				break // Достаточно одного правила
			}

			// Если ни одно правило не подошло и нет записей в trace, добавим общую запись
			if !ruleMatched && trace != nil && len(tmpl.Routing.RoutingRules) == 0 {
				trace.Entries = append(trace.Entries, RoutingTraceEntry{
					TemplateCode: tmpl.TemplateCode,
					Status:       "no_rules",
				})
			}
		}
	}

	return candidates
}

// compareCandidates сравнивает кандидатов по tie-break правилам
// Возвращает >0 если a лучше b, <0 если b лучше a, 0 если равны
func compareCandidates(a, b TemplateCandidate) int {
	// 1. visual_kinds_any совпадение
	if a.VisualMatched && !b.VisualMatched {
		return 1
	}
	if !a.VisualMatched && b.VisualMatched {
		return -1
	}

	// 2. больше сильных якорей
	if a.AnchorsMatched != b.AnchorsMatched {
		return a.AnchorsMatched - b.AnchorsMatched
	}

	// 3. выше score (включает routing_priority)
	if a.Score != b.Score {
		return a.Score - b.Score
	}

	// 4. стабильный порядок по template_code
	return strings.Compare(a.Template.TemplateCode, b.Template.TemplateCode)
}

// getTemplate выбирает шаблон и возвращает template_profile_core как JSON
func getTemplate(task types.ParseTask, items []types.ParseItem) string {
	ctx := buildRoutingContext(task, items)

	candidate, found := selectTemplate(ctx)
	if !found || candidate.Profile == nil {
		return ""
	}

	// Формируем template_profile_core (только необходимое для HINT)
	profileCore := map[string]interface{}{
		"template_id":         candidate.Template.TemplateID,
		"max_hints_default":   candidate.Profile.MaxHintsDefault,
		"age_language":        candidate.Profile.AgeLanguage,
		"teaching_pattern":    candidate.Profile.TeachingPattern,
		"common_mistakes":     candidate.Profile.CommonMistakes,
		"disclosure_defaults": candidate.Profile.DisclosureDefaults,
	}

	js, err := json.Marshal(profileCore)
	if err != nil {
		return ""
	}

	return string(js)
}

// TemplateRoutingResult содержит результат роутинга шаблона
type TemplateRoutingResult struct {
	TemplateID string
	Found      bool
	DebugInfo  map[string]interface{}
}

// getTemplateIDWithDebug возвращает ID шаблона и debug-информацию
func getTemplateIDWithDebug(task types.ParseTask, items []types.ParseItem) TemplateRoutingResult {
	ctx := buildRoutingContext(task, items)

	candidate, found := selectTemplate(ctx)
	if !found {
		return TemplateRoutingResult{
			TemplateID: "",
			Found:      false,
			DebugInfo: map[string]interface{}{
				"reason":    "no_template_found",
				"subject":   ctx.Subject,
				"task_type": ctx.TaskType,
				"format":    ctx.Format,
				"grade":     ctx.Grade,
				"text_preview": func() string {
					if len(ctx.TextAll) > 100 {
						return ctx.TextAll[:100] + "..."
					}
					return ctx.TextAll
				}(),
			},
		}
	}

	return TemplateRoutingResult{
		TemplateID: candidate.Template.TemplateID,
		Found:      true,
		DebugInfo: map[string]interface{}{
			"template_code": candidate.Template.TemplateCode,
			"matched_rule":  candidate.MatchedRuleID,
			"score":         candidate.Score,
		},
	}
}

// getTemplateID возвращает только ID выбранного шаблона (для обратной совместимости)
func getTemplateID(task types.ParseTask, items []types.ParseItem) string {
	result := getTemplateIDWithDebug(task, items)
	return result.TemplateID
}

// tryArithmeticFallback проверяет общие арифметические паттерны и возвращает fallback шаблон
func tryArithmeticFallback(ctx RoutingContext, registries []TemplateRegistry) *TemplateCandidate {
	// Паттерны, указывающие на арифметическую задачу
	arithmeticPatterns := []string{
		"вычисли",
		"посчитай",
		"найди значение",
		"реши пример",
		"сколько будет",
		"выполни действ",
		"\\d+\\s*[+\\-×·\\*:]\\s*\\d+", // числа с операциями
		"\\d+\\s*\\+\\s*\\d+",          // сложение
		"\\d+\\s*-\\s*\\d+",            // вычитание
		"\\d+\\s*[×·\\*]\\s*\\d+",      // умножение
		"\\d+\\s*[:/÷]\\s*\\d+",        // деление
		"сравни.*\\d+",                 // сравнение чисел
		"больше|меньше|равно",
		"сложи|вычти|умнож|раздели",
		"сумм|разност|произведен|частно",
	}

	matched := false
	for _, pattern := range arithmeticPatterns {
		re, err := regexp.Compile("(?i)" + pattern)
		if err == nil && re.MatchString(ctx.TextAll) {
			matched = true
			break
		}
	}

	if !matched {
		return nil
	}

	// Ищем T35 (порядок действий) или T11 (свойства) как fallback
	fallbackTemplates := []string{"T35", "T11", "T8"}

	for _, tmplCode := range fallbackTemplates {
		for i := range registries {
			reg := &registries[i]
			for j := range reg.Registry.Templates {
				tmpl := &reg.Registry.Templates[j]
				if tmpl.TemplateCode == tmplCode {
					// Найден fallback шаблон
					var profile *TemplateProfile
					if reg.Profiles != nil {
						if p, ok := reg.Profiles[tmpl.TemplateID]; ok {
							profile = &p
						}
					}
					return &TemplateCandidate{
						Template:      tmpl,
						Profile:       profile,
						Score:         10, // низкий score для fallback
						MatchedRuleID: "ARITHMETIC_FALLBACK",
					}
				}
			}
		}
	}

	return nil
}
