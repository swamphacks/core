package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type HackathonRepository struct {
	db *database.DB
}

func NewHackathonRepository(db *database.DB) *HackathonRepository {
	return &HackathonRepository{
		db: db,
	}
}

func (r *HackathonRepository) NewTx(tx pgx.Tx) *HackathonRepository {
	txDB := &database.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &HackathonRepository{
		db: txDB,
	}
}

func (r *HackathonRepository) CreateHackathon(ctx context.Context, params sqlc.CreateHackathonParams) (*sqlc.Hackathon, error) {
	hackathon, err := r.db.Query.CreateHackathon(ctx, params)
	return &hackathon, err
}

func (r *HackathonRepository) GetHackathon(ctx context.Context) (*sqlc.Hackathon, error) {
	hackathon, err := r.db.Query.GetHackathon(ctx)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEntityNotFound
		}
		return nil, err
	}

	return &hackathon, err
}

func (r *HackathonRepository) UpdateHackathon(ctx context.Context, params sqlc.UpdateHackathonParams) error {
	return r.db.Query.UpdateHackathon(ctx, params)
}

func (r *HackathonRepository) GetStaff(ctx context.Context) (*[]sqlc.GetStaffRow, error) {
	users, err := r.db.Query.GetStaff(ctx)
	return &users, err
}

func (r *HackathonRepository) GetAttendeesWithDiscord(ctx context.Context) (*[]sqlc.GetAttendeesWithDiscordRow, error) {
	attendees, err := r.db.Query.GetAttendeesWithDiscord(ctx)
	return &attendees, err
}

func (r *HackathonRepository) GetAttendeeUserIds(ctx context.Context) ([]uuid.UUID, error) {
	return r.db.Query.GetAttendeeUserIds(ctx)
}

func (r *HackathonRepository) GetAttendeeCount(ctx context.Context) (int64, error) {
	return r.db.Query.GetAttendeeCount(ctx)
}
