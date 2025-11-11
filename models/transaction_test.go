package models

import (
	"testing"
)

func TestCreateTransactionRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateTransactionRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:              "50.25",
			},
			wantErr: false,
		},
		{
			name: "invalid source_account_id (zero)",
			req: CreateTransactionRequest{
				SourceAccountID:      0,
				DestinationAccountID: 456,
				Amount:              "50.25",
			},
			wantErr: true,
		},
		{
			name: "invalid destination_account_id (zero)",
			req: CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 0,
				Amount:              "50.25",
			},
			wantErr: true,
		},
		{
			name: "same source and destination",
			req: CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 123,
				Amount:              "50.25",
			},
			wantErr: true,
		},
		{
			name: "missing amount",
			req: CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:              "",
			},
			wantErr: true,
		},
		{
			name: "invalid amount (not a number)",
			req: CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:              "not-a-number",
			},
			wantErr: true,
		},
		{
			name: "invalid amount (zero)",
			req: CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:              "0",
			},
			wantErr: true,
		},
		{
			name: "invalid amount (negative)",
			req: CreateTransactionRequest{
				SourceAccountID:      123,
				DestinationAccountID: 456,
				Amount:              "-10.00",
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

