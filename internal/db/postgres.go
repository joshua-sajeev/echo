// Package db
package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(url string) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	Pool = pool
	log.Println("database connected")
}
