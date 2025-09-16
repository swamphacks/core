-- +goose Up
-- +goose StatementBegin
ALTER TABLE auth.users
ADD COLUMN email_consent BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE auth.users
DROP COLUMN email_consent;
-- +goose StatementEnd
