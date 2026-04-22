package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"time"

	"child-bot/api/internal/store"

	"github.com/google/uuid"
)

// ReportService генерирует отчёты для родителей
type ReportService struct {
	store *store.Store
}

// NewReportService создаёт новый ReportService
func NewReportService(store *store.Store) *ReportService {
	return &ReportService{store: store}
}

// GetStore возвращает store для доступа к БД
func (s *ReportService) GetStore() *store.Store {
	return s.store
}

// WeeklyReportData данные для еженедельного отчёта
type WeeklyReportData struct {
	ChildName          string
	AvatarEmoji        string
	Grade              int
	Level              int
	XPTotal            int
	XPForNext          int
	CoinsBalance       int
	StreakDays         int
	TotalAttempts      int
	SuccessfulAttempts int
	AccuracyPercent    float64
	HintsUsed          int
	WeekAvgTimeMinutes float64
	NewAchievements    []AchievementData
	VillainBattles     []VillainBattleData
	AttemptsChange     *int
	AccuracyChange     *float64
	TimeChange         *int
	ReportWeekStart    time.Time
	ReportWeekEnd      time.Time
}

type AchievementData struct {
	Title            string
	Icon             string
	Count            int
	CurrentProgress  int
	RequirementValue int
	UnlockedAt       time.Time
}

type VillainBattleData struct {
	Name              string
	Icon              string
	CurrentHP         int
	MaxHP             int
	HPPercent         float64
	TotalDamageDealt  int
	CorrectTasksCount int
	Status            string
}

