package application

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

func TestAdmissionWaitlistIntegration(t *testing.T) {
	const localURL = "postgres://postgres:postgres@127.0.0.1:55432/waitlist_test?sslmode=disable"
	url := os.Getenv("WAITLIST_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set WAITLIST_TEST_DATABASE_URL to run local integration tests")
	}
	if url != localURL {
		t.Fatal("integration tests require the dedicated local waitlist_test database")
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

	var activeCount int
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM hackathons WHERE is_active",
	).Scan(&activeCount); err != nil {
		t.Fatal(err)
	}
	if activeCount != 0 {
		t.Fatal("test database already has an active hackathon; refusing to modify it")
	}

	hackathonID := "waitlist-test-" + uuid.NewString()
	exec(`INSERT INTO hackathons (
		id, name, application_open, application_close,
		start_time, end_time, rsvp_deadline, is_active
	) VALUES (
		$1, 'Waitlist integration test', now(), now() + interval '1 day',
		now() + interval '2 days', now() + interval '3 days',
		now() + interval '1 day', true
	)`, hackathonID)

	defer func() {
		_, _ = pool.Exec(context.Background(),
			"DELETE FROM hackathons WHERE id = $1", hackathonID)
	}()

	db := &database.DB{Pool: pool, Query: sqlc.New(pool)}
	service := &ApplicationService{
		db: db, txm: database.NewTransactionManager(db),
		logger: zerolog.New(os.Stderr),
	}

	newApplicant := func(t *testing.T, status, role string) uuid.UUID {
		t.Helper()
		id := uuid.New()
		exec(`INSERT INTO users (id, name, role)
			VALUES ($1, 'Integration test', $2::role)`, id, role)
		exec(`INSERT INTO applications (user_id, hackathon_id, status)
			VALUES ($1, $2, $3::application_status)`,
			id, hackathonID, status)
		t.Cleanup(func() {
			if _, err := pool.Exec(context.Background(),
				"DELETE FROM users WHERE id = $1", id); err != nil {
				t.Errorf("cleanup: %v", err)
			}
		})
		return id
	}

	assertState := func(t *testing.T, id uuid.UUID, wantStatus, wantRole string) {
		t.Helper()
		var status, role string
		err := pool.QueryRow(ctx, `
			SELECT a.status::text, u.role::text
			FROM applications a JOIN users u ON u.id = a.user_id
			WHERE a.user_id = $1 AND a.hackathon_id = $2
		`, id, hackathonID).Scan(&status, &role)
		if err != nil {
			t.Fatal(err)
		}
		if status != wantStatus || role != wantRole {
			t.Fatalf("got status=%s role=%s; want status=%s role=%s",
				status, role, wantStatus, wantRole)
		}
	}

	t.Run("duplicate join preserves one entry and timestamp", func(t *testing.T) {
		id := newApplicant(t, "rejected", "applicant")
		if err := service.JoinAdmissionWaitlist(ctx, id); err != nil {
			t.Fatal(err)
		}
		var first time.Time
		if err := pool.QueryRow(ctx,
			"SELECT created_at FROM waitlist WHERE user_id=$1", id,
		).Scan(&first); err != nil {
			t.Fatal(err)
		}
		if err := service.JoinAdmissionWaitlist(ctx, id); err != nil {
			t.Fatal(err)
		}
		var count int
		var timestamp time.Time
		if err := pool.QueryRow(ctx, `
			SELECT count(*), min(created_at)
			FROM waitlist WHERE user_id=$1
		`, id).Scan(&count, &timestamp); err != nil {
			t.Fatal(err)
		}
		if count != 1 || !timestamp.Equal(first) {
			t.Fatal("duplicate join changed the entry count or join timestamp")
		}
		assertState(t, id, "waitlisted", "applicant")
		application, err := service.GetApplicationByUserId(ctx, id)
		if err != nil || application == nil {
			t.Fatalf("portal lookup failed: %v", err)
		}
		if string(application.Status) != "waitlisted" {
			t.Fatal("portal lookup did not return the saved waitlist status")
		}
	})

	t.Run("join rejects ineligible statuses", func(t *testing.T) {
		for _, status := range []string{
			"started", "submitted", "under_review",
			"accepted", "confirmed", "withdrawn",
		} {
			t.Run(status, func(t *testing.T) {
				id := newApplicant(t, status, "applicant")
				if err := service.JoinAdmissionWaitlist(ctx, id); !errors.Is(err, ErrWaitlistEligibility) {
					t.Fatalf("expected eligibility error, got %v", err)
				}
				assertState(t, id, status, "applicant")
				var count int
				if err := pool.QueryRow(ctx,
					"SELECT count(*) FROM waitlist WHERE user_id=$1", id,
				).Scan(&count); err != nil {
					t.Fatal(err)
				}
				if count != 0 {
					t.Fatal("ineligible applicant was added to waitlist")
				}
			})
		}
	})

	t.Run("accepted and confirmed can withdraw repeatedly", func(t *testing.T) {
		for _, status := range []string{"accepted", "confirmed"} {
			t.Run(status, func(t *testing.T) {
				role := "applicant"
				if status == "confirmed" {
					role = "attendee"
				}
				id := newApplicant(t, status, role)
				for attempt := 0; attempt < 2; attempt++ {
					if err := service.WithdrawApplication(ctx, id); err != nil {
						t.Fatal(err)
					}
					assertState(t, id, "withdrawn", "applicant")
				}
			})
		}
	})

	t.Run("withdraw rejects ineligible statuses", func(t *testing.T) {
		for _, status := range []string{
			"started", "submitted", "under_review", "rejected", "waitlisted",
		} {
			t.Run(status, func(t *testing.T) {
				id := newApplicant(t, status, "applicant")
				if err := service.WithdrawApplication(ctx, id); !errors.Is(err, ErrWithdrawalEligibility) {
					t.Fatalf("expected eligibility error, got %v", err)
				}
				assertState(t, id, status, "applicant")
			})
		}
	})

	t.Run("confirmation cannot restore withdrawn attendance", func(t *testing.T) {
		for attempt := 0; attempt < 10; attempt++ {
			id := newApplicant(t, "accepted", "applicant")
			start := make(chan struct{})
			var wg sync.WaitGroup
			var confirmErr, withdrawErr error
			wg.Add(2)
			go func() {
				defer wg.Done()
				<-start
				confirmErr = service.ConfirmAttendance(ctx, id)
			}()
			go func() {
				defer wg.Done()
				<-start
				withdrawErr = service.WithdrawApplication(ctx, id)
			}()
			close(start)
			wg.Wait()
			if withdrawErr != nil {
				t.Fatal(withdrawErr)
			}
			if confirmErr != nil && !errors.Is(confirmErr, ErrConfirmAttendance) {
				t.Fatal(confirmErr)
			}
			assertState(t, id, "withdrawn", "applicant")
		}
	})
}
