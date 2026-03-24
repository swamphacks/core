package hackathon

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	. "github.com/swamphacks/core/apps/api/internal/parse"
)

func RegisterRoutes(hackathonHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "get-hackathon",
		Method:      http.MethodGet,
		Summary:     "Get Hackathon",
		Description: "Returns information of the hackathon",
		Tags:        []string{"Hackathon"},
		Path:        "/hackathon",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, hackathonHandler.handleGetHackathon)

	huma.Register(group, huma.Operation{
		OperationID:   "update-hackathon",
		Method:        http.MethodPatch,
		Summary:       "Update Hackathon",
		Description:   "Updates the information of the hackathon",
		Tags:          []string{"Hackathon"},
		Path:          "",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, hackathonHandler.handleUpdateHackathon)

	huma.Register(group, huma.Operation{
		OperationID: "get-hackathon-staff",
		Method:      http.MethodGet,
		Summary:     "Get Hackathon Staff",
		Description: "Returns the users who are part of the current staff of the hackathon",
		Tags:        []string{"Hackathon"},
		Path:        "/staff",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, hackathonHandler.handleGetStaff)

	huma.Register(group, huma.Operation{
		OperationID: "get-hackathon-attendees-with-discord",
		Method:      http.MethodGet,
		Summary:     "Get Hackathon Attendees with Discord",
		Description: "Returns all users with a discord account that is also attending the hackathon",
		Tags:        []string{"Hackathon"},
		Path:        "/attendees/discord",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, hackathonHandler.handleGetAttendeesWithDiscord)

	huma.Register(group, huma.Operation{
		OperationID: "get-hackathon-attendees-userids",
		Method:      http.MethodGet,
		Summary:     "Get Hackathon Attendees User Ids",
		Description: "Returns all users ids of users who are attending the hackathon",
		Tags:        []string{"Hackathon"},
		Path:        "/attendees/userids",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, hackathonHandler.handleGetAttendeeUserIds)

	huma.Register(group, huma.Operation{
		OperationID: "get-hackathon-attendees-count",
		Method:      http.MethodGet,
		Summary:     "Get Hackathon Attendees Count",
		Description: "Returns the number of users who is attending the hackathon",
		Tags:        []string{"Hackathon"},
		Path:        "/attendees/count",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, hackathonHandler.handleGetAttendeeCount)
}

type handler struct {
	hackathonService *HackathonService
	config           *config.Config
	logger           zerolog.Logger
}

func NewHandler(hackathonService *HackathonService, config *config.Config, logger zerolog.Logger) *handler {
	return &handler{
		hackathonService: hackathonService,
		config:           config,
		logger:           logger.With().Str("handler", "HackathonHandler").Str("domain", "hackathon").Logger(),
	}
}

type GetHackathonOutput struct {
	Body *sqlc.Hackathon
}

func (h *handler) handleGetHackathon(ctx context.Context, input *struct{}) (*GetHackathonOutput, error) {
	hackathon, err := h.hackathonService.GetHackathon(ctx)

	if err != nil {
		h.logger.Err(err).Msg("")
		if errors.Is(err, repository.ErrEntityNotFound) {
			return nil, huma.Error404NotFound("Hackathon not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get hackathon")
	}

	return &GetHackathonOutput{Body: hackathon}, nil
}

type UpdateHackathonOutput struct {
	Status int
}

type UpdateHackathonRequest struct {
	Name             OmittableNullable[string]     `json:"name,omitempty"`
	Description      OmittableNullable[*string]    `json:"description,omitempty"`
	Location         OmittableNullable[*string]    `json:"location,omitempty"`
	LocationUrl      OmittableNullable[*string]    `json:"location_url,omitempty"`
	MaxAttendees     OmittableNullable[*int32]     `json:"max_attendees,omitempty"`
	ApplicationOpen  OmittableNullable[time.Time]  `json:"application_open,omitempty"`
	ApplicationClose OmittableNullable[time.Time]  `json:"application_close,omitempty"`
	RsvpDeadline     OmittableNullable[*time.Time] `json:"rsvp_deadline,omitempty"`
	DecisionRelease  OmittableNullable[*time.Time] `json:"decision_release,omitempty"`
	StartTime        OmittableNullable[time.Time]  `json:"start_time,omitempty"`
	EndTime          OmittableNullable[time.Time]  `json:"end_time,omitempty"`
	WebsiteUrl       OmittableNullable[*string]    `json:"website_url,omitempty"`
	IsPublished      OmittableNullable[bool]       `json:"is_published,omitempty"`
}

func (h *handler) handleUpdateHackathon(ctx context.Context, input *struct {
	Body UpdateHackathonRequest
}) (*UpdateHackathonOutput, error) {
	params := sqlc.UpdateHackathonParams{
		NameDoUpdate: input.Body.Name.Sent,
		Name:         input.Body.Name.Value,

		DescriptionDoUpdate: input.Body.Description.Sent,
		Description:         input.Body.Description.Value,

		LocationDoUpdate: input.Body.Location.Sent,
		Location:         input.Body.Location.Value,

		LocationUrlDoUpdate: input.Body.LocationUrl.Sent,
		LocationUrl:         input.Body.LocationUrl.Value,

		MaxAttendeesDoUpdate: input.Body.MaxAttendees.Sent,
		MaxAttendees:         input.Body.MaxAttendees.Value,

		ApplicationOpenDoUpdate: input.Body.ApplicationOpen.Sent,
		ApplicationOpen:         input.Body.ApplicationOpen.Value,

		ApplicationCloseDoUpdate: input.Body.ApplicationClose.Sent,
		ApplicationClose:         input.Body.ApplicationClose.Value,

		RsvpDeadlineDoUpdate: input.Body.RsvpDeadline.Sent,
		RsvpDeadline:         input.Body.RsvpDeadline.Value,

		DecisionReleaseDoUpdate: input.Body.DecisionRelease.Sent,
		DecisionRelease:         input.Body.DecisionRelease.Value,

		StartTimeDoUpdate: input.Body.StartTime.Sent,
		StartTime:         input.Body.StartTime.Value,

		EndTimeDoUpdate: input.Body.EndTime.Sent,
		EndTime:         input.Body.EndTime.Value,

		WebsiteUrlDoUpdate: input.Body.WebsiteUrl.Sent,
		WebsiteUrl:         input.Body.WebsiteUrl.Value,

		IsPublishedDoUpdate: input.Body.IsPublished.Sent,
		IsPublished:         &input.Body.IsPublished.Value,

		BannerDoUpdate: false, // Banners are uploaded using a separate endpoint
		Banner:         nil,
	}

	err := h.hackathonService.UpdateHackathon(ctx, params)

	if err != nil {
		h.logger.Err(err).Msg("")
		return nil, huma.Error500InternalServerError("Failed to get update hackathon")
	}

	return &UpdateHackathonOutput{Status: http.StatusOK}, nil
}

type GetStaffOutput struct {
	Body []sqlc.GetStaffRow
}

func (h *handler) handleGetStaff(ctx context.Context, input *struct{}) (*GetStaffOutput, error) {
	staff, err := h.hackathonService.GetStaff(ctx)

	if err != nil {
		h.logger.Err(err).Msg("")
		return nil, huma.Error500InternalServerError("Failed to get staff")
	}

	return &GetStaffOutput{Body: staff}, nil
}

type GetAttendeesWithDiscordOutput struct {
	Body []sqlc.GetAttendeesWithDiscordRow
}

func (h *handler) handleGetAttendeesWithDiscord(ctx context.Context, input *struct{}) (*GetAttendeesWithDiscordOutput, error) {
	attendees, err := h.hackathonService.GetAttendeesWithDiscord(ctx)

	if err != nil {
		h.logger.Err(err).Msg("")
		return nil, huma.Error500InternalServerError("Failed to get attendees with discord")
	}

	return &GetAttendeesWithDiscordOutput{Body: attendees}, nil
}

type GetAttendeeUserIdsOutput struct {
	Body []uuid.UUID
}

func (h *handler) handleGetAttendeeUserIds(ctx context.Context, input *struct{}) (*GetAttendeeUserIdsOutput, error) {
	userIds, err := h.hackathonService.GetAttendeeUserIds(ctx)

	if err != nil {
		h.logger.Err(err).Msg("")
		return nil, huma.Error500InternalServerError("Failed to get attendee user ids")
	}

	return &GetAttendeeUserIdsOutput{Body: userIds}, nil
}

type GetAttendeeCountOutput struct {
	Body int64
}

func (h *handler) handleGetAttendeeCount(ctx context.Context, input *struct{}) (*GetAttendeeCountOutput, error) {
	count, err := h.hackathonService.GetAttendeeCount(ctx)

	if err != nil {
		h.logger.Err(err).Msg("")
		return nil, huma.Error500InternalServerError("Failed to get attendees count")
	}

	return &GetAttendeeCountOutput{Body: count}, nil
}
