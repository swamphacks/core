-- name: CreateInvitation :one
INSERT INTO team_invitations (
    team_id, invited_by_user_id, invited_email, expires_at
) VALUES (
    @team_id, @invited_by_user_id, @invited_email, @expires_at
)
RETURNING *;

-- name: GetInvitationByID :one
SELECT * FROM team_invitations 
WHERE id = @id;

-- name: ListPendingInvitationsByTeam :many
SELECT *
FROM team_invitations
WHERE team_id = @team_id AND status = 'PENDING'::invitation_status
ORDER BY created_at DESC;

-- name: UpdateInvitation :one
UPDATE team_invitations
SET
    invited_user_id = CASE WHEN @invited_user_id_do_update::boolean THEN @invited_user_id::uuid ELSE invited_user_id END,
    status = CASE WHEN @status_do_update::boolean THEN @status::invitation_status ELSE status END,
    expires_at = CASE WHEN @expires_at_do_update::boolean THEN @expires_at::timestamptz ELSE expires_at END,
    updated_at = NOW()
WHERE
    id = @id::uuid
RETURNING *;

-- name: DeleteInvitation :exec
DELETE FROM team_invitations 
WHERE id = @id;

-- name: AcceptInvitation :one
UPDATE team_invitations
SET
    status = 'ACCEPTED'::invitation_status,
    invited_user_id = @invited_user_id::uuid,
    updated_at = NOW()
WHERE
    id = @id AND 
    status = 'PENDING'::invitation_status AND 
    expires_at > NOW()
RETURNING *;

-- name: RejectInvitation :one
UPDATE team_invitations
SET
    status = 'REJECTED'::invitation_status,
    updated_at = NOW()
WHERE
    id = @id AND 
    status = 'PENDING'::invitation_status AND 
    expires_at > NOW()
RETURNING *;

-- name: ListInvitationsByInvitedUserIDAndStatus :many
SELECT *
FROM team_invitations
WHERE invited_user_id = @invited_user_id AND status = @status::invitation_status
ORDER BY created_at DESC;

-- name: GetPendingInvitationByEmailAndTeam :one
SELECT *
FROM team_invitations
WHERE invited_email = @invited_email 
  AND team_id = @team_id 
  AND status = 'PENDING'::invitation_status
  AND expires_at > NOW()
LIMIT 1;