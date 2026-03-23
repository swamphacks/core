-- +goose Up
ALTER TABLE bat_runs DROP COLUMN IF EXISTS event_id;

-- +goose Down
ALTER TABLE bat_runs ADD COLUMN event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE;