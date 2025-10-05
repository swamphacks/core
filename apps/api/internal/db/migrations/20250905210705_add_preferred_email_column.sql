-- +goose Up
-- +goose StatementBegin
ALTER TABLE auth.users
ADD COLUMN preferred_email TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE auth.users
DROP COLUMN preferred_email;
-- +goose StatementEnd
