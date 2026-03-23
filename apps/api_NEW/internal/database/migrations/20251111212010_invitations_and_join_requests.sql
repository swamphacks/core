-- +goose Up
-- +goose StatementBegin
CREATE TYPE invitation_status AS ENUM ('PENDING', 'ACCEPTED', 'EXPIRED', 'REJECTED');
CREATE TYPE join_request_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TABLE team_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    invited_by_user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    invited_email TEXT NOT NULL,
    invited_user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    status invitation_status NOT NULL DEFAULT 'PENDING',
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_team_invitations
BEFORE UPDATE ON team_invitations
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TABLE team_join_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    request_message TEXT,
    status join_request_status NOT NULL DEFAULT 'PENDING',
    processed_by_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_team_join_requests
BEFORE UPDATE ON team_join_requests
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- Ensure a user can have only one pending invitation per team
CREATE UNIQUE INDEX idx_unique_pending_request 
ON team_join_requests (team_id, user_id) 
WHERE status = 'PENDING';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_team_invitations ON team_invitations;
DROP TRIGGER IF EXISTS set_updated_at_team_join_requests ON team_join_requests;

DROP INDEX IF EXISTS idx_unique_pending_request;

DROP TABLE IF EXISTS team_invitations;
DROP TABLE IF EXISTS team_join_requests;

DROP TYPE IF EXISTS invitation_status;
DROP TYPE IF EXISTS join_request_status;
-- +goose StatementEnd
