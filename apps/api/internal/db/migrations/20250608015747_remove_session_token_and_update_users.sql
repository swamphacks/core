-- +goose Up
-- +goose StatementBegin
ALTER TABLE auth.sessions DROP COLUMN IF EXISTS token;

ALTER TABLE auth.users ALTER COLUMN email DROP NOT NULL;

ALTER TABLE auth.sessions ADD COLUMN last_used_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE auth.sessions ADD COLUMN token TEXT UNIQUE;

ALTER TABLE auth.users ALTER COLUMN email SET NOT NULL;

ALTER TABLE auth.sessions DROP COLUMN IF EXISTS last_used_at;
-- +goose StatementEnd