type WeeklyReport struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	ReportDate  time.Time  `json:"report_date"`
	HTMLContent string     `json:"html_content"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AchievementInfo информация о достижении
type AchievementInfo struct {
	Title      string
	Icon       string
	UnlockedAt string
}

// GetWeeklyReportData получает данные для отчёта за неделю
func (s *ReportService) GetWeeklyReportData(ctx context.Context, childProfileID string, weekStart time.Time) (*WeeklyReportData, error) {
	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		return nil, fmt.Errorf("invalid child_profile_id: %w", err)
	}

	weekEnd := weekStart.AddDate(0, 0, 6) // Sunday
	prevWeekStart := weekStart.AddDate(0, 0, -7)
	prevWeekEnd := weekEnd.AddDate(0, 0, -7)

	data := &WeeklyReportData{}

	// Get profile data
	var avatarID string
	query := `SELECT display_name, avatar_id, grade, level, xp_total, coins_balance, streak_days FROM child_profiles WHERE id = $1`
	err = s.store.DB.QueryRowContext(ctx, query, profileUUID).Scan(
		&data.ChildName, &avatarID, &data.Grade, &data.Level, &data.XPTotal, &data.CoinsBalance, &data.StreakDays,
	)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}

	// Calculate XP for next level
	data.XPForNext = store.XPForLevel(data.Level)

	// Avatar emoji
	avatarEmojis := map[string]string{
		"cat": "🐱", "dog": "🐶", "panda": "🐼", "fox": "🦊",
		"bear": "🐻", "lion": "🦁", "tiger": "🐯", "unicorn": "🦄",
		"robot": "🤖", "alien": "👽",
	}
	if emoji, ok := avatarEmojis[avatarID]; ok {
		data.AvatarEmoji = emoji
	} else {
		data.AvatarEmoji = "🦊"
	}

	// Get attempts data
	attemptsQuery := `
		SELECT COUNT(*) as total_attempts,
			   COALESCE(SUM(CASE WHEN is_correct THEN 1 ELSE 0 END), 0) as successful_attempts,
			   AVG(time_spent_seconds) as avg_time,
			   COALESCE(SUM(hints_used), 0) as total_hints
		FROM attempts
		WHERE child_profile_id = $1 AND created_at >= $2 AND created_at < $3`

	var totalAttempts, successfulAttempts, totalHints int
	var avgTime sql.NullFloat64
	err = s.store.DB.QueryRowContext(ctx, attemptsQuery, profileUUID, weekStart, weekEnd.AddDate(0, 0, 1)).Scan(
		&totalAttempts, &successfulAttempts, &avgTime, &totalHints)
	if err != nil {
		return nil, fmt.Errorf("failed to get attempts data: %w", err)
	}

	// Fill data fields
	data.TotalAttempts = totalAttempts
	data.SuccessfulAttempts = successfulAttempts
	data.HintsUsed = totalHints

	// Calculate accuracy
	if data.TotalAttempts > 0 {
		data.AccuracyPercent = float64(data.SuccessfulAttempts) / float64(data.TotalAttempts) * 100
	}

	// Calculate average time
	if avgTime.Valid && avgTime.Float64 > 0 {
		data.WeekAvgTimeMinutes = avgTime.Float64 / 60
	}

	// Get achievements
	newAchievements, err := s.getNewAchievements(ctx, profileUUID, weekStart, weekEnd.AddDate(0, 0, 1))
	if err != nil {
		return nil, fmt.Errorf("failed to get new achievements: %w", err)
	}
	data.NewAchievements = newAchievements

	// Get villain battles (active + defeated in current week)
	villainBattles, err := s.getVillainBattles(ctx, profileUUID, weekStart, weekEnd.AddDate(0, 0, 1))
	if err != nil {
		return nil, fmt.Errorf("failed to get villain battles: %w", err)
	}
	data.VillainBattles = villainBattles

	// Calculate comparison with previous week
	log.Printf("[ReportService] Calculating comparison: weekStart=%s, prevWeekStart=%s, totalAttempts=%d",
		weekStart.Format("2006-01-02"), prevWeekStart.Format("2006-01-02"), totalAttempts)

	// Получаем данные за предыдущую неделю
	var prevTotalAttempts, prevSuccessfulAttempts int
	var prevAvgTime sql.NullFloat64
	var prevHints int
	err = s.store.DB.QueryRowContext(ctx, attemptsQuery, profileUUID, prevWeekStart, prevWeekEnd.AddDate(0, 0, 1)).Scan(
		&prevTotalAttempts, &prevSuccessfulAttempts, &prevAvgTime, &prevHints)
	log.Printf("[ReportService] Previous week query: err=%v, prevTotalAttempts=%d, prevSuccessfulAttempts=%d",
		err, prevTotalAttempts, prevSuccessfulAttempts)

	// Сравниваем только если есть данные за обе недели
	if err == nil && prevTotalAttempts > 0 && totalAttempts > 0 {
		// Сравнение количества попыток
		attemptsChange := totalAttempts - prevTotalAttempts
		data.AttemptsChange = &attemptsChange
		log.Printf("[ReportService] Setting AttemptsChange to %d", attemptsChange)

		// Сравнение точности
		if prevSuccessfulAttempts > 0 && successfulAttempts > 0 {
			prevAccuracy := float64(prevSuccessfulAttempts) / float64(prevTotalAttempts) * 100
			currAccuracy := data.AccuracyPercent
			if currAccuracy != prevAccuracy {
				change := currAccuracy - prevAccuracy
				data.AccuracyChange = &change
				log.Printf("[ReportService] Setting AccuracyChange to %.1f", change)
			}
		}

		// Сравнение времени
		if avgTime.Valid && prevAvgTime.Valid && prevAvgTime.Float64 > 0 {
			prevTime := prevAvgTime.Float64 / 60
			currTime := data.WeekAvgTimeMinutes
			if currTime != prevTime {
				change := int((currTime - prevTime) * 60) // in seconds
				data.TimeChange = &change
				log.Printf("[ReportService] Setting TimeChange to %d seconds", change)
			}
		}
	}

	return data, nil
}

func (s *ReportService) getNewAchievements(ctx context.Context, childProfileID uuid.UUID, weekStart, weekEnd time.Time) ([]AchievementData, error) {
	query := `
		SELECT a.title, a.icon, a.requirement_value, ca.current_progress, COUNT(*) as count, MIN(ca.unlocked_at) as unlocked_at
		FROM achievements a
		JOIN child_achievements ca ON a.id = ca.achievement_id
		WHERE ca.child_profile_id = $1 AND ca.unlocked_at >= $2 AND ca.unlocked_at <= $3
		GROUP BY a.id, a.title, a.icon, a.requirement_value, ca.current_progress
		ORDER BY unlocked_at`

	rows, err := s.store.DB.QueryContext(ctx, query, childProfileID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []AchievementData
	for rows.Next() {
		var ach AchievementData
		err := rows.Scan(&ach.Title, &ach.Icon, &ach.RequirementValue, &ach.CurrentProgress, &ach.Count, &ach.UnlockedAt)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, ach)
	}

	return achievements, rows.Err()
}

func (s *ReportService) getVillainBattles(ctx context.Context, childProfileID uuid.UUID, weekStart, weekEnd time.Time) ([]VillainBattleData, error) {
	// Получаем активные битвы и побеждённых злодеев за текущую неделю
	// Используем emoji вместо image_url для совместимости с PDF
	villainEmojis := map[string]string{
		"count_error":         "👿",
		"baron_confusion":     "😵",
		"duchess_distraction": "💃",
		"sir_procrastination": "🦥",
		"madame_mistake":      "🎭",
		"lord_laziness":       "😴",
		"boss_week_chaos":     "👹",
	}

	query := `
		SELECT v.id, v.name, vb.current_hp, v.max_hp, vb.total_damage_dealt, vb.correct_tasks_count, vb.status
		FROM villain_battles vb
		JOIN villains v ON vb.villain_id = v.id
		WHERE vb.child_profile_id = $1
		  AND (vb.status = 'active' OR (vb.status = 'defeated' AND vb.defeated_at >= $2 AND vb.defeated_at < $3))
		ORDER BY
		  CASE vb.status
		    WHEN 'defeated' THEN 1
		    WHEN 'active' THEN 2
		  END,
		  vb.defeated_at DESC,
		  vb.started_at DESC`

	rows, err := s.store.DB.QueryContext(ctx, query, childProfileID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var battles []VillainBattleData
	for rows.Next() {
		var battle VillainBattleData
		var villainID string
		err := rows.Scan(&villainID, &battle.Name, &battle.CurrentHP, &battle.MaxHP, &battle.TotalDamageDealt, &battle.CorrectTasksCount, &battle.Status)
		if err != nil {
			return nil, err
		}
		battle.HPPercent = float64(battle.CurrentHP) / float64(battle.MaxHP) * 100

		// Устанавливаем emoji для злодея
		if emoji, ok := villainEmojis[villainID]; ok {
			battle.Icon = emoji
		} else {
			battle.Icon = "👹" // fallback
		}

		battles = append(battles, battle)
	}

	return battles, rows.Err()
}

func (s *ReportService) GenerateWeeklyReport(ctx context.Context, childProfileID uuid.UUID, reportDate time.Time) (*WeeklyReport, error) {
	// Calculate week start (Monday) and end (Sunday) for the report date
	daysSinceMonday := int(reportDate.Weekday() - time.Monday)
	if daysSinceMonday < 0 {
		daysSinceMonday += 7
	}
	weekStart := reportDate.AddDate(0, 0, -daysSinceMonday)
	weekEnd := weekStart.AddDate(0, 0, 6)

	data, err := s.GetWeeklyReportData(ctx, childProfileID.String(), weekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly report data: %w", err)
	}

	// Update data with calculated fields
	data.ReportWeekStart = weekStart
	data.ReportWeekEnd = weekEnd

	html, err := s.generateWeeklyReportHTML(data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML: %w", err)
	}

	report := &WeeklyReport{
		UserID:      childProfileID,
		ReportDate:  reportDate,
		HTMLContent: html,
	}

	return report, nil
}

func (s *ReportService) SaveWeeklyReport(ctx context.Context, report *WeeklyReport) error {
	query := `
		INSERT INTO weekly_reports (user_id, report_date, html_content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, report_date) DO UPDATE SET
			html_content = EXCLUDED.html_content,
			updated_at = EXCLUDED.updated_at
		RETURNING id`

	now := time.Now()
	err := s.store.DB.QueryRowContext(ctx, query,
		report.UserID, report.ReportDate, report.HTMLContent, now, now).Scan(&report.ID)
	if err != nil {
		return fmt.Errorf("failed to save weekly report: %w", err)
	}

	report.CreatedAt = now
	report.UpdatedAt = now
	return nil
}

func (s *ReportService) GetWeeklyReports(ctx context.Context, childProfileID uuid.UUID) ([]*WeeklyReport, error) {
	query := `
		SELECT id, user_id, report_date, html_content, sent_at, created_at, updated_at
		FROM weekly_reports
		WHERE user_id = $1
		ORDER BY report_date DESC`

	rows, err := s.store.DB.QueryContext(ctx, query, childProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*WeeklyReport
	for rows.Next() {
		var report WeeklyReport
		err := rows.Scan(&report.ID, &report.UserID, &report.ReportDate, &report.HTMLContent,
			&report.SentAt, &report.CreatedAt, &report.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, &report)
	}

	return reports, rows.Err()
}

func (s *ReportService) generateWeeklyReportHTML(data *WeeklyReportData) (string, error) {
	// Generate achievements HTML
	achievementsHTML := ""
	if len(data.NewAchievements) > 0 {
		for _, ach := range data.NewAchievements {
			// Для серийных достижений показываем прогресс (например "5/10")
			progressStr := ""
			if ach.RequirementValue > 1 {
				progressStr = fmt.Sprintf(" <span class=\"badge progress-badge\">%d/%d</span>", ach.CurrentProgress, ach.RequirementValue)
			}

			// Если одно и то же достижение разблокировано несколько раз (маловероятно, но проверяем)
			countStr := ""
			if ach.Count > 1 {
				countStr = fmt.Sprintf(" <span class=\"badge\">×%d</span>", ach.Count)
			}

			achievementsHTML += fmt.Sprintf(`
				<div class="achievement-item">
					<span class="achievement-icon">%s</span>
					<div class="achievement-info">
						<div class="achievement-title">%s%s%s</div>
						<div class="achievement-date">%s</div>
					</div>
				</div>`,
				ach.Icon, ach.Title, progressStr, countStr, ach.UnlockedAt.Format("02.01.2006"))
		}
	}

	// Generate villain battles HTML
	villainHTML := ""
	if len(data.VillainBattles) > 0 {
		for _, battle := range data.VillainBattles {
			// Для побеждённых злодеев показываем статус
			statusBadge := ""
			if battle.Status == "defeated" {
				statusBadge = " <span class=\"defeated-badge\">✓ Побеждён</span>"
			}

			villainHTML += fmt.Sprintf(`
				<div class="villain-card">
					<div class="villain-header">
						<span class="villain-icon">%s</span>
						<div class="villain-info">
							<div class="villain-name">%s%s</div>
							<div class="villain-hp">HP: %d/%d</div>
						</div>
					</div>
					<div class="progress-bar">
						<div class="progress-fill progress-fill-hp" style="width: %.1f%%;"></div>
					</div>
					<div class="villain-stats">
						<div class="villain-stat">
							<span class="stat-label">Урон нанесён</span>
							<span class="stat-value">%d</span>
						</div>
						<div class="villain-stat">
							<span class="stat-label">Задач решено</span>
							<span class="stat-value">%d</span>
						</div>
					</div>
				</div>`,
				battle.Icon, battle.Name, statusBadge, battle.CurrentHP, battle.MaxHP, battle.HPPercent,
				battle.TotalDamageDealt, battle.CorrectTasksCount)
		}
	}

	// Generate comparison HTML
	comparisonHTML := ""
	if data.AttemptsChange != nil || data.AccuracyChange != nil || data.TimeChange != nil {
		comparisonHTML = "<div class=\"comparison-grid\">"

		if data.AttemptsChange != nil {
			change := *data.AttemptsChange
			class := "neutral"
			if change > 0 {
				class = "positive"
			} else if change < 0 {
				class = "negative"
			}
			sign := ""
			if change > 0 {
				sign = "+"
			}
			comparisonHTML += fmt.Sprintf(`
				<div class="comparison-item">
					<div class="comparison-label">Попытки</div>
					<div class="comparison-value %s">%s%d</div>
				</div>`, class, sign, change)
		}

		if data.AccuracyChange != nil {
			change := *data.AccuracyChange
			class := "neutral"
			if change > 0 {
				class = "positive"
			} else if change < 0 {
				class = "negative"
			}
			sign := ""
			if change > 0 {
				sign = "+"
			}
			comparisonHTML += fmt.Sprintf(`
				<div class="comparison-item">
					<div class="comparison-label">Точность</div>
					<div class="comparison-value %s">%s%.1f%%%%</div>
				</div>`, class, sign, change)
		}

		if data.TimeChange != nil {
			change := *data.TimeChange
			class := "neutral"
			if change < 0 {
				class = "positive"
			} else if change > 0 {
				class = "negative"
			}
			sign := ""
			if change < 0 {
				sign = ""
			} else if change > 0 {
				sign = "+"
			}
			comparisonHTML += fmt.Sprintf(`
				<div class="comparison-item">
					<div class="comparison-label">Время</div>
					<div class="comparison-value %s">%s%d мин</div>
				</div>`, class, sign, int(math.Abs(float64(change))))
		}

		comparisonHTML += "</div>"
	} else {
		comparisonHTML = `<div class="empty-state">
			<p>📊 Недостаточно данных для сравнения с прошлой неделей</p>
		</div>`
	}

	// Calculate XP progress
	xpProgressPercent := float64(data.XPTotal%data.XPForNext) / float64(data.XPForNext) * 100

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Еженедельный отчёт — %s</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }

        /* Убираем заголовки и футеры при печати в PDF */
        @page {
            margin: 0;
            size: A4;
        }
        @media print {
            body {
                margin: 0;
                padding: 16px;
            }
        }

        :root {
            --background: #F0F4FF;
            --foreground: #2D3436;
            --card: #ffffff;
            --primary: #6C5CE7;
            --secondary: #A29BFE;
            --muted: #E8E4FF;
            --muted-foreground: #636E72;
            --success: #00B894;
            --destructive: #FF6B6B;
            --border: rgba(108, 92, 231, 0.15);
            --radius: 1rem;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            margin: 0;
            padding: 24px;
            background: linear-gradient(to bottom, var(--background), #E8EFFF);
            color: var(--foreground);
            line-height: 1.6;
        }
        .container { max-width: 900px; margin: 0 auto; }
        .header {
            background: linear-gradient(135deg, var(--primary) 0%%, var(--secondary) 100%%);
            color: white;
            padding: 32px;
            text-align: center;
            border-radius: var(--radius);
            margin-bottom: 24px;
            box-shadow: 0 4px 12px rgba(108, 92, 231, 0.3);
        }
        .header h1 { font-size: 28px; font-weight: 600; margin-bottom: 8px; }
        .header p { font-size: 16px; opacity: 0.95; }
        .section {
            background-color: var(--card);
            margin-bottom: 20px;
            padding: 24px;
            border-radius: var(--radius);
            box-shadow: 0 2px 8px rgba(0,0,0,0.08);
            border: 1px solid var(--border);
        }
        .section h2 { font-size: 20px; color: var(--primary); margin-bottom: 16px; font-weight: 600; }
        .profile { display: flex; align-items: flex-start; gap: 20px; }
        .profile .avatar {
            font-size: 72px;
            line-height: 1;
            background: linear-gradient(135deg, var(--muted) 0%%, #D8DDFF 100%%);
            width: 100px;
            height: 100px;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 16px;
            flex-shrink: 0;
        }
        .profile-info { flex: 1; }
        .profile-row {
            display: flex;
            gap: 12px;
            margin-bottom: 10px;
            font-size: 15px;
        }
        .profile-label { color: var(--muted-foreground); min-width: 100px; }
        .profile-value { font-weight: 600; color: var(--foreground); }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
            gap: 12px;
            margin-top: 16px;
        }
        .stat-card {
            text-align: center;
            padding: 20px 16px;
            background: linear-gradient(135deg, var(--muted) 0%%, #DFE4FF 100%%);
            border-radius: 12px;
        }
        .stat-card .stat-label {
            font-size: 13px;
            color: var(--muted-foreground);
            margin-bottom: 6px;
            display: block;
        }
        .stat-card .stat-value {
            font-size: 24px;
            font-weight: 700;
            color: var(--foreground);
        }
        .stat-card .stat-sub {
            font-size: 13px;
            color: var(--muted-foreground);
            margin-top: 4px;
        }
        .achievement-item {
            display: flex;
            align-items: center;
            gap: 16px;
            padding: 14px;
            background-color: var(--muted);
            border-radius: 12px;
            margin-bottom: 10px;
        }
        .achievement-icon { font-size: 32px; line-height: 1; }
        .achievement-info { flex: 1; }
        .achievement-title { font-weight: 600; font-size: 15px; color: var(--foreground); }
        .achievement-date { font-size: 13px; color: var(--muted-foreground); margin-top: 2px; }
        .badge {
            display: inline-block;
            background-color: var(--primary);
            color: white;
            padding: 2px 8px;
            border-radius: 8px;
            font-size: 12px;
            font-weight: 600;
            margin-left: 6px;
        }
        .progress-badge {
            background-color: var(--secondary);
            color: var(--foreground);
            font-size: 11px;
        }
        .defeated-badge {
            background-color: var(--success);
            color: white;
            font-size: 11px;
            padding: 3px 10px;
        }
        .progress-bar {
            background-color: var(--muted);
            border-radius: 20px;
            height: 12px;
            margin: 12px 0;
            overflow: hidden;
            position: relative;
        }
        .progress-fill {
            background: linear-gradient(90deg, var(--primary) 0%%, var(--secondary) 100%%);
            height: 100%%;
            border-radius: 20px;
            transition: width 0.3s ease;
        }
        .progress-fill-hp {
            background: linear-gradient(90deg, var(--destructive) 0%%, #FF8787 100%%);
        }
        .villain-card {
            padding: 18px;
            background-color: var(--muted);
            border-radius: 12px;
            margin-bottom: 12px;
        }
        .villain-header {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 12px;
        }
        .villain-icon { font-size: 40px; line-height: 1; }
        .villain-info { flex: 1; }
        .villain-name { font-weight: 600; font-size: 16px; color: var(--foreground); }
        .villain-hp { font-size: 13px; color: var(--muted-foreground); margin-top: 2px; }
        .villain-stats {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 12px;
            margin-top: 12px;
        }
        .villain-stat {
            text-align: center;
            padding: 10px;
            background-color: rgba(255, 255, 255, 0.5);
            border-radius: 8px;
        }
        .villain-stat .stat-label {
            font-size: 12px;
            color: var(--muted-foreground);
            display: block;
        }
        .villain-stat .stat-value {
            font-size: 20px;
            font-weight: 700;
            color: var(--foreground);
            margin-top: 4px;
        }
        .comparison-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
            gap: 12px;
        }
        .comparison-item {
            padding: 16px;
            background-color: var(--muted);
            border-radius: 12px;
            text-align: center;
        }
        .comparison-label {
            font-size: 14px;
            color: var(--muted-foreground);
            margin-bottom: 8px;
        }
        .comparison-value {
            font-size: 20px;
            font-weight: 700;
        }
        .positive { color: var(--success); }
        .negative { color: var(--destructive); }
        .neutral { color: var(--muted-foreground); }
        .empty-state {
            text-align: center;
            padding: 32px;
            color: var(--muted-foreground);
            font-size: 15px;
        }
        @media (max-width: 600px) {
            body { padding: 16px; }
            .header { padding: 24px 20px; }
            .header h1 { font-size: 22px; }
            .section { padding: 18px; }
            .profile { flex-direction: column; align-items: center; text-align: center; }
            .stats-grid { grid-template-columns: 1fr 1fr; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📊 Еженедельный отчёт</h1>
            <p>%s · Неделя с %s по %s</p>
        </div>

        <div class="section">
            <h2>👤 Профиль ребёнка</h2>
            <div class="profile">
                <div class="avatar">%s</div>
                <div class="profile-info">
                    <div class="profile-row">
                        <span class="profile-label">Имя:</span>
                        <span class="profile-value">%s</span>
                    </div>
                    <div class="profile-row">
                        <span class="profile-label">Класс:</span>
                        <span class="profile-value">%d</span>
                    </div>
                    <div class="profile-row">
                        <span class="profile-label">Уровень:</span>
                        <span class="profile-value">%d</span>
                    </div>
                    <div class="profile-row">
                        <span class="profile-label">XP:</span>
                        <span class="profile-value">%d (%d%%%%)</span>
                    </div>
                    <div class="progress-bar">
                        <div class="progress-fill" style="width: %d%%%%;"></div>
                    </div>
                    <div class="profile-row">
                        <span class="profile-label">Монеты:</span>
                        <span class="profile-value">%d 🪙</span>
                    </div>
                    <div class="profile-row">
                        <span class="profile-label">Стрик:</span>
                        <span class="profile-value">%d 🔥</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="section">
            <h2>📈 Активность за неделю</h2>
            <div class="stats-grid">
                <div class="stat-card">
                    <span class="stat-label">Попыток</span>
                    <div class="stat-value">%d</div>
                </div>
                <div class="stat-card">
                    <span class="stat-label">Успешных</span>
                    <div class="stat-value">%d</div>
                    <div class="stat-sub">%.1f%%</div>
                </div>
                <div class="stat-card">
                    <span class="stat-label">Подсказки</span>
                    <div class="stat-value">%d</div>
                </div>
                <div class="stat-card">
                    <span class="stat-label">Среднее время</span>
                    <div class="stat-value">%s</div>
                </div>
            </div>
        </div>

        %s

        %s

        <div class="section">
            <h2>📊 Сравнение с прошлой неделей</h2>
            %s
        </div>
    </div>
</body>
</html>`,
		data.ChildName,
		data.ChildName,
		data.ReportWeekStart.Format("02.01.2006"),
		data.ReportWeekEnd.Format("02.01.2006"),
		data.AvatarEmoji,
		data.ChildName,
		data.Grade,
		data.Level,
		data.XPTotal,
		int(xpProgressPercent),
		int(xpProgressPercent),
		data.CoinsBalance,
		data.StreakDays,
		data.TotalAttempts,
		data.SuccessfulAttempts,
		data.AccuracyPercent,
		data.HintsUsed,
		func() string {
			if data.WeekAvgTimeMinutes == 0.0 {
				return "—"
			}
			return fmt.Sprintf("%.1f мин", data.WeekAvgTimeMinutes)
		}(),
		func() string {
			if achievementsHTML != "" {
				return fmt.Sprintf(`<div class="section">
            <h2>🏆 Новые достижения</h2>
            %s
        </div>`, achievementsHTML)
			}
			return ""
		}(),
		func() string {
			if villainHTML != "" {
				return fmt.Sprintf(`<div class="section">
            <h2>⚔️ Битвы со злодеями</h2>
            %s
        </div>`, villainHTML)
			}
			return ""
		}(),
		comparisonHTML)

	return html, nil
}

