-- name: AddEmail :one
INSERT INTO emails (event_id, user_id, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteEmailByID :exec
DELETE FROM emails
WHERE id = $1;

-- name: GetEmailsByUser :many
SELECT * FROM emails
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetEmailsByEvent :many
SELECT * FROM emails
WHERE event_id = $1
ORDER BY created_at DESC;

-- name: GetEmailByID :one
SELECT * FROM emails
WHERE id = $1;

-- name: UpdateEmailAddress :one
UPDATE emails
SET email = $2, updated_at = now()
WHERE id = $1
RETURNING *;