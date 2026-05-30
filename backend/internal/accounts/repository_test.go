package accounts

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

func TestAccountRepo_Create(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewAccountRepository(testDB)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid name", input: "Cash", wantErr: false},
		{name: "another valid name", input: "Savings", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.Create(ctx, tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"Create(%q) error = %v, wantErr %v",
					tt.input,
					err,
					tt.wantErr,
				)
			}

			if !tt.wantErr {
				if id == 0 {
					t.Fatal("expected non-zero id")
				}

				accounts, err := repo.List(ctx)
				if err != nil {
					t.Fatalf("List() error = %v", err)
				}

				found := false

				for _, acc := range accounts {
					if acc.ID == id && acc.Name == tt.input {
						found = true
					}
				}

				if !found {
					t.Fatal("created account not found")
				}
			}
		})
	}
}

func TestAccountRepo_List(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewAccountRepository(testDB)

	cashID, _ := repo.Create(ctx, "Cash")
	savingsID, _ := repo.Create(ctx, "Savings")

	archivedID, _ := repo.Create(ctx, "Archived")
	_ = repo.Archive(ctx, archivedID)

	accounts, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts))
	}

	if accounts[0].ID != savingsID {
		t.Fatal("expected newest account first")
	}

	if accounts[1].ID != cashID {
		t.Fatal("expected oldest account second")
	}

	for _, acc := range accounts {
		if acc.IsArchived {
			t.Fatal("expected only active accounts")
		}
	}
}

func TestAccountRepo_ListWithBalances(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewAccountRepository(testDB)

	cashID, _ := repo.Create(ctx, "Cash")
	savingsID, _ := repo.Create(ctx, "Savings")

	accounts, err := repo.ListWithBalances(ctx)
	if err != nil {
		t.Fatalf("ListWithBalances() error = %v", err)
	}

	if len(accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts))
	}

	foundCash := false
	foundSavings := false

	for _, acc := range accounts {
		switch acc.ID {
		case cashID:
			foundCash = true

			if acc.Balance != 0 {
				t.Fatal("expected zero balance")
			}

		case savingsID:
			foundSavings = true

			if acc.Balance != 0 {
				t.Fatal("expected zero balance")
			}
		}
	}

	if !foundCash || !foundSavings {
		t.Fatal("missing expected accounts")
	}
}

func TestAccountRepo_ListArchivedWithBalances(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewAccountRepository(testDB)

	activeID, _ := repo.Create(ctx, "Cash")

	archivedID, _ := repo.Create(ctx, "Old Savings")
	_ = repo.Archive(ctx, archivedID)

	accounts, err := repo.ListArchivedWithBalances(ctx)
	if err != nil {
		t.Fatalf("ListArchivedWithBalances() error = %v", err)
	}

	if len(accounts) != 1 {
		t.Fatalf("expected 1 archived account, got %d", len(accounts))
	}

	acc := accounts[0]

	if acc.ID != archivedID {
		t.Fatalf(
			"expected archived account id %d, got %d",
			archivedID,
			acc.ID,
		)
	}

	if !acc.IsArchived {
		t.Fatal("expected archived account")
	}

	if acc.ID == activeID {
		t.Fatal("active account should not be returned")
	}
}

