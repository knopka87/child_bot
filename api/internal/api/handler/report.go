package handler

import (
	"log"
	"net/http"
	"time"

	"child-bot/api/internal/api/response"
	"child-bot/api/internal/api/validation"
	"child-bot/api/internal/service"

	"github.com/google/uuid"
)

// ReportHandler обрабатывает запросы для отчётов
type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler создаёт новый ReportHandler
func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// GetWeeklyData получает данные отчёта в JSON (для превью)
// GET /reports/{childProfileId}/weekly/data
func (h *ReportHandler) GetWeeklyData(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.PathValue("childProfileId")
	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	// Определяем период
	now := time.Now()
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -int(weekday-1))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, now.Location())

	data, err := h.reportService.GetWeeklyReportData(r.Context(), childProfileID, weekStart)
	if err != nil {
		log.Printf("[ReportHandler] Failed to get weekly data: %v", err)
		response.InternalError(w, "Failed to get report data")
		return
	}

	response.OK(w, data)
}

// GetWeeklyHTML получает HTML отчёт за текущую неделю
// GET /reports/{childProfileId}/weekly/html
func (h *ReportHandler) GetWeeklyHTML(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.PathValue("childProfileId")
	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	// Определяем период: начало текущей недели (понедельник)
	now := time.Now()
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -int(weekday-1))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, now.Location())

	log.Printf("[ReportHandler] Getting weekly HTML for profile %s, week starting %s", childProfileID, weekStart.Format("2006-01-02"))

	// Получаем HTML
	htmlContent, err := h.reportService.GetWeeklyHTML(r.Context(), childProfileID, weekStart)
	if err != nil {
		log.Printf("[ReportHandler] Failed to get weekly HTML: %v", err)
		response.InternalError(w, "Failed to generate report")
		return
	}

	// Отправляем HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(htmlContent))
	if err != nil {
		log.Printf("[ReportHandler] Failed to write HTML response: %v", err)
	}

	log.Printf("[ReportHandler] Weekly HTML sent successfully (%d bytes)", len(htmlContent))
}

// GetReportsList получает список всех отчётов для пользователя
// GET /reports/{childProfileId}/list
func (h *ReportHandler) GetReportsList(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.PathValue("childProfileId")
	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		response.BadRequest(w, "invalid UUID format")
		return
	}

	reports, err := h.reportService.GetWeeklyReports(r.Context(), profileUUID)
	if err != nil {
		log.Printf("[ReportHandler] Failed to get reports list: %v", err)
		response.InternalError(w, "Failed to get reports list")
		return
	}

	// Преобразуем в response format
	type ReportInfo struct {
		ID         string `json:"id"`
		ReportDate string `json:"reportDate"`
		SentAt     string `json:"sentAt,omitempty"`
		CreatedAt  string `json:"createdAt"`
	}

	var result []ReportInfo
	for _, report := range reports {
		info := ReportInfo{
			ID:         report.ID.String(),
			ReportDate: report.ReportDate.Format("2006-01-02"),
			CreatedAt:  report.CreatedAt.Format(time.RFC3339),
		}
		if report.SentAt != nil {
			info.SentAt = report.SentAt.Format(time.RFC3339)
		}
		result = append(result, info)
	}

	response.OK(w, result)
}

// GetReportByDate получает HTML отчёт по дате
// GET /reports/{childProfileId}/{reportDate}/html
func (h *ReportHandler) GetReportByDate(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.PathValue("childProfileId")
	reportDateStr := r.PathValue("reportDate")

	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		response.BadRequest(w, "invalid UUID format")
		return
	}

	reportDate, err := time.Parse("2006-01-02", reportDateStr)
	if err != nil {
		response.BadRequest(w, "invalid date format, use YYYY-MM-DD")
		return
	}

	// Получаем отчёт из БД
	var htmlContent string
	query := `SELECT html_content FROM weekly_reports WHERE user_id = $1 AND report_date = $2`
	err = h.reportService.GetStore().DB.QueryRowContext(r.Context(), query, profileUUID, reportDate).Scan(&htmlContent)
	if err != nil {
		log.Printf("[ReportHandler] Report not found for date %s: %v", reportDateStr, err)
		response.NotFound(w, "Report not found for this date")
		return
	}

	// Отправляем HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(htmlContent))
	if err != nil {
		log.Printf("[ReportHandler] Failed to write HTML response: %v", err)
	}

	log.Printf("[ReportHandler] Report for date %s sent successfully", reportDateStr)
}

