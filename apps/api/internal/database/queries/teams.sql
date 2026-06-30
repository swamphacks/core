-- name: CreateTeam :one
INSERT INTO teams (name, owner_id) VALUES (@name, @owner_id) RETURNING *;

-- name: GetTeamById :one
SELECT *
FROM teams
WHERE id = @id;

-- name: GetTeamMembers :many
SELECT tm.user_id, users.image, users.name FROM team_members tm
JOIN users ON users.id = tm.user_id
WHERE tm.team_id = @team_id;

-- name: GetTeamDetails :one
SELECT 
    t.*, 
    COALESCE(
        json_agg(
            json_build_object(
                'id', tm.user_id,
                'name', users.name,
                'image', users.image,
                'joinedAt', tm.joined_at
            )
        ) FILTER (WHERE tm.user_id IS NOT NULL),
        '[]'
    ) AS members
FROM teams t
LEFT JOIN team_members tm ON tm.team_id = t.id
LEFT JOIN users ON users.id = tm.user_id
WHERE t.id = @id
GROUP BY t.id;

-- name: GetTeamByUserId :one
SELECT
    t.id,
    t.name,
    t.owner_id
FROM teams t
JOIN team_members tm ON t.id = tm.team_id
WHERE tm.user_id = @user_id;

-- name: UpdateTeamById :one
UPDATE teams
SET
    owner_id = CASE WHEN @owner_id_do_update::boolean THEN @owner_id ELSE owner_id END,
    name = CASE WHEN @name_do_update::boolean THEN @name ELSE name END
WHERE
    id = @id
RETURNING *;

-- name: DeleteTeamById :exec
DELETE FROM teams WHERE id = @id;

-- name: ListTeamsWithMembers :many
SELECT
    t.id,
    t.name,
    t.owner_id,
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
FROM teams t
LEFT JOIN team_members tm ON t.id = tm.team_id
LEFT JOIN users u ON tm.user_id = u.id
GROUP BY t.id
ORDER BY t.created_at DESC
LIMIT $1 OFFSET $2;

-- name: AddUserToTeam :one
INSERT INTO team_members (team_id, user_id) VALUES (@team_id, @user_id) RETURNING *;

-- name: RemoveUserFromTeam :exec
DELETE FROM team_members WHERE user_id = @user_id AND team_id = @team_id;

-- name: CreateInvitation :one
INSERT INTO team_invitations (
    team_id, inviter_id, expires_at
) VALUES (
    @team_id, @inviter_id, @expires_at
)
RETURNING *;

-- name: GetInvitationByID :one
SELECT * FROM team_invitations WHERE id = @id;

-- name: DeleteInvitation :exec
DELETE FROM team_invitations WHERE id = @id;