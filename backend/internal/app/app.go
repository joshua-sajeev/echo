// Package app handles the app wiring
package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joshu-sajeev/echo/internal/accounts"
	"github.com/joshu-sajeev/echo/internal/auth"
	"github.com/joshu-sajeev/echo/internal/jars"
	"github.com/joshu-sajeev/echo/internal/router"
	"github.com/joshu-sajeev/echo/internal/transactions"
	"github.com/pressly/goose"
)

type App struct {
	Router http.Handler
	DB     *pgxpool.Pool
}

func New(ctx context.Context, dbConnString string) (*App, error) {
	pool, err := pgxpool.New(ctx, dbConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	sqlDB := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		slog.Error("failed running migrations", "err", err)
		os.Exit(1)
	}

	accountRepo := accounts.NewAccountRepository(pool)
	accountService := accounts.NewAccountService(accountRepo)
	accountHandler := accounts.NewAccountHandler(accountService)

	txRepo := transactions.NewTransactionRepository(pool)
	txService := transactions.NewTransactionService(txRepo)
	txHandler := transactions.NewTransactionHandler(txService)

	jarRepo := jars.NewJarRepository(pool)
	jarService := jars.NewJarService(jarRepo, txRepo)
	jarHandler := jars.NewJarHandler(jarService)

	store := auth.NewStore()

	authHandler := auth.NewHandler(store)
	appRouter := router.New(router.Config{
		AccountHandler:     accountHandler,
		JarHandler:         jarHandler,
		TransactionHandler: txHandler,
		AuthHandler:        authHandler,
	})

	return &App{
		Router: appRouter,
		DB:     pool,
	}, nil
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
