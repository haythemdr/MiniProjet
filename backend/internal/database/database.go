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

	// Products table
	_, err := DB.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price TEXT,
		image TEXT,
		url TEXT UNIQUE,
		store TEXT,
		last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return err
	}

	// Search cache table
	_, err = DB.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS search_cache (
		id SERIAL PRIMARY KEY,
		query TEXT NOT NULL,
		store TEXT NOT NULL,
		product_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		CONSTRAINT fk_product
			FOREIGN KEY (product_url)
			REFERENCES products(url)
			ON DELETE CASCADE,

		CONSTRAINT unique_search
			UNIQUE(query, store, product_url)
	);
	`)
	if err != nil {
		return err
	}
	// Search statistics table
	_, err = DB.Exec(context.Background(), `
CREATE TABLE IF NOT EXISTS search_queries (
	id SERIAL PRIMARY KEY,
	query TEXT UNIQUE NOT NULL,
	search_count INT DEFAULT 0,
	last_search TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	last_refresh TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`)
	if err != nil {
		return err
	}

	return nil
}
