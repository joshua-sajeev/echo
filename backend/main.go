package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/joshu-sajeev/echo/internal/db"
	"github.com/joshu-sajeev/echo/internal/router"
)

func main() {
	_ = godotenv.Load()

	pool, err := db.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("database connected")
	r := router.New(pool)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
