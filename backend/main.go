package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/joshu-sajeev/echo/internal/db"
	"github.com/joshu-sajeev/echo/internal/router"
	"github.com/pressly/goose"
)

func main() {
	_ = godotenv.Load()

	pool, err := db.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("database connected")

	sqlDB := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatalf("failed running migrations: %v", err)
	}

	r := router.New(pool)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server running on :" + port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
