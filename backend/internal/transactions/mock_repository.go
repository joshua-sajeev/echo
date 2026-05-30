// Package transactions
package transactions

import "context"

type MockTransactionRepo struct {
	CreateFunc  func(context.Context, Transaction) (int64, error)
	UpdateFunc  func(context.Context, Transaction) error
	DeleteFunc  func(context.Context, int64) error
	ListFunc    func(context.Context) ([]Transaction, error)
	GetByIDFunc func(context.Context, int64) (*Transaction, error)
}

func (m *MockTransactionRepo) Create(ctx context.Context, tx Transaction) (int64, error) {
	return m.CreateFunc(ctx, tx)
}

func (m *MockTransactionRepo) List(ctx context.Context) ([]Transaction, error) {
	return m.ListFunc(ctx)
}

func (m *MockTransactionRepo) GetByID(ctx context.Context, id int64) (*Transaction, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}

	return nil, nil
}

func (m *MockTransactionRepo) Update(ctx context.Context, tx Transaction) error {
	return m.UpdateFunc(ctx, tx)
}

func (m *MockTransactionRepo) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
