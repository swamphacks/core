-- +goose Up
CREATE TABLE confirmed_qr_email_outbox (
    application_id uuid PRIMARY KEY REFERENCES applications(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    hackathon_id text NOT NULL REFERENCES hackathons(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    sent_at timestamptz,
    cancelled_at timestamptz,
    attempts integer NOT NULL DEFAULT 0,
    next_attempt_at timestamptz NOT NULL DEFAULT now(),
    last_error text,
    message_id text
);

CREATE INDEX confirmed_qr_email_outbox_pending
ON confirmed_qr_email_outbox (next_attempt_at, created_at)
WHERE sent_at IS NULL AND cancelled_at IS NULL;

-- +goose StatementBegin
CREATE FUNCTION queue_confirmed_qr_email()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF OLD.status::text = 'accepted'
       AND NEW.status::text = 'confirmed'
       AND NEW.hackathon_id = 'xii'
       AND NOT NEW.is_fake
       AND EXISTS (
           SELECT 1 FROM users
           WHERE id=NEW.user_id AND NOT is_fake
       ) THEN
        INSERT INTO confirmed_qr_email_outbox (
            application_id, user_id, hackathon_id
        ) VALUES (NEW.id, NEW.user_id, NEW.hackathon_id)
        ON CONFLICT (application_id) DO NOTHING;
    END IF;
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

CREATE TRIGGER application_confirmed_qr_email
AFTER UPDATE OF status ON applications
FOR EACH ROW EXECUTE FUNCTION queue_confirmed_qr_email();

-- +goose Down
DROP TRIGGER application_confirmed_qr_email ON applications;
DROP FUNCTION queue_confirmed_qr_email();
DROP TABLE confirmed_qr_email_outbox;
