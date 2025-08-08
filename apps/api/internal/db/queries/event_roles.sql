-- name: GetEventStaff :many
SELECT u.*
FROM auth.users u
JOIN event_roles er ON u.id = er.user_id
WHERE er.event_id = $1
  AND er.role IN ('admin', 'staff');
