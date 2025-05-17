package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

type DB struct {
	Pool  *pgxpool.Pool
	Query *sqlc.Queries
}

func NewDB(connStr string) *DB {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	query := sqlc.New(pool)

	db := DB{
		Pool:  pool,
		Query: query,
	}

	return &db
}

func (d *DB) Close() {
	d.Pool.Close()
}
