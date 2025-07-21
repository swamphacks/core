-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    location TEXT,
    location_url TEXT,
    max_attendees INT,

    -- Application phases
    application_open TIMESTAMPTZ NOT NULL,
    application_close TIMESTAMPTZ NOT NULL,
    rsvp_deadline TIMESTAMPTZ,
    decision_release TIMESTAMPTZ,

    -- Event phase
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,

    -- Metadata
    website_url TEXT,
    is_published BOOLEAN DEFAULT FALSE,
    saved_at TIMESTAMPTZ DEFAULT NOW(),

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TYPE event_role_type AS ENUM ('admin', 'staff', 'attendee', 'applicant');

CREATE TABLE event_roles (
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    role event_role_type NOT NULL,
    assigned_at TIMESTAMPTZ DEFAULT NOW(),

    PRIMARY KEY (user_id, event_id)
);

-- Triggers for update_at
CREATE TRIGGER set_updated_at_events
BEFORE UPDATE ON events
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_updated_at_events ON events;
DROP TABLE IF EXISTS event_roles;
DROP TYPE IF EXISTS event_role_type;
DROP TABLE IF EXISTS events;

-- +goose StatementEnd
