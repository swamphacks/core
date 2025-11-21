-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
ADD COLUMN application_review_started BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE applications
ADD COLUMN experience_rating INTEGER;
ALTER TABLE applications
ADD COLUMN passion_rating INTEGER;
ALTER TABLE applications
ADD COLUMN assigned_reviewer UUID REFERENCES auth.users(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
DROP COLUMN application_review_started;
ALTER TABLE applications
DROP COLUMN experience_rating;
ALTER TABLE applications
DROP COLUMN passion_rating;
ALTER TABLE applications
DROP COLUMN assigned_reviewer;
-- +goose StatementEnd
