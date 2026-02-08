package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Initialize variables for total amount and transaction details
	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	// Loop through items to get product details, calculate total amount and transaction details
	for _, item := range items {
		var productPrice, stock int
		var productName string

		// Get product details from database
		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Calculate subtotal and add to total amount
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// Update product stock
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		// Add transaction details to details slice
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Insert transaction into database
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// Insert transaction details into database using bulk insert
	if len(details) > 0 {
		args := make([]interface{}, 0, len(details)*4)
		valueStrings := make([]string, 0, len(details))

		for i, detail := range details {
			details[i].TransactionID = transactionID
			valueStrings = append(valueStrings,
				fmt.Sprintf("($%d, $%d, $%d, $%d)",
					i*4+1, i*4+2, i*4+3, i*4+4))
			args = append(args, transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
		}

		query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES " +
			strings.Join(valueStrings, ", ")
		_, err = tx.Exec(query, args...)
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return transaction
	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
