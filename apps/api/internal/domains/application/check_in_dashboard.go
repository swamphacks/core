package application

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/swamphacks/core/apps/api/internal/api/cookie"
	"github.com/swamphacks/core/apps/api/internal/api/middleware"
)

type CheckInDashboardInput struct {
	HackathonID  string `query:"hackathonId"`
	Search       string `query:"search"`
	Status       string `query:"status"`
	CheckedIn    string `query:"checkedIn" enum:",yes,no"`
	RedeemableID string `query:"redeemableId"`
}

type CheckInDashboardRow struct {
	UserID         uuid.UUID       `json:"userId"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	Status         string          `json:"status"`
	CheckedInAt    *time.Time      `json:"checkedInAt"`
	RFID           *string         `json:"rfid"`
	SignupSource   string          `json:"signupSource"`
	StandbyArrival *time.Time      `json:"standbyArrival"`
	Redemptions    json.RawMessage `json:"redemptions"`
}

type CheckInDashboardOutput struct {
	Body struct {
		HackathonID        string                `json:"hackathonId"`
		Confirmed          int64                 `json:"confirmed"`
		ConfirmedCheckedIn int64                 `json:"confirmedCheckedIn"`
		DayOfSignups       int64                 `json:"dayOfSignups"`
		StandbyWaiting     int64                 `json:"standbyWaiting"`
		Rows               []CheckInDashboardRow `json:"rows"`
	}
}

func registerCheckInDashboardRoutes(h *handler, group huma.API, mw *middleware.Middleware) {
	huma.Register(group, huma.Operation{
		OperationID: "get-check-in-dashboard",
		Method:      http.MethodGet,
		Path:        "/check-in-dashboard",
		Summary:     "Get Staff Check-In Dashboard",
		Tags:        []string{"Application"},
		Middlewares: huma.Middlewares{
			mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma,
		},
		Parameters: []*huma.Param{cookie.SessionCookieHumaParam},
		Errors:     []int{400, 401, 403, 404, 500},
	}, h.handleCheckInDashboard)
}

func (h *handler) handleCheckInDashboard(ctx context.Context, input *CheckInDashboardInput) (*CheckInDashboardOutput, error) {
	var redeemableID *uuid.UUID
	if input.RedeemableID != "" {
		id, err := uuid.Parse(input.RedeemableID)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid redeemable ID")
		}
		redeemableID = &id
	}

	output := &CheckInDashboardOutput{}
	output.Body.Rows = []CheckInDashboardRow{}
	err := h.applicationService.db.Pool.QueryRow(ctx, `
		SELECT id FROM hackathons
		WHERE ($1='' AND is_active) OR ($1<>'' AND id=$1)
	`, input.HackathonID).Scan(&output.Body.HackathonID)
	if err != nil {
		return nil, huma.Error404NotFound("Hackathon not found")
	}

	err = h.applicationService.db.Pool.QueryRow(ctx, `
		SELECT
			count(*) FILTER (WHERE a.status::text='confirmed'),
			count(*) FILTER (WHERE a.status::text='confirmed' AND u.checked_in_at IS NOT NULL),
			count(*) FILTER (WHERE w.signup_source='day_of'),
			count(*) FILTER (WHERE a.status::text='waitlist_confirmed')
		FROM applications a
		JOIN users u ON u.id=a.user_id
		LEFT JOIN waitlist w ON w.hackathon_id=a.hackathon_id AND w.user_id=a.user_id
		WHERE a.hackathon_id=$1 AND NOT a.is_fake AND NOT u.is_fake
	`, output.Body.HackathonID).Scan(
		&output.Body.Confirmed, &output.Body.ConfirmedCheckedIn,
		&output.Body.DayOfSignups, &output.Body.StandbyWaiting,
	)
	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to load check-in totals")
	}

	rows, err := h.applicationService.db.Pool.Query(ctx, `
		SELECT u.id,
			COALESCE(NULLIF(trim(concat_ws(' ',a.application->>'firstName',
				a.application->>'lastName')),''),u.name),
			COALESCE(NULLIF(u.preferred_email,''),u.email),
			a.status::text, u.checked_in_at, u.rfid,
			COALESCE(w.signup_source,'preregistered'),w.in_person_joined_at,
			COALESCE((
				SELECT jsonb_agg(jsonb_build_object(
					'redeemableId',r.id,'name',r.name,'amount',ur.amount
				) ORDER BY r.name,r.id)
				FROM user_redemptions ur JOIN redeemables r ON r.id=ur.redeemable_id
				WHERE ur.user_id=u.id AND ur.hackathon_id=a.hackathon_id
					AND r.hackathon_id=a.hackathon_id AND ur.amount>0
			),'[]'::jsonb)
		FROM applications a JOIN users u ON u.id=a.user_id
		LEFT JOIN waitlist w ON w.hackathon_id=a.hackathon_id AND w.user_id=a.user_id
		WHERE a.hackathon_id=$1 AND NOT a.is_fake AND NOT u.is_fake
			AND ($2='' OR strpos(lower(concat_ws(' ',u.name,
				a.application->>'firstName',a.application->>'lastName',
				u.email,u.preferred_email)),lower($2))>0)
			AND ($3='' OR a.status::text=$3)
			AND ($4='' OR ($4='yes' AND u.checked_in_at IS NOT NULL)
				OR ($4='no' AND u.checked_in_at IS NULL))
			AND ($5::uuid IS NULL OR EXISTS (
				SELECT 1 FROM user_redemptions ur JOIN redeemables r ON r.id=ur.redeemable_id
				WHERE ur.user_id=u.id AND ur.hackathon_id=a.hackathon_id
					AND r.hackathon_id=a.hackathon_id
					AND ur.redeemable_id=$5 AND ur.amount>0
			))
		ORDER BY CASE WHEN a.status::text='waitlist_confirmed' THEN 0 ELSE 1 END,
			CASE WHEN w.signup_source='day_of' THEN 1 ELSE 0 END,
			w.in_person_joined_at NULLS LAST,w.created_at NULLS LAST,u.id
	`, output.Body.HackathonID, input.Search, input.Status, input.CheckedIn, redeemableID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to load check-in table")
	}
	defer rows.Close()
	for rows.Next() {
		var row CheckInDashboardRow
		var redemptions []byte
		if err := rows.Scan(&row.UserID, &row.Name, &row.Email, &row.Status,
			&row.CheckedInAt, &row.RFID, &row.SignupSource, &row.StandbyArrival,
			&redemptions); err != nil {
			return nil, huma.Error500InternalServerError("Unable to read check-in table")
		}
		row.Redemptions = json.RawMessage(redemptions)
		output.Body.Rows = append(output.Body.Rows, row)
	}
	if rows.Err() != nil {
		return nil, huma.Error500InternalServerError("Unable to read check-in table")
	}
	return output, nil
}
