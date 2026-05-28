package accounts

import (
	"context"
)

type MockAccountService struct {
	CreateFn                   func(ctx context.Context, name string) (int64, error)
	ListFn                     func(ctx context.Context) ([]Account, error)
	ListWithBalancesFn         func(ctx context.Context) ([]AccountWithBalance, error)
	ListArchivedWithBalancesFn func(ctx context.Context) ([]AccountWithBalance, error)
	RenameFn                   func(ctx context.Context, id int64, name string) error
	ArchiveFn                  func(ctx context.Context, id int64) error
	UnarchiveFn                func(ctx context.Context, id int64) error
}

func (m *MockAccountService) Create(
	ctx context.Context,
	name string,
) (int64, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, name)
	}

	return 0, nil
}

func (m *MockAccountService) List(
	ctx context.Context,
) ([]Account, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx)
	}

	return nil, nil
}

func (m *MockAccountService) ListWithBalances(
	ctx context.Context,
) ([]AccountWithBalance, error) {
	if m.ListWithBalancesFn != nil {
		return m.ListWithBalancesFn(ctx)
	}

	return nil, nil
}

func (m *MockAccountService) ListArchivedWithBalances(
	ctx context.Context,
) ([]AccountWithBalance, error) {
	if m.ListArchivedWithBalancesFn != nil {
		return m.ListArchivedWithBalancesFn(ctx)
	}

	return nil, nil
}

func (m *MockAccountService) Rename(
	ctx context.Context,
	id int64,
	name string,
) error {
	if m.RenameFn != nil {
		return m.RenameFn(ctx, id, name)
	}

	return nil
}

func (m *MockAccountService) Archive(
	ctx context.Context,
	id int64,
) error {
	if m.ArchiveFn != nil {
		return m.ArchiveFn(ctx, id)
	}

	return nil
}

func (m *MockAccountService) Unarchive(
	ctx context.Context,
	id int64,
) error {
	if m.UnarchiveFn != nil {
		return m.UnarchiveFn(ctx, id)
	}

	return nil
}
