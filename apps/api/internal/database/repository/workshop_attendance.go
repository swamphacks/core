package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrWorkshopAttendanceNotFound = errors.New("workshop or hacker not found")
var ErrWorkshopAttendanceIneligible = errors.New("hacker must be confirmed and checked in")

type WorkshopAttendance struct {
	Name            string `json:"name"`
	WorkshopTitle   string `json:"workshopTitle"`
	AlreadyRecorded bool   `json:"alreadyRecorded"`
}

type AttendanceWorkshop struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	Type           string    `json:"type"`
	StartTime      time.Time `json:"startTime"`
	CheckedInCount int64     `json:"checkedInCount"`
}

func (r *WorkshopsRepository) AttendanceWorkshops(ctx context.Context) ([]AttendanceWorkshop, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT w.id,w.title,w.type,w.start_time,count(DISTINCT a.user_id)
FROM workshops w
LEFT JOIN nfc_tags_workshops a ON a.workshop_id=w.id
GROUP BY w.id,w.title,w.type,w.start_time
ORDER BY w.start_time,w.title,w.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []AttendanceWorkshop{}
	for rows.Next() {
		var w AttendanceWorkshop
		if err := rows.Scan(&w.ID, &w.Title, &w.Type, &w.StartTime, &w.CheckedInCount); err != nil {
			return nil, err
		}
		result = append(result, w)
	}
	return result, rows.Err()
}

// QR and NFC attendance share the same (user_id, workshop_id) uniqueness rule.
// Registration/RSVP counts are deliberately not changed by an attendance scan.
func (r *WorkshopsRepository) RecordAttendance(ctx context.Context, workshopID, userID uuid.UUID) (*WorkshopAttendance, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	result := &WorkshopAttendance{}
	err = tx.QueryRow(ctx, `SELECT title FROM workshops WHERE id=$1 FOR KEY SHARE`, workshopID).Scan(&result.WorkshopTitle)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrWorkshopAttendanceNotFound
	}
	if err != nil {
		return nil, err
	}
	var eligible bool
	err = tx.QueryRow(ctx, `SELECT COALESCE(NULLIF(trim(concat_ws(' ',a.application->>'firstName',a.application->>'lastName')),''),u.name),
 a.status::text='confirmed' AND NOT a.is_fake AND NOT u.is_fake
 AND (u.checked_in_at IS NOT NULL OR EXISTS(SELECT 1 FROM nfc_tags_user n WHERE n.user_id=u.id))
 FROM users u JOIN applications a ON a.user_id=u.id JOIN hackathons h ON h.id=a.hackathon_id AND h.is_active
 WHERE u.id=$1 FOR UPDATE OF u,a`, userID).Scan(&result.Name, &eligible)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrWorkshopAttendanceNotFound
	}
	if err != nil {
		return nil, err
	}
	if !eligible {
		return nil, ErrWorkshopAttendanceIneligible
	}
	tag, err := tx.Exec(ctx, `INSERT INTO nfc_tags_workshops (user_id,workshop_id) VALUES ($1,$2) ON CONFLICT (user_id,workshop_id) DO NOTHING`, userID, workshopID)
	if err != nil {
		return nil, err
	}
	result.AlreadyRecorded = tag.RowsAffected() == 0
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return result, nil
}
