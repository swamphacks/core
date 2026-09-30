package application

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
)

type InPersonAdmissionWaitlistInput struct {
	Body struct {
		UserID uuid.UUID `json:"userId" required:"true"`
	}
}

func registerInPersonAdmissionWaitlistRoutes(
	h *handler, group huma.API, mw *middleware.Middleware,
) {
	huma.Register(group, huma.Operation{
		OperationID: "record-in-person-admission-waitlist",
		Method:      http.MethodPost,
		Path:        "/waitlist/in-person",
		Summary:     "Record In-Person Waitlist Arrival",
		Tags:        []string{"Application"},
		Middlewares: huma.Middlewares{
			mw.Auth.RequireAuthHuma, mw.Auth.RequireAdminHuma,
		},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		Errors:        []int{401, 403, 409, 500},
		DefaultStatus: http.StatusOK,
	}, h.handleRecordInPersonAdmissionWaitlist)
}

func (h *handler) handleRecordInPersonAdmissionWaitlist(
	ctx context.Context, input *InPersonAdmissionWaitlistInput,
) (*AdmissionWaitlistOutput, error) {
	user := ctxutils.GetUserFromCtx(ctx)
	if user == nil {
		return nil, huma.Error401Unauthorized("Authentication required")
	}
	err := h.applicationService.RecordInPersonAdmissionWaitlist(
		ctx, input.Body.UserID, user.UserID, time.Now().UTC(),
	)
	if errors.Is(err, ErrInPersonWaitlistPermission) {
		return nil, huma.Error403Forbidden(err.Error())
	}
	if errors.Is(err, ErrInPersonWaitlistEligibility) ||
		errors.Is(err, ErrInPersonWaitlistClosed) {
		return nil, huma.Error409Conflict(err.Error())
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to record in-person arrival")
	}
	output := &AdmissionWaitlistOutput{}
	output.Body.Status = "waitlisted"
	return output, nil
}
