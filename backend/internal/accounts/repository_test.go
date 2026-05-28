package accounts

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v4"
	"github.com/pressly/goose"
)

const (
	testDBUser     = "postgres"
	testDBPassword = "secret"
	testDBName     = "testdb"
	testDBPort     = "5432"
)

var globalDBPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Initialize dockertest pool
	pool, err := dockertest.NewPool(ctx, "")
	if err != nil {
		log.Fatalf("could not connect to docker: %v", err)
	}

	// 2. Run the Postgres container ONCE
	postgres, err := pool.Run(
		ctx,
		"postgres",
		dockertest.WithTag("14"),
		dockertest.WithEnv([]string{
			"POSTGRES_PASSWORD=" + testDBPassword,
			"POSTGRES_DB=" + testDBName,
		}),
	)
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	hostPort := postgres.GetHostPort(testDBPort + "/tcp")
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		testDBUser,
		testDBPassword,
		hostPort,
		testDBName,
	)

	// 3. Wait for Postgres to be ready
	err = pool.Retry(ctx, 30*time.Second, func() error {
		var err error
		globalDBPool, err = pgxpool.New(ctx, databaseURL)
		if err != nil {
			return err
		}
		return globalDBPool.Ping(ctx)
	})
	if err != nil {
		log.Fatalf("could not connect to postgres: %v", err)
	}

	// 4. Run migrations ONCE
	sqlDB := stdlib.OpenDBFromPool(globalDBPool)
	if err := goose.Up(sqlDB, "../../migrations"); err != nil {
		log.Fatalf("failed running migrations: %v", err)
	}

	// 5. Run all tests in the package
	code := m.Run()

	// 6. Global Teardown after all tests finish
	sqlDB.Close()
	globalDBPool.Close()
	_ = postgres.Close(ctx) // stop & remove container

	os.Exit(code)
}

func TestAccountRepo_Create(t *testing.T) {
	ctx := context.Background()
	repo := NewAccountRepository(globalDBPool)

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
				t.Fatalf("Create(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && id == 0 {
				t.Fatal("expected non-zero id")
			}
		})
	}
}

func TestAccountRepo_Rename(t *testing.T) {
	ctx := context.Background()
	repo := NewAccountRepository(globalDBPool)

	activeID, _ := repo.Create(ctx, "Cash")
	archivedID, _ := repo.Create(ctx, "Old Savings")
	if err := repo.Archive(ctx, archivedID); err != nil {
		t.Fatal("unexpected err", err)
	}

	tests := []struct {
		name    string
		id      int64
		input   string
		wantErr bool
	}{
		{name: "valid rename", id: activeID, input: "Savings", wantErr: false},
		{name: "empty name", id: activeID, input: "", wantErr: true},
		{name: "invalid id", id: 0, input: "X", wantErr: true},
		{name: "non-existent id", id: 99999, input: "X", wantErr: true},
		{name: "archived account", id: archivedID, input: "New Name", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Rename(ctx, tt.id, tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Rename(%d, %q) error = %v, wantErr %v", tt.id, tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestAccountRepo_Archive(t *testing.T) {
	ctx := context.Background()
	repo := NewAccountRepository(globalDBPool)

	activeID, _ := repo.Create(ctx, "Cash")
	alreadyArchivedID, _ := repo.Create(ctx, "Old")
	if err := repo.Archive(ctx, alreadyArchivedID); err != nil {
		t.Fatal("unexpected err", err)
	}

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{name: "archives active account", id: activeID, wantErr: false},
		{name: "already archived", id: alreadyArchivedID, wantErr: true},
		{name: "invalid id", id: 0, wantErr: true},
		{name: "non-existent id", id: 99999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Archive(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Archive(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

func TestAccountRepo_Unarchive(t *testing.T) {
	ctx := context.Background()
	repo := NewAccountRepository(globalDBPool)

	archivedID, _ := repo.Create(ctx, "Old Savings")
	if err := repo.Archive(ctx, archivedID); err != nil {
		t.Fatal("unexpected err", err)
	}
	activeID, _ := repo.Create(ctx, "Cash")

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{name: "unarchives archived account", id: archivedID, wantErr: false},
		{name: "active account", id: activeID, wantErr: true},
		{name: "invalid id", id: 0, wantErr: true},
		{name: "non-existent id", id: 99999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Unarchive(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unarchive(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}
