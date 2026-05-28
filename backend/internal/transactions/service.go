package transactions

import (
	"context"
	"fmt"
)

type TransactionService struct {
	repo TransactionRepositoryInterface
}

type TransactionServiceInterface interface {
	Create(ctx context.Context, tx Transaction) (int64, error)
	List(ctx context.Context) ([]Transaction, error)
	Update(ctx context.Context, tx Transaction) error
	Delete(ctx context.Context, id int64) error
}

var _ TransactionServiceInterface = (*TransactionService)(nil)

func NewTransactionService(repo TransactionRepositoryInterface) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Create(ctx context.Context, tx Transaction) (int64, error) {
	if tx.Name == "" {
		return 0, fmt.Errorf("name is required")
	}

	if tx.Type == "" {
		return 0, fmt.Errorf("type is required")
	}

	if tx.Amount <= 0 {
		return 0, fmt.Errorf("amount must be greater than 0")
	}

	// optional rule: prevent invalid account mapping
	if tx.FromAccountID != nil && tx.ToAccountID != nil &&
		*tx.FromAccountID == *tx.ToAccountID {
		return 0, fmt.Errorf("from and to account cannot be same")
	}

	return s.repo.Create(ctx, tx)
}

func (s *TransactionService) List(ctx context.Context) ([]Transaction, error) {
	return s.repo.List(ctx)
}

func (s *TransactionService) Update(ctx context.Context, tx Transaction) error {
	if tx.ID <= 0 {
		return fmt.Errorf("invalid id")
	}

	if tx.Name == "" {
		return fmt.Errorf("name is required")
	}

	if tx.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	return s.repo.Update(ctx, tx)
}

func (s *TransactionService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid id")
	}

	return s.repo.Delete(ctx, id)
}
