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

-- name: GetEventRoleByUserID :one
SELECT event_id, role
FROM event_roles
WHERE user_id = $1
ORDER BY assigned_at DESC
LIMIT 1;