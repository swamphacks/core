package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	ctxu "github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/storage"
	"golang.org/x/sync/errgroup"
)

var (
	ErrFailedToCreateEvent  = errors.New("failed to create event")
	ErrFailedToGetEvent     = errors.New("failed to get event")
	ErrFailedToUpdateEvent  = errors.New("failed to update event")
	ErrFailedToDeleteEvent  = errors.New("failed to delete event")
	ErrFailedToParseUUID    = errors.New("failed to parse uuid")
	ErrMissingFields        = errors.New("missing fields")
	ErrMissingPerms         = errors.New("missing perms")
	ErrFailedToUploadBanner = errors.New("failed to upload banner")
	ErrUnexpectedFileType   = errors.New("did not expect this file type")

	ErrFailedToSubmitApplication = errors.New("failed to submit application")
	ErrGetEventOverview          = errors.New("failed to aggregate event stats")
)

type EventService struct {
	eventRepo *repository.EventRepository
	userRepo  *repository.UserRepository
	storage   storage.Storage
	buckets   *config.CoreBuckets
	logger    zerolog.Logger
}

func NewEventService(eventRepo *repository.EventRepository, userRepo *repository.UserRepository, storage storage.Storage, buckets *config.CoreBuckets, logger zerolog.Logger) *EventService {
	return &EventService{
		eventRepo: eventRepo,
		userRepo:  userRepo,
		storage:   storage,
		buckets:   buckets,
		logger:    logger.With().Str("service", "EventService").Str("component", "events").Logger(),
	}
}

func (s *EventService) CreateEvent(ctx context.Context, params sqlc.CreateEventParams) (*sqlc.Event, error) {
	event, err := s.eventRepo.CreateEvent(ctx, params)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			s.logger.Err(err).Msg(repository.ErrEventNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, ErrFailedToCreateEvent
	}

	return event, nil
}

func (s *EventService) GetEventByID(ctx context.Context, id uuid.UUID) (*sqlc.Event, error) {
	event, err := s.eventRepo.GetEventByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			s.logger.Err(err).Msg(repository.ErrEventNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, ErrFailedToGetEvent
	}

	return event, nil
}

func (s *EventService) UpdateEventById(ctx context.Context, params sqlc.UpdateEventByIdParams) (*sqlc.Event, error) {
	err := s.eventRepo.UpdateEventById(ctx, params)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			s.logger.Err(err).Msg(repository.ErrEventNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, ErrFailedToUpdateEvent
	}

	event, err := s.eventRepo.GetEventByID(ctx, params.ID)

	return event, err
}

func (s *EventService) DeleteEventById(ctx context.Context, id uuid.UUID) error {
	err := s.eventRepo.DeleteEventById(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrEventNotFound:
			s.logger.Err(err).Msg(repository.ErrEventNotFound.Error())
		case repository.ErrNoEventsDeleted:
			s.logger.Err(err).Msg(repository.ErrEventNotFound.Error())
		case repository.ErrMultipleEventsDeleted:
			s.logger.Err(err).Msg(repository.ErrMultipleEventsDeleted.Error())
		default:
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return ErrFailedToDeleteEvent
	}

	return err
}

func (s *EventService) GetEvents(ctx context.Context, scope sqlc.GetEventScopeType) (*[]sqlc.GetEventsWithUserInfoRow, error) {
	isSuperuser := ctxu.IsSuperuser(ctx)
	userId := ctxu.GetUserIdFromCtx(ctx)

	// Non-superusers can't get all unpublished events
	if !isSuperuser && scope == "all" {
		return nil, ErrMissingPerms
	}

	return s.eventRepo.GetEventsWithRoles(ctx, userId, scope)

}

