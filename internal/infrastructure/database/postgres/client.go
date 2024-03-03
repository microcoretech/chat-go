package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgres(ctx context.Context, connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
