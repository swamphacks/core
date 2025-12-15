-- +goose Up
-- +goose StatementBegin
CREATE TYPE bat_result_status AS ENUM ('running','completed','failed');

CREATE TABLE bat_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    accepted_applicants UUID[] DEFAULT '{}',
    rejected_applicants UUID[] DEFAULT '{}',
    status bat_result_status DEFAULT 'running',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bat_results;
DROP TYPE IF EXISTS bat_result_status;
-- +goose StatementEnd
