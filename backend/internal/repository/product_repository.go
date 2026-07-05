package repository

import (
	"context"
	"fmt"
	"strings"

	"tunisianet-scraper/internal/database"
	"tunisianet-scraper/internal/models"
)

func SearchByQuery(search string) ([]models.Product, error) {

	search = strings.ToLower(strings.TrimSpace(search))

	words := strings.Fields(search)

	baseQuery := `
		SELECT name, price, image, url, store
		FROM products
		WHERE `

	var (
		args       []interface{}
		conditions []string
	)

	for _, word := range words {

		if len(word) < 2 {
			continue
		}

		args = append(args, "%"+word+"%")

		conditions = append(
			conditions,
			fmt.Sprintf("LOWER(name) LIKE $%d", len(args)),
		)
	}

	if len(conditions) == 0 {
		return []models.Product{}, nil
	}

	query := baseQuery + strings.Join(conditions, " AND ")

	rows, err := database.DB.Query(
		context.Background(),
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []models.Product

	for rows.Next() {

		var product models.Product

		err := rows.Scan(
			&product.Name,
			&product.Price,
			&product.Image,
			&product.URL,
			&product.Store,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}
func SearchByStoreAndQuery(
	store string,
	search string,
) ([]models.Product, error) {

	rows, err := database.DB.Query(
		context.Background(),
		`
		SELECT
			p.name,
			p.price,
			p.image,
			p.url,
			p.store
		FROM products p
		INNER JOIN search_cache sc
			ON p.url = sc.product_url
		WHERE
			LOWER(sc.query) = LOWER($1)
			AND LOWER(sc.store) = LOWER($2)
		`,
		search,
		store,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []models.Product

	for rows.Next() {

		var product models.Product

		err := rows.Scan(
			&product.Name,
			&product.Price,
			&product.Image,
			&product.URL,
			&product.Store,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}
func SaveProducts(query string, products []models.Product) error {

	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	for _, product := range products {

		// Save product
		_, err = tx.Exec(
			context.Background(),
			`
			INSERT INTO products
			(name, price, image, url, store, last_updated)
			VALUES ($1,$2,$3,$4,$5,NOW())

			ON CONFLICT(url)

			DO UPDATE SET
				name = EXCLUDED.name,
				price = EXCLUDED.price,
				image = EXCLUDED.image,
				store = EXCLUDED.store,
				last_updated = NOW()
			`,
			product.Name,
			product.Price,
			product.Image,
			product.URL,
			product.Store,
		)

		if err != nil {
			return err
		}

		// Save search mapping
		_, err = tx.Exec(
			context.Background(),
			`
			INSERT INTO search_cache
			(query, store, product_url)

			VALUES ($1,$2,$3)

			ON CONFLICT(query, store, product_url)

			DO NOTHING
			`,
			query,
			product.Store,
			product.URL,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}
