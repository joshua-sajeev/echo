// Package testutils used for inital test setup with Docker
package testutils

import (
	"context"
	"fmt"
	"log"
	"sync"
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
