-- +goose Up
ALTER TABLE redeemables DROP COLUMN IF EXISTS event_id;
-- +goose Down
ALTER TABLE redeemables ADD COLUMN event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE;
