package transactions

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

func TestTransactionRepositoryCreate(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewTransactionRepository(testDB)

	tx := Transaction{
		Type:   "income",
		Name:   "Salary",
		Amount: 50000,
	}

	id, err := repo.Create(ctx, tx)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if id == 0 {
		t.Fatal("expected valid id")
	}

	// verify using list
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	found := false
	for _, v := range list {
		if v.ID == id && v.Name == "Salary" {
			found = true
		}
	}

	if !found {
		t.Fatal("created transaction not found")
	}
}

func TestTransactionRepositoryList(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewTransactionRepository(testDB)

	// seed data
	inputs := []Transaction{
		{Type: "income", Name: "A", Amount: 100},
		{Type: "expense", Name: "B", Amount: 200},
	}

	for _, tx := range inputs {
		_, err := repo.Create(ctx, tx)
		if err != nil {
			t.Fatalf("setup create failed: %v", err)
		}
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(list) < len(inputs) {
		t.Fatalf("expected at least %d got %d", len(inputs), len(list))
	}

	names := map[string]bool{}
	for _, l := range list {
		names[l.Name] = true
	}

	for _, in := range inputs {
		if !names[in.Name] {
			t.Fatalf("missing transaction %s", in.Name)
		}
	}
}

func TestTransactionRepositoryUpdate(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewTransactionRepository(testDB)

	id, err := repo.Create(ctx, Transaction{
		Type:   "income",
		Name:   "Old",
		Amount: 1000,
	})
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	err = repo.Update(ctx, Transaction{
		ID:     id,
		Type:   "income",
		Name:   "Updated",
		Amount: 2000,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	found := false
	for _, tx := range list {
		if tx.ID == id && tx.Name == "Updated" && tx.Amount == 2000 {
			found = true
		}
	}

	if !found {
		t.Fatal("update not reflected in DB")
	}
}

func TestTransactionRepositoryDelete(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewTransactionRepository(testDB)

	id, err := repo.Create(ctx, Transaction{
		Type:   "expense",
		Name:   "ToDelete",
		Amount: 300,
	})
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	err = repo.Delete(ctx, id)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	for _, tx := range list {
		if tx.ID == id {
			t.Fatal("transaction was not deleted")
		}
	}
}
