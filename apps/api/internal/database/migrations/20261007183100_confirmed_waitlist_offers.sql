-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION record_waitlist_offer_deadline()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    offer_time timestamptz := clock_timestamp();
BEGIN
    IF OLD.status::text IN ('waitlisted', 'waitlist_confirmed')
       AND NEW.status::text = 'accepted' THEN
        INSERT INTO application_waitlist_offers (
            application_id, offered_at, confirmation_deadline
        ) VALUES (
            NEW.id, offer_time, offer_time + interval '48 hours'
        )
        ON CONFLICT (application_id) DO UPDATE
        SET offered_at = EXCLUDED.offered_at,
            confirmation_deadline = EXCLUDED.confirmation_deadline;
    END IF;
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION record_waitlist_offer_deadline()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    offer_time timestamptz := clock_timestamp();
BEGIN
    IF OLD.status = 'waitlisted' AND NEW.status = 'accepted' THEN
        INSERT INTO application_waitlist_offers (
            application_id, offered_at, confirmation_deadline
        ) VALUES (
            NEW.id, offer_time, offer_time + interval '48 hours'
        )
        ON CONFLICT (application_id) DO UPDATE
        SET offered_at = EXCLUDED.offered_at,
            confirmation_deadline = EXCLUDED.confirmation_deadline;
    END IF;
    RETURN NEW;
END;
$$;
-- +goose StatementEnd
