-- +goose Up
-- +goose StatementBegin
ALTER TABLE applications DROP COLUMN IF EXISTS resume_url;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE applications ADD COLUMN resume_url TEXT;
-- +goose StatementEnd
