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

-- name: ListTeamsWithMembersByEvent :many
SELECT
    t.id,
    t.name,
    t.owner_id,
    t.event_id,
    -- Step 1: Cast the aggregated JSON array to JSONB
    (COALESCE(
        json_agg(
            json_build_object(
                'user_id', u.id,
                'name', u.name,
                'email', u.email,
                'image', u.image,
                'joined_at', tm.joined_at
            )
        ) FILTER (WHERE u.id IS NOT NULL),
        '[]'::json
    ))::jsonb AS members -- <--- Explicitly cast to jsonb
FROM
    teams t
LEFT JOIN
    team_members tm ON t.id = tm.team_id
LEFT JOIN
    auth.users u ON tm.user_id = u.id
WHERE
    t.event_id = $1
GROUP BY
    t.id
ORDER BY
    t.created_at DESC
LIMIT $2
OFFSET $3;