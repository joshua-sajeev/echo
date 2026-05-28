package jars

import "context"

type MockJarService struct {
	CreateJarFunc func(ctx context.Context, jar Jar) (int64, error)
	ListJarsFunc  func(ctx context.Context) ([]Jar, error)
	UpdateJarFunc func(ctx context.Context, jar Jar) error
	DeleteJarFunc func(ctx context.Context, id int64) error
}

func (m *MockJarService) CreateJar(ctx context.Context, jar Jar) (int64, error) {
	return m.CreateJarFunc(ctx, jar)
}

func (m *MockJarService) ListJars(ctx context.Context) ([]Jar, error) {
	return m.ListJarsFunc(ctx)
}

func (m *MockJarService) UpdateJar(ctx context.Context, jar Jar) error {
	return m.UpdateJarFunc(ctx, jar)
}

func (m *MockJarService) DeleteJar(ctx context.Context, id int64) error {
	return m.DeleteJarFunc(ctx, id)
}
