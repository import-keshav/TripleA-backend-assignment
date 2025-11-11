package handlers

import (
	"strings"
)

// isValidationError checks if an error is a validation error
func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "validation error") ||
		strings.Contains(err.Error(), "must be") ||
		strings.Contains(err.Error(), "is required") ||
		strings.Contains(err.Error(), "cannot be")
}

// isAccountNotFoundError checks if an error is an account not found error
func isAccountNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "account not found") ||
		strings.Contains(err.Error(), "Account not found")
}

// isAccountExistsError checks if an error is an account already exists error
func isAccountExistsError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "already exists")
}

// isInsufficientBalanceError checks if an error is an insufficient balance error
func isInsufficientBalanceError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "insufficient balance")
}

