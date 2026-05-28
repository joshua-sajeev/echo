package jars

import (
	"context"
	"testing"
)

func TestJarService_CreateJar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		input   Jar
		mock    func(*MockJarRepository)
		wantErr bool
	}{
		{
			name: "empty name",
			input: Jar{
				Name:           "",
				AllocationType: AllocationFixed,
				Value:          100,
			},
			mock:    func(m *MockJarRepository) {},
			wantErr: true,
		},
		{
			name: "percentage must be positive",
			input: Jar{
				Name:           "Invest",
				AllocationType: AllocationPercentage,
				Value:          0,
			},
			mock:    func(m *MockJarRepository) {},
			wantErr: true,
		},
		{
			name: "percentage exceeds 100",
			input: Jar{
				Name:           "New",
				AllocationType: AllocationPercentage,
				Value:          60,
			},
			mock: func(m *MockJarRepository) {
				m.ListFunc = func(ctx context.Context) ([]Jar, error) {
					return []Jar{
						{AllocationType: AllocationPercentage, Value: 50},
					}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "valid fixed jar",
			input: Jar{
				Name:           "Rent",
				AllocationType: AllocationFixed,
				Value:          1000,
			},
			mock: func(m *MockJarRepository) {
				m.CreateFunc = func(ctx context.Context, jar Jar) (int64, error) {
					return 1, nil
				}
			},
			wantErr: false,
		},
		{
			name: "valid percentage jar",
			input: Jar{
				Name:           "Savings",
				AllocationType: AllocationPercentage,
				Value:          20,
			},
			mock: func(m *MockJarRepository) {
				m.ListFunc = func(ctx context.Context) ([]Jar, error) {
					return []Jar{
						{AllocationType: AllocationPercentage, Value: 30},
					}, nil
				}

				m.CreateFunc = func(ctx context.Context, jar Jar) (int64, error) {
					return 10, nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockJarRepository{}

			// safe defaults (avoid nil panic)
			mockRepo.CreateFunc = func(ctx context.Context, jar Jar) (int64, error) {
				return 0, nil
			}
			mockRepo.ListFunc = func(ctx context.Context) ([]Jar, error) {
				return []Jar{}, nil
			}

			tt.mock(mockRepo)

			service := NewJarService(mockRepo)

			id, err := service.CreateJar(ctx, tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected err=%v got %v", tt.wantErr, err)
			}

			if !tt.wantErr && id == 0 {
				t.Fatal("expected valid id")
			}
		})
	}
}

func TestJarService_ListJars(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockJarRepository{
		ListFunc: func(ctx context.Context) ([]Jar, error) {
			return []Jar{
				{Name: "A"},
				{Name: "B"},
			}, nil
		},
	}

	service := NewJarService(mockRepo)

	jars, err := service.ListJars(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(jars) != 2 {
		t.Fatalf("expected 2 jars got %d", len(jars))
	}
}

func TestJarService_UpdateJar(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockJarRepository{
		ListFunc: func(ctx context.Context) ([]Jar, error) {
			return []Jar{
				{ID: 2, AllocationType: AllocationPercentage, Value: 30},
			}, nil
		},
		UpdateFunc: func(ctx context.Context, jar Jar) error {
			return nil
		},
	}

	service := NewJarService(mockRepo)

	err := service.UpdateJar(ctx, Jar{
		ID:             1,
		Name:           "Updated",
		AllocationType: AllocationFixed,
		Value:          100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJarService_DeleteJar(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockJarRepository{
		DeleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	service := NewJarService(mockRepo)

	err := service.DeleteJar(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
