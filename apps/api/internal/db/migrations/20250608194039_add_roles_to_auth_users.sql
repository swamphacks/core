-- +goose Up
-- +goose StatementBegin
CREATE TYPE auth_user_role AS ENUM ('user', 'superuser');

ALTER TABLE auth.users ADD COLUMN role auth_user_role NOT NULL DEFAULT 'user';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE auth.users DROP COLUMN IF EXISTS role;
DROP TYPE IF EXISTS auth_user_role;
-- +goose StatementEnd
