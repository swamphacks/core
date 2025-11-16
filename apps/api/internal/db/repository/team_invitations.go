package repository

import (
	"context"
	"errors"
	"time"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/db"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrInvitationNotFound = errors.New("invitation not found")
)

type TeamInvitationRepository struct {
	db *db.DB
}

func NewTeamInvitationRespository(db *db.DB) *TeamInvitationRepository {
	return &TeamInvitationRepository{
		db: db,
	}
}

func (r *TeamInvitationRepository) NewTx(tx pgx.Tx) *TeamInvitationRepository {
	txDB := &db.DB{
		Pool:  r.db.Pool,
		Query: sqlc.New(tx),
	}

	return &TeamInvitationRepository{
		db: txDB,
	}
}

func (r *TeamInvitationRepository) Create(ctx context.Context, teamId, invitedByUserId uuid.UUID, invitedEmail string, expiresAt *time.Time) (*sqlc.TeamInvitation, error) {
	invitation, err := r.db.Query.CreateInvitation(ctx, sqlc.CreateInvitationParams{
		TeamID:          teamId,
		InvitedByUserID: invitedByUserId,
		InvitedEmail:    invitedEmail,
		ExpiresAt:       expiresAt,
	})
	if err != nil {
		return nil, err
	}

	return &invitation, nil
}

func (r *TeamInvitationRepository) GetByID(ctx context.Context, id uuid.UUID) (*sqlc.TeamInvitation, error) {
	invitation, err := r.db.Query.GetInvitationByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}

	return &invitation, nil
}

func (r *TeamInvitationRepository) AcceptInvitation(ctx context.Context, id, invitedUserId uuid.UUID) (*sqlc.TeamInvitation, error) {
	invitation, err := r.db.Query.AcceptInvitation(ctx, sqlc.AcceptInvitationParams{
		ID:            id,
		InvitedUserID: invitedUserId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}

	return &invitation, nil
}

func (r *TeamInvitationRepository) RejectInvitation(ctx context.Context, id uuid.UUID) (*sqlc.TeamInvitation, error) {
	invitation, err := r.db.Query.RejectInvitation(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}

	return &invitation, nil
}

func (r *TeamInvitationRepository) GetByEmailAndTeam(ctx context.Context, email string, teamId uuid.UUID) (*sqlc.TeamInvitation, error) {
	invitation, err := r.db.Query.GetPendingInvitationByEmailAndTeam(ctx, sqlc.GetPendingInvitationByEmailAndTeamParams{
		InvitedEmail: email,
		TeamID:       teamId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}

	return &invitation, nil
}

func (r *TeamInvitationRepository) Update(ctx context.Context, id uuid.UUID, invitedUserId *uuid.UUID, status *sqlc.InvitationStatus) (*sqlc.TeamInvitation, error) {
	params := sqlc.UpdateInvitationParams{
		ID: id,
	}

	if invitedUserId != nil {
		params.InvitedUserIDDoUpdate = true
		params.InvitedUserID = *invitedUserId
	}

	if status != nil {
		params.StatusDoUpdate = true
		params.Status = *status
	}

	invitation, err := r.db.Query.UpdateInvitation(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}

	return &invitation, nil
}
