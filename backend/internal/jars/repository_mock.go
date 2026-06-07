package jars

import "context"

type MockJarRepository struct {
	CreateFunc                  func(ctx context.Context, jar Jar) (int64, error)
	ListFunc                    func(ctx context.Context) ([]Jar, error)
	GetByIDFunc                 func(ctx context.Context, id int64) (Jar, error)
	UpdateFunc                  func(ctx context.Context, jar Jar) error
	DeleteFunc                  func(ctx context.Context, id int64) error
	GetAllJarBalancesFunc       func(ctx context.Context) (map[int64]int64, error)
	GetSpentThisMonthPerJarFunc func(ctx context.Context) (map[int64]int64, error)
}

func (m *MockJarRepository) Create(ctx context.Context, jar Jar) (int64, error) {
	return m.CreateFunc(ctx, jar)
}

func (m *MockJarRepository) List(ctx context.Context) ([]Jar, error) {
	return m.ListFunc(ctx)
}

func (m *MockJarRepository) GetByID(ctx context.Context, id int64) (Jar, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}

	// Fallback: search inside ListFunc results so old tests still work.
	if m.ListFunc != nil {
		jars, err := m.ListFunc(ctx)
		if err == nil {
			for _, j := range jars {
				if j.ID == id {
					return j, nil
				}
			}
			return Jar{}, ErrJarNotFound
		}
	}

	return Jar{ID: id}, nil
}

func (m *MockJarRepository) Update(ctx context.Context, jar Jar) error {
	return m.UpdateFunc(ctx, jar)
}

func (m *MockJarRepository) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockJarRepository) GetAllJarBalances(ctx context.Context) (map[int64]int64, error) {
	if m.GetAllJarBalancesFunc != nil {
		return m.GetAllJarBalancesFunc(ctx)
	}
	// Default: return zero balance for every jar returned by ListFunc.
	balances := make(map[int64]int64)
	if m.ListFunc != nil {
		jars, err := m.ListFunc(ctx)
		if err == nil {
			for _, j := range jars {
				balances[j.ID] = 0
			}
		}
	}
	return balances, nil
}

func (m *MockJarRepository) GetSpentThisMonthPerJar(ctx context.Context) (map[int64]int64, error) {
	if m.GetSpentThisMonthPerJarFunc != nil {
		return m.GetSpentThisMonthPerJarFunc(ctx)
	}
	// Default: no spending this month.
	return make(map[int64]int64), nil
}
