-- name: GetStaff :many
SELECT u.*, er.role AS event_role
FROM auth.users u
JOIN event_roles er ON u.id = er.user_id
WHERE er.role IN ('admin', 'staff');

-- name: AssignRole :exec
INSERT INTO event_roles (user_id, role)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveRole :exec
DELETE FROM event_roles
WHERE user_id = $1;

-- name: GetUsesWithRoles :many
SELECT u.*, er.role AS event_role
FROM auth.users u
JOIN event_roles er ON u.id = er.user_id;

-- name: UpdateRole :exec
UPDATE event_roles
SET role = $2
WHERE user_id = $1;

-- name: GetAttendeesWithDiscord :many
SELECT 
    a.account_id as discord_id,
    u.id as user_id,
    u.name,
    u.email
FROM auth.users u
JOIN event_roles er ON u.id = er.user_id
JOIN auth.accounts a ON u.id = a.user_id
WHERE er.role = 'attendee'
    AND a.provider_id = 'discord';

-- name: GetRoleByDiscordID :one
SELECT er.role
FROM event_roles er
JOIN auth.accounts a ON er.user_id = a.user_id
WHERE a.provider_id = 'discord'
    AND a.account_id = $1;

-- name: UpdateRoleByUserId :exec
UPDATE event_roles
SET
  role = CASE WHEN @role_do_update::boolean THEN @role ELSE role END,
  rfid = CASE WHEN @rfid_do_update::boolean THEN @rfid ELSE rfid END,
  checked_in_at = CASE WHEN @checked_in_at_do_update::boolean THEN @checked_in_at ELSE checked_in_at END
WHERE user_id = @user_id;

-- name: GetAttendeeCount :one
SELECT COUNT(*) FROM event_roles AS er
WHERE er.role = 'attendee';

-- name: GetAttendeeUserIds :many
SELECT er.user_id FROM event_roles AS er
WHERE er.role = 'attendee';

-- name: GetUserByRFID :one
SELECT u.*
FROM auth.users u
JOIN event_roles er ON u.id = er.user_id
WHERE er.rfid = $1;

-- name: GetCheckedInStatusByUserIds :one
SELECT EXISTS (
    SELECT 1 
    FROM event_roles 
    WHERE user_id = $1 
      AND checked_in_at IS NOT NULL
)::bool;
