package transactions

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joshu-sajeev/echo/internal/utils"
)

type TransactionService struct {
	repo TransactionRepositoryInterface
}

type TransactionServiceInterface interface {
	Create(ctx context.Context, request CreateTransactionRequest) (int64, error)
	List(ctx context.Context) ([]Transaction, error)
	Update(ctx context.Context, id int64, request UpdateTransactionRequest) error
	GetByID(ctx context.Context, id int64) (*Transaction, error)
	Delete(ctx context.Context, id int64) error
}

var _ TransactionServiceInterface = (*TransactionService)(nil)

func NewTransactionService(repo TransactionRepositoryInterface) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Create(ctx context.Context, request CreateTransactionRequest) (int64, error) {
	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" {
		return 0, ErrTransactionNameRequired
	}
	if request.Type == "" {
		return 0, ErrTransactionTypeRequired
	}
	if request.Amount <= 0 {
		return 0, ErrTransactionAmountInvalid
	}
	if request.FromAccountID != nil && request.ToAccountID != nil && *request.FromAccountID == *request.ToAccountID {
		return 0, ErrTransactionSameAccount
	}
	if request.Date.IsZero() {
		return 0, ErrTransactionDateRequired
	}

	tx := Transaction{
		Name:           request.Name,
		Type:           request.Type,
		Amount:         request.Amount,
		Date:           request.Date,
		FromAccountID:  request.FromAccountID,
		ToAccountID:    request.ToAccountID,
		Category:       request.Category,
		JarID:          request.JarID,
		IsMasterIncome: request.IsMasterIncome,
	}

	id, err := s.repo.Create(ctx, tx)
	if err != nil {
		utils.LogError(ctx, "TransactionService.Create", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			switch pgErr.ConstraintName {
			case "transactions_jar_id_fkey":
				return 0, ErrJarNotFound
			case "transactions_from_account_id_fkey", "transactions_to_account_id_fkey":
				return 0, ErrAccountNotFound
			}
		}
		return 0, err
	}
	return id, nil
}

func (s *TransactionService) List(ctx context.Context) ([]Transaction, error) {
	transactions, err := s.repo.List(ctx)
	if err != nil {
		utils.LogError(ctx, "TransactionService.List", err)
		return nil, err
	}
	return transactions, nil
}

func (s *TransactionService) Update(ctx context.Context, id int64, request UpdateTransactionRequest) error {
	if id <= 0 {
		return ErrInvalidTransactionID
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			return ErrTransactionNotFound
		}
		utils.LogError(ctx, "TransactionService.Update (GetByID)", err)
		return err
	}

	if request.Name != nil {
		trimmedName := strings.TrimSpace(*request.Name)
		if trimmedName == "" {
			return ErrTransactionNameRequired
		}
		existing.Name = trimmedName
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

	// Safely evaluate optional nullable references
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

	if request.Type != nil {
		existing.Type = *request.Type

		switch *request.Type {

		case "expense":
			// expense should NOT have destination account
			existing.ToAccountID = nil
			existing.IsMasterIncome = false

		case "income":
			// income should NOT have source account
			existing.FromAccountID = nil

		case "transfer":
			// transfer should NOT have jar or income flags
			existing.JarID = nil
			existing.IsMasterIncome = false
		}
	}

	if existing.Type == "transfer" {
		if existing.FromAccountID != nil &&
			existing.ToAccountID != nil &&
			*existing.FromAccountID == *existing.ToAccountID {
			return ErrTransactionSameAccount
		}
	}

	err = s.repo.Update(ctx, *existing)
	if err != nil {
		utils.LogError(ctx, "TransactionService.Update (Update)", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			switch pgErr.ConstraintName {
			case "transactions_jar_id_fkey":
				return ErrJarNotFound
			case "transactions_from_account_id_fkey", "transactions_to_account_id_fkey":
				return ErrAccountNotFound
			}
		}
		return err
	}
	return nil
}

func (s *TransactionService) GetByID(ctx context.Context, id int64) (*Transaction, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TransactionService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidTransactionID
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			return ErrTransactionNotFound
		}
		utils.LogError(ctx, "TransactionService.Delete", err)
		return err
	}
	return nil
}