func (s *EventService) GetEventRoleByIds(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (*sqlc.EventRole, error) {
	eventRole, err := s.eventRepo.GetEventRoleByIds(ctx, userId, eventId)
	if err != nil {
		if errors.Is(err, repository.ErrEventRoleNotFound) {
			s.logger.Err(err).Msg(repository.ErrEventRoleNotFound.Error())
		} else {
			s.logger.Err(err).Msg(repository.ErrUnknown.Error())
		}
		return nil, err
	}

	return eventRole, err
}

func (s *EventService) GetEventStaffUsers(ctx context.Context, eventId uuid.UUID) (*[]sqlc.GetEventStaffRow, error) {
	return s.eventRepo.GetEventStaff(ctx, eventId)
}

func (s *EventService) AssignEventRole(
	ctx context.Context,
	userId *uuid.UUID,
	email *string,
	eventId uuid.UUID,
	role sqlc.EventRoleType,
) error {
	if userId == nil && email == nil {
		return errors.New("must provide either userId or email")
	}

	var selectedUser *sqlc.AuthUser
	var err error

	if userId != nil {
		selectedUser, err = s.userRepo.GetByID(ctx, *userId)
		// Do not return if user not found, the query needs to fallback to other optiosn
		if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
			s.logger.Err(err).Msg("Something went wrong getting by id")
			return err
		}
	}

	if selectedUser == nil && email != nil {
		selectedUser, err = s.userRepo.GetByEmail(ctx, *email)
		if err != nil {
			s.logger.Err(err).Msg("Something went wrong getting by email")
			return err
		}
	}

	// Just a double safety check (should usually be caught by queries above)
	if selectedUser == nil {
		s.logger.Warn().Msg(("User not found from email OR id"))
		return repository.ErrUserNotFound
	}

	// Now assign the event role
	err = s.eventRepo.AssignRole(ctx, sqlc.AssignRoleParams{
		EventID: eventId,
		UserID:  selectedUser.ID,
		Role:    role,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *EventService) RevokeEventRole(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) error {
	return s.eventRepo.RevokeRole(ctx, userId, eventId)
}

func (s *EventService) IsApplicationsOpen(ctx context.Context, eventId uuid.UUID) (bool, error) {
	event, err := s.GetEventByID(ctx, eventId)
	if err != nil {
		s.logger.Err(err).Msg("ApplicationOpen check error: " + err.Error())
		return false, err
	}

	if !*event.IsPublished {
		return false, nil
	}

	now := time.Now()
	open := now.After(event.ApplicationOpen) && now.Before(event.ApplicationClose)
	return open, nil
}

func (s *EventService) UploadBanner(ctx context.Context, eventId uuid.UUID, banner multipart.File, header *multipart.FileHeader) (*string, error) {
	bannerFileBuffer := bytes.NewBuffer(nil)

	fileName := header.Filename
	fileExt := strings.ToLower(filepath.Ext(fileName))

	s.logger.Info().Str("Filetype", fileExt).Msg("The file type")

	switch fileExt {
	case ".jpg", ".png", ".jpeg":
		// Do nothing
	default:
		return nil, ErrUnexpectedFileType
	}

	fileType := mime.TypeByExtension(fileExt)

	if fileType == "" {
		return nil, ErrFailedToUploadBanner
	}

	if _, err := io.Copy(bannerFileBuffer, banner); err != nil {
		return nil, ErrFailedToUploadBanner
	}

	uploadKey := fmt.Sprintf("%s/banner%s", eventId, fileExt)

	err := s.storage.Store(ctx, s.buckets.EventAssets, uploadKey, bannerFileBuffer.Bytes(), &fileType)
	if err != nil {
		return nil, ErrFailedToUploadBanner
	}

	// Reconstrust URL with cache buster
	url := fmt.Sprintf("%s/%s?t=%d", s.buckets.EventAssetsBaseUrl, uploadKey, time.Now().Unix())

	err = s.eventRepo.UpdateEventById(ctx, sqlc.UpdateEventByIdParams{
		ID:             eventId,
		BannerDoUpdate: true,
		Banner:         &url,
	})
	if err != nil {
		return nil, ErrFailedToUpdateEvent
	}

	return &url, nil

}

func (s *EventService) DeleteBanner(ctx context.Context, eventId uuid.UUID) error {
	// For now its a soft delete, not actually deleting banner is easiest, just set to null
	return s.eventRepo.UpdateEventById(ctx, sqlc.UpdateEventByIdParams{
		ID:             eventId,
		BannerDoUpdate: true,
		Banner:         nil,
	})
}

type SubmissionTimesStatistics struct {
	Day   time.Time `json:"day" format:"date-time"`
	Count int64     `json:"count"`
}

type EventOverview struct {
	EventDetails                sqlc.Event                        `json:"event_details"`
	ApplicationStatusStatistics sqlc.GetApplicationStatusSplitRow `json:"application_status_stats"`
	SubmissionTimesStatistics   []SubmissionTimesStatistics       `json:"application_submission_stats"`
}

func (s *EventService) GetEventOverview(ctx context.Context, eventId uuid.UUID) (*EventOverview, error) {
	g, ctx := errgroup.WithContext(ctx)

	var eventDetails *sqlc.Event
	var statusStats sqlc.GetApplicationStatusSplitRow
	var submissionTimesStats []SubmissionTimesStatistics

	g.Go(func() error {
		var err error
		eventDetails, err = s.GetEventByID(ctx, eventId)
		return err
	})

	g.Go(func() error {
		var err error
		statusStats, err = s.eventRepo.GetApplicationStatuses(ctx, eventId)
		return err
	})

	g.Go(func() error {
		var err error
		temp, err := s.eventRepo.GetSubmissionTimes(ctx, eventId)

		for _, v := range temp {
			submissionTimesStats = append(submissionTimesStats, SubmissionTimesStatistics{
				Count: v.Count,
				Day:   v.Day.Time,
			})
		}

		fmt.Println(submissionTimesStats)
		return err
	})

	if err := g.Wait(); err != nil {
		s.logger.Err(err).Msg("Something went wrong while getting event overview stats")
		return nil, ErrGetEventOverview
	}

	return &EventOverview{
		EventDetails:                *eventDetails,
		ApplicationStatusStatistics: statusStats,
		SubmissionTimesStatistics:   submissionTimesStats,
	}, nil
}
