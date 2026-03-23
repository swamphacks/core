-- +goose Up
ALTER TABLE event_roles DROP CONSTRAINT event_roles_pkey;
ALTER TABLE event_roles DROP COLUMN IF EXISTS event_id;
ALTER TABLE event_roles ADD PRIMARY KEY (user_id);

-- +goose Down
ALTER TABLE event_roles ADD COLUMN event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE;
ALTER TABLE event_roles DROP CONSTRAINT event_roles_pkey;
ALTER TABLE event_roles ADD PRIMARY KEY (user_id, event_id);