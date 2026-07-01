package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(databaseURL string) {
	var err error

	DB, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Connected to PostgreSQL")
}
func CreateTables() error {
	_, err := DB.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
        query TEXT NOT NULL,
        name TEXT NOT NULL,
        price TEXT,
        image TEXT,
        url TEXT UNIQUE,
        store TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)

	return err
}
