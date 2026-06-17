package transactions

import "context"

type MockTransactionService struct {
	CreateFunc func(ctx context.Context, request CreateTransactionRequest) (int64, error)
	ListFunc   func(ctx context.Context) ([]TransactionListItem, error)
	UpdateFunc func(ctx context.Context, id int64, request UpdateTransactionRequest) error

	GetByIDFunc func(ctx context.Context, id int64) (*Transaction, error)
	DeleteFunc  func(ctx context.Context, id int64) error
}

var _ TransactionServiceInterface = (*MockTransactionService)(nil)

func (m *MockTransactionService) Create(ctx context.Context, request CreateTransactionRequest) (int64, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, request)
	}
	return 0, nil
}

func (m *MockTransactionService) List(ctx context.Context) ([]TransactionListItem, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

func (m *MockTransactionService) Update(ctx context.Context, id int64, request UpdateTransactionRequest) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, request)
	}
	return nil
}

func (m *MockTransactionService) GetByID(ctx context.Context, id int64) (*Transaction, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTransactionService) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
