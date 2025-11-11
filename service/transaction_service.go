package service

import (
	"database/sql"
	"fmt"

	"triplea-backend-assignment/database"
	"triplea-backend-assignment/models"
	"triplea-backend-assignment/repository"
)

// TransactionService handles business logic for transactions
type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	transactionRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

// ProcessTransaction processes a transfer between two accounts
// This method ensures atomicity by using a database transaction
func (s *TransactionService) ProcessTransaction(req *models.CreateTransactionRequest) error {
	// Validate request
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Start database transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if source account exists
	sourceExists, err := s.accountRepo.Exists(req.SourceAccountID)
	if err != nil {
		return fmt.Errorf("failed to check source account existence: %w", err)
	}
	if !sourceExists {
		return fmt.Errorf("source account %d not found", req.SourceAccountID)
	}

	// Check if destination account exists
	destExists, err := s.accountRepo.Exists(req.DestinationAccountID)
	if err != nil {
		return fmt.Errorf("failed to check destination account existence: %w", err)
	}
	if !destExists {
		return fmt.Errorf("destination account %d not found", req.DestinationAccountID)
	}

	// Get source account balance
	sourceAccount, err := s.accountRepo.GetByID(req.SourceAccountID)
	if err != nil {
		return fmt.Errorf("failed to get source account: %w", err)
	}

	// Get destination account balance
	destAccount, err := s.accountRepo.GetByID(req.DestinationAccountID)
	if err != nil {
		return fmt.Errorf("failed to get destination account: %w", err)
	}

	// Convert amounts to float64 for calculation
	amount, err := models.Decimal(req.Amount).Float64()
	if err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	sourceBalance, err := sourceAccount.Balance.Float64()
	if err != nil {
		return fmt.Errorf("invalid source balance format: %w", err)
	}

	destBalance, err := destAccount.Balance.Float64()
	if err != nil {
		return fmt.Errorf("invalid destination balance format: %w", err)
	}

	// Check if source account has sufficient balance
	if sourceBalance < amount {
		return fmt.Errorf("insufficient balance in source account %d", req.SourceAccountID)
	}

	// Calculate new balances
	newSourceBalance := sourceBalance - amount
	newDestBalance := destBalance + amount

	// Create transaction record
	transaction, err := s.transactionRepo.Create(tx, req.SourceAccountID, req.DestinationAccountID, models.Decimal(req.Amount))
	if err != nil {
		return fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Update source account balance
	sourceBalanceStr := fmt.Sprintf("%.10f", newSourceBalance)
	if err := s.updateAccountBalanceInTx(tx, req.SourceAccountID, models.Decimal(sourceBalanceStr)); err != nil {
		return fmt.Errorf("failed to update source account balance: %w", err)
	}

	// Update destination account balance
	destBalanceStr := fmt.Sprintf("%.10f", newDestBalance)
	if err := s.updateAccountBalanceInTx(tx, req.DestinationAccountID, models.Decimal(destBalanceStr)); err != nil {
		return fmt.Errorf("failed to update destination account balance: %w", err)
	}

	// Update transaction status to completed
	if err := s.transactionRepo.UpdateStatus(tx, transaction.ID, models.TransactionStatusCompleted); err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// updateAccountBalanceInTx updates account balance within a database transaction
func (s *TransactionService) updateAccountBalanceInTx(tx *sql.Tx, accountID int64, newBalance models.Decimal) error {
	query := `UPDATE accounts SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE account_id = $2`
	result, err := tx.Exec(query, newBalance, accountID)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}
	return nil
}

