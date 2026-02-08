package repositories

import (
	"database/sql"
	"go-kasir-api/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// 1. Insert Transaction
	// We only insert date and total. ID is auto-increment.
	// Assuming the table structure matches.
	query := "INSERT INTO transactions (date, total_amount) VALUES ($1, $2) RETURNING id"
	// Assuming date is current time if not set, or we use the passed date.
	if transaction.Date.IsZero() {
		transaction.Date = time.Now()
	}

	var id int
	err = tx.QueryRow(query, transaction.Date, transaction.Total).Scan(&id)
	if err != nil {
		tx.Rollback()
		return err
	}
	// transaction.ID = int(id) // id is already int
	transaction.ID = id

	// 2. Insert Transaction Details
	for _, detail := range transaction.Details {
		detailQuery := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)"
		// Ensure formatting matches DB types.
		// detail.TransactionID is set to the new ID.
		_, err := tx.Exec(detailQuery, transaction.ID, detail.ProductID, detail.Quantity, detail.Subtotal)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *TransactionRepository) GetDailyReport(date time.Time) (map[string]interface{}, error) {
	// Format date for query, assuming DATETIME or DATE column
	startDate := date.Format("2006-01-02") + " 00:00:00"
	endDate := date.Format("2006-01-02") + " 23:59:59"

	var totalRevenue int
	var totalTransactions int

	// Query Total Revenue and Total Transactions
	query := `
		SELECT 
			COALESCE(SUM(total_amount), 0), 
			COUNT(id) 
		FROM transactions 
		WHERE date BETWEEN $1 AND $2`
	
	err := r.db.QueryRow(query, startDate, endDate).Scan(&totalRevenue, &totalTransactions)
	if err != nil {
		return nil, err
	}

	// Query Best Selling Product
	// This requires joining with products table to get name, or just returning ID
	// Let's assume we want product name.
	// Adjust table names if necessary (transaction_details vs transaction_detail?)
	// Based on common naming conventions: transaction_details
	bestSellingQuery := `
		SELECT 
			p.name, 
			SUM(td.quantity) as total_qty 
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.date BETWEEN $1 AND $2
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`
	
	var bestProductName string
	var bestProductQty int

	err = r.db.QueryRow(bestSellingQuery, startDate, endDate).Scan(&bestProductName, &bestProductQty)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	
	// Construct response
	result := map[string]interface{}{
		"total_revenue":     totalRevenue,
		"total_transaksi":   totalTransactions,
		"produk_terlaris": map[string]interface{}{
			"nama":        bestProductName,
			"qty_terjual": bestProductQty,
		},
	}

	return result, nil
}
