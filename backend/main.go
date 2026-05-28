package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/joshu-sajeev/echo/internal/app"
	"github.com/lmittmann/tint"
)

func main() {
	_ = godotenv.Load()

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen, // Simplifies timestamp to "7:15PM"
	}))
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	application, err := app.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("Failed to initialize application", "err", err)
		os.Exit(1)
	}
	defer application.Close()

	slog.Info("Database connected successfully")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server running", "port", port)
	if err := http.ListenAndServe(":"+port, application.Router); err != nil {
		slog.Error("Server forced to shutdown", "err", err)
		os.Exit(1)
	}
}
