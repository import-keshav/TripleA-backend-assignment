package service

import (
	"fmt"

	"triplea-backend-assignment/models"
	"triplea-backend-assignment/repository"
)

type AccountService struct {
	accountRepo *repository.AccountRepository
}

func NewAccountService(accountRepo *repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

func (s *AccountService) CreateAccount(req *models.CreateAccountRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	exists, err := s.accountRepo.Exists(req.AccountID)
	if err != nil {
		return fmt.Errorf("failed to check account existence: %w", err)
	}
	if exists {
		return fmt.Errorf("account with ID %d already exists", req.AccountID)
	}

	balance := models.Decimal(req.InitialBalance)
	if err := s.accountRepo.Create(req.AccountID, balance); err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

func (s *AccountService) GetAccount(accountID int64) (*models.Account, error) {
	if accountID <= 0 {
		return nil, fmt.Errorf("account_id must be a positive integer")
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

