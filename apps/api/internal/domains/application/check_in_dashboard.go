package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	ApplicationID  uuid.UUID       `json:"applicationId"`
	UserID         uuid.UUID       `json:"userId"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	Status         string          `json:"status"`
	CheckedInAt    *time.Time      `json:"checkedInAt"`
	NFCLinked      bool            `json:"nfcLinked"`
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

	huma.Register(group, huma.Operation{
		OperationID: "manual-check-in",
		Method:      http.MethodPost,
		Path:        "/check-in",
		Summary:     "Check In Hacker Without NFC",
		Description: "Staff route for checking in a confirmed hacker for the active hackathon without linking an NFC nametag. The nametag can still be linked afterwards from the mobile app.",
		Tags:        []string{"Application"},
		Middlewares: huma.Middlewares{
			mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma,
		},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		Errors:        []int{401, 403, 404, 409, 500},
		DefaultStatus: http.StatusOK,
	}, h.handleManualCheckIn)

	huma.Register(group, huma.Operation{
		OperationID: "undo-manual-check-in",
		Method:      http.MethodPost,
		Path:        "/check-in/undo",
		Summary:     "Undo Check-In Without NFC",
		Description: "Staff route for clearing a check-in made without an NFC nametag, for example a mis-tap. Check-ins with a linked NFC nametag cannot be undone here.",
		Tags:        []string{"Application"},
		Middlewares: huma.Middlewares{
			mw.Auth.RequireAuthHuma, mw.Auth.RequireStaffHuma,
		},
		Parameters:    []*huma.Param{cookie.SessionCookieHumaParam},
		Errors:        []int{401, 403, 404, 409, 500},
		DefaultStatus: http.StatusOK,
	}, h.handleUndoManualCheckIn)
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
        count(*) FILTER (WHERE a.status::text='confirmed'
            AND (ci.user_id IS NOT NULL OR u.checked_in_at IS NOT NULL)),
        count(*) FILTER (WHERE w.signup_source='day_of'),
        count(*) FILTER (WHERE a.status::text='waitlist_confirmed')
    FROM applications a
    JOIN users u ON u.id=a.user_id
    LEFT JOIN LATERAL (
        SELECT ntu.user_id, ntu.created_at
        FROM nfc_tags_user ntu
        WHERE ntu.user_id=u.id
        ORDER BY ntu.created_at ASC NULLS LAST
        LIMIT 1
    ) ci ON true
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
    SELECT a.id, u.id,
        COALESCE(NULLIF(trim(concat_ws(' ',a.application->>'firstName',
            a.application->>'lastName')),''),u.name),
        COALESCE(NULLIF(u.preferred_email,''),u.email),
        a.status::text, LEAST(ci.created_at, u.checked_in_at),
        ci.user_id IS NOT NULL, u.rfid,
        COALESCE(w.signup_source,'preregistered'),w.in_person_joined_at,
        COALESCE((
            SELECT jsonb_agg(jsonb_build_object(
                'redeemableId',x.id,'name',x.name,'amount',x.amount
            ) ORDER BY x.name,x.id)
            FROM (
                SELECT r.id, r.name, count(*) AS amount
                FROM nfc_tags_redeemables ntr
                JOIN nfc_tags_user ntu ON ntu.tag_id=ntr.tag_id
                JOIN redeemables r ON r.id=ntr.redeemable_id
                WHERE ntu.user_id=u.id AND r.hackathon_id=a.hackathon_id
                GROUP BY r.id, r.name
            ) x
        ),'[]'::jsonb)
    FROM applications a
    JOIN users u ON u.id=a.user_id
    LEFT JOIN LATERAL (
        SELECT ntu.user_id, ntu.created_at
        FROM nfc_tags_user ntu
        WHERE ntu.user_id=u.id
        ORDER BY ntu.created_at ASC NULLS LAST
        LIMIT 1
    ) ci ON true
    LEFT JOIN waitlist w ON w.hackathon_id=a.hackathon_id AND w.user_id=a.user_id
    WHERE a.hackathon_id=$1 AND NOT a.is_fake AND NOT u.is_fake
        AND ($2='' OR strpos(lower(concat_ws(' ',u.name,
            a.application->>'firstName',a.application->>'lastName',
            u.email,u.preferred_email)),lower($2))>0)
        AND ($3='' OR a.status::text=$3)
        AND ($4='' OR ($4='yes' AND (ci.user_id IS NOT NULL OR u.checked_in_at IS NOT NULL))
            OR ($4='no' AND ci.user_id IS NULL AND u.checked_in_at IS NULL))
        AND ($5::uuid IS NULL OR EXISTS (
            SELECT 1 FROM nfc_tags_redeemables ntr
            JOIN nfc_tags_user ntu ON ntu.tag_id=ntr.tag_id
            JOIN redeemables r ON r.id=ntr.redeemable_id
            WHERE ntu.user_id=u.id AND r.hackathon_id=a.hackathon_id
                AND ntr.redeemable_id=$5
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
		if err := rows.Scan(
			&row.ApplicationID, &row.UserID, &row.Name, &row.Email, &row.Status,
			&row.CheckedInAt, &row.NFCLinked, &row.RFID, &row.SignupSource, &row.StandbyArrival,
			&redemptions,
		); err != nil {
			fmt.Print(err)
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

type ManualCheckInInput struct {
	Body struct {
		UserID uuid.UUID `json:"userId" required:"true"`
	}
}

type ManualCheckInOutput struct {
	Body struct {
		CheckedInAt time.Time `json:"checkedInAt"`
	}
}

// A hacker counts as checked in once either an NFC nametag is linked
// (nfc_tags_user) or staff checked them in here (users.checked_in_at). The
// update is conditional so a double-tap or two staff scanning the same
// hacker cannot both succeed.
func (h *handler) handleManualCheckIn(ctx context.Context, input *ManualCheckInInput) (*ManualCheckInOutput, error) {
	output := &ManualCheckInOutput{}
	err := h.applicationService.db.Pool.QueryRow(ctx, `
    UPDATE users u SET checked_in_at=now(), updated_at=now()
    WHERE u.id=$1 AND u.checked_in_at IS NULL
        AND NOT EXISTS (SELECT 1 FROM nfc_tags_user ntu WHERE ntu.user_id=u.id)
        AND EXISTS (
            SELECT 1 FROM applications a
            JOIN hackathons hk ON hk.id=a.hackathon_id AND hk.is_active
            WHERE a.user_id=u.id AND a.status::text='confirmed'
        )
    RETURNING u.checked_in_at
