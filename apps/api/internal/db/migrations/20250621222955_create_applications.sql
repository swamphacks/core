-- +goose Up
-- +goose StatementBegin

CREATE TYPE application_status AS ENUM ('started', 'submitted', 'under_review', 'accepted', 'rejected', 'waitlisted', 'withdrawn');

CREATE TABLE applications (
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    -- event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    -- TODO: remove this
    event_id UUID, 
    status application_status DEFAULT 'started',
    application JSONB NOT NULL DEFAULT '{}'::JSONB,
    resume_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    saved_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, event_id)
);

CREATE INDEX idx_applications_status ON applications(status);
CREATE INDEX idx_applications_event_id ON applications(event_id);

-- Create trigger to update application updates
CREATE TRIGGER set_updated_at_applications
BEFORE UPDATE ON applications
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


DROP TRIGGER IF EXISTS set_updated_at_applications ON applications;
DROP TABLE IF EXISTS applications;
DROP TYPE IF EXISTS application_status;

-- +goose StatementEnd