func TestAccountRepo_Rename(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewAccountRepository(testDB)

	activeID, _ := repo.Create(ctx, "Cash")

	archivedID, _ := repo.Create(ctx, "Old Savings")
	if err := repo.Archive(ctx, archivedID); err != nil {
		t.Fatal("unexpected error:", err)
	}

	tests := []struct {
		name        string
		id          int64
		input       string
		wantErr     bool
		verifyName  bool
		expectedNew string
	}{
		{
			name:        "valid rename",
			id:          activeID,
			input:       "Savings",
			wantErr:     false,
			verifyName:  true,
			expectedNew: "Savings",
		},
		{
			name:    "empty name",
			id:      activeID,
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid id",
			id:      0,
			input:   "X",
			wantErr: true,
		},
		{
			name:    "non-existent id",
			id:      99999,
			input:   "X",
			wantErr: true,
		},
		{
			name:    "archived account",
			id:      archivedID,
			input:   "New Name",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Rename(ctx, tt.id, tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"Rename(%d, %q) error = %v, wantErr %v",
					tt.id,
					tt.input,
					err,
					tt.wantErr,
				)
			}

			if tt.verifyName {
				accounts, err := repo.List(ctx)
				if err != nil {
					t.Fatalf("List() error = %v", err)
				}

				found := false

				for _, acc := range accounts {
					if acc.ID == tt.id {
						found = true

						if acc.Name != tt.expectedNew {
							t.Fatalf(
								"expected renamed account %q, got %q",
								tt.expectedNew,
								acc.Name,
							)
						}
					}
				}

				if !found {
					t.Fatal("renamed account not found")
				}
			}
		})
	}
}

func TestAccountRepo_Archive(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewAccountRepository(testDB)

	activeID, _ := repo.Create(ctx, "Cash")

	alreadyArchivedID, _ := repo.Create(ctx, "Old")
	if err := repo.Archive(ctx, alreadyArchivedID); err != nil {
		t.Fatal("unexpected error:", err)
	}

	tests := []struct {
		name       string
		id         int64
		wantErr    bool
		verifyGone bool
	}{
		{
			name:       "archives active account",
			id:         activeID,
			wantErr:    false,
			verifyGone: true,
		},
		{
			name:    "already archived",
			id:      alreadyArchivedID,
			wantErr: true,
		},
		{
			name:    "invalid id",
			id:      0,
			wantErr: true,
		},
		{
			name:    "non-existent id",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Archive(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"Archive(%d) error = %v, wantErr %v",
					tt.id,
					err,
					tt.wantErr,
				)
			}

			if tt.verifyGone {
				accounts, err := repo.List(ctx)
				if err != nil {
					t.Fatalf("List() error = %v", err)
				}

				for _, acc := range accounts {
					if acc.ID == tt.id {
						t.Fatal("archived account still appears in active list")
					}
				}
			}
		})
	}
}

func TestAccountRepo_Unarchive(t *testing.T) {
	utils.ResetTables()

	ctx := context.Background()
	repo := NewAccountRepository(testDB)

	archivedID, _ := repo.Create(ctx, "Old Savings")
	if err := repo.Archive(ctx, archivedID); err != nil {
		t.Fatal("unexpected error:", err)
	}

	activeID, _ := repo.Create(ctx, "Cash")

	tests := []struct {
		name          string
		id            int64
		wantErr       bool
		verifyPresent bool
	}{
		{
			name:          "unarchives archived account",
			id:            archivedID,
			wantErr:       false,
			verifyPresent: true,
		},
		{
			name:    "active account",
			id:      activeID,
			wantErr: true,
		},
		{
			name:    "invalid id",
			id:      0,
			wantErr: true,
		},
		{
			name:    "non-existent id",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Unarchive(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"Unarchive(%d) error = %v, wantErr %v",
					tt.id,
					err,
					tt.wantErr,
				)
			}

			if tt.verifyPresent {
				accounts, err := repo.List(ctx)
				if err != nil {
					t.Fatalf("List() error = %v", err)
				}

				found := false

				for _, acc := range accounts {
					if acc.ID == tt.id {
						found = true
					}
				}

				if !found {
					t.Fatal("unarchived account not found in active list")
				}
			}
		})
	}
}

func TestAccountRepo_Exists(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		accountID int64
		want      bool
	}{
		{
			name:      "account exists",
			accountID: 1,
			want:      true,
		},
		{
			name:      "account does not exist",
			accountID: 99999,
			want:      false,
		},
	}

	repo := NewAccountRepository(testDB)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Exists(ctx, tt.accountID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf(
					"expected %v, got %v",
					tt.want,
					got,
				)
			}
		})
	}
}
