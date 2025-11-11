package models

import (
	"testing"
)

func TestCreateAccountRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateAccountRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateAccountRequest{
				AccountID:     123,
				InitialBalance: "100.50",
			},
			wantErr: false,
		},
		{
			name: "invalid account_id (zero)",
			req: CreateAccountRequest{
				AccountID:     0,
				InitialBalance: "100.50",
			},
			wantErr: true,
		},
		{
			name: "invalid account_id (negative)",
			req: CreateAccountRequest{
				AccountID:     -1,
				InitialBalance: "100.50",
			},
			wantErr: true,
		},
		{
			name: "missing initial_balance",
			req: CreateAccountRequest{
				AccountID:     123,
				InitialBalance: "",
			},
			wantErr: true,
		},
		{
			name: "invalid initial_balance (not a number)",
			req: CreateAccountRequest{
				AccountID:     123,
				InitialBalance: "not-a-number",
			},
			wantErr: true,
		},
		{
			name: "invalid initial_balance (negative)",
			req: CreateAccountRequest{
				AccountID:     123,
				InitialBalance: "-10.00",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

