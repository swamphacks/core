package application

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
	"github.com/swamphacks/core/apps/api/internal/ctxutils"
)

var errStandbyClosed = errors.New("standby admission is not currently open")
var errStandbyPriority = errors.New("accept the next hacker in priority order")
var errStandbyEmpty = errors.New("no confirmed standby hackers are waiting")

type AcceptStandbyInput struct {
	Body struct {
		HackathonID string    `json:"hackathonId" required:"true"`
		UserID      uuid.UUID `json:"userId" required:"true"`
	}
}

func registerStandbyAcceptRoutes(h *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "accept-standby-hacker",
		Method:      http.MethodPost,
		Path:        "/waitlist/accept-standby",
		Summary:     "Accept Next Standby Hacker",
		Tags:        []string{"Application"},
		Middlewares: huma.Middlewares{
			mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma,
		},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
		Errors:     []int{401, 403, 409, 500},
	}, h.handleAcceptStandby)
}

func (h *handler) handleAcceptStandby(ctx context.Context, input *AcceptStandbyInput) (*AdmissionWaitlistOutput, error) {
	staff := ctxutils.GetUserFromCtx(ctx)
	if staff == nil {
		return nil, huma.Error401Unauthorized("Authentication required")
	}
	err := h.applicationService.acceptStandby(
		ctx, input.Body.HackathonID, input.Body.UserID, staff.UserID, time.Now().UTC(),
	)
	switch {
	case errors.Is(err, ErrInPersonWaitlistPermission):
		return nil, huma.Error403Forbidden(err.Error())
	case errors.Is(err, errStandbyClosed), errors.Is(err, errStandbyPriority), errors.Is(err, errStandbyEmpty):
		return nil, huma.Error409Conflict(err.Error())
	case err != nil:
		return nil, huma.Error500InternalServerError("Unable to accept standby hacker")
	}
	output := &AdmissionWaitlistOutput{}
	output.Body.Status = "accepted"
	return output, nil
}

func (s *ApplicationService) acceptStandby(ctx context.Context, eventID string, userID, staffID uuid.UUID, now time.Time) error {
	return s.txm.WithTx(ctx, func(tx pgx.Tx) error {
		var permitted bool
		if err := tx.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM users
				WHERE id=$1 AND role IN ('admin','staff') AND NOT is_fake)
		`, staffID).Scan(&permitted); err != nil {
			return err
		}
		if !permitted {
			return ErrInPersonWaitlistPermission
		}

		var opensAt, closesAt time.Time
		err := tx.QueryRow(ctx, `
			SELECT p.in_person_opens_at,p.invitations_close_at
			FROM hackathons h JOIN waitlist_dispatch_policies p ON p.hackathon_id=h.id
			WHERE h.id=$1 AND h.is_active AND p.enabled
			FOR UPDATE OF h,p
		`, eventID).Scan(&opensAt, &closesAt)
		if errors.Is(err, pgx.ErrNoRows) {
			return errStandbyClosed
		}
		if err != nil {
			return err
		}
		if now.Before(opensAt) || !now.Before(closesAt) {
			return errStandbyClosed
		}

		var nextUser, applicationID uuid.UUID
		err = tx.QueryRow(ctx, `
			SELECT a.user_id,a.id
			FROM applications a
			JOIN users u ON u.id=a.user_id
			JOIN waitlist w ON w.hackathon_id=a.hackathon_id AND w.user_id=a.user_id
			WHERE a.hackathon_id=$1 AND a.status::text='waitlist_confirmed'
				AND u.role='applicant' AND NOT a.is_fake AND NOT u.is_fake
				AND w.in_person_joined_at IS NOT NULL
			ORDER BY CASE WHEN w.signup_source='preregistered' THEN 0 ELSE 1 END,
				w.in_person_joined_at,w.created_at,w.user_id
			LIMIT 1 FOR UPDATE OF a,w
		`, eventID).Scan(&nextUser, &applicationID)
		if errors.Is(err, pgx.ErrNoRows) {
			return errStandbyEmpty
		}
		if err != nil {
			return err
		}
		if nextUser != userID {
			return errStandbyPriority
		}
		_, err = tx.Exec(ctx, `
			UPDATE applications SET status='accepted',updated_at=now()
			WHERE id=$1 AND status::text='waitlist_confirmed'
		`, applicationID)
		return err
	})
}
