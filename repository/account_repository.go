package repository

import (
	"database/sql"
	"fmt"

	"triplea-backend-assignment/database"
	"triplea-backend-assignment/models"
)

type AccountRepository struct{}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{}
}

func (r *AccountRepository) Create(accountID int64, balance models.Decimal) error {
	query := `INSERT INTO accounts (account_id, balance) VALUES ($1, $2)`
	_, err := database.DB.Exec(query, accountID, balance)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	return nil
}

func (r *AccountRepository) GetByID(accountID int64) (*models.Account, error) {
	query := `SELECT account_id, balance FROM accounts WHERE account_id = $1`
	account := &models.Account{}
	err := database.DB.QueryRow(query, accountID).Scan(&account.AccountID, &account.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

func (r *AccountRepository) UpdateBalance(accountID int64, newBalance models.Decimal) error {
	query := `UPDATE accounts SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE account_id = $2`
	result, err := database.DB.Exec(query, newBalance, accountID)
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

func (r *AccountRepository) Exists(accountID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM accounts WHERE account_id = $1)`
	var exists bool
	err := database.DB.QueryRow(query, accountID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check account existence: %w", err)
	}
	return exists, nil
}

func (r *AccountRepository) GetByIDWithLock(tx *sql.Tx, accountID int64) (*models.Account, error) {
	query := `SELECT account_id, balance FROM accounts WHERE account_id = $1 FOR UPDATE`
	account := &models.Account{}
	err := tx.QueryRow(query, accountID).Scan(&account.AccountID, &account.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

