-- +goose Up
CREATE TABLE event_interest_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    source TEXT
);

CREATE INDEX idx_event_interest_event_id ON event_interest_submissions (event_id);
CREATE UNIQUE INDEX uniq_event_email ON event_interest_submissions (event_id, email);

-- +goose Down
DROP INDEX IF EXISTS uniq_event_email;
DROP INDEX IF EXISTS idx_event_interest_event_id;
DROP TABLE IF EXISTS event_interest_submissions;
