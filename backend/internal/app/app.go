package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joshu-sajeev/echo/internal/accounts"
	"github.com/joshu-sajeev/echo/internal/jars"
	"github.com/joshu-sajeev/echo/internal/router"
	"github.com/joshu-sajeev/echo/internal/transactions"
	"github.com/pressly/goose"
)

type App struct {
	Router http.Handler
	DB     *pgxpool.Pool
}

// New initializes all infrastructure, domains, and routes.
func New(ctx context.Context, dbConnString string) (*App, error) {
	// 1. Initialize Database
	pool, err := pgxpool.New(ctx, dbConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	sqlDB := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatalf("failed running migrations: %v", err)
	}
	// 2. Wire up the domains (The "heavy lifting" lives here now)
	accountRepo := accounts.NewAccountRepository(pool)
	accountService := accounts.NewAccountService(accountRepo)
	accountHandler := accounts.NewAccountHandler(accountService)

	jarRepo := jars.NewJarRepository(pool)
	jarService := jars.NewJarService(jarRepo)
	jarHandler := jars.NewJarHandler(jarService)

	txRepo := transactions.NewTransactionRepository(pool)
	txService := transactions.NewTransactionService(txRepo)
	txHandler := transactions.NewTransactionHandler(txService)

	// 3. Build the router config
	appRouter := router.New(router.Config{
		AccountHandler:     accountHandler,
		JarHandler:         jarHandler,
		TransactionHandler: txHandler,
	})

	return &App{
		Router: appRouter,
		DB:     pool,
	}, nil
}

// Close ensures resources like the database pool are cleaned up gracefully
func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
