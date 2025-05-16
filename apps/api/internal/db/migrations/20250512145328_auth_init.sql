-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    onboarded BOOLEAN NOT NULL DEFAULT FALSE,
    image TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE auth.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE auth.accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    provider_id TEXT NOT NULL,
    account_id TEXT NOT NULL,
    hashed_password TEXT,
    access_token TEXT,
    refresh_token TEXT,
    id_token TEXT,
    access_token_expires_at TIMESTAMPTZ,
    refresh_token_expires_at TIMESTAMPTZ,
    scope TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider_id, account_id)
);

-- Index for quick lookup of sessions by user
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON auth.sessions (user_id);

-- Index for cleaning up expired sessions quickly
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON auth.sessions (expires_at);

-- Index for quick lookup of accounts by user
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON auth.accounts (user_id);

-- Index for querying provider + account
CREATE INDEX IF NOT EXISTS idx_accounts_provider_account ON auth.accounts (provider_id, account_id);

-- Updated at function
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = clock_timestamp();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for users
CREATE TRIGGER set_updated_at_users
BEFORE UPDATE ON auth.users
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- Create trigger for sessions
CREATE TRIGGER set_updated_at_sessions
BEFORE UPDATE ON auth.sessions
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- Create trigger for accounts
CREATE TRIGGER set_updated_at_accounts
BEFORE UPDATE ON auth.accounts
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_users;
DROP TRIGGER IF EXISTS set_updated_at_sessions;
DROP TRIGGER IF EXISTS set_updated_at_accounts;

DROP FUNCTION IF EXISTS update_modified_column;


DROP INDEX IF EXISTS auth.idx_accounts_provider_account;
DROP INDEX IF EXISTS auth.idx_accounts_user_id;

DROP INDEX IF EXISTS auth.idx_sessions_expires_at;
DROP INDEX IF EXISTS auth.idx_sessions_user_id;

DROP TABLE IF EXISTS auth.accounts;
DROP TABLE IF EXISTS auth.sessions;
DROP TABLE IF EXISTS auth.users;

DROP SCHEMA IF EXISTS auth;
-- +goose StatementEnd
