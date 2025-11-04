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

-- name: GetTeamMembers :many
SELECT
    u.id AS user_id,
    u.email,
    u.image,
    u.name,
    tm.joined_at
FROM
    team_members tm
JOIN
    auth.users u ON tm.user_id = u.id
WHERE
    tm.team_id = $1;
