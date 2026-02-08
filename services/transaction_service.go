package services

import (
	"go-kasir-api/models"
	"go-kasir-api/repositories"
	"time"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) CreateTransaction(transaction *models.Transaction) error {
	// Calculate total amount if not provided
	if transaction.Total == 0 {
		for _, detail := range transaction.Details {
			transaction.Total += detail.Subtotal
		}
	}
	return s.repo.CreateTransaction(transaction)
}

func (s *TransactionService) GetDailyReport(date time.Time) (map[string]interface{}, error) {
	return s.repo.GetDailyReport(date)
}
