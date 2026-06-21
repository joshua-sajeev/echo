package goals

import "context"

type MockGoalService struct {
	CreateFunc              func(ctx context.Context, request CreateGoalRequest) (int64, error)
	CreateWithRebalanceFunc func(ctx context.Context, request CreateGoalWithRebalanceRequest) (*RebalanceSummary, error)
	ListFunc                func(ctx context.Context) ([]GoalWithProgress, error)
	GetByIDFunc             func(ctx context.Context, id int64) (*GoalWithProgress, error)
	UpdateFunc              func(ctx context.Context, id int64, request UpdateGoalRequest) error
	AddProgressFunc         func(ctx context.Context, id int64, amount int64) error
	ArchiveFunc             func(ctx context.Context, id int64) error
	RestoreFunc             func(ctx context.Context, id int64) error
}

var _ GoalServiceInterface = (*MockGoalService)(nil)

func (m *MockGoalService) Create(ctx context.Context, request CreateGoalRequest) (int64, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, request)
	}
	return 0, nil
}

func (m *MockGoalService) CreateWithRebalance(ctx context.Context, request CreateGoalWithRebalanceRequest) (*RebalanceSummary, error) {
	if m.CreateWithRebalanceFunc != nil {
		return m.CreateWithRebalanceFunc(ctx, request)
	}
	return nil, nil
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

func (m *MockGoalService) Archive(ctx context.Context, id int64) error {
	if m.ArchiveFunc != nil {
		return m.ArchiveFunc(ctx, id)
	}
	return nil
}

func (m *MockGoalService) Restore(ctx context.Context, id int64) error {
	if m.RestoreFunc != nil {
		return m.RestoreFunc(ctx, id)
	}
	return nil
}
