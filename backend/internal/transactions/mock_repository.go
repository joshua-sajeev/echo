// Package transactions
package transactions

import "context"

type MockTransactionRepo struct {
	CreateFunc func(ctx context.Context, tx Transaction) (int64, error)
	ListFunc   func(ctx context.Context) ([]Transaction, error)
	UpdateFunc func(ctx context.Context, tx Transaction) error
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockTransactionRepo) Create(ctx context.Context, tx Transaction) (int64, error) {
	return m.CreateFunc(ctx, tx)
}

func (m *MockTransactionRepo) List(ctx context.Context) ([]Transaction, error) {
	return m.ListFunc(ctx)
}

func (m *MockTransactionRepo) Update(ctx context.Context, tx Transaction) error {
	return m.UpdateFunc(ctx, tx)
}

func (m *MockTransactionRepo) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
