package repository

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swamphacks/core/apps/api/internal/database"
)

// Uses an isolated schema; never points at the event database by default.
func TestWorkshopAttendanceIntegration(t *testing.T) {
	uri := os.Getenv("WORKSHOP_TEST_DATABASE_URL")
	if uri == "" {
		t.Skip("set WORKSHOP_TEST_DATABASE_URL to a local disposable PostgreSQL database")
	}
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		t.Fatal(err)
	}
	if config.ConnConfig.Host != "127.0.0.1" && config.ConnConfig.Host != "localhost" {
		t.Fatal("workshop tests require a local database")
	}
	setup, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer setup.Close()
	schema := "workshop_test_" + uuid.NewString()[:8]
	if _, err = setup.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	defer func() { _, _ = setup.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE") }()
	config.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	_, err = pool.Exec(ctx, `CREATE TABLE users(id uuid PRIMARY KEY,name text,checked_in_at timestamptz,is_fake boolean DEFAULT false);
 CREATE TABLE hackathons(id text PRIMARY KEY,is_active boolean);
 CREATE TABLE applications(user_id uuid,hackathon_id text,status text,application jsonb DEFAULT '{}',is_fake boolean DEFAULT false);
 CREATE TABLE workshops(id uuid PRIMARY KEY,title text,type text,start_time timestamptz,num_attendees integer DEFAULT 0);
 CREATE TABLE nfc_tags_user(tag_id text PRIMARY KEY,user_id uuid UNIQUE);
 CREATE TABLE nfc_tags_workshops(user_id uuid,workshop_id uuid,tag_id text,UNIQUE(user_id,workshop_id));
 INSERT INTO hackathons VALUES('xii',true);`)
	if err != nil {
		t.Fatal(err)
	}
	r := NewWorkshopsRepository(&database.DB{Pool: pool})
	workshop := uuid.New()
	if _, err = pool.Exec(ctx, `INSERT INTO workshops(id,title,type,start_time) VALUES($1,'Workshop','workshop',now())`, workshop); err != nil {
		t.Fatal(err)
	}
	addUser := func(status string, checkin bool, nfc bool) uuid.UUID {
		t.Helper()
		id := uuid.New()
		_, err := pool.Exec(ctx, `INSERT INTO users(id,name,checked_in_at) VALUES($1,'Hacker',CASE WHEN $2 THEN now() ELSE NULL END)`, id, checkin)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = pool.Exec(ctx, `INSERT INTO applications(user_id,hackathon_id,status) VALUES($1,'xii',$2)`, id, status); err != nil {
			t.Fatal(err)
		}
		if nfc {
			if _, err = pool.Exec(ctx, `INSERT INTO nfc_tags_user VALUES($1,$2)`, id.String(), id); err != nil {
				t.Fatal(err)
			}
		}
		return id
	}
	t.Run("QR without NFC and repeat", func(t *testing.T) {
		id := addUser("confirmed", true, false)
		got, err := r.RecordAttendance(ctx, workshop, id)
		if err != nil || got.AlreadyRecorded {
			t.Fatalf("first=%+v err=%v", got, err)
		}
		got, err = r.RecordAttendance(ctx, workshop, id)
		if err != nil || !got.AlreadyRecorded {
			t.Fatalf("repeat=%+v err=%v", got, err)
		}
	})
	t.Run("NFC-only checkin qualifies", func(t *testing.T) {
		id := addUser("confirmed", false, true)
		if _, err := r.RecordAttendance(ctx, workshop, id); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("existing NFC attendance deduplicates", func(t *testing.T) {
		id := addUser("confirmed", false, true)
		if _, err := pool.Exec(ctx, `INSERT INTO nfc_tags_workshops VALUES($1,$2,$3)`, id, workshop, id.String()); err != nil {
			t.Fatal(err)
		}
		got, err := r.RecordAttendance(ctx, workshop, id)
		if err != nil || !got.AlreadyRecorded {
			t.Fatalf("got=%+v err=%v", got, err)
		}
	})
	for _, status := range []string{"accepted", "waitlist_confirmed", "rejected", "started"} {
		t.Run("reject "+status, func(t *testing.T) {
			id := addUser(status, true, false)
			_, err := r.RecordAttendance(ctx, workshop, id)
			if !errors.Is(err, ErrWorkshopAttendanceIneligible) {
				t.Fatal(err)
			}
		})
	}
	t.Run("reject not checked in", func(t *testing.T) {
		id := addUser("confirmed", false, false)
		_, err := r.RecordAttendance(ctx, workshop, id)
		if !errors.Is(err, ErrWorkshopAttendanceIneligible) {
			t.Fatal(err)
		}
	})
	t.Run("missing workshop", func(t *testing.T) {
		_, err := r.RecordAttendance(ctx, uuid.New(), uuid.New())
		if !errors.Is(err, ErrWorkshopAttendanceNotFound) {
			t.Fatal(err)
		}
	})
	t.Run("concurrent scans", func(t *testing.T) {
		id := addUser("confirmed", true, false)
		var wg sync.WaitGroup
		errs := make(chan error, 8)
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); _, err := r.RecordAttendance(ctx, workshop, id); errs <- err }()
		}
		wg.Wait()
		close(errs)
		for err := range errs {
			if err != nil {
				t.Fatal(err)
			}
		}
		var count int
		if err := pool.QueryRow(ctx, `SELECT count(*) FROM nfc_tags_workshops WHERE user_id=$1 AND workshop_id=$2`, id, workshop).Scan(&count); err != nil || count != 1 {
			t.Fatalf("count=%d err=%v", count, err)
		}
	})
	var registrations int
	if err := pool.QueryRow(ctx, `SELECT num_attendees FROM workshops WHERE id=$1`, workshop).Scan(&registrations); err != nil || registrations != 0 {
		t.Fatalf("RSVP count changed: %d err=%v", registrations, err)
	}
}
