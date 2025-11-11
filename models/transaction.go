package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Transaction represents a transaction in the system
type Transaction struct {
	ID                 int64     `json:"id" db:"id"`
	SourceAccountID    int64     `json:"source_account_id" db:"source_account_id"`
	DestinationAccountID int64   `json:"destination_account_id" db:"destination_account_id"`
	Amount             Decimal   `json:"amount" db:"amount"`
	Status             string    `json:"status" db:"status"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// TransactionStatus constants
const (
	TransactionStatusPending   = "pending"
	TransactionStatusCompleted = "completed"
	TransactionStatusFailed    = "failed"
)

// CreateTransactionRequest represents the request body for creating a transaction
type CreateTransactionRequest struct {
	SourceAccountID      int64  `json:"source_account_id"`
	DestinationAccountID int64  `json:"destination_account_id"`
	Amount               string `json:"amount"`
}

// Validate validates the create transaction request
func (r *CreateTransactionRequest) Validate() error {
	if r.SourceAccountID <= 0 {
		return errors.New("source_account_id must be a positive integer")
	}
	if r.DestinationAccountID <= 0 {
		return errors.New("destination_account_id must be a positive integer")
	}
	if r.SourceAccountID == r.DestinationAccountID {
		return errors.New("source_account_id and destination_account_id cannot be the same")
	}
	if r.Amount == "" {
		return errors.New("amount is required")
	}
	// Validate that amount is a valid decimal
	amount, err := strconv.ParseFloat(r.Amount, 64)
	if err != nil {
		return fmt.Errorf("amount must be a valid decimal number: %w", err)
	}
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return nil
}

