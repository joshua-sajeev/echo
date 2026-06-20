package goals

import "context"

type MockGoalRepository struct {
	CreateFunc      func(ctx context.Context, goal Goal) (int64, error)
	GetByIDFunc     func(ctx context.Context, id int64) (*Goal, error)
	ListFunc        func(ctx context.Context) ([]Goal, error)
	UpdateFunc      func(ctx context.Context, goal Goal) error
	DeleteFunc      func(ctx context.Context, id int64) error
	AddProgressFunc func(ctx context.Context, id int64, amount int64) error
	ExistsFunc      func(ctx context.Context, id int64) (bool, error)
}

func (m *MockGoalRepository) Create(ctx context.Context, goal Goal) (int64, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, goal)
	}
	return 0, nil
}

func (m *MockGoalRepository) GetByID(ctx context.Context, id int64) (*Goal, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockGoalRepository) List(ctx context.Context) ([]Goal, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

func (m *MockGoalRepository) Update(ctx context.Context, goal Goal) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, goal)
	}
	return nil
}

func (m *MockGoalRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockGoalRepository) AddProgress(ctx context.Context, id int64, amount int64) error {
	if m.AddProgressFunc != nil {
		return m.AddProgressFunc(ctx, id, amount)
	}
	return nil
}

func (m *MockGoalRepository) Exists(ctx context.Context, id int64) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, id)
	}
	return false, nil
}
