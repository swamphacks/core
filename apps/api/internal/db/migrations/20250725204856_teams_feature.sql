-- +goose Up
-- +goose StatementBegin

CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL,

    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    admin_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE team_members (
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (team_id, user_id)
);

CREATE TABLE team_invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    inviter_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    invited_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    invited_email TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    accepted_at TIMESTAMPTZ,

    UNIQUE (team_id, invited_email)
);

ALTER TABLE events ADD COLUMN max_team_size INTEGER NOT NULL DEFAULT 4;

-- Triggers for updated at
CREATE TRIGGER set_updated_at_teams
BEFORE UPDATE ON teams
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_teams on teams;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS team_invites;
ALTER TABLE events DROP COLUMN IF EXISTS max_team_size;
-- +goose StatementEnd