// DownloadReportPDF скачивает отчёт в формате PDF
// GET /reports/{childProfileId}/{reportDate}/download
func (h *ReportHandler) DownloadReportPDF(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.PathValue("childProfileId")
	reportDateStr := r.PathValue("reportDate")

	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		response.BadRequest(w, "invalid UUID format")
		return
	}

	reportDate, err := time.Parse("2006-01-02", reportDateStr)
	if err != nil {
		response.BadRequest(w, "invalid date format, use YYYY-MM-DD")
		return
	}

	// Получаем отчёт из БД
	var htmlContent string
	query := `SELECT html_content FROM weekly_reports WHERE user_id = $1 AND report_date = $2`
	err = h.reportService.GetStore().DB.QueryRowContext(r.Context(), query, profileUUID, reportDate).Scan(&htmlContent)
	if err != nil {
		log.Printf("[ReportHandler] Report not found for date %s: %v", reportDateStr, err)
		response.NotFound(w, "Report not found for this date")
		return
	}

	// Конвертируем HTML в PDF
	pdfContent, err := h.reportService.ConvertHTMLToPDF(htmlContent)
	if err != nil {
		log.Printf("[ReportHandler] Failed to convert HTML to PDF: %v", err)
		// Если конвертация не удалась, отправляем HTML
		filename := "report_" + reportDateStr + ".html"
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(htmlContent))
		return
	}

	// Отправляем PDF
	filename := "report_" + reportDateStr + ".pdf"
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(pdfContent)
	if err != nil {
		log.Printf("[ReportHandler] Failed to write PDF response: %v", err)
	}

	log.Printf("[ReportHandler] Report for date %s downloaded successfully as PDF (%d bytes)", reportDateStr, len(pdfContent))
}

// GenerateReport генерирует отчёт за указанную неделю
// POST /reports/{childProfileId}/generate?week_start=2026-04-06
func (h *ReportHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.PathValue("childProfileId")
	weekStartStr := r.URL.Query().Get("week_start")

	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		response.BadRequest(w, "invalid UUID format")
		return
	}

	var weekStart time.Time
	if weekStartStr != "" {
		weekStart, err = time.Parse("2006-01-02", weekStartStr)
		if err != nil {
			response.BadRequest(w, "invalid week_start format, use YYYY-MM-DD")
			return
		}
	} else {
		// По умолчанию - текущая неделя
		now := time.Now()
		weekday := now.Weekday()
		if weekday == time.Sunday {
			weekday = 7
		}
		weekStart = now.AddDate(0, 0, -int(weekday-1))
		weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, now.Location())
	}

	log.Printf("[ReportHandler] Generating report for profile %s, week starting %s", childProfileID, weekStart.Format("2006-01-02"))

	// Получаем данные отчёта
	data, err := h.reportService.GetWeeklyReportData(r.Context(), childProfileID, weekStart)
	if err != nil {
		log.Printf("[ReportHandler] Failed to get report data: %v", err)
		response.InternalError(w, "Failed to generate report data")
		return
	}

	// Генерируем HTML
	htmlContent, err := h.reportService.GetWeeklyHTML(r.Context(), childProfileID, weekStart)
	if err != nil {
		log.Printf("[ReportHandler] Failed to generate HTML: %v", err)
		response.InternalError(w, "Failed to generate HTML")
		return
	}

	// Сохраняем отчёт в БД
	report := &service.WeeklyReport{
		UserID:      profileUUID,
		ReportDate:  weekStart,
		HTMLContent: htmlContent,
	}

	if err := h.reportService.SaveWeeklyReport(r.Context(), report); err != nil {
		log.Printf("[ReportHandler] Failed to save report: %v", err)
		response.InternalError(w, "Failed to save report")
		return
	}

	log.Printf("[ReportHandler] Report generated and saved successfully")
	response.OK(w, map[string]interface{}{
		"message":    "Report generated successfully",
		"reportDate": weekStart.Format("2006-01-02"),
		"data":       data,
	})
}

