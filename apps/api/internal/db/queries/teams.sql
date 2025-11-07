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

-- name: GetUserEventTeam :one
SELECT
    t.id,
    t.name,
    t.owner_id,
    t.event_id
FROM
    teams t
JOIN
    team_members tm ON t.id = tm.team_id
WHERE
    t.event_id = $1
    AND tm.user_id = $2
LIMIT 1;

-- name: GetTeamById :one
SELECT *
FROM teams
WHERE id = $1;

-- name: UpdateTeamById :one
UPDATE teams
SET
    owner_id = CASE WHEN @owner_id_do_update::boolean THEN @owner_id ELSE owner_id END,
    name = CASE WHEN @name_do_update::boolean THEN @name ELSE name END
WHERE
    id = @id::uuid
RETURNING *;
