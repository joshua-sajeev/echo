package jars

import "context"

type MockJarRepository struct {
	CreateFunc func(ctx context.Context, jar Jar) (int64, error)

	ListFunc    func(ctx context.Context) ([]Jar, error)
	GetByIDFunc func(ctx context.Context, id int64) (Jar, error)
	UpdateFunc  func(ctx context.Context, jar Jar) error

	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockJarRepository) Create(ctx context.Context, jar Jar) (int64, error) {
	return m.CreateFunc(ctx, jar)
}

func (m *MockJarRepository) List(ctx context.Context) ([]Jar, error) {
	return m.ListFunc(ctx)
}

func (m *MockJarRepository) GetByID(ctx context.Context, id int64) (Jar, error) {
	// 1. If the test explicitly defined a GetByIDFunc, use it.
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}

	// 2. Fallback: If GetByIDFunc is nil, look for the jar inside ListFunc
	// to reuse the data your old tests already defined.
	if m.ListFunc != nil {
		jars, err := m.ListFunc(ctx)
		if err == nil {
			for _, j := range jars {
				if j.ID == id {
					return j, nil
				}
			}
			// If ListFunc is defined but doesn't have the ID, mimic a genuine DB Not Found error
			return Jar{}, ErrJarNotFound
		}
	}

	// 3. Last resort safely formatted default to prevent panic
	return Jar{ID: id}, nil
}

func (m *MockJarRepository) Update(ctx context.Context, jar Jar) error {
	return m.UpdateFunc(ctx, jar)
}

func (m *MockJarRepository) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
