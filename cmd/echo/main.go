package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/joshu-sajeev/echo/internal/db"
	"github.com/joshu-sajeev/echo/internal/router"
	"github.com/joshu-sajeev/echo/ui"
)

func main() {
	_ = godotenv.Load()
	pool := db.Connect(os.Getenv("DATABASE_URL"))
	tmpl := ui.LoadTemplates()

	r := router.New(tmpl, pool)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
