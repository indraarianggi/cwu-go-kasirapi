package repositories

import (
	"database/sql"
	"time"

	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetSalesSummary(startDate, endDate time.Time) (*models.SalesSummary, error) {
	summary := &models.SalesSummary{}

	// Build WHERE clause dynamically based on provided dates
	var whereClause string
	var args []interface{}

	if startDate.IsZero() && endDate.IsZero() {
		// No date filter - return all-time data
		whereClause = ""
	} else if startDate.IsZero() {
		// Only end date provided
		endDatePlusOne := endDate.AddDate(0, 0, 1)
		whereClause = "WHERE created_at < $1"
		args = append(args, endDatePlusOne)
	} else if endDate.IsZero() {
		// Only start date provided
		whereClause = "WHERE created_at >= $1"
		args = append(args, startDate)
	} else {
		// Both dates provided
		endDatePlusOne := endDate.AddDate(0, 0, 1)
		whereClause = "WHERE created_at >= $1 AND created_at < $2"
		args = append(args, startDate, endDatePlusOne)
	}

	// Query 1: Get total revenue and total transactions
	revenueQuery := `
		SELECT
			COALESCE(SUM(total_amount), 0) AS total_revenue,
			COUNT(*) AS total_transaction
		FROM transactions
		` + whereClause

	err := repo.db.QueryRow(revenueQuery, args...).Scan(
		&summary.TotalRevenue,
		&summary.TotalTransaction,
	)
	if err != nil {
		return nil, err
	}

	// Query 2: Get best selling product
	// Adjust WHERE clause for joined table (t.created_at instead of created_at)
	joinedWhereClause := whereClause
	if whereClause != "" {
		joinedWhereClause = "WHERE t.created_at" + whereClause[len("WHERE created_at"):]
	}

	bestSellingQuery := `
		SELECT
			td.product_id,
			p.name AS product_name,
			SUM(td.quantity) AS count
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		` + joinedWhereClause + `
		GROUP BY td.product_id, p.name
		ORDER BY count DESC
		LIMIT 1
	`

	var bestSelling models.BestSellingProduct
	err = repo.db.QueryRow(bestSellingQuery, args...).Scan(
		&bestSelling.ProductID,
		&bestSelling.ProductName,
		&bestSelling.Count,
	)

	if err == sql.ErrNoRows {
		// No transactions in date range - best_selling_product will be nil
		summary.BestSellingProduct = nil
	} else if err != nil {
		return nil, err
	} else {
		summary.BestSellingProduct = &bestSelling
	}

	return summary, nil
}
