-- +goose Up
ALTER TABLE teams DROP COLUMN IF EXISTS event_id;

-- +goose Down
ALTER TABLE teams ADD COLUMN event_id UUID REFERENCES events(id) ON DELETE CASCADE;
