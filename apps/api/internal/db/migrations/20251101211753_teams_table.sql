-- +goose Up
-- +goose StatementBegin
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    owner_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_teams
BEFORE UPDATE ON auth.users
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TABLE team_members (
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT NOW(),

    PRIMARY KEY (user_id, team_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS team_members;
DROP TRIGGER IF EXISTS set_updated_at_teams ON auth.users;
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
