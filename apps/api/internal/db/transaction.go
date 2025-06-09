package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type TransactionManager struct {
	*DB
}

func NewTransactionManager(db *DB) *TransactionManager {
	return &TransactionManager{
		DB: db,
	}
}

func (txm *TransactionManager) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := txm.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Err(err).Msg("Failed to begin new transaction")
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err = fn(tx); err != nil {
		log.Err(err).Msg("Failed to execute transaction functions")
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to commit and close transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
