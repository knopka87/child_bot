package telegram

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

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
	HintStyleProfile   string            `json:"hint_style_profile"`
	MaxHintsDefault    int               `json:"max_hints_default"`
	AgeLanguage        AgeLanguage       `json:"age_language"`
	TeachingPattern    TeachingPattern   `json:"teaching_pattern"`
	CommonMistakes     []string          `json:"common_mistakes"`
	TerminologyRules   []string          `json:"terminology_rules"`
	DisclosureDefaults map[string]string `json:"disclosure_defaults"`
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

var (
	templatesCache     []TemplateRegistry
	templatesCacheOnce sync.Once
	templatesDir       = "api/internal/v2/templates"
)

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
				continue
			}
			var reg TemplateRegistry
			if err := json.Unmarshal(data, &reg); err != nil {
				continue
			}
			templatesCache = append(templatesCache, reg)
		}
	})
	return templatesCache
}

// normalizeText нормализует текст для сравнения: lower-case, ё→е, убрать лишние пробелы
func normalizeText(s string) string {
	s = strings.ToLower(s)
	// ё → е
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	result = strings.ReplaceAll(result, "ё", "е")
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

// checkMustHave проверяет, выполняются ли must_have условия
func checkMustHave(ctx RoutingContext, patterns RulePatterns) (bool, int, bool) {
	anchorsMatched := 0
	visualMatched := false

	// Проверяем text_patterns_any (OR)
	textMatched := len(patterns.TextPatternsAny) == 0
	for _, pattern := range patterns.TextPatternsAny {
		if strings.Contains(ctx.TextAll, normalizeText(pattern)) {
			textMatched = true
			anchorsMatched++
		}
	}

	// Проверяем visual_kinds_any (OR)
	visualKindsMatched := len(patterns.VisualKindsAny) == 0
	for _, kind := range patterns.VisualKindsAny {
		if ctx.VisualKinds[strings.ToLower(kind)] {
			visualKindsMatched = true
			visualMatched = true
			break
		}
	}

	return textMatched && visualKindsMatched, anchorsMatched, visualMatched
}

// checkMustNot проверяет, нарушены ли must_not условия
func checkMustNot(ctx RoutingContext, patterns RulePatterns) bool {
	// Проверяем text_patterns_any
	for _, pattern := range patterns.TextPatternsAny {
		if strings.Contains(ctx.TextAll, normalizeText(pattern)) {
			return true // нарушено
		}
	}

	// Проверяем visual_kinds_any
	for _, kind := range patterns.VisualKindsAny {
		if ctx.VisualKinds[strings.ToLower(kind)] {
			return true // нарушено
		}
	}

	return false
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

	return score
}

// selectTemplate выбирает лучший шаблон по алгоритму из ТЗ
func selectTemplate(ctx RoutingContext) (*TemplateCandidate, bool) {
	// Только для math
	if ctx.Subject != types.SubjectMath {
		return nil, false
	}

	registries := loadTemplates()
	var candidates []TemplateCandidate

	for i := range registries {
		reg := &registries[i]
		for j := range reg.Registry.Templates {
			tmpl := &reg.Registry.Templates[j]

			// Проверяем grade
			if ctx.Grade > 0 && (ctx.Grade < tmpl.GradeMin || ctx.Grade > tmpl.GradeMax) {
				continue
			}

			// Проверяем match_keys (task_type)
			if ctx.TaskType != "" && tmpl.Routing.MatchKeys.TaskType != "" {
				if ctx.TaskType != tmpl.Routing.MatchKeys.TaskType {
					continue
				}
			}

			// Проверяем routing_rules (OR)
			for k := range tmpl.Routing.RoutingRules {
				rule := &tmpl.Routing.RoutingRules[k]

				// Проверяем must_not
				if checkMustNot(ctx, rule.MustNot) {
					continue
				}

				// Проверяем must_have
				matched, anchorsMatched, visualMatched := checkMustHave(ctx, rule.MustHave)
				if !matched {
					continue
				}

				// Вычисляем score
				score := scoreCandidate(ctx, tmpl, rule, anchorsMatched, visualMatched)

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
				break // Достаточно одного правила
			}
		}
	}

	if len(candidates) == 0 {
		return nil, false
	}

	// Сортировка по tie-break правилам
	best := candidates[0]
	for _, c := range candidates[1:] {
		if compareCandidates(c, best) > 0 {
			best = c
		}
	}

	return &best, true
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
