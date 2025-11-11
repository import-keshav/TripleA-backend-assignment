package repository

import (
	"database/sql"
	"fmt"

	"triplea-backend-assignment/database"
	"triplea-backend-assignment/models"
)

type TransactionRepository struct{}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) Create(tx *sql.Tx, sourceAccountID, destinationAccountID int64, amount models.Decimal) (*models.Transaction, error) {
	query := `INSERT INTO transactions (source_account_id, destination_account_id, amount, status)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, source_account_id, destination_account_id, amount, status, created_at, updated_at`
	
	transaction := &models.Transaction{}
	err := tx.QueryRow(query, sourceAccountID, destinationAccountID, amount, models.TransactionStatusPending).
		Scan(&transaction.ID, &transaction.SourceAccountID, &transaction.DestinationAccountID,
			&transaction.Amount, &transaction.Status, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return transaction, nil
}

func (r *TransactionRepository) UpdateStatus(tx *sql.Tx, transactionID int64, status string) error {
	query := `UPDATE transactions SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := tx.Exec(query, status, transactionID)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}
	return nil
}

func (r *TransactionRepository) GetByID(transactionID int64) (*models.Transaction, error) {
	query := `SELECT id, source_account_id, destination_account_id, amount, status, created_at, updated_at
			  FROM transactions WHERE id = $1`
	transaction := &models.Transaction{}
	err := database.DB.QueryRow(query, transactionID).
		Scan(&transaction.ID, &transaction.SourceAccountID, &transaction.DestinationAccountID,
			&transaction.Amount, &transaction.Status, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return transaction, nil
}

