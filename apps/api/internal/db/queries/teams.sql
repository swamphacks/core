-- name: CreateTeam :one
INSERT INTO teams (
    name,
    owner_id,
    event_id
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams
WHERE id = $1;
