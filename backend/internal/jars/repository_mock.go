package jars

import "context"

type MockJarRepository struct {
	CreateFunc func(ctx context.Context, jar Jar) (int64, error)

	ListFunc func(ctx context.Context) ([]Jar, error)

	UpdateFunc func(ctx context.Context, jar Jar) error

	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockJarRepository) Create(ctx context.Context, jar Jar) (int64, error) {
	return m.CreateFunc(ctx, jar)
}

func (m *MockJarRepository) List(ctx context.Context) ([]Jar, error) {
	return m.ListFunc(ctx)
}

func (m *MockJarRepository) Update(ctx context.Context, jar Jar) error {
	return m.UpdateFunc(ctx, jar)
}

func (m *MockJarRepository) Delete(ctx context.Context, id int64) error {
	return m.DeleteFunc(ctx, id)
}
