package workshops

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
)

func registerAttendanceRoutes(h *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "list-attendance-workshops", Method: http.MethodGet, Path: "/attendance-options", Summary: "List workshops and socials for staff scanning", Tags: []string{"Workshops"},
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma}, Parameters: []*huma.Param{cookie.SessionCookieHumaParam}, Errors: []int{401, 403, 500},
	}, h.handleAttendanceOptions)
	huma.Register(group, huma.Operation{
		OperationID: "record-workshop-attendance", Method: http.MethodPost, Path: "/{workshopId}/attendance", Summary: "Record workshop attendance by hacker QR", Tags: []string{"Workshops"},
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma}, Parameters: []*huma.Param{cookie.SessionCookieHumaParam}, Errors: []int{401, 403, 404, 409, 500},
	}, h.handleAttendance)
}

type AttendanceOptionsOutput struct {
	Body []repository.AttendanceWorkshop
}

func (h *handler) handleAttendanceOptions(ctx context.Context, _ *struct{}) (*AttendanceOptionsOutput, error) {
	rows, err := h.workshopService.workshopsRepo.AttendanceWorkshops(ctx)
	if err != nil {
		h.logger.Err(err).Msg("Unable to list attendance workshops")
		return nil, huma.Error500InternalServerError("Unable to load workshops")
	}
	return &AttendanceOptionsOutput{Body: rows}, nil
}

type AttendanceInput struct {
	WorkshopID uuid.UUID `path:"workshopId"`
	Body       struct {
		UserID uuid.UUID `json:"userId" required:"true"`
	}
}
type AttendanceOutput struct {
	Body *repository.WorkshopAttendance
}

func (h *handler) handleAttendance(ctx context.Context, input *AttendanceInput) (*AttendanceOutput, error) {
	row, err := h.workshopService.workshopsRepo.RecordAttendance(ctx, input.WorkshopID, input.Body.UserID)
	if errors.Is(err, repository.ErrWorkshopAttendanceNotFound) {
		return nil, huma.Error404NotFound("Workshop or hacker application not found")
	}
	if errors.Is(err, repository.ErrWorkshopAttendanceIneligible) {
		return nil, huma.Error409Conflict("Hacker must be confirmed and checked in first")
	}
	if err != nil {
		h.logger.Err(err).Msg("Unable to record workshop attendance")
		return nil, huma.Error500InternalServerError("Unable to record attendance")
	}
	h.logger.Info().Str("user_id", input.Body.UserID.String()).Str("workshop_id", input.WorkshopID.String()).Bool("already_recorded", row.AlreadyRecorded).Msg("Workshop QR attendance")
	return &AttendanceOutput{Body: row}, nil
}
