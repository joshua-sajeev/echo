package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/joshu-sajeev/echo/internal/db"
)

func main() {
	_ = godotenv.Load()
	fmt.Println("Hello World!")
	db.Connect(os.Getenv("DATABASE_URL"))
}
