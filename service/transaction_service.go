package service

import (
	"database/sql"
	"fmt"

	"triplea-backend-assignment/database"
	"triplea-backend-assignment/models"
	"triplea-backend-assignment/repository"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
}

func NewTransactionService(
	transactionRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (s *TransactionService) ProcessTransaction(req *models.CreateTransactionRequest) error {
	if err := s.validateBeforeTxn(req); err != nil {
		return err
	}

	amount, err := models.Decimal(req.Amount).Float64()
	if err != nil {
		return fmt.Errorf("invalid amount format: %w", err)
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	sourceAccount, err := s.accountRepo.GetByIDWithLock(tx, req.SourceAccountID)
	if err != nil {
		return fmt.Errorf("failed to get source account: %w", err)
	}

	destAccount, err := s.accountRepo.GetByIDWithLock(tx, req.DestinationAccountID)
	if err != nil {
		return fmt.Errorf("failed to get destination account: %w", err)
	}

	sourceBalance, err := sourceAccount.Balance.Float64()
	if err != nil {
		return fmt.Errorf("invalid source balance format: %w", err)
	}

	destBalance, err := destAccount.Balance.Float64()
	if err != nil {
		return fmt.Errorf("invalid destination balance format: %w", err)
	}

	if sourceBalance < amount {
		return fmt.Errorf("insufficient balance in source account %d", req.SourceAccountID)
	}

	newSourceBalance := sourceBalance - amount
	newDestBalance := destBalance + amount

	transaction, err := s.transactionRepo.Create(tx, req.SourceAccountID, req.DestinationAccountID, models.Decimal(req.Amount))
	if err != nil {
		return fmt.Errorf("failed to create transaction record: %w", err)
	}

	sourceBalanceStr := fmt.Sprintf("%.10f", newSourceBalance)
	if err := s.updateAccountBalanceInTx(tx, req.SourceAccountID, models.Decimal(sourceBalanceStr)); err != nil {
		return fmt.Errorf("failed to update source account balance: %w", err)
	}

	destBalanceStr := fmt.Sprintf("%.10f", newDestBalance)
	if err := s.updateAccountBalanceInTx(tx, req.DestinationAccountID, models.Decimal(destBalanceStr)); err != nil {
		return fmt.Errorf("failed to update destination account balance: %w", err)
	}

	if err := s.transactionRepo.UpdateStatus(tx, transaction.ID, models.TransactionStatusCompleted); err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *TransactionService) validateBeforeTxn(req *models.CreateTransactionRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	sourceExists, err := s.accountRepo.Exists(req.SourceAccountID)
	if err != nil {
		return fmt.Errorf("failed to check source account existence: %w", err)
	}
	if !sourceExists {
		return fmt.Errorf("source account %d not found", req.SourceAccountID)
	}

	destExists, err := s.accountRepo.Exists(req.DestinationAccountID)
	if err != nil {
		return fmt.Errorf("failed to check destination account existence: %w", err)
	}
	if !destExists {
		return fmt.Errorf("destination account %d not found", req.DestinationAccountID)
	}

	return nil
}

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

