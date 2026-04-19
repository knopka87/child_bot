package handler

import (
	"log"
	"net/http"
	"time"

	"child-bot/api/internal/api/response"
	"child-bot/api/internal/api/validation"
	"child-bot/api/internal/service"
)

// ReportHandler обрабатывает запросы для отчётов
type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler создаёт новый ReportHandler
func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// GetWeeklyPDF генерирует PDF-отчёт за текущую неделю
// GET /reports/{childProfileId}/weekly/pdf
func (h *ReportHandler) GetWeeklyPDF(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("[ReportHandler] Generating weekly PDF for profile %s, week starting %s", childProfileID, weekStart.Format("2006-01-02"))

	// Генерируем PDF
	pdfData, err := h.reportService.GenerateWeeklyPDF(r.Context(), childProfileID, weekStart)
	if err != nil {
		log.Printf("[ReportHandler] Failed to generate weekly PDF: %v", err)
		response.InternalError(w, "Failed to generate report")
		return
	}

	// Отправляем PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=homework_report_"+weekStart.Format("2006-01-02")+".pdf")
	w.Header().Set("Content-Length", "0") // Dynamic

	_, err = w.Write(pdfData)
	if err != nil {
		log.Printf("[ReportHandler] Failed to write PDF response: %v", err)
	}

	log.Printf("[ReportHandler] Weekly PDF sent successfully (%d bytes)", len(pdfData))
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
