package announcements

import (
	"time"

	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type CreateAnnouncementRequest struct {
	HackathonID string     `json:"hackathonId" required:"true"`
	Title       string     `json:"title" minLength:"1"`
	Body        string     `json:"body" minLength:"1"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}

type UpdateAnnouncementRequest struct {
	Title     *string    `json:"title,omitempty"`
	Body      *string    `json:"body,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

type AnnouncementOutput struct {
	Body *sqlc.Announcement
}

type ListAnnouncementsOutput struct {
	Body []sqlc.Announcement
}

type DeleteAnnouncementOutput struct {
	Status int
}

type DismissAnnouncementOutput struct {
	Status int
}

type ListDismissedAnnouncementsOutput struct {
	Body []uuid.UUID
}