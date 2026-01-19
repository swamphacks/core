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

-- name: UpdateEventRoleByIds :exec
UPDATE event_roles
SET
  role = CASE WHEN @role_do_update::boolean THEN @role ELSE role END,
  rfid = CASE WHEN @rfid_do_update::boolean THEN @rfid ELSE rfid END,
  checked_in_at = CASE WHEN @checked_in_at_do_update::boolean THEN @checked_in_at ELSE checked_in_at END
WHERE user_id = @user_id
  AND event_id = @event_id;

-- name: GetAttendeeCountByEventId :one
SELECT COUNT(*) FROM event_roles AS er
WHERE er.event_id = @event_id::uuid
  AND er.role = 'attendee';

-- name: GetUserByRFID :one
SELECT u.*
FROM auth.users u
JOIN event_roles er ON u.id = er.user_id
WHERE er.event_id = $1
  AND er.rfid = $2;