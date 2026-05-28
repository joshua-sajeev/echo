package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/joshu-sajeev/echo/internal/models"
	"github.com/joshu-sajeev/echo/internal/repository"
)

var (
	now = time.Now()

	stubAccounts = []models.Account{
		{ID: 1, Name: "Checking", IsArchived: false, CreatedAt: now},
		{ID: 2, Name: "Savings", IsArchived: false, CreatedAt: now},
	}

	stubAccountsWithBalances = []models.AccountWithBalance{
		{
			Account: models.Account{
				ID:   1,
				Name: "Checking",
			},
			Balance: 500,
		},
		{
			Account: models.Account{
				ID:   2,
				Name: "Savings",
			},
			Balance: 1200,
		},
	}
)

// --- Create ---
func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		repo    *repository.MockAccountRepo
		wantID  int64
		wantErr error
	}{
		{
			name:  "success",
			input: "  Checking  ",
			repo: &repository.MockAccountRepo{
				CreateFn: func(_ context.Context, name string) (int64, error) {
					return 42, nil
				},
			},
			wantID: 42,
		},
		{
			name:    "empty name",
			input:   "   ",
			repo:    &repository.MockAccountRepo{},
			wantErr: ErrInvalidAccountName,
		},
		{
			name:  "repo error",
			input: "Checking",
			repo: &repository.MockAccountRepo{
				CreateFn: func(_ context.Context, name string) (int64, error) {
					return 0, errors.New("db error")
				},
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAccountService(tt.repo)

			id, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				if err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if id != tt.wantID {
				t.Fatalf("expected id %d, got %d", tt.wantID, id)
			}
		})
	}
}

// --- List ---
func TestAccountService_List(t *testing.T) {
	repoErr := errors.New("db error")

	tests := []struct {
		name    string
		repo    *repository.MockAccountRepo
		want    []models.Account
		wantErr error
	}{
		{
			name: "success",
			repo: &repository.MockAccountRepo{
				ListFn: func(_ context.Context) ([]models.Account, error) {
					return stubAccounts, nil
				},
			},
			want: stubAccounts,
		},
		{
			name: "repo error",
			repo: &repository.MockAccountRepo{
				ListFn: func(_ context.Context) ([]models.Account, error) {
					return nil, repoErr
				},
			},
			wantErr: repoErr,
		},
	}
	for _, tt := range tests {
		svc := NewAccountService(tt.repo)

		got, err := svc.List(context.Background())
		if !errors.Is(err, tt.wantErr) {
			t.Fatalf("expected error %v, got %v", tt.wantErr, err)
		}

		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestList_Success(t *testing.T) {
	svc := NewAccountService(&repository.MockAccountRepo{
		ListFn: func(_ context.Context) ([]models.Account, error) {
			return stubAccounts, nil
		},
	})

	accounts, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(accounts) != 2 {
		t.Errorf("expected 2 accounts, got %d", len(accounts))
	}
}

// --- ListWithBalances ---

func TestAccountService_ListWithBalances(t *testing.T) {
	repoErr := errors.New("db error")

	tests := []struct {
		name    string
		repo    *repository.MockAccountRepo
		want    []models.AccountWithBalance
		wantErr error
	}{
		{
			name: "success",
			repo: &repository.MockAccountRepo{
				ListWithBalancesFn: func(_ context.Context) ([]models.AccountWithBalance, error) {
					return stubAccountsWithBalances, nil
				},
			},
			want: stubAccountsWithBalances,
		},
		{
			name: "repo error",
			repo: &repository.MockAccountRepo{
				ListWithBalancesFn: func(_ context.Context) ([]models.AccountWithBalance, error) {
					return nil, repoErr
				},
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAccountService(tt.repo)

			got, err := svc.ListWithBalances(context.Background())

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAccountService_ListArchivedWithBalances(t *testing.T) {
	repoErr := errors.New("db error")

	archived := []models.AccountWithBalance{
		{
			Account: models.Account{
				ID:         3,
				Name:       "Old Account",
				IsArchived: true,
			},
			Balance: 0,
		},
	}

	tests := []struct {
		name    string
		repo    *repository.MockAccountRepo
		want    []models.AccountWithBalance
		wantErr error
	}{
		{
			name: "success",
			repo: &repository.MockAccountRepo{
				ListArchivedWithBalancesFn: func(_ context.Context) ([]models.AccountWithBalance, error) {
					return archived, nil
				},
			},
			want: archived,
		},
		{
			name: "repo error",
			repo: &repository.MockAccountRepo{
				ListArchivedWithBalancesFn: func(_ context.Context) ([]models.AccountWithBalance, error) {
					return nil, repoErr
				},
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAccountService(tt.repo)

			got, err := svc.ListArchivedWithBalances(context.Background())

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAccountService_Rename(t *testing.T) {
	repoErr := errors.New("db error")

	tests := []struct {
		name    string
		id      int64
		input   string
		repo    *repository.MockAccountRepo
		wantErr error
	}{
		{
			name:  "success",
			id:    1,
			input: "  New Name  ",
			repo: &repository.MockAccountRepo{
				RenameFn: func(_ context.Context, id int64, name string) error {
					if id != 1 {
						t.Fatalf("expected id 1, got %d", id)
					}

					if name != "New Name" {
						t.Fatalf("expected trimmed name 'New Name', got %q", name)
					}

					return nil
				},
			},
		},
		{
			name:    "invalid id",
			id:      0,
			input:   "name",
			repo:    &repository.MockAccountRepo{},
			wantErr: ErrInvalidAccountID,
		},
		{
			name:    "empty name",
			id:      1,
			input:   "   ",
			repo:    &repository.MockAccountRepo{},
			wantErr: ErrInvalidAccountName,
		},
		{
			name:  "repo error",
			id:    1,
			input: "Checking",
			repo: &repository.MockAccountRepo{
				RenameFn: func(_ context.Context, id int64, name string) error {
					return repoErr
				},
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAccountService(tt.repo)

			err := svc.Rename(context.Background(), tt.id, tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestAccountService_Archive(t *testing.T) {
	repoErr := errors.New("archive error")

	tests := []struct {
		name    string
		id      int64
		repo    *repository.MockAccountRepo
		wantErr error
	}{
		{
			name: "success",
			id:   1,
			repo: &repository.MockAccountRepo{
				ArchiveFn: func(_ context.Context, id int64) error {
					return nil
				},
			},
		},
		{
			name:    "invalid id",
			id:      -1,
			repo:    &repository.MockAccountRepo{},
			wantErr: ErrInvalidAccountID,
		},
		{
			name: "repo error",
			id:   99,
			repo: &repository.MockAccountRepo{
				ArchiveFn: func(_ context.Context, id int64) error {
					return repoErr
				},
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAccountService(tt.repo)

			err := svc.Archive(context.Background(), tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestAccountService_Unarchive(t *testing.T) {
	repoErr := errors.New("unarchive error")

	tests := []struct {
		name    string
		id      int64
		repo    *repository.MockAccountRepo
		wantErr error
	}{
		{
			name: "success",
			id:   1,
			repo: &repository.MockAccountRepo{
				UnarchiveFn: func(_ context.Context, id int64) error {
					return nil
				},
			},
		},
		{
			name:    "invalid id",
			id:      0,
			repo:    &repository.MockAccountRepo{},
			wantErr: ErrInvalidAccountID,
		},
		{
			name: "repo error",
			id:   99,
			repo: &repository.MockAccountRepo{
				UnarchiveFn: func(_ context.Context, id int64) error {
					return repoErr
				},
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAccountService(tt.repo)

			err := svc.Unarchive(context.Background(), tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
