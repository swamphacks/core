-- +goose Up
CREATE TABLE application_waitlist_offers (
    application_id uuid PRIMARY KEY
        REFERENCES applications(id) ON DELETE CASCADE,
    offered_at timestamptz NOT NULL,
    confirmation_deadline timestamptz NOT NULL,
    CHECK (confirmation_deadline > offered_at)
);

-- +goose StatementBegin
CREATE FUNCTION record_waitlist_offer_deadline()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    offer_time timestamptz := clock_timestamp();
BEGIN
    IF OLD.status = 'waitlisted' AND NEW.status = 'accepted' THEN
        INSERT INTO application_waitlist_offers (
            application_id, offered_at, confirmation_deadline
        )
        VALUES (
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

CREATE TRIGGER record_waitlist_offer_deadline
AFTER UPDATE OF status ON applications
FOR EACH ROW
EXECUTE FUNCTION record_waitlist_offer_deadline();

-- +goose Down
DROP TRIGGER record_waitlist_offer_deadline ON applications;
DROP FUNCTION record_waitlist_offer_deadline();
DROP TABLE application_waitlist_offers;
