package handlers

import (
	"strings"
)

func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "validation error") ||
		strings.Contains(err.Error(), "must be") ||
		strings.Contains(err.Error(), "is required") ||
		strings.Contains(err.Error(), "cannot be")
}

func isAccountNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "account not found") ||
		strings.Contains(err.Error(), "Account not found")
}

func isAccountExistsError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "already exists")
}

func isInsufficientBalanceError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "insufficient balance")
}

