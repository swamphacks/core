-- +goose Up
ALTER TABLE applications DROP CONSTRAINT applications_pkey;
ALTER TABLE applications DROP COLUMN IF EXISTS id;

ALTER TABLE applications ADD PRIMARY KEY (user_id);

-- +goose Down
ALTER TABLE applications DROP CONSTRAINT applications_pkey;

ALTER TABLE applications ADD COLUMN id UUID DEFAULT gen_random_uuid();
UPDATE applications SET id = gen_random_uuid() WHERE id IS NULL;
ALTER TABLE applications ALTER COLUMN id SET NOT NULL;
ALTER TABLE applications ADD PRIMARY KEY (id);