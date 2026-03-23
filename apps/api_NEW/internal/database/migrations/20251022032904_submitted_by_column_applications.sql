-- +goose Up
-- +goose StatementBegin
ALTER TABLE applications
ADD COLUMN submitted_at TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE applications
DROP COLUMN submitted_at;
-- +goose StatementEnd
