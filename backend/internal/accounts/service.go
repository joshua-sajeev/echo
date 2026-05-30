package accounts

import (
	"context"
	"strings"
)

type AccountServiceInterface interface {
	Create(ctx context.Context, name string) (int64, error)

	List(ctx context.Context) ([]Account, error)

	ListWithBalances(ctx context.Context) ([]AccountWithBalance, error)

	ListArchivedWithBalances(ctx context.Context) ([]AccountWithBalance, error)
	Rename(ctx context.Context, id int64, name string) error

	Archive(ctx context.Context, id int64) error
	Unarchive(ctx context.Context, id int64) error
}

type AccountService struct {
	repo AccountRepositoryInterface
}

// NewAccountService creates a new account service
func NewAccountService(repo AccountRepositoryInterface) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

// Create creates a new account
func (s *AccountService) Create(ctx context.Context, name string) (int64, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return 0, ErrInvalidAccountName
	}

	return s.repo.Create(ctx, name)
}

// List returns all active accounts
func (s *AccountService) List(ctx context.Context) ([]Account, error) {
	return s.repo.List(ctx)
}

// ListWithBalances returns active accounts with balances
func (s *AccountService) ListWithBalances(ctx context.Context) ([]AccountWithBalance, error) {
	return s.repo.ListWithBalances(ctx)
}

// ListArchivedWithBalances returns archived accounts with balances
func (s *AccountService) ListArchivedWithBalances(ctx context.Context) ([]AccountWithBalance, error) {
	return s.repo.ListArchivedWithBalances(ctx)
}

// Rename updates an account name
func (s *AccountService) Rename(ctx context.Context, id int64, name string) error {
	if id <= 0 {
		return ErrInvalidAccountID
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return ErrInvalidAccountName
	}

	return s.repo.Rename(ctx, id, name)
}

// Archive archives an account
func (s *AccountService) Archive(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidAccountID
	}

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return ErrAccountNotFound
	}

	return s.repo.Archive(ctx, id)
}

// Unarchive restores an archived account
func (s *AccountService) Unarchive(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidAccountID
	}

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return ErrAccountNotFound
	}

	return s.repo.Unarchive(ctx, id)
}
