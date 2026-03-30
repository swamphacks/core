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
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
	"github.com/swamphacks/core/apps/api/internal/emailutils"
	. "github.com/swamphacks/core/apps/api/internal/parse"
)

func RegisterRoutes(hackathonHandler *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "get-hackathon",
		Method:      http.MethodGet,
		Summary:     "Get Hackathon",
		Description: "Returns information of the hackathon",
		Tags:        []string{"Hackathon"},
		Path:        "",
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
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
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
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
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
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
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
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
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
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Errors:      []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
	}, hackathonHandler.handleGetAttendeeCount)

	huma.Register(group, huma.Operation{
		OperationID:   "check-in",
		Method:        http.MethodPost,
		Summary:       "Check In User",
		Description:   "Staff route for checking a user to an event. The user to check in must be an attendee and have never been checked in yet.",
		Tags:          []string{"Hackathon"},
		Path:          "/checkin",
		Middlewares:   huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma},
		Errors:        []int{http.StatusUnauthorized, http.StatusInternalServerError},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		DefaultStatus: http.StatusOK,
	}, hackathonHandler.handleCheckIn)

	huma.Register(group, huma.Operation{
		OperationID:   "submit-interest-email",
		Method:        http.MethodPost,
		Summary:       "Submit Interest Email",
		Description:   "Submits an email to interest/mailing list for the hackathon",
		Tags:          []string{"Hackathon"},
		Path:          "/interest", // public route
		Errors:        []int{http.StatusBadRequest, http.StatusInternalServerError},
		DefaultStatus: http.StatusOK,
	}, hackathonHandler.handleSubmitInterestEmail)

	huma.Register(group, huma.Operation{
		OperationID: "upload-banner",
		Method:      http.MethodPost,
		Summary:     "Upload Banner",
		Description: "Uploads an image to be used as the banner for the hackathon",
		Tags:        []string{"Hackathon"},
		Path:        "/banner",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
		Errors:      []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, hackathonHandler.handleUploadBanner)

	huma.Register(group, huma.Operation{
		OperationID: "delete-banner",
		Method:      http.MethodDelete,
		Summary:     "Delete Banner",
		Description: "Deletes the banner",
		Tags:        []string{"Hackathon"},
		Path:        "/banner",
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
		Errors:      []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, hackathonHandler.handleDeleteBanner)
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
		if errors.Is(err, database.ErrEntityNotFound) {
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
	LocationUrl      OmittableNullable[*string]    `json:"locationUrl,omitempty"`
	MaxAttendees     OmittableNullable[*int32]     `json:"maxAttendees,omitempty"`
	ApplicationOpen  OmittableNullable[time.Time]  `json:"applicationOpen,omitempty"`
	ApplicationClose OmittableNullable[time.Time]  `json:"applicationClose,omitempty"`
	RsvpDeadline     OmittableNullable[*time.Time] `json:"rsvpDeadline,omitempty"`
	DecisionRelease  OmittableNullable[*time.Time] `json:"decisionRelease,omitempty"`
	StartTime        OmittableNullable[time.Time]  `json:"startTime,omitempty"`
	EndTime          OmittableNullable[time.Time]  `json:"endTime,omitempty"`
	WebsiteUrl       OmittableNullable[*string]    `json:"websiteUrl,omitempty"`
	IsPublished      OmittableNullable[bool]       `json:"isPublished,omitempty"`
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
	Body []sqlc.User
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

type CheckInRequest struct {
	UserID uuid.UUID `json:"userId"`
	RFID   *string   `json:"rfid"`
}

type CheckInOutput struct {
	Status int
}

func (h *handler) handleCheckIn(ctx context.Context, input *struct {
	Body CheckInRequest
}) (*CheckInOutput, error) {
	if input.Body.RFID != nil {
		if *input.Body.RFID == "" {
			input.Body.RFID = nil
		}
	}

	err := h.hackathonService.CheckInAttendee(ctx, input.Body.UserID, input.Body.RFID)

	if err != nil {
		h.logger.Err(err).Msg("check in user failed")
		return nil, huma.Error500InternalServerError("Failed to check in user")
	}

	return &CheckInOutput{Status: http.StatusOK}, nil
}

type SubmitInterestEmailRequest struct {
	Email  string  `json:"email"`
	Source *string `json:"source"`
}

type SubmitInterestEmailOutput struct {
	Status int
}

func (h *handler) handleSubmitInterestEmail(ctx context.Context, input *struct {
	Body SubmitInterestEmailRequest
}) (*SubmitInterestEmailOutput, error) {
	if !emailutils.IsValidEmail(input.Body.Email) {
		return nil, huma.Error400BadRequest("Invalid email")
	}

	_, err := h.hackathonService.SubmitInterestEmail(ctx, input.Body.Email, input.Body.Source)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to submit interest email")
	}

	return &SubmitInterestEmailOutput{Status: http.StatusOK}, nil
}

type UploadBannerOutput struct {
	Body *string
}

func (h *handler) handleUploadBanner(ctx context.Context, input *struct {
	RawBody huma.MultipartFormFiles[struct {
		Image huma.FormFile `form:"image" contentType:"image/png, image/jpeg, image/jpg" required:"true"`
	}]
}) (*UploadBannerOutput, error) {
	fileHeader := input.RawBody.Form.File["image"][0]

	if fileHeader.Size > 5*1024*1024 { // 5 MiB
		return nil, huma.Error400BadRequest("File too large")
	}

	file, err := fileHeader.Open()

	if err != nil {
		return nil, huma.Error400BadRequest("Failed to parse uploaded banner image")
	}

	url, err := h.hackathonService.UploadBanner(ctx, file, fileHeader)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to upload banner")
	}

	return &UploadBannerOutput{Body: url}, nil
}

type DeleteBannerOutput struct {
	Status int
}

func (h *handler) handleDeleteBanner(ctx context.Context, input *struct{}) (*DeleteBannerOutput, error) {
	err := h.hackathonService.DeleteBanner(ctx)

	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete banner")
	}

	return &DeleteBannerOutput{Status: http.StatusOK}, nil
}
