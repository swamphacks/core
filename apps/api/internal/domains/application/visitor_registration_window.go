package application

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
)

type VisitorRegistrationWindowOutput struct {
	Body struct {
		Open  bool `json:"open"`
		DayOf bool `json:"dayOf"`
	}
}

func registerVisitorRegistrationWindowRoutes(h *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "get-visitor-registration-window",
		Method:      http.MethodGet,
		Path:        "/visitor-registration-window",
		Summary:     "Get Visitor Registration Availability",
		Tags:        []string{"Application"},
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
		Errors:      []int{401, 500},
	}, h.handleVisitorRegistrationWindow)
}

func (h *handler) handleVisitorRegistrationWindow(ctx context.Context, input *struct{}) (*VisitorRegistrationWindowOutput, error) {
	output := &VisitorRegistrationWindowOutput{}
	now := time.Now().UTC()
	var active bool
	err := h.applicationService.db.Pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM hackathons WHERE is_active),
			EXISTS(
				SELECT 1 FROM hackathons h
				JOIN waitlist_dispatch_policies p ON p.hackathon_id=h.id
				WHERE h.is_active AND p.enabled
					AND p.in_person_opens_at <= $1 AND p.invitations_close_at > $1
			)
	`, now).Scan(&active, &output.Body.DayOf)
	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to load registration availability")
	}
	output.Body.Open = active && (output.Body.DayOf || admissionWaitlistIsOpen(now))
	return output, nil
}
