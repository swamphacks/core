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