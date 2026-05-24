package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/joshu-sajeev/echo/internal/db"
)

func main() {
	_ = godotenv.Load()

	pool, err := db.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("database connected")
	fmt.Print(pool)
}
