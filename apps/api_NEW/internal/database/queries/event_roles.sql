-- name: AssignRole :exec
INSERT INTO event_roles (user_id, role)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveRole :exec
DELETE FROM event_roles
WHERE user_id = $1;

-- name: UpdateRole :exec
UPDATE event_roles
SET role = $2
WHERE user_id = $1;

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
