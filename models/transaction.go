package models

import "time"

type Transaction struct {
	ID        int                 `json:"id"`
	Date      time.Time           `json:"date"`
	Total     int                 `json:"total"`
	Details   []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            int `json:"id"`
	TransactionID int `json:"transaction_id"`
	ProductID     int `json:"product_id"`
	Quantity      int `json:"quantity"`
	Subtotal      int `json:"subtotal"`
}
