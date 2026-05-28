package transactions

import "context"

type MockTransactionService struct {
	CreateFunc func(ctx context.Context, tx Transaction) (int64, error)
	ListFunc   func(ctx context.Context) ([]Transaction, error)
	UpdateFunc func(ctx context.Context, tx Transaction) error
	DeleteFunc func(ctx context.Context, id int64) error
}

var _ TransactionServiceInterface = (*MockTransactionService)(nil)

func (m *MockTransactionService) Create(ctx context.Context, tx Transaction) (int64, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, tx)
	}
	return 0, nil
}

func (m *MockTransactionService) List(ctx context.Context) ([]Transaction, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

func (m *MockTransactionService) Update(ctx context.Context, tx Transaction) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, tx)
	}
	return nil
}

func (m *MockTransactionService) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
