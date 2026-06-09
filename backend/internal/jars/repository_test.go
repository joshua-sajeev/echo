package jars

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/utils"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	testDB = utils.GetTestDB()

	code := m.Run()

	utils.CleanupTestDB()

	os.Exit(code)
}

func TestJarRepository_Create(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewJarRepository(testDB)

	tests := []struct {
		name    string
		input   Jar
		wantErr bool
	}{
		{
			name: "valid jar percentage",
			input: Jar{
				Name:           "Savings",
				AllocationType: AllocationPercentage,
				Value:          50,
			},
			wantErr: false,
		},
		{
			name: "valid jar remainder",
			input: Jar{
				Name:           "Necessities",
				AllocationType: AllocationRemainder,
				Value:          0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.ResetTables()

			id, err := repo.Create(ctx, tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if id == 0 {
					t.Fatal("expected non-zero id")
				}

				jars, err := repo.List(ctx)
				if err != nil {
					t.Fatalf("List() error = %v", err)
				}

				found := false
				for _, j := range jars {
					if j.ID == id && j.Name == tt.input.Name {
						found = true
					}
				}

				if !found {
					t.Fatal("created jar not found")
				}
			}
		})
	}
}

func TestJarRepository_List(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewJarRepository(testDB)

	inputs := []Jar{
		{
			Name:           "Savings",
			AllocationType: AllocationPercentage,
			Value:          30,
		},
		{
			Name:           "Investments",
			AllocationType: AllocationPercentage,
			Value:          20,
		},
		{
			Name:           "Necessities",
			AllocationType: AllocationRemainder,
			Value:          0,
		},
	}

	var createdIDs []int64

	for _, j := range inputs {
		id, err := repo.Create(ctx, j)
		if err != nil {
			t.Fatalf("setup create failed: %v", err)
		}
		createdIDs = append(createdIDs, id)
	}

	jars, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(jars) < len(inputs) {
		t.Fatalf("expected at least %d jars, got %d", len(inputs), len(jars))
	}

	// Verify all inserted jars exist in result
	for i, input := range inputs {
		found := false

		for _, j := range jars {
			if j.ID == createdIDs[i] &&
				j.Name == input.Name &&
				j.AllocationType == input.AllocationType &&
				j.Value == input.Value {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("jar not found in list: %+v", input)
		}
	}

	// Verify ordering is by id ASC (as per the repository query)
	for i := 1; i < len(jars); i++ {
		if jars[i-1].ID > jars[i].ID {
			t.Fatal("list ordering broken: id not sorted ASC")
		}
	}
}

func TestJarRepository_Update(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewJarRepository(testDB)

	id, err := repo.Create(ctx, Jar{
		Name:           "Old",
		AllocationType: AllocationPercentage,
		Value:          10,
	})
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	err = repo.Update(ctx, Jar{
		ID:             id,
		Name:           "Updated",
		AllocationType: AllocationPercentage,
		Value:          20,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	jars, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	found := false
	for _, j := range jars {
		if j.ID == id {
			if j.Name != "Updated" {
				t.Fatalf("expected name 'Updated', got %q", j.Name)
			}
			if j.Value != 20 {
				t.Fatalf("expected value 20, got %d", j.Value)
			}
			found = true
		}
	}

	if !found {
		t.Fatal("updated jar not reflected in DB")
	}
}

func TestJarRepository_Delete(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewJarRepository(testDB)

	id, err := repo.Create(ctx, Jar{
		Name:           "ToDelete",
		AllocationType: AllocationPercentage,
		Value:          10,
	})
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	err = repo.Delete(ctx, id)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	jars, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	for _, j := range jars {
		if j.ID == id {
			t.Fatal("jar was not deleted")
		}
	}
}
