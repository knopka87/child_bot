package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
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
	Title      string
	Icon       string
	Count      int
	UnlockedAt time.Time
}

type VillainBattleData struct {
	Name              string
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
			   SUM(CASE WHEN is_correct THEN 1 ELSE 0 END) as successful_attempts,
			   AVG(time_spent_seconds) as avg_time,
			   SUM(hints_used) as total_hints
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

	// Get villain battles
	villainBattles, err := s.getVillainBattles(ctx, profileUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get villain battles: %w", err)
	}
	data.VillainBattles = villainBattles

	// Calculate comparison with previous week
	if totalAttempts > 0 {
		var prevTotalAttempts, prevSuccessfulAttempts int
		err = s.store.DB.QueryRowContext(ctx, attemptsQuery, profileUUID, prevWeekStart, prevWeekEnd.AddDate(0, 0, 1)).Scan(
			&prevTotalAttempts, &prevSuccessfulAttempts, nil, nil)
		if err == nil && prevTotalAttempts > 0 {
			attemptsChange := totalAttempts - prevTotalAttempts
			data.AttemptsChange = &attemptsChange

			if prevSuccessfulAttempts > 0 && successfulAttempts > 0 {
				prevAccuracy := float64(prevSuccessfulAttempts) / float64(prevTotalAttempts) * 100
				currAccuracy := data.AccuracyPercent
				if currAccuracy != prevAccuracy {
					change := currAccuracy - prevAccuracy
					data.AccuracyChange = &change
				}
			}

			if avgTime.Valid && prevTotalAttempts > 0 {
				var prevAvgTime sql.NullFloat64
				err = s.store.DB.QueryRowContext(ctx, attemptsQuery, profileUUID, prevWeekStart, prevWeekEnd.AddDate(0, 0, 1)).Scan(
					nil, nil, &prevAvgTime, nil)
				if err == nil && prevAvgTime.Valid && prevAvgTime.Float64 > 0 {
					prevTime := prevAvgTime.Float64 / 60
					currTime := data.WeekAvgTimeMinutes
					if currTime != prevTime {
						change := int((currTime - prevTime) * 60) // in minutes
						data.TimeChange = &change
					}
				}
			}
		}
	}

	return data, nil
}

func (s *ReportService) getNewAchievements(ctx context.Context, childProfileID uuid.UUID, weekStart, weekEnd time.Time) ([]AchievementData, error) {
	query := `
		SELECT a.title, a.icon, COUNT(*) as count, MIN(ca.unlocked_at) as unlocked_at
		FROM achievements a
		JOIN child_achievements ca ON a.id = ca.achievement_id
		WHERE ca.child_profile_id = $1 AND ca.unlocked_at >= $2 AND ca.unlocked_at <= $3
		GROUP BY a.id, a.title, a.icon
		ORDER BY unlocked_at`

	rows, err := s.store.DB.QueryContext(ctx, query, childProfileID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []AchievementData
	for rows.Next() {
		var ach AchievementData
		err := rows.Scan(&ach.Title, &ach.Icon, &ach.Count, &ach.UnlockedAt)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, ach)
	}

	return achievements, rows.Err()
}

