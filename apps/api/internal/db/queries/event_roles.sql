-- name: GetEventStaff :many
SELECT u.*, er.role AS event_role
FROM auth.users u
JOIN event_roles er ON u.id = er.user_id
WHERE er.event_id = $1
  AND er.role IN ('admin', 'staff');

-- name: AssignRole :exec
INSERT INTO event_roles (event_id, user_id, role)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: RemoveRole :exec
DELETE FROM event_roles
WHERE event_id = $1
  AND user_id = $2;

-- name: GetEventUsers :many
SELECT u.*, er.role AS event_role
FROM auth.users u
JOIN event_roles er ON u.id = er.user_id
WHERE er.event_id = $1;

-- name: UpdateRole :exec
UPDATE event_roles
SET role = $3
WHERE event_id = $1 AND user_id = $2;

-- name: GetEventAttendeesWithDiscord :many
SELECT 
    a.account_id as discord_id,
    u.id as user_id,
    u.name,
    u.email
FROM auth.users u
JOIN event_roles er ON u.id = er.user_id
JOIN auth.accounts a ON u.id = a.user_id
WHERE er.event_id = $1
    AND er.role = 'attendee'
    AND a.provider_id = 'discord';

-- name: GetEventRoleByDiscordIDAndEventId :one
SELECT er.event_id, er.role
FROM event_roles er
JOIN auth.accounts a ON er.user_id = a.user_id
WHERE a.provider_id = 'discord'
    AND a.account_id = $1
    AND er.event_id = $2;