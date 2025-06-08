-- name: CreateSession :one
INSERT INTO auth.sessions (user_id, expires_at, ip_address, user_agent)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM auth.sessions
WHERE id = $1;

-- name: GetSessionsByUserID :many
SELECT * FROM auth.sessions
WHERE user_id = $1;

-- name: UpdateSessionExpiration :exec
UPDATE auth.sessions
SET expires_at = $2
WHERE id = $1;

-- name: DeleteExpiredSession :exec
DELETE FROM auth.sessions
WHERE expires_at < NOW();
