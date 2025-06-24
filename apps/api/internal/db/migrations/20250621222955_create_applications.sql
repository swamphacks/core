-- +goose Up
-- +goose StatementBegin

CREATE TYPE application_status AS ENUM ('open', 'submitted', 'accepted', 'rejected', 'waitlisted');

CREATE TABLE applications (
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    status application_status DEFAULT 'open',
    application JSONB,
    resume_url TEXT,
    opened_at TIMESTAMPTZ DEFAULT NOW(),
    submitted_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    

    PRIMARY KEY (user_id, event_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS applications;
DROP TYPE IF EXISTS application_status;

-- +goose StatementEnd
