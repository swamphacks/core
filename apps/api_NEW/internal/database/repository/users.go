package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Call this to create a copy with transactional queries
func (r *UserRepository) NewTx(tx pgx.Tx) *UserRepository {
	txDB := &database.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &UserRepository{db: txDB}
}

func (r *UserRepository) CreateUser(ctx context.Context, params sqlc.CreateUserParams) (*sqlc.User, error) {
	user, err := r.db.Query.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, params sqlc.UpdateUserParams) error {
	err := r.db.Query.UpdateUser(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrUserNotFound
		}
	}
	return err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*sqlc.User, error) {
	user, err := r.db.Query.GetUserByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*sqlc.User, error) {
	user, err := r.db.Query.GetUserByEmail(ctx, &email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserEmailInfoById(ctx context.Context, id uuid.UUID) (*sqlc.GetUserEmailInfoByIdRow, error) {
	row, err := r.db.Query.GetUserEmailInfoById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *UserRepository) GetUserByRFID(ctx context.Context, rfid string) (*sqlc.User, error) {
	user, err := r.db.Query.GetUserByRFID(ctx, &rfid)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context, search *string, limit, offset int32) ([]sqlc.User, error) {
	params := sqlc.GetUsersParams{
		Search: search,
		Limit:  limit,
		Offset: offset,
	}

	return r.db.Query.GetUsers(ctx, params)
}

func (r *UserRepository) UpdateRole(ctx context.Context, updateRoleParams sqlc.UpdateRoleParams) error {
	return r.db.Query.UpdateRole(ctx, updateRoleParams)
}

func (r *UserRepository) RemoveRole(ctx context.Context, userId uuid.UUID) error {
	return r.db.Query.RemoveRole(ctx, userId)
}

func (r *UserRepository) UpdateCheckInTime(ctx context.Context, updateCheckInParams sqlc.UpdateCheckInTimeParams) error {
	return r.db.Query.UpdateCheckInTime(ctx, updateCheckInParams)
}

func (r *UserRepository) UpdateRFID(ctx context.Context, updateRFIDParams sqlc.UpdateRFIDParams) error {
	return r.db.Query.UpdateRFID(ctx, updateRFIDParams)
}
