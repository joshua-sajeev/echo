package jars

import "context"

type MockJarService struct {
	CreateJarFunc func(ctx context.Context, jar CreateJarRequest) (int64, error)
	ListJarsFunc  func(ctx context.Context) ([]Jar, error)
	UpdateJarFunc func(ctx context.Context, id int64, jar UpdateJarRequest) error
	DeleteJarFunc func(ctx context.Context, id int64) error
}

func (m *MockJarService) CreateJar(ctx context.Context, jar CreateJarRequest) (int64, error) {
	return m.CreateJarFunc(ctx, jar)
}

func (m *MockJarService) ListJars(ctx context.Context) ([]Jar, error) {
	return m.ListJarsFunc(ctx)
}

func (m *MockJarService) UpdateJar(ctx context.Context, id int64, jar UpdateJarRequest) error {
	return m.UpdateJarFunc(ctx, id, jar)
}

func (m *MockJarService) DeleteJar(ctx context.Context, id int64) error {
	return m.DeleteJarFunc(ctx, id)
}
