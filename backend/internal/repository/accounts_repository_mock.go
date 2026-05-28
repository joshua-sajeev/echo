// Package repository provides mock implementations for testing
package repository

import (
	"context"

	"github.com/joshu-sajeev/echo/internal/models"
)

type MockAccountRepo struct {
	CreateFn                   func(ctx context.Context, name string) (int64, error)
	ListFn                     func(ctx context.Context) ([]models.Account, error)
	ListWithBalancesFn         func(ctx context.Context) ([]models.AccountWithBalance, error)
	ListArchivedWithBalancesFn func(ctx context.Context) ([]models.AccountWithBalance, error)
	RenameFn                   func(ctx context.Context, id int64, name string) error
	ArchiveFn                  func(ctx context.Context, id int64) error
	UnarchiveFn                func(ctx context.Context, id int64) error
}

func (m *MockAccountRepo) Create(ctx context.Context, name string) (int64, error) {
	return m.CreateFn(ctx, name)
}

func (m *MockAccountRepo) List(ctx context.Context) ([]models.Account, error) {
	return m.ListFn(ctx)
}

func (m *MockAccountRepo) ListWithBalances(ctx context.Context) ([]models.AccountWithBalance, error) {
	return m.ListWithBalancesFn(ctx)
}

func (m *MockAccountRepo) ListArchivedWithBalances(ctx context.Context) ([]models.AccountWithBalance, error) {
	return m.ListArchivedWithBalancesFn(ctx)
}

func (m *MockAccountRepo) Rename(ctx context.Context, id int64, name string) error {
	return m.RenameFn(ctx, id, name)
}

func (m *MockAccountRepo) Archive(ctx context.Context, id int64) error {
	return m.ArchiveFn(ctx, id)
}

func (m *MockAccountRepo) Unarchive(ctx context.Context, id int64) error {
	return m.UnarchiveFn(ctx, id)
}
