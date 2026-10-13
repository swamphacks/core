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
	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

func TestConfirmedQRIntegration(t *testing.T) {
	const localURL = "postgres://postgres:postgres@127.0.0.1:55432/waitlist_test?sslmode=disable"
	url := os.Getenv("WAITLIST_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set WAITLIST_TEST_DATABASE_URL")
	}
	if url != localURL {
		t.Fatal("requires dedicated local test database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
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
	var active, existing int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM hackathons WHERE is_active").Scan(&active); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM hackathons WHERE id='xii'").Scan(&existing); err != nil {
		t.Fatal(err)
	}
	if active != 0 || existing != 0 {
		t.Fatal("local test requires no active event and no existing xii fixture")
	}
	exec(`INSERT INTO hackathons (
		id,name,application_open,application_close,start_time,end_time,rsvp_deadline,is_active
	) VALUES (
		'xii','QR email integration test',now()-interval '2 days',
		now()-interval '1 day',now()+interval '1 day',now()+interval '3 days',
		now()+interval '12 hours',true
	)`)
	defer func() {
		if _, err := pool.Exec(context.Background(), "DELETE FROM hackathons WHERE id='xii'"); err != nil {
			t.Errorf("fixture cleanup failed: %v", err)
		}
	}()

	db := &database.DB{Pool: pool, Query: sqlc.New(pool)}
	service := &ApplicationService{
		db: db, txm: database.NewTransactionManager(db), logger: zerolog.Nop(),
	}

	t.Run("confirmation queues once and retries with application email fallback", func(t *testing.T) {
		id := uuid.New()
		exec("INSERT INTO users (id,name,role) VALUES ($1,'Preview QR','applicant')", id)
		defer func() {
			if _, err := pool.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id); err != nil {
				t.Errorf("fixture cleanup failed: %v", err)
			}
		}()
		exec(`INSERT INTO applications (user_id,hackathon_id,status,application)
			VALUES ($1,'xii','accepted',
			'{"firstName":"Preview","preferredEmail":"preview@example.test"}')`, id)
		if err := service.ConfirmAttendance(ctx, id); err != nil {
			t.Fatal(err)
		}
		// A repeated status update must not enqueue a second message.
		exec("UPDATE applications SET status='confirmed' WHERE user_id=$1 AND hackathon_id='xii'", id)
		var queued int
		if err := pool.QueryRow(ctx,
			"SELECT count(*) FROM confirmed_qr_email_outbox WHERE user_id=$1", id,
		).Scan(&queued); err != nil || queued != 1 {
			t.Fatalf("queue count=%d, error=%v", queued, err)
		}

		calls := 0
		failure := func(_ context.Context, got uuid.UUID, email, name string) (string, error) {
			calls++
			if got != id || email != "preview@example.test" || name != "Preview" {
				t.Fatalf("wrong personalized recipient: %s %s %s", got, email, name)
			}
			return "", errors.New("simulated SES failure")
		}
		if err := service.DeliverConfirmedQREmails(ctx, failure); err != nil {
			t.Fatal(err)
		}
		var attempts int
		var sent *time.Time
		var retry time.Time
		if err := pool.QueryRow(ctx, `
			SELECT attempts,sent_at,next_attempt_at
			FROM confirmed_qr_email_outbox WHERE user_id=$1
		`, id).Scan(&attempts, &sent, &retry); err != nil {
			t.Fatal(err)
		}
		if calls != 1 || attempts != 1 || sent != nil || !retry.After(time.Now()) {
			t.Fatal("failed message was not retained for delayed retry")
		}
		if err := service.DeliverConfirmedQREmails(ctx, failure); err != nil {
			t.Fatal(err)
		}
		if calls != 1 {
			t.Fatal("retried before scheduled retry time")
		}
		exec("UPDATE confirmed_qr_email_outbox SET next_attempt_at=now() WHERE user_id=$1", id)
		success := func(_ context.Context, got uuid.UUID, email, name string) (string, error) {
			calls++
			if got != id || email != "preview@example.test" || name != "Preview" {
				t.Fatal("retry changed recipient")
			}
			return "test-message-id", nil
		}
		if err := service.DeliverConfirmedQREmails(ctx, success); err != nil {
			t.Fatal(err)
		}
		if err := service.DeliverConfirmedQREmails(ctx, success); err != nil {
			t.Fatal(err)
		}
		if calls != 2 {
			t.Fatal("successful message was sent again")
		}
	})

	t.Run("withdrawal cancels pending email", func(t *testing.T) {
		id := uuid.New()
		exec("INSERT INTO users (id,name,role) VALUES ($1,'Withdraw test','applicant')", id)
		defer func() {
			if _, err := pool.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id); err != nil {
				t.Errorf("fixture cleanup failed: %v", err)
			}
		}()
		exec("INSERT INTO applications (user_id,hackathon_id,status) VALUES ($1,'xii','accepted')", id)
		if err := service.ConfirmAttendance(ctx, id); err != nil {
			t.Fatal(err)
		}
		if err := service.WithdrawApplication(ctx, id); err != nil {
			t.Fatal(err)
		}
		err := service.DeliverConfirmedQREmails(ctx,
			func(context.Context, uuid.UUID, string, string) (string, error) {
				t.Fatal("sent email to withdrawn hacker")
				return "", nil
			})
		if err != nil {
			t.Fatal(err)
		}
		var cancelled bool
		if err := pool.QueryRow(ctx,
			"SELECT cancelled_at IS NOT NULL FROM confirmed_qr_email_outbox WHERE user_id=$1", id,
		).Scan(&cancelled); err != nil || !cancelled {
			t.Fatalf("cancelled=%v, error=%v", cancelled, err)
		}
	})
}
