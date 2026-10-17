package redeemables

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

// Re-exported so the handler checks the service-level names, matching
// ErrEmailCampaignNotFound in the email package.
var (
	ErrNoActiveHackathon  = repository.ErrNoActiveHackathon
	ErrRedeemableNotFound = repository.ErrRedeemableNotFound
)

// Why a QR redemption was refused. RedeemError wraps one of these.
var (
	ErrHackerNotFound     = errors.New("hacker not found")
	ErrNoApplication      = errors.New("no application for this hackathon")
	ErrNotConfirmed       = errors.New("application not confirmed")
	ErrNotCheckedIn       = errors.New("hacker not checked in")
	ErrRedeemLimitReached = errors.New("redeem limit reached")
)

// RedeemError carries a message staff can act on alongside the reason.
type RedeemError struct {
	Reason  error
	Message string
}

func (e *RedeemError) Error() string { return e.Message }
func (e *RedeemError) Unwrap() error { return e.Reason }

type RedeemablesService struct {
	redeemablesRepo *repository.RedeemablesRepository
	logger          zerolog.Logger
}

func NewService(redeemablesRepo *repository.RedeemablesRepository, logger zerolog.Logger) *RedeemablesService {
	return &RedeemablesService{
		redeemablesRepo: redeemablesRepo,
		logger:          logger,
	}
}

func (s *RedeemablesService) GetRedeemables(ctx context.Context) (*[]sqlc.GetRedeemablesRow, error) {
	redeemables, err := s.redeemablesRepo.GetRedeemables(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get redeemables by event ID")
		return nil, err
	}
	return redeemables, nil
}

func (s *RedeemablesService) CreateRedeemable(ctx context.Context, name string, maxUserAmount int, redeemableType string) (*sqlc.Redeemable, error) {
	params := sqlc.CreateRedeemableParams{
		Name:          name,
		MaxUserAmount: int32(maxUserAmount),
		Type:          redeemableType,
	}
	redeemable, err := s.redeemablesRepo.CreateRedeemable(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create redeemable")
		return nil, err
	}
	return redeemable, nil
}

func (s *RedeemablesService) DeleteRedeemable(ctx context.Context, redeemableID uuid.UUID) error {
	err := s.redeemablesRepo.DeleteRedeemable(ctx, redeemableID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to delete redeemable")
		return err
	}
	return nil
}

func (s *RedeemablesService) UpdateRedeemable(ctx context.Context, redeemableID uuid.UUID, name *string, maxUserAmount *int, redeemableType *string) (*sqlc.Redeemable, error) {

	var maxUserAmount32 *int32
	if maxUserAmount != nil {
		v := int32(*maxUserAmount)
		maxUserAmount32 = &v
	}

	redeemable, err := s.redeemablesRepo.UpdateRedeemable(ctx, sqlc.UpdateRedeemableParams{
		ID:            redeemableID,
		Name:          name,
		MaxUserAmount: maxUserAmount32,
		Type:          redeemableType,
	})

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to update redeemable")
		return nil, err
	}
	return redeemable, nil
}

// RedeemForUser records one item handed to a hacker whose QR code staff
// scanned. The insert is conditional, so when it is refused a second lookup
// works out why and returns a *RedeemError.
func (s *RedeemablesService) RedeemForUser(ctx context.Context, redeemableID uuid.UUID, userID uuid.UUID) (*sqlc.RedeemForUserRow, error) {
	row, err := s.redeemablesRepo.RedeemForUser(ctx, sqlc.RedeemForUserParams{
		RedeemableID: redeemableID,
		UserID:       userID,
	})
	if err == nil {
		s.logger.Info().Str("type", "redeemable").Str("user_id", userID.String()).Str("redeemable_id", redeemableID.String()).
			Msgf("Redeemed by QR - Name: %s - Redeemable: %s - Items remaining: %d", row.Name, row.RedeemableName, row.ItemsRemaining)
		return row, nil
	}

	// check_redeemable_count_trigger raises P0001 when the hacker is at the limit.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "P0001" {
		eligibility, lookupErr := s.redeemablesRepo.GetRedeemEligibility(ctx, sqlc.GetRedeemEligibilityParams{
			RedeemableID: redeemableID,
			UserID:       userID,
		})
		if lookupErr != nil {
			return nil, &RedeemError{Reason: ErrRedeemLimitReached, Message: "Hacker already redeemed the limit for this item"}
		}
		return nil, &RedeemError{
			Reason:  ErrRedeemLimitReached,
			Message: fmt.Sprintf("Hacker already redeemed %s (limit %d per hacker)", eligibility.RedeemableName, eligibility.MaxUserAmount),
		}
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		s.logger.Error().Err(err).Msg("failed to redeem by QR")
		return nil, err
	}

	eligibility, err := s.redeemablesRepo.GetRedeemEligibility(ctx, sqlc.GetRedeemEligibilityParams{
		RedeemableID: redeemableID,
		UserID:       userID,
	})
	if err != nil {
		if !errors.Is(err, ErrRedeemableNotFound) {
			s.logger.Error().Err(err).Msg("failed to look up redeem eligibility")
		}
		return nil, err
	}

	switch {
	case !eligibility.UserExists:
		return nil, &RedeemError{Reason: ErrHackerNotFound, Message: "No hacker matches this QR code"}
	case eligibility.ApplicationStatus == "":
		return nil, &RedeemError{Reason: ErrNoApplication, Message: "Hacker has no application for this hackathon"}
	case eligibility.ApplicationStatus != string(sqlc.ApplicationStatusConfirmed):
		return nil, &RedeemError{
			Reason:  ErrNotConfirmed,
			Message: fmt.Sprintf("Only confirmed hackers can redeem items (status: %s)", strings.ReplaceAll(eligibility.ApplicationStatus, "_", " ")),
		}
	case !eligibility.CheckedIn:
		return nil, &RedeemError{Reason: ErrNotCheckedIn, Message: "Hacker is not checked in yet. Check them in first"}
	}

	// Eligible now, so the state changed between the two queries (for
	// example a check-in landed). Asking staff to scan again is enough.
	return nil, &RedeemError{Reason: ErrNotCheckedIn, Message: "Hacker's status just changed. Scan again"}
}
