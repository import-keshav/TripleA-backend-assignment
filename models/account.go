package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Decimal represents a decimal number as a string to maintain precision
type Decimal string

// Value implements the driver.Valuer interface for database storage
func (d Decimal) Value() (driver.Value, error) {
	if d == "" {
		return nil, nil
	}
	// Validate that it's a valid decimal number
	_, err := strconv.ParseFloat(string(d), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid decimal value: %w", err)
	}
	return string(d), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (d *Decimal) Scan(value interface{}) error {
	if value == nil {
		*d = ""
		return nil
	}

	switch v := value.(type) {
	case string:
		*d = Decimal(v)
	case []byte:
		*d = Decimal(v)
	case float64:
		*d = Decimal(strconv.FormatFloat(v, 'f', -1, 64))
	default:
		return errors.New("cannot scan non-string value into Decimal")
	}
	return nil
}

// String returns the string representation of the decimal
func (d Decimal) String() string {
	return string(d)
}

// Float64 converts the decimal to float64 (may lose precision)
func (d Decimal) Float64() (float64, error) {
	return strconv.ParseFloat(string(d), 64)
}

// MarshalJSON implements json.Marshaler
func (d Decimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(d))
}

// UnmarshalJSON implements json.Unmarshaler
func (d *Decimal) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*d = Decimal(s)
	return nil
}

// Account represents an account in the system
type Account struct {
	AccountID int64   `json:"account_id" db:"account_id"`
	Balance   Decimal `json:"balance" db:"balance"`
}

// CreateAccountRequest represents the request body for creating an account
type CreateAccountRequest struct {
	AccountID     int64  `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

// Validate validates the create account request
func (r *CreateAccountRequest) Validate() error {
	if r.AccountID <= 0 {
		return errors.New("account_id must be a positive integer")
	}
	if r.InitialBalance == "" {
		return errors.New("initial_balance is required")
	}
	// Validate that initial_balance is a valid decimal
	_, err := strconv.ParseFloat(r.InitialBalance, 64)
	if err != nil {
		return fmt.Errorf("initial_balance must be a valid decimal number: %w", err)
	}
	balance, _ := strconv.ParseFloat(r.InitialBalance, 64)
	if balance < 0 {
		return errors.New("initial_balance cannot be negative")
	}
	return nil
}

