-- +goose Up
-- +goose StatementBegin

CREATE TYPE application_status AS ENUM ('open', 'submitted', 'under_review', accepted', 'rejected', 'waitlisted');

CREATE TABLE applications (
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    status application_status DEFAULT 'open',
    application JSONB,
    resume_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    PRIMARY KEY (user_id, event_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS applications;
DROP TYPE IF EXISTS application_status;

-- +goose StatementEnd
