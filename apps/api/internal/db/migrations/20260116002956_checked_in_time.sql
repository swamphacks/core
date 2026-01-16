-- +goose Up
-- +goose StatementBegin
ALTER TABLE event_roles
ADD COLUMN checked_in_at TIMESTAMPTZ,
ADD COLUMN rfid TEXT;

CREATE UNIQUE INDEX idx_event_roles_rfid ON event_roles(rfid)
WHERE rfid IS NOT NULL; -- Only index non-nulls to optimize.
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event_roles
DROP COLUMN checked_in_at,
DROP COLUMN rfid;

DROP INDEX IF EXISTS idx_event_roles_rfid;
-- +goose StatementEnd