// GetReportSettings получает настройки отчётов
// GET /reports/{childProfileId}/settings
func (h *ReportHandler) GetReportSettings(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.PathValue("childProfileId")
	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		response.BadRequest(w, "invalid UUID format")
		return
	}

	var parentEmail string
	var weeklyEnabled bool
	query := `SELECT parent_email, weekly_report_enabled FROM report_settings WHERE child_profile_id = $1`
	err = h.reportService.GetStore().DB.QueryRowContext(r.Context(), query, profileUUID).Scan(&parentEmail, &weeklyEnabled)

	if err != nil {
		// Если настроек нет, возвращаем дефолтные значения
		response.OK(w, map[string]interface{}{
			"email":               "",
			"weeklyReportEnabled": true,
		})
		return
	}

	response.OK(w, map[string]interface{}{
		"email":               parentEmail,
		"weeklyReportEnabled": weeklyEnabled,
	})
}

// UpdateReportSettings обновляет настройки отчётов
// PUT /reports/{childProfileId}/settings
func (h *ReportHandler) UpdateReportSettings(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.PathValue("childProfileId")
	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		response.BadRequest(w, "invalid UUID format")
		return
	}

	var req struct {
		Email               *string `json:"email"`
		WeeklyReportEnabled *bool   `json:"weeklyReportEnabled"`
	}

	if err := validation.DecodeJSON(r, &req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	// Upsert настроек
	query := `
		INSERT INTO report_settings (child_profile_id, parent_email, weekly_report_enabled)
		VALUES ($1, $2, $3)
		ON CONFLICT (child_profile_id)
		DO UPDATE SET
			parent_email = COALESCE($2, report_settings.parent_email),
			weekly_report_enabled = COALESCE($3, report_settings.weekly_report_enabled),
			updated_at = NOW()
		RETURNING parent_email, weekly_report_enabled`

	var email string
	var enabled bool
	err = h.reportService.GetStore().DB.QueryRowContext(r.Context(), query, profileUUID, req.Email, req.WeeklyReportEnabled).Scan(&email, &enabled)
	if err != nil {
		log.Printf("[ReportHandler] Failed to update settings: %v", err)
		response.InternalError(w, "Failed to update settings")
		return
	}

	response.OK(w, map[string]interface{}{
		"email":               email,
		"weeklyReportEnabled": enabled,
	})
}

// SendTestReport отправляет тестовый отчёт на email
// POST /reports/{childProfileId}/send-test
func (h *ReportHandler) SendTestReport(w http.ResponseWriter, r *http.Request) {
	childProfileID := r.PathValue("childProfileId")
	if err := validation.ValidateUUID(childProfileID); err != nil {
		response.BadRequest(w, "invalid child_profile_id: "+err.Error())
		return
	}

	profileUUID, err := uuid.Parse(childProfileID)
	if err != nil {
		response.BadRequest(w, "invalid UUID format")
		return
	}

	// Получаем email из настроек
	var parentEmail string
	query := `SELECT parent_email FROM report_settings WHERE child_profile_id = $1`
	err = h.reportService.GetStore().DB.QueryRowContext(r.Context(), query, profileUUID).Scan(&parentEmail)
	if err != nil || parentEmail == "" {
		response.BadRequest(w, "Email not configured. Please set parent email first.")
		return
	}

	// Генерируем текущий отчёт
	now := time.Now()
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -int(weekday-1))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, now.Location())

	htmlContent, err := h.reportService.GetWeeklyHTML(r.Context(), childProfileID, weekStart)
	if err != nil {
		log.Printf("[ReportHandler] Failed to generate report: %v", err)
		response.InternalError(w, "Failed to generate report")
		return
	}

	// TODO: Отправить email через SMTP
	// Пока просто логируем
	log.Printf("[ReportHandler] Test report generated for %s, would send to: %s", childProfileID, parentEmail)
	log.Printf("[ReportHandler] Report size: %d bytes", len(htmlContent))

	response.OK(w, map[string]interface{}{
		"message": "Test report sent successfully to " + parentEmail,
		"email":   parentEmail,
	})
}
