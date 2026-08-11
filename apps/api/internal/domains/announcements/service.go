package announcements

import (
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
)

type AnnouncementsService struct {
	announcementsRepo *repository.AnnouncementsRepository
	logger            zerolog.Logger
}

func NewService(
	announcementsRepo *repository.AnnouncementsRepository,
	logger zerolog.Logger,
) *AnnouncementsService {
	return &AnnouncementsService{
		announcementsRepo: announcementsRepo,
		logger: logger.With().Str("service", "AnnouncementsService").Str("domain", "announcements").Logger(),
	}
}