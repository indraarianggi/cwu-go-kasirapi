package services

import (
	"errors"
	"time"

	"kasir-api/models"
	"kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// parseDateString parses a date string in "yyyy-MM-dd" format
func parseDateString(dateStr string) (time.Time, error) {
	layout := "2006-01-02" // Go's reference date format for yyyy-MM-dd
	return time.Parse(layout, dateStr)
}

func (s *ReportService) GetSalesSummary(startDateStr, endDateStr string) (*models.SalesSummary, error) {
	var startDate, endDate time.Time
	var err error

	// Handle optional start_date
	if startDateStr == "" {
		// Default to beginning of time (zero time)
		startDate = time.Time{}
	} else {
		startDate, err = parseDateString(startDateStr)
		if err != nil {
			return nil, errors.New("invalid start_date format, expected yyyy-MM-dd")
		}
	}

	// Handle optional end_date
	if endDateStr == "" {
		// Default to current time
		endDate = time.Now()
	} else {
		endDate, err = parseDateString(endDateStr)
		if err != nil {
			return nil, errors.New("invalid end_date format, expected yyyy-MM-dd")
		}
	}

	// Validate date range (only if both dates are not zero/default)
	if !startDate.IsZero() && endDate.Before(startDate) {
		return nil, errors.New("end_date must be after or equal to start_date")
	}

	return s.repo.GetSalesSummary(startDate, endDate)
}
