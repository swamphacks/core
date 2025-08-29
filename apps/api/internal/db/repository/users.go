package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db *db.DB
}

func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Call this to create a copy with transactional queries
func (r *UserRepository) NewTx(tx pgx.Tx) *UserRepository {
	txDB := &db.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &UserRepository{db: txDB}
}

func (r *UserRepository) Create(ctx context.Context, params sqlc.CreateUserParams) (*sqlc.AuthUser, error) {
	user, err := r.db.Query.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*sqlc.AuthUser, error) {
	user, err := r.db.Query.GetUserByID(ctx, id)
	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*sqlc.AuthUser, error) {
	user, err := r.db.Query.GetUserByEmail(ctx, &email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
