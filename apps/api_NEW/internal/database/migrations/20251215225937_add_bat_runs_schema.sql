-- +goose Up
-- +goose StatementBegin
CREATE TYPE bat_run_status AS ENUM ('running','completed','failed');

CREATE TABLE bat_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    accepted_applicants UUID[] DEFAULT '{}',
    rejected_applicants UUID[] DEFAULT '{}',
    status bat_run_status DEFAULT 'running',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bat_runs;
DROP TYPE IF EXISTS bat_run_status;
-- +goose StatementEnd