func (s *ReportService) getVillainBattles(ctx context.Context, childProfileID uuid.UUID) ([]VillainBattleData, error) {
	query := `
		SELECT v.name, vb.current_hp, v.max_hp, vb.total_damage_dealt, vb.correct_tasks_count, vb.status
		FROM villain_battles vb
		JOIN villains v ON vb.villain_id = v.id
		WHERE vb.child_profile_id = $1 AND vb.status IN ('active', 'defeated')
		ORDER BY vb.started_at DESC`

	rows, err := s.store.DB.QueryContext(ctx, query, childProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var battles []VillainBattleData
	for rows.Next() {
		var battle VillainBattleData
		err := rows.Scan(&battle.Name, &battle.CurrentHP, &battle.MaxHP, &battle.TotalDamageDealt, &battle.CorrectTasksCount, &battle.Status)
		if err != nil {
			return nil, err
		}
		battle.HPPercent = float64(battle.CurrentHP) / float64(battle.MaxHP) * 100
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
		achievementsHTML = "<ul class=\"achievements\">"
		for _, ach := range data.NewAchievements {
			countStr := ""
			if ach.Count > 1 {
				countStr = fmt.Sprintf(" (%d)", ach.Count)
			}
			achievementsHTML += fmt.Sprintf("<li>%s <strong>%s%s</strong> (%s)</li>",
				ach.Icon, ach.Title, countStr, ach.UnlockedAt.Format("02.01.2006"))
		}
		achievementsHTML += "</ul>"
	}

	// Generate villain battles HTML
	villainHTML := ""
	if len(data.VillainBattles) > 0 {
		for _, battle := range data.VillainBattles {
			if battle.TotalDamageDealt > 0 {
				villainHTML += fmt.Sprintf(`
					<div class="villain">
						<p><strong>%s</strong></p>
						<p><strong>HP:</strong> %d/%d (%d%%)</p>
						<div class="progress-bar"><div class="progress-fill" style="width: %d%%;"></div></div>
						<p><strong>Общий урон:</strong> %d</p>
						<p><strong>Правильные задачи:</strong> %d</p>
					</div>`,
					battle.Name, battle.CurrentHP, battle.MaxHP, int(battle.HPPercent),
					int(battle.HPPercent), battle.TotalDamageDealt, battle.CorrectTasksCount)
			}
		}
	}

	// Generate comparison HTML
	comparisonHTML := ""
	if data.AttemptsChange != nil || data.AccuracyChange != nil || data.TimeChange != nil {
		comparisonHTML = "<div class=\"comparison\">"

		if data.AttemptsChange != nil {
			change := *data.AttemptsChange
			class := ""
			if change > 0 {
				class = "positive"
			} else if change < 0 {
				class = "negative"
			}
			sign := ""
			if change > 0 {
				sign = "+"
			}
			comparisonHTML += fmt.Sprintf("<div><strong>Попытки:</strong> <span class=\"%s\">%s%d</span></div>", class, sign, change)
		}

		if data.AccuracyChange != nil {
			change := *data.AccuracyChange
			class := ""
			if change > 0 {
				class = "positive"
			} else if change < 0 {
				class = "negative"
			}
			sign := ""
			if change > 0 {
				sign = "+"
			}
			comparisonHTML += fmt.Sprintf("<div><strong>Точность:</strong> <span class=\"%s\">%s%.1f%%</span></div>", class, sign, change)
		}

		if data.TimeChange != nil {
			change := *data.TimeChange
			class := ""
			if change < 0 {
				class = "positive"
			} else if change > 0 {
				class = "negative"
			}
			comparisonHTML += fmt.Sprintf("<div><strong>Время:</strong> <span class=\"%s\">%s%d мин</span></div>", class,
				func() string {
					if change < 0 {
						return "-"
					} else if change > 0 {
						return "+"
					}
					return ""
				}(), int(math.Abs(float64(change))))
		}

		comparisonHTML += "</div>"
	} else {
		comparisonHTML = "<div class=\"comparison\"><div>Недостаточно данных для сравнения</div></div>"
	}

	// Calculate XP progress
	xpProgressPercent := float64(data.XPTotal%data.XPForNext) / float64(data.XPForNext) * 100

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Еженедельный отчёт</title>
    <style>
        :root {
            --background: #F0F4FF;
            --foreground: #2D3436;
            --card: #ffffff;
            --primary: #6C5CE7;
            --secondary: #A29BFE;
            --muted: #E8E4FF;
            --success: #00B894;
            --destructive: #FF6B6B;
            --border: rgba(108, 92, 231, 0.15);
            --radius: 1rem;
        }
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: var(--background); color: var(--foreground); }
        .header { background-color: var(--primary); color: white; padding: 20px; text-align: center; border-radius: var(--radius); margin-bottom: 20px; }
        .section { background-color: var(--card); margin-bottom: 20px; padding: 20px; border-radius: var(--radius); box-shadow: 0 2px 4px rgba(0,0,0,0.1); border: 1px solid var(--border); }
        .profile { display: flex; align-items: center; }
        .profile .avatar { font-size: 80px; margin-right: 20px; }
        .stats { display: flex; justify-content: space-around; margin-top: 20px; flex-wrap: wrap; }
        .stat { text-align: center; padding: 15px; background-color: var(--muted); border-radius: calc(var(--radius) - 4px); margin: 5px; flex: 1; min-width: 120px; }
        .achievements { list-style: none; padding: 0; }
        .achievements li { padding: 10px; background-color: var(--muted); margin-bottom: 5px; border-radius: calc(var(--radius) - 4px); }
        .progress-bar { background-color: var(--muted); border-radius: 20px; height: 20px; margin: 10px 0; overflow: hidden; }
        .progress-fill { background-color: var(--primary); height: 100%%; border-radius: 20px; }
        .comparison { display: flex; justify-content: space-between; flex-wrap: wrap; }
        .comparison div { padding: 10px; background-color: var(--muted); border-radius: calc(var(--radius) - 4px); margin: 5px; flex: 1; text-align: center; }
        .positive { color: var(--success); }
        .negative { color: var(--destructive); }
        .villain { padding: 15px; background-color: var(--muted); border-radius: var(--radius); margin-bottom: 10px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Еженедельный отчёт</h1>
        <p>Для %s | Неделя с %s по %s</p>
    </div>

    <div class="section">
        <h2>Профиль ребёнка</h2>
        <div class="profile">
            <div class="avatar">%s</div>
            <div>
                <p><strong>Имя:</strong> %s</p>
                <p><strong>Класс:</strong> %d</p>
                <p><strong>Уровень:</strong> %d</p>
                <p><strong>XP:</strong> %d (предполагаемый прогресс: %d%%)</p>
                <div class="progress-bar"><div class="progress-fill" style="width: %d%%;"></div></div>
                <p><strong>Монеты:</strong> %d 🪙</p>
                <p><strong>Стрик дней:</strong> %d 🔥</p>
            </div>
        </div>
    </div>

    <div class="section">
        <h2>Активность за неделю</h2>
        <div class="stats">
            <div class="stat"><strong>Попыток:</strong><br>%d</div>
            <div class="stat"><strong>Успешных:</strong><br>%d (%.2f%%)</div>
            <div class="stat"><strong>Подсказки:</strong><br>%d</div>
            <div class="stat"><strong>Среднее время:</strong><br>%s</div>
        </div>
    </div>

    %s

    %s

    %s
</body>
</html>`,
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
				return "Неизвестно"
			}
			return fmt.Sprintf("%.1f мин", data.WeekAvgTimeMinutes)
		}(),
		func() string {
			if achievementsHTML != "" {
				return fmt.Sprintf(`<div class="section">
        <h2>Новые достижения</h2>
        %s
    </div>`, achievementsHTML)
			}
			return ""
		}(),
		func() string {
			if villainHTML != "" {
				return fmt.Sprintf(`<div class="section">
        <h2>Битва со злодеями</h2>
        %s
    </div>`, villainHTML)
			}
			return ""
		}(),
		func() string {
			return fmt.Sprintf(`<div class="section">
        <h2>Сравнение с прошлой неделей</h2>
        %s
    </div>`, comparisonHTML)
		}())

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
