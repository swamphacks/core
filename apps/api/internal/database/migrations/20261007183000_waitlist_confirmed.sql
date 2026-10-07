-- +goose Up
ALTER TYPE application_status
    ADD VALUE IF NOT EXISTS 'waitlist_confirmed';

ALTER TABLE waitlist
    ADD COLUMN signup_source text NOT NULL DEFAULT 'preregistered'
    CHECK (signup_source IN ('preregistered', 'day_of'));

-- +goose Down
ALTER TABLE waitlist DROP COLUMN signup_source;
-- PostgreSQL enum values cannot be removed without rebuilding the enum.
