package database

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type DB struct {
	Pool  *pgxpool.Pool
	Query *sqlc.Queries
}

func NewDB(connStr string) *DB {
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal(err)
	}

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return registerCustomTypes(ctx, conn)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
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

func (d *DB) NewTX(tx pgx.Tx) *DB {
	txDB := DB{
		Pool:  d.Pool,
		Query: sqlc.New(tx),
	}

	return &txDB
}

func (d *DB) Close() {
	d.Pool.Close()
}

func registerCustomTypes(ctx context.Context, conn *pgx.Conn) error {
	typeNames := []string{
		"email_campaign_format",
		"email_campaign_status",
		"email_recipient_type",
		"_email_recipient_type",
	}

	for _, typeName := range typeNames {
		dataType, err := conn.LoadType(ctx, typeName)
		if isUndefinedType(err) {
			continue
		}
		if err != nil {
			return err
		}

		conn.TypeMap().RegisterType(dataType)
	}

	return nil
}

func isUndefinedType(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "42704"
}
