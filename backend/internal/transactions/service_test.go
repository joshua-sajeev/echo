package transactions

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTransactionService_Create(t *testing.T) {
	tests := []struct {
		name       string
		input      CreateTransactionRequest
		mockReturn int64
		mockErr    error
		wantErr    bool
		expectID   int64
	}{
		{
			name: "success create",
			input: CreateTransactionRequest{
				Name:   "Salary",
				Type:   "income",
				Amount: 1000,
				Date:   time.Now(),
			},
			mockReturn: 1,
			expectID:   1,
		},
		{
			name: "missing name",
			input: CreateTransactionRequest{
				Type:   "income",
				Amount: 1000,
				Date:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing type",
			input: CreateTransactionRequest{
				Name:   "Salary",
				Amount: 1000,
				Date:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid amount",
			input: CreateTransactionRequest{
				Name:   "Salary",
				Type:   "income",
				Amount: 0,
				Date:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "same from and to account",
			input: CreateTransactionRequest{
				Name:          "Transfer",
				Type:          "transfer",
				Amount:        100,
				Date:          time.Now(),
				FromAccountID: new(int64(1)),
				ToAccountID:   new(int64(1)),
			},
			wantErr: true,
		},
		{
			name: "repo error",
			input: CreateTransactionRequest{
				Name:   "Salary",
				Type:   "income",
				Amount: 1000,
				Date:   time.Now(),
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockTransactionRepo{
				CreateFunc: func(ctx context.Context, tx Transaction) (int64, error) {
					return tt.mockReturn, tt.mockErr
				},
			}

			service := NewTransactionService(mock)

			id, err := service.Create(context.Background(), tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if id != tt.expectID {
				t.Fatalf("expected id %d, got %d", tt.expectID, id)
			}
		})
	}
}

func TestTransactionService_Update(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		input   UpdateTransactionRequest
		mockErr error
		wantErr bool
	}{
		{
			name: "success update",
			id:   1,
			input: UpdateTransactionRequest{
				Name:   new("Updated"),
				Amount: new(int64(200)),
			},
		},
		{
			name:    "invalid id",
			id:      0,
			input:   UpdateTransactionRequest{},
			wantErr: true,
		},
		{
			name: "empty name",
			id:   1,
			input: UpdateTransactionRequest{
				Name: new(""),
			},
			wantErr: true,
		},
		{
			name: "invalid amount",
			id:   1,
			input: UpdateTransactionRequest{
				Amount: new(int64(0)),
			},
			wantErr: true,
		},
		{
			name: "repo error",
			id:   1,
			input: UpdateTransactionRequest{
				Name: new("Updated"),
			},
			mockErr: errors.New("db fail"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockTransactionRepo{
				GetByIDFunc: func(ctx context.Context, id int64) (*Transaction, error) {
					return &Transaction{
						ID:     id,
						Name:   "Original",
						Type:   "income",
						Amount: 100,
					}, nil
				},
				UpdateFunc: func(ctx context.Context, tx Transaction) error {
					return tt.mockErr
				},
			}

			service := NewTransactionService(mock)

			err := service.Update(
				context.Background(),
				tt.id,
				tt.input,
			)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTransactionService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		mockErr error
		wantErr bool
	}{
		{
			name: "success delete",
			id:   1,
		},
		{
			name:    "invalid id",
			id:      0,
			wantErr: true,
		},
		{
			name:    "repo error",
			id:      1,
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockTransactionRepo{
				DeleteFunc: func(ctx context.Context, id int64) error {
					return tt.mockErr
				},
			}

			service := NewTransactionService(mock)

			err := service.Delete(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTransactionService_List(t *testing.T) {
	expected := []Transaction{
		{ID: 1, Name: "A", Amount: 100},
	}

	mock := &MockTransactionRepo{
		ListFunc: func(ctx context.Context) ([]Transaction, error) {
			return expected, nil
		},
	}

	service := NewTransactionService(mock)

	res, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != len(expected) {
		t.Fatalf("expected %d, got %d", len(expected), len(res))
	}
}
