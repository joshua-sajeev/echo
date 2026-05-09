// Package db
package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func Connect(url string) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	Conn = conn
	log.Println("database connected")
}
