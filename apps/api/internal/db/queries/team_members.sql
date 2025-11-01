-- name: AddTeamMember :one
INSERT INTO team_members (
    user_id,
    team_id
) VALUES (
    $1,
    $2
)
RETURNING *;

-- name: RemoveTeamMember :exec
DELETE FROM team_members
WHERE user_id = $1
  AND team_id = $2;
