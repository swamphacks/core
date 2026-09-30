package application

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type AdmissionWaitlistOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

func registerAdmissionWaitlistRoutes(h *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "join-admission-waitlist",
		Method:      http.MethodPost,
		Path:        "/join-waitlist",
		Summary:     "Join Waitlist",
		Description: "Join the active hackathon waitlist after rejection.",
		Tags:        []string{"Application"},
		Middlewares: huma.Middlewares{mw.Auth.RequireAuthHuma},
		Parameters:  []*huma.Param{cookie.SessionCookieHumaParam},
		Errors: []int{
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusConflict,
			http.StatusInternalServerError,
		},
		DefaultStatus: http.StatusOK,
	}, h.handleJoinAdmissionWaitlist)
}

func (h *handler) handleJoinAdmissionWaitlist(
	ctx context.Context, input *struct{},
) (*AdmissionWaitlistOutput, error) {
	user := ctxutils.GetUserFromCtx(ctx)
	if user == nil {
		return nil, huma.Error401Unauthorized("Authentication required")
	}
	if user.Role != sqlc.RoleApplicant {
		return nil, huma.Error403Forbidden("Only applicants can join the waitlist")
	}
	err := h.applicationService.JoinAdmissionWaitlist(ctx, user.UserID)
	if errors.Is(err, ErrWaitlistEligibility) {
		return nil, huma.Error409Conflict(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to join waitlist")
	}
	output := &AdmissionWaitlistOutput{}
	output.Body.Status = "waitlisted"
	return output, nil
}
