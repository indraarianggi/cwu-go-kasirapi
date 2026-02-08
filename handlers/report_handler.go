package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"net/http"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// HandleReport - GET /api/report
func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetSalesSummary(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetSalesSummary godoc
// @Summary Get sales summary report
// @Description Get sales summary including total revenue, transaction count, and best selling product for a date range. If no dates provided, returns all-time data.
// @Tags reports
// @Produce json
// @Param start_date query string false "Start date (yyyy-MM-dd format, e.g., 2026-01-01). If omitted, includes all transactions from the beginning."
// @Param end_date query string false "End date (yyyy-MM-dd format, e.g., 2026-01-31). If omitted, includes all transactions up to now."
// @Success 200 {object} models.SalesSummary
// @Failure 400 {string} string "Invalid request parameters"
// @Failure 500 {string} string "Internal server error"
// @Router /api/report [get]
func (h *ReportHandler) GetSalesSummary(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	summary, err := h.service.GetSalesSummary(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
