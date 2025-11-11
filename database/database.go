package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"triplea-backend-assignment/config"
)

var DB *sql.DB

func Connect(cfg *config.Config) error {
	dsn := cfg.GetDSN()
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS accounts (
			account_id BIGINT PRIMARY KEY,
			balance DECIMAL(20, 10) NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id BIGSERIAL PRIMARY KEY,
			source_account_id BIGINT NOT NULL,
			destination_account_id BIGINT NOT NULL,
			amount DECIMAL(20, 10) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (source_account_id) REFERENCES accounts(account_id),
			FOREIGN KEY (destination_account_id) REFERENCES accounts(account_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_source_account ON transactions(source_account_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_destination_account ON transactions(destination_account_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}

