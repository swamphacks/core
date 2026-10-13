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

func TestWaitlistReceiptIntegration(t *testing.T) {
	const localURL = "postgres://postgres:postgres@127.0.0.1:55432/waitlist_test?sslmode=disable"
	url := os.Getenv("WAITLIST_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set WAITLIST_TEST_DATABASE_URL")
	}
	if url != localURL {
		t.Fatal("requires dedicated local test database")
	}
	if !admissionWaitlistIsOpen(time.Now()) {
		t.Skip("online signup window has closed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
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
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM hackathons WHERE is_active",
	).Scan(&active); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx,
		"SELECT count(*) FROM hackathons WHERE id='xii'",
	).Scan(&existing); err != nil {
		t.Fatal(err)
	}
	if active != 0 || existing != 0 {
		t.Fatal("requires no active event and no existing xii fixture")
	}

	exec(`INSERT INTO hackathons (
		id,name,application_open,application_close,
		start_time,end_time,rsvp_deadline,is_active
	) VALUES (
		'xii','Receipt test',now()-interval '2 days',
		now()-interval '1 day',now()+interval '1 day',
		now()+interval '3 days',now()+interval '12 hours',true
	)`)
	defer func() {
		if _, err := pool.Exec(context.Background(),
			"DELETE FROM hackathons WHERE id='xii'",
		); err != nil {
			t.Errorf("event cleanup: %v", err)
		}
	}()
	exec(`INSERT INTO waitlist_dispatch_policies (
		hackathon_id,enabled,invitations_open_at,online_join_closes_at,
		in_person_opens_at,invitations_close_at
	) VALUES (
		'xii',true,now()-interval '3 hours',now()-interval '2 hours',
		now()-interval '1 hour',now()+interval '1 day'
	)`)

	db := &database.DB{Pool: pool, Query: sqlc.New(pool)}
	service := &ApplicationService{
		db: db, txm: database.NewTransactionManager(db), logger: zerolog.Nop(),
	}
	newUser := func(t *testing.T) uuid.UUID {
		t.Helper()
		id := uuid.New()
		exec("INSERT INTO users (id,name,role) VALUES ($1,'Receipt Test','applicant')", id)
		exec(`INSERT INTO applications (
			user_id,hackathon_id,status,application
		) VALUES (
			$1,'xii','rejected',
			'{"firstName":"Receipt","preferredEmail":"receipt@example.test"}'
		)`, id)
		t.Cleanup(func() {
			if _, err := pool.Exec(context.Background(),
				"DELETE FROM users WHERE id=$1", id,
			); err != nil {
				t.Errorf("user cleanup: %v", err)
			}
		})
		return id
	}
	join := func(t *testing.T, id uuid.UUID) {
		t.Helper()
		if err := service.JoinAdmissionWaitlist(ctx, id); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("signup queues once and retries without resending", func(t *testing.T) {
		id := newUser(t)
		join(t, id)
		join(t, id)
		var queued int
		if err := pool.QueryRow(ctx, `
			SELECT count(*) FROM waitlist_receipt_email_outbox
			WHERE user_id=$1 AND kind='signup'
		`, id).Scan(&queued); err != nil || queued != 1 {
			t.Fatalf("queued=%d error=%v", queued, err)
		}

		calls := 0
		fail := func(_ context.Context, email, name, kind, source string) error {
			calls++
			if email != "receipt@example.test" || name != "Receipt" ||
				kind != "signup" || source != "preregistered" {
				t.Fatal("wrong signup recipient or receipt")
			}
			return errors.New("simulated send failure")
		}
		if err := service.DeliverWaitlistReceipts(ctx, fail); err != nil {
			t.Fatal(err)
		}
		if err := service.DeliverWaitlistReceipts(ctx, fail); err != nil {
			t.Fatal(err)
		}
		if calls != 1 {
			t.Fatal("retry ran before its scheduled time")
		}
		var attempts int
		var retry time.Time
		if err := pool.QueryRow(ctx, `
			SELECT attempts,next_attempt_at
			FROM waitlist_receipt_email_outbox WHERE user_id=$1
		`, id).Scan(&attempts, &retry); err != nil {
			t.Fatal(err)
		}
		if attempts != 1 || !retry.After(time.Now()) {
			t.Fatal("failed receipt was not scheduled for retry")
		}
		exec("UPDATE waitlist_receipt_email_outbox SET next_attempt_at=now() WHERE user_id=$1", id)
		success := func(context.Context, string, string, string, string) error {
			calls++
			return nil
		}
		if err := service.DeliverWaitlistReceipts(ctx, success); err != nil {
			t.Fatal(err)
		}
		if err := service.DeliverWaitlistReceipts(ctx, success); err != nil {
			t.Fatal(err)
		}
		var sent bool
		if err := pool.QueryRow(ctx, `
			SELECT sent_at IS NOT NULL
			FROM waitlist_receipt_email_outbox WHERE user_id=$1
		`, id).Scan(&sent); err != nil || !sent || calls != 2 {
			t.Fatalf("sent=%v calls=%d error=%v", sent, calls, err)
		}
	})

	t.Run("repeated arrival queues once and supersedes pending signup", func(t *testing.T) {
		id := newUser(t)
		join(t, id)
		staff := uuid.New()
		exec("INSERT INTO users (id,name,role) VALUES ($1,'Receipt Staff','staff')", staff)
		t.Cleanup(func() {
			if _, err := pool.Exec(context.Background(),
				"DELETE FROM users WHERE id=$1", staff,
			); err != nil {
				t.Errorf("staff cleanup: %v", err)
			}
		})
		now := time.Now().UTC()
		for i := 0; i < 2; i++ {
			if err := service.RecordInPersonAdmissionWaitlist(ctx, id, staff, now); err != nil {
				t.Fatal(err)
			}
		}
		var queued int
		if err := pool.QueryRow(ctx, `
			SELECT count(*) FROM waitlist_receipt_email_outbox
			WHERE user_id=$1 AND kind='arrival'
		`, id).Scan(&queued); err != nil || queued != 1 {
			t.Fatalf("arrival queued=%d error=%v", queued, err)
		}
		calls := 0
		if err := service.DeliverWaitlistReceipts(ctx,
			func(_ context.Context, email, name, kind, source string) error {
				calls++
				if kind != "arrival" || email != "receipt@example.test" {
					t.Fatal("sent stale signup instead of arrival")
				}
				return nil
			},
		); err != nil {
			t.Fatal(err)
		}
		if calls != 1 {
			t.Fatalf("arrival sends=%d; want 1", calls)
		}
	})

	t.Run("day-of registration selects day-of receipt", func(t *testing.T) {
		id := newUser(t)
		exec("UPDATE applications SET status='waitlisted' WHERE user_id=$1", id)
		exec(`INSERT INTO waitlist (hackathon_id,user_id,signup_source)
			VALUES ('xii',$1,'day_of')`, id)
		calls := 0
		if err := service.DeliverWaitlistReceipts(ctx,
			func(_ context.Context, email, name, kind, source string) error {
				calls++
				if kind != "signup" || source != "day_of" {
					t.Fatal("wrong day-of receipt")
				}
				return nil
			},
		); err != nil {
			t.Fatal(err)
		}
		if calls != 1 {
			t.Fatalf("day-of sends=%d; want 1", calls)
		}
	})

	t.Run("leaving cancels pending receipt", func(t *testing.T) {
		id := newUser(t)
		join(t, id)
		if err := service.LeaveAdmissionWaitlist(ctx, id); err != nil {
			t.Fatal(err)
		}
		if err := service.DeliverWaitlistReceipts(ctx,
			func(context.Context, string, string, string, string) error {
				t.Fatal("sent receipt after leaving")
				return nil
			},
		); err != nil {
			t.Fatal(err)
		}
		var cancelled bool
		if err := pool.QueryRow(ctx, `
			SELECT cancelled_at IS NOT NULL
			FROM waitlist_receipt_email_outbox WHERE user_id=$1
		`, id).Scan(&cancelled); err != nil || !cancelled {
			t.Fatalf("cancelled=%v error=%v", cancelled, err)
		}
	})
}
