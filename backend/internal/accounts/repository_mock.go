package accounts

import (
	"context"
)

type MockAccountRepo struct {
	CreateFn                   func(context.Context, string) (int64, error)
	ListFn                     func(context.Context) ([]Account, error)
	ListWithBalancesFn         func(context.Context) ([]AccountWithBalance, error)
	ListArchivedWithBalancesFn func(context.Context) ([]AccountWithBalance, error)
	RenameFn                   func(context.Context, int64, string) error
	ArchiveFn                  func(context.Context, int64) error
	UnarchiveFn                func(context.Context, int64) error
	ExistsFn                   func(context.Context, int64) (bool, error)
}

func (m *MockAccountRepo) Create(ctx context.Context, name string) (int64, error) {
	return m.CreateFn(ctx, name)
}

func (m *MockAccountRepo) List(ctx context.Context) ([]Account, error) {
	return m.ListFn(ctx)
}

func (m *MockAccountRepo) ListWithBalances(ctx context.Context) ([]AccountWithBalance, error) {
	return m.ListWithBalancesFn(ctx)
}

func (m *MockAccountRepo) ListArchivedWithBalances(ctx context.Context) ([]AccountWithBalance, error) {
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

func (m *MockAccountRepo) Exists(ctx context.Context, id int64) (bool, error) {
	return m.ExistsFn(ctx, id)
}
