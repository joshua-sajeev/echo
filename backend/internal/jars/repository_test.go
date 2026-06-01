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
			name: "valid jar fixed",
			input: Jar{
				Name:           "Rent",
				AllocationType: AllocationFixed,
				Value:          1000,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

	// seed data
	inputs := []Jar{
		{
			Name:           "Savings",
			AllocationType: AllocationPercentage,
			Value:          30,
		},
		{
			Name:           "Rent",
			AllocationType: AllocationPercentage,
			Value:          10,
		},
		{
			Name:           "Investment",
			AllocationType: AllocationPercentage,
			Value:          20,
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

	// call list
	jars, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	// basic validation: count match
	if len(jars) < len(inputs) {
		t.Fatalf("expected at least %d jars, got %d", len(inputs), len(jars))
	}

	// verify all inserted jars exist in result
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

	// optional: ensure ordering is correct (priority ASC, id ASC)
	for i := 1; i < len(jars); i++ {
		prev := jars[i-1]
		curr := jars[i]

		if prev.Priority > curr.Priority {
			t.Fatal("list ordering broken: priority not sorted ASC")
		}
	}
}

func TestJarRepository_Update(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewJarRepository(testDB)

	// create initial jar
	id, err := repo.Create(ctx, Jar{
		Name:           "Old",
		AllocationType: AllocationFixed,
		Value:          500,
		Priority:       1,
	})
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	// update it
	err = repo.Update(ctx, Jar{
		ID:             id,
		Name:           "Updated",
		AllocationType: AllocationPercentage,
		Value:          60,
		Priority:       5,
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
			found = j.Name == "Updated" && j.Priority == 5
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
		AllocationType: AllocationFixed,
		Value:          200,
		Priority:       1,
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
