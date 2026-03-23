-- +goose Up
ALTER TABLE events DROP CONSTRAINT events_pkey;
ALTER TABLE events DROP COLUMN IF EXISTS id;

ALTER TABLE events ADD COLUMN onerow_id bool DEFAULT true;
ALTER TABLE events ADD PRIMARY KEY (onerow_id);

ALTER TABLE events ADD CONSTRAINT onerow_uni CHECK (onerow_id);

ALTER TABLE events RENAME TO hackathon_config;

REVOKE DELETE, TRUNCATE ON hackathon_config FROM public;

-- +goose Down
ALTER TABLE hackathon_config RENAME TO events;

GRANT DELETE, TRUNCATE ON events TO public;

ALTER TABLE events DROP CONSTRAINT onerow_uni;
ALTER TABLE events DROP CONSTRAINT events_pkey;
ALTER TABLE events DROP COLUMN onerow_id;

ALTER TABLE events ADD COLUMN id UUID DEFAULT gen_random_uuid();
UPDATE events SET id = gen_random_uuid() WHERE id IS NULL;
ALTER TABLE events ALTER COLUMN id SET NOT NULL;
ALTER TABLE events ADD PRIMARY KEY (id);
