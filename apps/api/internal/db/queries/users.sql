-- name: CreateUser :one
INSERT INTO auth.users (name, email, image)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM auth.users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM auth.users
WHERE email = $1;

-- name: UpdateUserOnboarded :exec
UPDATE auth.users
SET onboarded = TRUE
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE auth.users
SET
    name = CASE WHEN @name_do_update::boolean THEN @name ELSE name END,
    email = CASE WHEN @email_do_update::boolean THEN @email ELSE email END,
    email_verified = CASE WHEN @email_verified_do_update::boolean THEN @email_verified ELSE email_verified END,
    preferred_email = CASE WHEN @preferred_email_do_update::boolean THEN @preferred_email ELSE preferred_email END,
    onboarded = CASE WHEN @onboarded_do_update::boolean THEN @onboarded ELSE onboarded END,
    image = CASE WHEN @image_do_update::boolean THEN @image ELSE image END,
    email_consent = CASE WHEN @email_consent_do_update::boolean THEN @email_consent ELSE email_consent END,
    updated_at = NOW()
WHERE
    id = @id::uuid;

-- name: DeleteUser :exec
DELETE FROM auth.users
WHERE id = $1;
