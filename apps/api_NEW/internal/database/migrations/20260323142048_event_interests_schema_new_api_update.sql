-- +goose Up
ALTER TABLE event_interest_submissions DROP COLUMN IF EXISTS event_id;

DROP INDEX IF EXISTS idx_event_interest_event_id;
DROP INDEX IF EXISTS uniq_event_email;

-- +goose Down
ALTER TABLE event_interest_submissions ADD COLUMN event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE;
CREATE INDEX idx_event_interest_event_id ON event_interest_submissions (event_id);
CREATE UNIQUE INDEX uniq_event_email ON event_interest_submissions (event_id, email);