`, input.Body.UserID).Scan(&output.Body.CheckedInAt)
	if err == nil {
		h.logger.Info().Str("type", "checkin").Str("user_id", input.Body.UserID.String()).Msg("User checked in without NFC")
		return output, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		h.logger.Err(err).Msg("manual check-in failed")
		return nil, huma.Error500InternalServerError("Unable to check in hacker")
	}

	// Nothing was updated; work out why so staff get a useful message.
	var status string
	var checkedIn bool
	err = h.applicationService.db.Pool.QueryRow(ctx, `
    SELECT a.status::text,
        u.checked_in_at IS NOT NULL
            OR EXISTS (SELECT 1 FROM nfc_tags_user ntu WHERE ntu.user_id=u.id)
    FROM users u
    JOIN applications a ON a.user_id=u.id
    JOIN hackathons hk ON hk.id=a.hackathon_id AND hk.is_active
    WHERE u.id=$1
`, input.Body.UserID).Scan(&status, &checkedIn)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, huma.Error404NotFound("No application for the active hackathon")
	}
	if err != nil {
		h.logger.Err(err).Msg("manual check-in lookup failed")
		return nil, huma.Error500InternalServerError("Unable to check in hacker")
	}
	if checkedIn {
		return nil, huma.Error409Conflict("Hacker is already checked in")
	}
	return nil, huma.Error409Conflict(fmt.Sprintf(
		"Only confirmed hackers can be checked in (status: %s)",
		strings.ReplaceAll(status, "_", " "),
	))
}

type UndoManualCheckInOutput struct {
	Body struct {
		UserID uuid.UUID `json:"userId"`
	}
}

// Only clears check-ins made without NFC. A hacker with a linked nametag would
// still count as checked in, so undoing their timestamp here would change
// nothing staff can see.
func (h *handler) handleUndoManualCheckIn(ctx context.Context, input *ManualCheckInInput) (*UndoManualCheckInOutput, error) {
	output := &UndoManualCheckInOutput{}
	err := h.applicationService.db.Pool.QueryRow(ctx, `
    UPDATE users u SET checked_in_at=NULL, updated_at=now()
    WHERE u.id=$1 AND u.checked_in_at IS NOT NULL
        AND NOT EXISTS (SELECT 1 FROM nfc_tags_user ntu WHERE ntu.user_id=u.id)
    RETURNING u.id
`, input.Body.UserID).Scan(&output.Body.UserID)
	if err == nil {
		h.logger.Info().Str("type", "checkin").Str("user_id", input.Body.UserID.String()).Msg("Check-in without NFC undone")
		return output, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		h.logger.Err(err).Msg("undo check-in failed")
		return nil, huma.Error500InternalServerError("Unable to undo check-in")
	}

	// Nothing was cleared; work out why so staff get a useful message.
	var nfcLinked bool
	err = h.applicationService.db.Pool.QueryRow(ctx, `
    SELECT EXISTS (SELECT 1 FROM nfc_tags_user ntu WHERE ntu.user_id=u.id)
    FROM users u WHERE u.id=$1
`, input.Body.UserID).Scan(&nfcLinked)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, huma.Error404NotFound("Hacker not found")
	}
	if err != nil {
		h.logger.Err(err).Msg("undo check-in lookup failed")
		return nil, huma.Error500InternalServerError("Unable to undo check-in")
	}
	if nfcLinked {
		return nil, huma.Error409Conflict("Hacker checked in with an NFC nametag, which can't be undone from the web")
	}
	return nil, huma.Error409Conflict("Hacker is not checked in")
}
