-- name: CreateUser :one
INSERT INTO users (name, email, image)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserEmailInfoById :one
SELECT
    id,
    name,
    email_consent,
    CASE
        WHEN preferred_email IS NOT NULL AND preferred_email != '' THEN preferred_email
        ELSE email
    END AS contact_email
FROM users
WHERE id = $1;

-- name: UpdateUserOnboarded :exec
UPDATE users
SET onboarded = TRUE
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET
    name = CASE WHEN @name_do_update::boolean THEN @name ELSE name END,
    email = CASE WHEN @email_do_update::boolean THEN @email ELSE email END,
    email_verified = CASE WHEN @email_verified_do_update::boolean THEN @email_verified ELSE email_verified END,
    preferred_email = CASE WHEN @preferred_email_do_update::boolean THEN @preferred_email ELSE preferred_email END,
    onboarded = CASE WHEN @onboarded_do_update::boolean THEN @onboarded ELSE onboarded END,
    image = CASE WHEN @image_do_update::boolean THEN @image ELSE image END,
    email_consent = CASE WHEN @email_consent_do_update::boolean THEN @email_consent ELSE email_consent END,
    checked_in_at = CASE WHEN @checked_in_at_do_update::boolean THEN @checked_in_at ELSE checked_in_at END,
    rfid = CASE WHEN @rfid_do_update::boolean THEN @rfid ELSE rfid END,
    role = CASE WHEN @role_do_update::boolean THEN @role ELSE role END,
    role_assigned_at = CASE WHEN @role_do_update::boolean THEN NOW() ELSE role_assigned_at END,
    updated_at = NOW()
WHERE
    id = @id::uuid;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUsers :many
SELECT *
FROM users
WHERE LOWER(name) LIKE LOWER('%' || COALESCE(sqlc.arg('search'), '') || '%')
   OR LOWER(email) LIKE LOWER('%' || COALESCE(sqlc.arg('search'), '') || '%')
ORDER BY name
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetUserByRFID :one
SELECT * FROM users
WHERE rfid = $1;

-- name: UpdateRole :exec
UPDATE users
SET role = @role::role_type,
    role_assigned_at = NOW()
WHERE id = @user_id::uuid;

-- name: RemoveRole :exec
UPDATE users
SET role = NULL,
    role_assigned_at = NOW()
WHERE id = @user_id::uuid;

-- name: UpdateCheckInTime :exec
UPDATE users
SET checked_in_at = @checked_in_at
WHERE id = @user_id::uuid;

-- name: UpdateRFID :exec
UPDATE users
SET rfid = @rfid
WHERE id = @user_id::uuid;