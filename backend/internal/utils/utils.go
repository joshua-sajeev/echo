// Package utils used for inital test setup with Docker
package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

var (
	dbPool   *pgxpool.Pool
	postgres dockertest.ClosableResource
	once     sync.Once
)

func GetTestDB() *pgxpool.Pool {
	once.Do(func() {
		ctx := context.Background()

		pool, err := dockertest.NewPool(ctx, "")
		if err != nil {
			log.Fatalf("could not connect to docker: %v", err)
		}

		postgres, err = pool.Run(
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

		err = pool.Retry(ctx, 30*time.Second, func() error {
			var err error

			dbPool, err = pgxpool.New(ctx, databaseURL)
			if err != nil {
				return err
			}

			return dbPool.Ping(ctx)
		})
		if err != nil {
			log.Fatalf("could not connect to postgres: %v", err)
		}

		sqlDB := stdlib.OpenDBFromPool(dbPool)

		if err := goose.Up(sqlDB, "../../migrations"); err != nil {
			log.Fatalf("failed running migrations: %v", err)
		}

		if err := sqlDB.Close(); err != nil {
			log.Fatalf("failed closing sqlDB: %v", err)
		}
	})

	return dbPool
}

func CleanupTestDB() {
	ctx := context.Background()

	if dbPool != nil {
		dbPool.Close()
	}

	if postgres != nil {
		_ = postgres.Close(ctx)
	}
}

func ResetTables() {
	ctx := context.Background()

	_, err := dbPool.Exec(
		ctx,
		`TRUNCATE TABLE accounts,jars,transactions RESTART IDENTITY CASCADE`,
	)
	if err != nil {
		log.Fatalf("failed truncating tables: %v", err)
	}
}

func LogError(ctx context.Context, op string, err error) {
	if err == nil {
		return
	}

	// Row not found → not an error in most APIs
	if errors.Is(err, pgx.ErrNoRows) {
		slog.InfoContext(ctx,
			"db record not found",
			"op", op,
		)
		return
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {

		// Clean grouped log
		slog.ErrorContext(ctx,
			"database operation failed",
			"op", op,
			"code", pgErr.Code,
			"type", classifyPgError(pgErr.Code),
			"message", pgErr.Message,
			"table", pgErr.TableName,
			"column", pgErr.ColumnName,
			"constraint", pgErr.ConstraintName,
		)

		return
	}

	// fallback
	slog.ErrorContext(ctx,
		"unexpected database error",
		"op", op,
		"error", err.Error(),
	)
}

func classifyPgError(code string) string {
	switch code {
	case "23505":
		return "unique_violation"
	case "23514":
		return "check_violation"
	case "23503":
		return "foreign_key_violation"
	case "08006":
		return "connection_failure"
	default:
		return "postgres_error"
	}
}
