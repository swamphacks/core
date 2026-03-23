-- +goose Up
ALTER TABLE applications ADD COLUMN id UUID DEFAULT gen_random_uuid();
-- Populate it for existing rows, if any
UPDATE applications SET id = gen_random_uuid() WHERE id IS NULL;
ALTER TABLE applications ALTER COLUMN id SET NOT NULL;

ALTER TABLE applications DROP CONSTRAINT applications_pkey;
ALTER TABLE applications ADD PRIMARY KEY (id);

ALTER TABLE applications DROP COLUMN IF EXISTS event_id;

DROP INDEX IF EXISTS idx_applications_event_id;

ALTER TABLE applications ADD COLUMN hackathon_iteration TEXT NOT NULL;

-- +goose Down
ALTER TABLE applications ADD COLUMN event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE;

ALTER TABLE applications DROP CONSTRAINT applications_pkey;

ALTER TABLE applications ADD PRIMARY KEY (user_id, event_id);

ALTER TABLE applications DROP COLUMN id;

ALTER TABLE applications DROP COLUMN hackathon_iteration;

CREATE INDEX idx_applications_event_id ON applications(event_id);