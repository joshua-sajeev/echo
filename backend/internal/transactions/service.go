package transactions

import (
	"context"
	"errors"
)

var (
	ErrTransactionNameRequired  = errors.New("name is required")
	ErrTransactionTypeRequired  = errors.New("type is required")
	ErrTransactionAmountInvalid = errors.New("amount must be greater than 0")
	ErrTransactionSameAccount   = errors.New("from and to account cannot be the same")
	ErrInvalidTransactionID     = errors.New("invalid transaction id")
	ErrTransactionNotFound      = errors.New("transaction not found")
)

type TransactionService struct {
	repo TransactionRepositoryInterface
}

type TransactionServiceInterface interface {
	Create(ctx context.Context, request CreateTransactionRequest) (int64, error)
	List(ctx context.Context) ([]Transaction, error)
	Update(ctx context.Context, id int64, request UpdateTransactionRequest) error
	Delete(ctx context.Context, id int64) error
}

var _ TransactionServiceInterface = (*TransactionService)(nil)

func NewTransactionService(repo TransactionRepositoryInterface) *TransactionService {
	return &TransactionService{repo: repo}
}

// Create — accept DTO, map to Transaction
func (s *TransactionService) Create(ctx context.Context, request CreateTransactionRequest) (int64, error) {
	if request.Name == "" {
		return 0, ErrTransactionNameRequired
	}
	if request.Type == "" {
		return 0, ErrTransactionTypeRequired
	}
	if request.Amount <= 0 {
		return 0, ErrTransactionAmountInvalid
	}
	if request.FromAccountID != nil && request.ToAccountID != nil &&
		*request.FromAccountID == *request.ToAccountID {
		return 0, ErrTransactionSameAccount
	}
	tx := Transaction{
		Name:           request.Name,
		Type:           request.Type, // plain string, no TransactionType cast
		Amount:         request.Amount,
		Date:           request.Date,
		FromAccountID:  request.FromAccountID,
		ToAccountID:    request.ToAccountID,
		Category:       request.Category,
		JarID:          request.JarID,
		IsMasterIncome: request.IsMasterIncome,
	}
	return s.repo.Create(ctx, tx)
}

func (s *TransactionService) List(ctx context.Context) ([]Transaction, error) {
	return s.repo.List(ctx)
}

func (s *TransactionService) Update(ctx context.Context, id int64, request UpdateTransactionRequest) error {
	if id <= 0 {
		return ErrInvalidTransactionID
	}

	// fetch existing so we only overwrite provided fields
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if request.Name != nil {
		if *request.Name == "" {
			return ErrTransactionNameRequired
		}
		existing.Name = *request.Name
	}
	if request.Type != nil {
		existing.Type = *request.Type
	}
	if request.Amount != nil {
		if *request.Amount <= 0 {
			return ErrTransactionAmountInvalid
		}
		existing.Amount = *request.Amount
	}
	if request.Date != nil {
		existing.Date = *request.Date
	}
	if request.FromAccountID != nil {
		existing.FromAccountID = request.FromAccountID
	}
	if request.ToAccountID != nil {
		existing.ToAccountID = request.ToAccountID
	}
	if request.Category != nil {
		existing.Category = request.Category
	}
	if request.JarID != nil {
		existing.JarID = request.JarID
	}
	if request.IsMasterIncome != nil {
		existing.IsMasterIncome = *request.IsMasterIncome
	}

	if existing.FromAccountID != nil && existing.ToAccountID != nil &&
		*existing.FromAccountID == *existing.ToAccountID {
		return ErrTransactionSameAccount
	}

	return s.repo.Update(ctx, *existing)
}

func (s *TransactionService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidTransactionID
	}
	return s.repo.Delete(ctx, id)
}