func (s *ReportService) StartWeeklyReportScheduler(ctx context.Context) {
	log.Printf("Starting weekly report scheduler")
	ticker := time.NewTicker(24 * time.Hour) // Check daily
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Weekly report scheduler stopped")
			return
		case <-ticker.C:
			s.generateAndSaveWeeklyReports(ctx)
		}
	}
}

func (s *ReportService) generateAndSaveWeeklyReports(ctx context.Context) {
	// Generate reports for the previous week (Monday to Sunday)
	now := time.Now()
	// Find previous Sunday (end of previous week)
	daysSinceSunday := int(now.Weekday())
	if daysSinceSunday == 0 { // Sunday
		daysSinceSunday = 7
	}
	weekEnd := now.AddDate(0, 0, -daysSinceSunday)
	weekStart := weekEnd.AddDate(0, 0, -6) // Monday of that week

	// Get all child profile IDs
	rows, err := s.store.DB.QueryContext(ctx, "SELECT id FROM child_profiles")
	if err != nil {
		log.Printf("Failed to query profiles for weekly reports: %v", err)
		return
	}
	defer rows.Close()

	var profileIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			log.Printf("Failed to scan profile ID: %v", err)
			continue
		}
		profileIDs = append(profileIDs, id)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating profiles: %v", err)
		return
	}

	for _, profileID := range profileIDs {
		report, err := s.GenerateWeeklyReport(ctx, profileID, weekEnd)
		if err != nil {
			log.Printf("Failed to generate weekly report for user %s: %v", profileID, err)
			continue
		}

		if err := s.SaveWeeklyReport(ctx, report); err != nil {
			log.Printf("Failed to save weekly report for user %s: %v", profileID, err)
			continue
		}

		log.Printf("Generated weekly report for user %s (week: %s to %s)", profileID,
			weekStart.Format("02.01.2006"), weekEnd.Format("02.01.2006"))

		// TODO: Send email with report
	}
}

