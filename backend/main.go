package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/joshu-sajeev/echo/internal/app"
)

func main() {
	_ = godotenv.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	application, err := app.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer application.Close()

	log.Println("Database connected successfully")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on :" + port)
	if err := http.ListenAndServe(":"+port, application.Router); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
