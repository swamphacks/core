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

	t.Run("leave is repeatable and rejoin gets a new position", func(t *testing.T) {
		id := newApplicant(t, "rejected", "applicant")
		if err := service.JoinAdmissionWaitlist(ctx, id); err != nil {
			t.Fatal(err)
		}
		exec("UPDATE waitlist SET created_at = now() - interval '1 day' WHERE user_id = $1", id)
		for i := 0; i < 2; i++ {
			if err := service.LeaveAdmissionWaitlist(ctx, id); err != nil {
				t.Fatal(err)
			}
		}
		assertState(t, id, "rejected", "applicant")
		var count int
		if err := pool.QueryRow(ctx, "SELECT count(*) FROM waitlist WHERE user_id=$1", id).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Fatal("leaving retained a queue entry")
		}
		if err := service.JoinAdmissionWaitlist(ctx, id); err != nil {
			t.Fatal(err)
		}
		var recent bool
		if err := pool.QueryRow(ctx, `
            SELECT created_at > now() - interval '1 minute'
            FROM waitlist WHERE user_id=$1
        `, id).Scan(&recent); err != nil {
			t.Fatal(err)
		}
		if !recent {
			t.Fatal("rejoin retained the old queue position")
		}
	})

	t.Run("leave cannot cancel an accepted offer", func(t *testing.T) {
		id := newApplicant(t, "accepted", "applicant")
		if err := service.LeaveAdmissionWaitlist(ctx, id); !errors.Is(err, ErrLeaveWaitlistEligibility) {
			t.Fatalf("expected eligibility error, got %v", err)
		}
		assertState(t, id, "accepted", "applicant")
	})

	t.Run("dispatch preserves order and reserves capacity", func(t *testing.T) {
		first := newApplicant(t, "rejected", "applicant")
		second := newApplicant(t, "rejected", "applicant")
		for _, id := range []uuid.UUID{first, second} {
			exec("UPDATE users SET email=$2 WHERE id=$1", id, id.String()+"@example.test")
			if err := service.JoinAdmissionWaitlist(ctx, id); err != nil {
				t.Fatal(err)
			}
		}
		exec("UPDATE waitlist SET created_at=now()-interval '1 day' WHERE user_id=$1", first)
		exec("UPDATE hackathons SET max_attendees=1 WHERE id=$1", hackathonID)
		now := time.Now().UTC()
		exec(`
            INSERT INTO waitlist_dispatch_policies
              (hackathon_id, enabled, invitations_open_at, online_join_closes_at,
               in_person_opens_at, invitations_close_at)
            VALUES ($1, true, $2, $3, $4, $5)
        `, hackathonID, now.Add(-time.Hour), now.Add(time.Hour),
			now.Add(2*time.Hour), now.Add(3*time.Hour))
		defer exec("DELETE FROM waitlist_dispatch_policies WHERE hackathon_id=$1", hackathonID)

		result, err := service.DispatchAdmissionWaitlist(ctx, now)
		if err != nil {
			t.Fatal(err)
		}
		if result.Offered != 1 {
			t.Fatalf("offered %d, want 1", result.Offered)
		}
		assertState(t, first, "accepted", "applicant")
		assertState(t, second, "waitlisted", "applicant")
		result, err = service.DispatchAdmissionWaitlist(ctx, now.Add(time.Minute))
		if err != nil {
			t.Fatal(err)
		}
		if result.Offered != 0 {
			t.Fatal("dispatcher overbooked the event")
		}

		var count int
		if err := pool.QueryRow(ctx, `
            SELECT count(*) FROM waitlist_invitation_outbox o
            JOIN applications a ON a.id=o.application_id
            WHERE a.hackathon_id=$1
        `, hackathonID).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("outbox contains %d invitations, want 1", count)
		}

		failed, err := service.DeliverAdmissionWaitlistInvitations(
			ctx, func(context.Context, string, string, time.Time) error {
				return errors.New("fake rate limit")
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		if failed.Failed != 1 || failed.Sent != 0 {
			t.Fatalf("unexpected failed-delivery result: %+v", failed)
		}
		exec(`
            UPDATE waitlist_invitation_outbox o SET next_attempt_at=now()
            FROM applications a
            WHERE o.application_id=a.id AND a.hackathon_id=$1
        `, hackathonID)

		calls := 0
		delivered, err := service.DeliverAdmissionWaitlistInvitations(
			ctx, func(_ context.Context, recipient, name string, deadline time.Time) error {
				calls++
				if recipient != first.String()+"@example.test" {
					t.Errorf("unexpected recipient: %s", recipient)
				}
				remaining := time.Until(deadline)
				if remaining < 48*time.Hour-time.Minute || remaining > 48*time.Hour {
					t.Errorf("invalid confirmation window: %v", remaining)
				}
				return nil
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		if delivered.Sent != 1 || calls != 1 {
			t.Fatalf("unexpected delivery: %+v, calls=%d", delivered, calls)
		}

		_, err = service.DeliverAdmissionWaitlistInvitations(
			ctx, func(context.Context, string, string, time.Time) error {
				t.Error("already sent invitation was sent again")
				return nil
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		// Closing invitations must preserve an already sent confirmation window.
		closedSweep, err := service.DispatchAdmissionWaitlist(ctx, now.Add(3*time.Hour))
		if err != nil {
			t.Fatal(err)
		}
		if closedSweep.Expired != 0 || closedSweep.Offered != 0 {
			t.Fatalf("closure changed a sent offer: %+v", closedSweep)
		}
		assertState(t, first, "accepted", "applicant")

		exec(`
            UPDATE application_waitlist_offers o
            SET offered_at=now()-interval '3 days',
                confirmation_deadline=now()-interval '1 day'
            FROM applications a
            WHERE o.application_id=a.id AND a.user_id=$1
        `, first)

		result, err = service.DispatchAdmissionWaitlist(ctx, time.Now().UTC())
		if err != nil {
			t.Fatal(err)
		}
		if result.Expired != 1 || result.Offered != 1 {
			t.Fatalf("expiration did not advance the queue: %+v", result)
		}
		assertState(t, first, "withdrawn", "applicant")
		assertState(t, second, "accepted", "applicant")

		// A pending invitation must never be sent after the delivery window closes.
		exec(`
            UPDATE waitlist_dispatch_policies
            SET invitations_open_at=now()-interval '4 hours',
                online_join_closes_at=now()-interval '3 hours',
                in_person_opens_at=now()-interval '2 hours',
                invitations_close_at=now()-interval '1 hour'
            WHERE hackathon_id=$1
        `, hackathonID)
		closedDelivery, err := service.DeliverAdmissionWaitlistInvitations(
			ctx, func(context.Context, string, string, time.Time) error {
				t.Error("invitation sent after closure")
				return nil
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		if closedDelivery.Sent != 0 {
			t.Fatalf("sent after closure: %+v", closedDelivery)
		}
		assertState(t, second, "withdrawn", "applicant")
		var cancelled bool
		if err := pool.QueryRow(ctx, `
            SELECT o.cancelled_at IS NOT NULL
            FROM waitlist_invitation_outbox o
            JOIN applications a ON a.id=o.application_id
            WHERE a.user_id=$1 AND a.hackathon_id=$2
        `, second, hackathonID).Scan(&cancelled); err != nil {
			t.Fatal(err)
		}
		if !cancelled {
			t.Fatal("closed invitation remained pending")
		}

	})

	t.Run("standby requires staff acceptance and preserves priority", func(t *testing.T) {
		online := newApplicant(t, "rejected", "applicant")
		present := newApplicant(t, "rejected", "applicant")
		dayOf := newApplicant(t, "rejected", "applicant")
		staff := newApplicant(t, "started", "staff")

		for _, id := range []uuid.UUID{online, present, dayOf} {
			if err := service.JoinAdmissionWaitlist(ctx, id); err != nil {
				t.Fatal(err)
			}
		}
		exec("UPDATE waitlist SET signup_source='day_of' WHERE user_id=$1", dayOf)

		now := time.Now().UTC()
		exec(`
			INSERT INTO waitlist_dispatch_policies
				(hackathon_id,enabled,invitations_open_at,online_join_closes_at,
				in_person_opens_at,invitations_close_at)
			VALUES ($1,true,$2,$3,$4,$5)
		`, hackathonID, now.Add(-3*time.Hour), now.Add(-2*time.Hour),
			now.Add(-time.Hour), now.Add(time.Hour))
		defer exec("DELETE FROM waitlist_dispatch_policies WHERE hackathon_id=$1", hackathonID)

		if err := service.RecordInPersonAdmissionWaitlist(ctx, present, online, now); !errors.Is(err, ErrInPersonWaitlistPermission) {
			t.Fatalf("non-staff arrival recording: %v", err)
		}
		if err := service.RecordInPersonAdmissionWaitlist(ctx, present, staff, now.Add(-2*time.Hour)); !errors.Is(err, ErrInPersonWaitlistClosed) {
			t.Fatalf("early arrival recording: %v", err)
		}

		// A day-of signup arrives first; preregistered hackers still take priority.
		if err := service.RecordInPersonAdmissionWaitlist(ctx, dayOf, staff, now.Add(-time.Minute)); err != nil {
			t.Fatal(err)
		}
		if err := service.RecordInPersonAdmissionWaitlist(ctx, present, staff, now); err != nil {
			t.Fatal(err)
		}
		if err := service.RecordInPersonAdmissionWaitlist(ctx, present, staff, now.Add(time.Minute)); err != nil {
			t.Fatal(err)
		}
		var arrival time.Time
		if err := pool.QueryRow(ctx,
			"SELECT in_person_joined_at FROM waitlist WHERE user_id=$1", present,
		).Scan(&arrival); err != nil {
			t.Fatal(err)
		}
		if arrival.Sub(now) > time.Millisecond || now.Sub(arrival) > time.Millisecond {
			t.Fatal("repeat recording changed arrival time")
		}
		assertState(t, present, "waitlist_confirmed", "applicant")
		assertState(t, dayOf, "waitlist_confirmed", "applicant")

		// Concurrent scheduler runs must not automatically admit anyone day-of.
		failures := make(chan error, 8)
		results := make(chan WaitlistDispatchResult, 8)
		var wg sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				result, err := service.DispatchAdmissionWaitlist(ctx, now)
				if err != nil {
					failures <- err
					return
				}
				results <- result
			}()
		}
		wg.Wait()
		close(failures)
		close(results)
		for err := range failures {
			t.Errorf("dispatch: %v", err)
		}
		for result := range results {
			if result.Offered != 0 {
				t.Fatalf("automatic day-of offers: %d", result.Offered)
			}
		}
		assertState(t, present, "waitlist_confirmed", "applicant")

		if err := service.acceptStandby(ctx, hackathonID, present, online, now); !errors.Is(err, ErrInPersonWaitlistPermission) {
			t.Fatalf("non-staff acceptance: %v", err)
		}
		if err := service.acceptStandby(ctx, hackathonID, present, staff, now.Add(-2*time.Hour)); !errors.Is(err, errStandbyClosed) {
			t.Fatalf("acceptance before opening: %v", err)
		}
		if err := service.acceptStandby(ctx, hackathonID, dayOf, staff, now); !errors.Is(err, errStandbyPriority) {
			t.Fatalf("day-of signup bypassed preregistered priority: %v", err)
		}

		// Two staff requests for the same hacker can only accept them once.
		acceptResults := make(chan error, 2)
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				acceptResults <- service.acceptStandby(ctx, hackathonID, present, staff, now)
			}()
		}
		wg.Wait()
		close(acceptResults)
		successes := 0
		for err := range acceptResults {
			if err == nil {
				successes++
				continue
			}
			if !errors.Is(err, errStandbyPriority) {
				t.Errorf("concurrent acceptance: %v", err)
			}
		}
		if successes != 1 {
			t.Fatalf("successful acceptances=%d; want 1", successes)
		}
		assertState(t, present, "accepted", "applicant")
		assertState(t, dayOf, "waitlist_confirmed", "applicant")
		assertState(t, online, "waitlisted", "applicant")

		var deadlineExists bool
		if err := pool.QueryRow(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM application_waitlist_offers o
				JOIN applications a ON a.id=o.application_id
				WHERE a.user_id=$1 AND a.hackathon_id=$2
					AND o.confirmation_deadline>o.offered_at
			)
		`, present, hackathonID).Scan(&deadlineExists); err != nil {
			t.Fatal(err)
		}
		if !deadlineExists {
			t.Fatal("manual acceptance did not create confirmation deadline")
		}

		if err := service.acceptStandby(ctx, hackathonID, dayOf, staff, now); err != nil {
			t.Fatal(err)
		}
		assertState(t, dayOf, "accepted", "applicant")
		if err := service.acceptStandby(ctx, hackathonID, online, staff, now); !errors.Is(err, errStandbyEmpty) {
			t.Fatalf("absent hacker was accepted: %v", err)
		}

		var source string
		if err := pool.QueryRow(ctx,
			"SELECT signup_source FROM waitlist WHERE user_id=$1", dayOf,
		).Scan(&source); err != nil {
			t.Fatal(err)
		}
		if source != "day_of" {
			t.Fatal("acceptance lost day-of signup provenance")
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

	t.Run("waitlist offer overrides expired shared deadline", func(t *testing.T) {
		exec("UPDATE hackathons SET rsvp_deadline = now() - interval '1 day' WHERE id = $1", hackathonID)
		defer exec("UPDATE hackathons SET rsvp_deadline = now() + interval '1 day' WHERE id = $1", hackathonID)

		id := newApplicant(t, "waitlisted", "applicant")
		exec("UPDATE applications SET status = 'accepted' WHERE user_id = $1 AND hackathon_id = $2", id, hackathonID)

		var applicationID uuid.UUID
		var offered, deadline time.Time
		err := pool.QueryRow(ctx, `
            SELECT a.id, o.offered_at, o.confirmation_deadline
            FROM applications a
            JOIN application_waitlist_offers o ON o.application_id = a.id
            WHERE a.user_id = $1 AND a.hackathon_id = $2
        `, id, hackathonID).Scan(&applicationID, &offered, &deadline)
		if err != nil {
			t.Fatal(err)
		}
		if deadline.Sub(offered) != 48*time.Hour {
			t.Fatal("offer does not provide exactly 48 hours")
		}

		exec("UPDATE applications SET status = 'accepted' WHERE id = $1", applicationID)
		saved, err := service.GetApplicationConfirmationDeadline(ctx, applicationID)
		if err != nil || saved == nil || !saved.Equal(deadline) {
			t.Fatalf("repeated acceptance changed deadline: %v", err)
		}
		if err := service.ConfirmAttendance(ctx, id); err != nil {
			t.Fatal(err)
		}
		assertState(t, id, "confirmed", "attendee")
	})

	t.Run("expired waitlist offer cannot confirm", func(t *testing.T) {
		id := newApplicant(t, "waitlisted", "applicant")
		exec("UPDATE applications SET status = 'accepted' WHERE user_id = $1 AND hackathon_id = $2", id, hackathonID)
		exec(`
            UPDATE application_waitlist_offers
            SET offered_at = now() - interval '3 days',
                confirmation_deadline = now() - interval '1 day'
            WHERE application_id = (
                SELECT id FROM applications WHERE user_id = $1 AND hackathon_id = $2
            )
        `, id, hackathonID)
		if err := service.ConfirmAttendance(ctx, id); !errors.Is(err, ErrConfirmAttendance) {
			t.Fatalf("expected expired confirmation error, got %v", err)
		}
		assertState(t, id, "accepted", "applicant")
	})

	t.Run("regular acceptance keeps shared deadline", func(t *testing.T) {
		exec("UPDATE hackathons SET rsvp_deadline = now() - interval '1 day' WHERE id = $1", hackathonID)
		defer exec("UPDATE hackathons SET rsvp_deadline = now() + interval '1 day' WHERE id = $1", hackathonID)

		id := newApplicant(t, "under_review", "applicant")
		exec("UPDATE applications SET status = 'accepted' WHERE user_id = $1 AND hackathon_id = $2", id, hackathonID)
		var count int
		if err := pool.QueryRow(ctx, `
            SELECT count(*) FROM application_waitlist_offers o
            JOIN applications a ON a.id = o.application_id
            WHERE a.user_id = $1
        `, id).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Fatal("regular acceptance unexpectedly received a waitlist offer")
		}
		if err := service.ConfirmAttendance(ctx, id); !errors.Is(err, ErrConfirmAttendance) {
			t.Fatalf("expected shared deadline rejection, got %v", err)
		}
		assertState(t, id, "accepted", "applicant")
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
