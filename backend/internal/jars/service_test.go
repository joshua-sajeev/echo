package jars

import (
	"context"
	"testing"

	"github.com/joshu-sajeev/echo/internal/transactions"
)

// newTestJarService builds a JarService with a stub txRepo that returns zero
// income (sufficient for all unit tests that don't exercise allocation math).
func newTestJarService(mockRepo *MockJarRepository) *JarService {
	txRepo := &transactions.MockTransactionRepo{
		GetCurrentMonthIncomeFunc: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
	}
	return NewJarService(mockRepo, txRepo)
}

func TestJarService_CreateJar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		input   CreateJarRequest
		mock    func(*MockJarRepository)
		wantErr bool
	}{
		{
			name: "empty name",
			input: CreateJarRequest{
				Name:           "",
				AllocationType: string(AllocationPercentage),
				Value:          10,
			},
			mock:    func(m *MockJarRepository) {},
			wantErr: true,
		},
		{
			name: "percentage must be positive",
			input: CreateJarRequest{
				Name:           "Invest",
				AllocationType: string(AllocationPercentage),
				Value:          0,
			},
			mock:    func(m *MockJarRepository) {},
			wantErr: true,
		},
		{
			name: "percentage exceeds 100",
			input: CreateJarRequest{
				Name:           "New",
				AllocationType: string(AllocationPercentage),
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
			name: "remainder value must be zero",
			input: CreateJarRequest{
				Name:           "Necessities",
				AllocationType: string(AllocationRemainder),
				Value:          10,
			},
			mock:    func(m *MockJarRepository) {},
			wantErr: true,
		},
		{
			name: "duplicate remainder jar",
			input: CreateJarRequest{
				Name:           "Second Remainder",
				AllocationType: string(AllocationRemainder),
				Value:          0,
			},
			mock: func(m *MockJarRepository) {
				m.ListFunc = func(ctx context.Context) ([]Jar, error) {
					return []Jar{
						{ID: 1, AllocationType: AllocationRemainder, Value: 0},
					}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "valid percentage jar",
			input: CreateJarRequest{
				Name:           "Savings",
				AllocationType: string(AllocationPercentage),
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
		{
			name: "valid remainder jar",
			input: CreateJarRequest{
				Name:           "Necessities",
				AllocationType: string(AllocationRemainder),
				Value:          0,
			},
			mock: func(m *MockJarRepository) {
				m.ListFunc = func(ctx context.Context) ([]Jar, error) {
					return []Jar{}, nil
				}
				m.CreateFunc = func(ctx context.Context, jar Jar) (int64, error) {
					return 1, nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockJarRepository{}

			// Safe defaults to avoid nil panics
			mockRepo.CreateFunc = func(ctx context.Context, jar Jar) (int64, error) {
				return 0, nil
			}
			mockRepo.ListFunc = func(ctx context.Context) ([]Jar, error) {
				return []Jar{}, nil
			}

			tt.mock(mockRepo)

			service := newTestJarService(mockRepo)

			id, err := service.CreateJar(ctx, tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected err=%v got %v (err: %v)", tt.wantErr, err != nil, err)
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

	service := newTestJarService(mockRepo)

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

	newName := "Updated"
	newType := string(AllocationPercentage)
	newValue := int64(20)

	mockRepo := &MockJarRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (Jar, error) {
			return Jar{
				ID:             1,
				Name:           "Emergency",
				AllocationType: AllocationPercentage,
				Value:          20,
			}, nil
		},
		ListFunc: func(ctx context.Context) ([]Jar, error) {
			// Return the other jar (id=2) so percentage check works
			return []Jar{
				{
					ID:             2,
					AllocationType: AllocationPercentage,
					Value:          30,
				},
			}, nil
		},
		UpdateFunc: func(ctx context.Context, jar Jar) error {
			if jar.ID != 1 {
				t.Fatalf("expected id 1, got %d", jar.ID)
			}
			if jar.Name != "Updated" {
				t.Fatalf("expected name 'Updated', got %q", jar.Name)
			}
			if jar.Value != 20 {
				t.Fatalf("expected value 20, got %d", jar.Value)
			}
			return nil
		},
	}

	service := newTestJarService(mockRepo)

	err := service.UpdateJar(ctx, 1, UpdateJarRequest{
		Name:           &newName,
		AllocationType: &newType,
		Value:          &newValue,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJarService_UpdateJar_InvalidID(t *testing.T) {
	service := newTestJarService(&MockJarRepository{})

	err := service.UpdateJar(context.Background(), 0, UpdateJarRequest{})
	if err != ErrInvalidJarID {
		t.Fatalf("expected ErrInvalidJarID, got %v", err)
	}
}

func TestJarService_DeleteJar(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockJarRepository{
		DeleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	service := newTestJarService(mockRepo)

	err := service.DeleteJar(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJarService_DeleteJar_InvalidID(t *testing.T) {
	service := newTestJarService(&MockJarRepository{})

	err := service.DeleteJar(context.Background(), 0)
	if err != ErrInvalidJarID {
		t.Fatalf("expected ErrInvalidJarID, got %v", err)
	}
}
