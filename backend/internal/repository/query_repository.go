package repository

import (
	"context"

	"tunisianet-scraper/internal/database"
)

func RecordSearch(query string) error {

	_, err := database.DB.Exec(
		context.Background(),
		`
		INSERT INTO search_queries
		(query, search_count, last_search)

		VALUES ($1,1,NOW())

		ON CONFLICT(query)

		DO UPDATE SET

			search_count = search_queries.search_count + 1,
			last_search = NOW()
		`,
		query,
	)

	return err
}

func GetQueriesToRefresh(limit int) ([]string, error) {

	rows, err := database.DB.Query(
		context.Background(),
		`
		SELECT query

		FROM search_queries

		WHERE last_refresh < NOW() - INTERVAL '1 hour'

		ORDER BY search_count DESC

		LIMIT $1
		`,
		limit,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var queries []string

	for rows.Next() {

		var query string

		if err := rows.Scan(&query); err != nil {
			return nil, err
		}

		queries = append(queries, query)
	}

	return queries, nil
}

func UpdateRefreshTime(query string) error {

	_, err := database.DB.Exec(
		context.Background(),
		`
		UPDATE search_queries

		SET last_refresh = NOW()

		WHERE query = $1
		`,
		query,
	)

	return err
}
