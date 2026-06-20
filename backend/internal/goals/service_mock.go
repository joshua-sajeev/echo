package goals

import "context"

type MockGoalService struct {
	CreateFunc      func(ctx context.Context, request CreateGoalRequest) (int64, error)
	ListFunc        func(ctx context.Context) ([]GoalWithProgress, error)
	GetByIDFunc     func(ctx context.Context, id int64) (*GoalWithProgress, error)
	UpdateFunc      func(ctx context.Context, id int64, request UpdateGoalRequest) error
	AddProgressFunc func(ctx context.Context, id int64, amount int64) error
	DeleteFunc      func(ctx context.Context, id int64) error
}

var _ GoalServiceInterface = (*MockGoalService)(nil)

func (m *MockGoalService) Create(ctx context.Context, request CreateGoalRequest) (int64, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, request)
	}
	return 0, nil
}

func (m *MockGoalService) List(ctx context.Context) ([]GoalWithProgress, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

func (m *MockGoalService) GetByID(ctx context.Context, id int64) (*GoalWithProgress, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockGoalService) Update(ctx context.Context, id int64, request UpdateGoalRequest) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, request)
	}
	return nil
}

func (m *MockGoalService) AddProgress(ctx context.Context, id int64, amount int64) error {
	if m.AddProgressFunc != nil {
		return m.AddProgressFunc(ctx, id, amount)
	}
	return nil
}

func (m *MockGoalService) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
