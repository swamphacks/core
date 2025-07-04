package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
)

var (
	ErrEmailNotFound       = errors.New("email not found")
	ErrEmailConflict       = errors.New("email already exists in this mailing list")
	ErrFailedToAddEmail    = errors.New("failed to add email")
	ErrFailedToUpdateEmail = errors.New("failed to update email")
	ErrFailedToDeleteEmail = errors.New("failed to delete email")
	ErrFailedToGetEmail    = errors.New("failed to get email(s)")
)

type MailingListService struct {
	mailingRepo *repository.MailingListRepository
	logger      zerolog.Logger
}

func NewMailingListService(mailingRepo *repository.MailingListRepository, logger zerolog.Logger) *MailingListService {
	return &MailingListService{
		mailingRepo: mailingRepo,
		logger:      logger.With().Str("service", "MailingListService").Str("component", "mailing_list").Logger(),
	}
}

// AddEmail adds a new email to the mailing list.
func (s *MailingListService) AddEmail(ctx context.Context, eventID, userID uuid.UUID, email string) (sqlc.MailingListEmail, error) {
	params := sqlc.AddEmailParams{
		EventID: eventID,
		UserID:  userID,
		Email:   email,
	}

	result, err := s.mailingRepo.AddEmail(ctx, params)
	if err != nil {
		// Check for a unique constraint violation error (PostgreSQL error code 23505)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			s.logger.Warn().Err(err).Msg("attempted to add a duplicate email")
			return sqlc.MailingListEmail{}, ErrEmailConflict
		}
		s.logger.Error().Err(err).Msg("failed to add email")
		return sqlc.MailingListEmail{}, ErrFailedToAddEmail
	}

	return result, nil
}

// UpdateEmailByID updates an email's address by its primary key ID.
func (s *MailingListService) UpdateEmailByID(ctx context.Context, emailID uuid.UUID, newEmail string) (sqlc.MailingListEmail, error) {
	params := sqlc.UpdateEmailByIDParams{
		ID:    emailID,
		Email: newEmail,
	}
	result, err := s.mailingRepo.UpdateEmailByID(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.MailingListEmail{}, ErrEmailNotFound
		}
		s.logger.Error().Err(err).Msg("failed to update email")
		return sqlc.MailingListEmail{}, ErrFailedToUpdateEmail
	}
	return result, nil
}

// DeleteEmailByID deletes an email by its primary key ID.
func (s *MailingListService) DeleteEmailByID(ctx context.Context, emailID uuid.UUID) error {
	err := s.mailingRepo.DeleteEmailByID(ctx, emailID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to delete email by id")
		return ErrFailedToDeleteEmail
	}
	return nil
}

// DeleteEmailByUserAndEvent deletes an email using the logical user and event IDs.
func (s *MailingListService) DeleteEmailByUserAndEvent(ctx context.Context, eventID, userID uuid.UUID) error {
	params := sqlc.DeleteEmailByUserAndEventParams{
		EventID: eventID,
		UserID:  userID,
	}
	err := s.mailingRepo.DeleteEmailByUserAndEvent(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to delete email by user and event")
		return ErrFailedToDeleteEmail
	}
	return nil
}

// GetEmailsByEvent retrieves all emails for a specific event.
func (s *MailingListService) GetEmailsByEvent(ctx context.Context, eventID uuid.UUID) ([]sqlc.MailingListEmail, error) {
	emails, err := s.mailingRepo.GetEmailsByEvent(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get emails by event")
		return nil, ErrFailedToGetEmail
	}
	return emails, nil
}

// GetEmailsByUser retrieves all mailing list entries for a user.
func (s *MailingListService) GetEmailsByUser(ctx context.Context, userID uuid.UUID) ([]sqlc.MailingListEmail, error) {
	emails, err := s.mailingRepo.GetEmailsByUser(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get emails by user")
		return nil, ErrFailedToGetEmail
	}
	return emails, nil
}
