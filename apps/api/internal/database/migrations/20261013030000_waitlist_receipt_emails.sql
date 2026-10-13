-- +goose Up
CREATE TABLE waitlist_receipt_email_outbox (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    hackathon_id text NOT NULL REFERENCES hackathons(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kind text NOT NULL CHECK (kind IN ('signup', 'arrival')),
    signup_source text NOT NULL,
    action_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    sent_at timestamptz,
    cancelled_at timestamptz,
    attempts integer NOT NULL DEFAULT 0,
    next_attempt_at timestamptz NOT NULL DEFAULT now(),
    last_error text,
    UNIQUE (hackathon_id, user_id, kind, action_at)
);

CREATE INDEX waitlist_receipt_email_outbox_pending
ON waitlist_receipt_email_outbox (next_attempt_at, created_at)
WHERE sent_at IS NULL AND cancelled_at IS NULL;

-- +goose StatementBegin
CREATE FUNCTION queue_waitlist_receipt_email()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    receipt_kind text;
    receipt_time timestamptz;
BEGIN
    IF NEW.hackathon_id <> 'xii' THEN
        RETURN NEW;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM users u
        JOIN applications a ON a.user_id=u.id
        WHERE u.id=NEW.user_id
          AND a.hackathon_id=NEW.hackathon_id
          AND NOT u.is_fake AND NOT a.is_fake
    ) THEN
        RETURN NEW;
    END IF;

    IF TG_OP = 'INSERT' THEN
        IF NEW.in_person_joined_at IS NULL THEN
            receipt_kind := 'signup';
            receipt_time := NEW.created_at;
        ELSE
            receipt_kind := 'arrival';
            receipt_time := NEW.in_person_joined_at;
        END IF;
    ELSIF OLD.in_person_joined_at IS NULL
          AND NEW.in_person_joined_at IS NOT NULL THEN
        receipt_kind := 'arrival';
        receipt_time := NEW.in_person_joined_at;
    ELSE
        RETURN NEW;
    END IF;

    INSERT INTO waitlist_receipt_email_outbox (
        hackathon_id, user_id, kind, signup_source, action_at
    ) VALUES (
        NEW.hackathon_id, NEW.user_id, receipt_kind,
        NEW.signup_source, receipt_time
    )
    ON CONFLICT (hackathon_id, user_id, kind, action_at) DO NOTHING;

    RETURN NEW;
END;
$$;
-- +goose StatementEnd

CREATE TRIGGER waitlist_receipt_email
AFTER INSERT OR UPDATE OF in_person_joined_at ON waitlist
FOR EACH ROW EXECUTE FUNCTION queue_waitlist_receipt_email();

-- +goose Down
DROP TRIGGER waitlist_receipt_email ON waitlist;
DROP FUNCTION queue_waitlist_receipt_email();
DROP TABLE waitlist_receipt_email_outbox;
