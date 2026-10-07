package application

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/config"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

type visitorTestStorage struct {
	calls int
	err   error
}

func (s *visitorTestStorage) Store(
	_ context.Context, _, _ string, _ []byte, _ *string,
) error {
	s.calls++
	return s.err
}
func (*visitorTestStorage) Retrieve(context.Context, string, string) ([]byte, error) {
	return nil, nil
}
func (*visitorTestStorage) Delete(context.Context, string, string) error {
	return nil
}
func (*visitorTestStorage) Close() error { return nil }

func TestVisitorWaitlistIntegration(t *testing.T) {
	t.Run("advance", func(t *testing.T) {
		runVisitorWaitlistIntegration(t, false)
	})
	t.Run("day-of", func(t *testing.T) {
		runVisitorWaitlistIntegration(t, true)
	})
}

func runVisitorWaitlistIntegration(t *testing.T, dayOf bool) {
	const localURL = "postgres://postgres:postgres@127.0.0.1:55432/waitlist_test?sslmode=disable"
	url := os.Getenv("WAITLIST_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set WAITLIST_TEST_DATABASE_URL to run local integration tests")
	}
	if url != localURL {
		t.Fatal("requires the dedicated local waitlist_test database")
	}
	if !dayOf && !admissionWaitlistIsOpen(time.Now()) {
		t.Skip("visitor registration window has closed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, query, args...); err != nil {
			t.Fatal(err)
		}
	}
	var active int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM hackathons WHERE is_active").Scan(&active); err != nil {
		t.Fatal(err)
	}
	if active != 0 {
		t.Fatal("test database must have no active hackathon")
	}

	eventID := "visitor-test-" + uuid.NewString()
	exec(`INSERT INTO hackathons (
		id, name, application_open, application_close,
		start_time, end_time, rsvp_deadline, is_active
	) VALUES (
		$1, 'Visitor integration test', now()-interval '2 days',
		now()-interval '1 day', now()+interval '2 days',
		now()+interval '3 days', now()+interval '1 day', true
	)`, eventID)
	defer func() {
		_, _ = pool.Exec(context.Background(),
			"DELETE FROM hackathons WHERE id=$1", eventID)
	}()

	wantMarker, expectedStatus := "visitor-waitlist", "waitlisted"
	if dayOf {
		wantMarker, expectedStatus = "day-of", "waitlist_confirmed"
		exec(`
            INSERT INTO waitlist_dispatch_policies
                (hackathon_id,enabled,invitations_open_at,online_join_closes_at,
                in_person_opens_at,invitations_close_at)
            VALUES ($1,true,now()-interval '3 hours',now()-interval '2 hours',
                now()-interval '1 hour',now()+interval '1 hour')
        `, eventID)
	}

	db := &database.DB{Pool: pool, Query: sqlc.New(pool)}
	uploadErr := errors.New("test upload failed")
	data := ApplicationSubmissionFields{PreferredEmail: "visitor@example.com"}

	for _, tc := range []struct {
		name, role, status string
		failUpload         bool
	}{
		{"new visitor", "visitor", "", false},
		{"started application", "visitor", "started", false},
		{"new visitor upload failure", "visitor", "", true},
		{"started application upload failure", "visitor", "started", true},
		{"staff is ineligible", "staff", "", false},
		{"completed application is ineligible", "visitor", "rejected", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.New()
			exec("INSERT INTO users (id,name,role) VALUES ($1,'Test visitor',$2::role)",
				id, tc.role)
			defer func() {
				_, _ = pool.Exec(context.Background(),
					"DELETE FROM users WHERE id=$1", id)
			}()
			if tc.status != "" {
				exec(`INSERT INTO applications (user_id,hackathon_id,status)
					VALUES ($1,$2,$3::application_status)`, id, eventID, tc.status)
			}

			store := &visitorTestStorage{}
			if tc.failUpload {
				store.err = uploadErr
			}
			service := &ApplicationService{
				db: db, txm: database.NewTransactionManager(db),
				storage: store,
				buckets: &config.CoreBuckets{ApplicationResumes: "test-resumes"},
				logger:  zerolog.Nop(),
			}
			stamp, err := service.SubmitVisitorWaitlist(
				ctx, data, []byte("%PDF-test"), id)
			ineligible := tc.role != "visitor" ||
				(tc.status != "" && tc.status != "started")
			switch {
			case ineligible:
				if !errors.Is(err, ErrVisitorWaitlistEligibility) || store.calls != 0 {
					t.Fatalf("expected eligibility rejection before upload: %v", err)
				}
			case tc.failUpload:
				if !errors.Is(err, uploadErr) || store.calls != 1 {
					t.Fatalf("expected upload failure: %v", err)
				}
			default:
				if err != nil || stamp == nil {
					t.Fatalf("submission failed: %v", err)
				}
				var saved time.Time
				var marker, email string
				if err := pool.QueryRow(ctx, `
					SELECT a.submitted_at, a.application->>'registrationType',
						u.preferred_email
					FROM applications a JOIN users u ON u.id=a.user_id
					WHERE a.user_id=$1 AND a.hackathon_id=$2
				`, id, eventID).Scan(&saved, &marker, &email); err != nil {
					t.Fatal(err)
				}
				if marker != wantMarker || email != data.PreferredEmail {
					t.Fatal("registration metadata or preferred email not saved")
				}
				var source string
				var arrival *time.Time
				if err := pool.QueryRow(ctx, `
                    SELECT signup_source,in_person_joined_at
                    FROM waitlist WHERE user_id=$1 AND hackathon_id=$2
                `, id, eventID).Scan(&source, &arrival); err != nil {
					t.Fatal(err)
				}
				if dayOf {
					if source != "day_of" || arrival == nil || !arrival.Equal(saved) {
						t.Fatal("day-of source or arrival timestamp missing")
					}
				} else if source != "preregistered" || arrival != nil {
					t.Fatal("advance registration was recorded as day-of arrival")
				}

				again, err := service.SubmitVisitorWaitlist(
					ctx, data, []byte("%PDF-retry"), id)
				if err != nil || again == nil || !again.Equal(saved) || store.calls != 1 {
					t.Fatalf("retry changed timestamp or repeated upload: %v", err)
				}
			}

			var role string
			var apps, queued int
			if err := pool.QueryRow(ctx,
				"SELECT role::text FROM users WHERE id=$1", id).Scan(&role); err != nil {
				t.Fatal(err)
			}
			if err := pool.QueryRow(ctx,
				"SELECT count(*) FROM applications WHERE user_id=$1 AND hackathon_id=$2",
				id, eventID).Scan(&apps); err != nil {
				t.Fatal(err)
			}
			if err := pool.QueryRow(ctx,
				"SELECT count(*) FROM waitlist WHERE user_id=$1 AND hackathon_id=$2",
				id, eventID).Scan(&queued); err != nil {
				t.Fatal(err)
			}

			wantRole, wantApps, wantQueued := tc.role, 0, 0
			if tc.status != "" {
				wantApps = 1
			}
			if !ineligible && !tc.failUpload {
				wantRole, wantApps, wantQueued = "applicant", 1, 1
			}
			if role != wantRole || apps != wantApps || queued != wantQueued {
				t.Fatalf("role/apps/queue = %s/%d/%d; want %s/%d/%d",
					role, apps, queued, wantRole, wantApps, wantQueued)
			}
			if apps == 1 {
				var status string
				if err := pool.QueryRow(ctx,
					"SELECT status::text FROM applications WHERE user_id=$1 AND hackathon_id=$2",
					id, eventID).Scan(&status); err != nil {
					t.Fatal(err)
				}
				wantStatus := tc.status
				if !ineligible && !tc.failUpload {
					wantStatus = expectedStatus
				}
				if status != wantStatus {
					t.Fatalf("status = %s; want %s", status, wantStatus)
				}
			}
		})
	}
}
