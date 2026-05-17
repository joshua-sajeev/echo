// Package service provides business logic for account operations.
package service

import (
	"context"
	"errors"
	"strings"

	"github.com/joshu-sajeev/echo/internal/models"
	"github.com/joshu-sajeev/echo/internal/repository"
)

var ErrInvalidAccountName = errors.New("invalid account name")

// AccountService handles account business logic.
type AccountService struct {
	repo repository.AccountRepositoryInterface
}

// NewAccountService creates a new AccountService instance.
func NewAccountService(repo repository.AccountRepositoryInterface) AccountServiceInterface {
	return &AccountService{
		repo: repo,
	}
}

// Create validates and creates a new account.
func (s *AccountService) Create(ctx context.Context, name string) (int64, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return 0, ErrInvalidAccountName
	}

	return s.repo.Create(ctx, name)
}

// List returns all active accounts.
func (s *AccountService) List(ctx context.Context) ([]models.Account, error) {
	return s.repo.List(ctx)
}

// ListWithBalances returns all active accounts
// with their calculated balances.
func (s *AccountService) ListWithBalances(ctx context.Context) ([]models.AccountWithBalance, error) {
	return s.repo.ListWithBalances(ctx)
}

// ListArchivedWithBalances returns archived accounts
// with their calculated balances.
func (s *AccountService) ListArchivedWithBalances(ctx context.Context) ([]models.AccountWithBalance, error) {
	return s.repo.ListArchivedWithBalances(ctx)
}

// Rename validates and updates an account name.
func (s *AccountService) Rename(ctx context.Context, id int64, name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrInvalidAccountName
	}

	return s.repo.Rename(ctx, id, name)
}

// Archive archives an account.
func (s *AccountService) Archive(ctx context.Context, id int64) error {
	return s.repo.Archive(ctx, id)
}

// Unarchive restores an archived account.
func (s *AccountService) Unarchive(ctx context.Context, id int64) error {
	return s.repo.Unarchive(ctx, id)
}
