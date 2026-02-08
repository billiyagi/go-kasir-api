package handlers

import (
	"encoding/json"
	"go-kasir-api/models"
	"go-kasir-api/services"
	"net/http"
	"time"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// HandleCreateTransaction creates a new transaction
// @Summary Create a new transaction
// @Description Create a new transaction with details
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal Server Error"
// @Router /transactions [post]
func (h *TransactionHandler) HandleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTransaction(&transaction); err != nil {
		http.Error(w, "Failed to create transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Transaction created successfully",
		"data":    transaction, // Returns transaction with ID populated
	})
}

// HandleDailyReport gets the daily sales report
// @Summary Get daily sales report
// @Description Get total revenue, total transactions, and best selling product for today
// @Tags reports
// @Produce json
// @Success 200 {object} models.DailyReport
// @Failure 500 {string} string "Internal Server Error"
// @Router /report/hari-ini [get]
func (h *TransactionHandler) HandleDailyReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For "hari-ini", we use time.Now()
	report, err := h.service.GetDailyReport(time.Now())
	if err != nil {
		http.Error(w, "Failed to get daily report: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