// GetWeeklyHTML возвращает HTML отчёт за текущую неделю (генерирует новый или берёт из БД)
func (s *ReportService) GetWeeklyHTML(ctx context.Context, childProfileID string, weekStart time.Time) (string, error) {
	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		return "", fmt.Errorf("invalid child_profile_id: %w", err)
	}

	// Пытаемся получить существующий отчёт из БД
	var htmlContent string
	query := `SELECT html_content FROM weekly_reports WHERE user_id = $1 AND report_date = $2`
	err = s.store.DB.QueryRowContext(ctx, query, profileUUID, weekStart).Scan(&htmlContent)
	if err == nil {
		// Отчёт найден в БД
		return htmlContent, nil
	}

	// Отчёта нет - генерируем новый
	log.Printf("[ReportService] No existing report found, generating new one for %s", childProfileID)
	report, err := s.GenerateWeeklyReport(ctx, profileUUID, weekStart)
	if err != nil {
		return "", fmt.Errorf("failed to generate weekly report: %w", err)
	}

	// Сохраняем в БД
	if err := s.SaveWeeklyReport(ctx, report); err != nil {
		log.Printf("[ReportService] Warning: failed to save report to DB: %v", err)
		// Продолжаем, возвращаем HTML даже если не удалось сохранить
	}

	return report.HTMLContent, nil
}

// ConvertHTMLToPDF конвертирует HTML в PDF используя headless chromium
func (s *ReportService) ConvertHTMLToPDF(htmlContent string) ([]byte, error) {
	// Создаем временные файлы
	tmpHTMLFile, err := os.CreateTemp("", "report-*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp HTML file: %w", err)
	}
	defer os.Remove(tmpHTMLFile.Name())

	tmpPDFFile, err := os.CreateTemp("", "report-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp PDF file: %w", err)
	}
	defer os.Remove(tmpPDFFile.Name())

	// Записываем HTML
	if _, err := tmpHTMLFile.WriteString(htmlContent); err != nil {
		return nil, fmt.Errorf("failed to write HTML: %w", err)
	}
	tmpHTMLFile.Close()

	// Запускаем chromium в headless режиме для генерации PDF
	cmd := exec.Command("chromium-browser",
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"--disable-software-rasterizer",
		"--print-to-pdf="+tmpPDFFile.Name(),
		"--print-to-pdf-no-header", // Убирает URL и дату из header/footer
		"--no-pdf-header-footer",   // Убирает все header/footer
		"--font-render-hinting=none",
		"file://"+tmpHTMLFile.Name(),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("chromium failed: %w, output: %s", err, string(output))
	}

	// Читаем PDF
	pdfContent, err := os.ReadFile(tmpPDFFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF: %w", err)
	}

	log.Printf("[ReportService] Successfully converted HTML to PDF (%d bytes)", len(pdfContent))
	return pdfContent, nil
}
