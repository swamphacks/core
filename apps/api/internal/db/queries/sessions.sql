-- name: CreateSession :one
INSERT INTO auth.sessions (user_id, token, expires_at, ip_address, user_agent)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM auth.sessions
WHERE TOKEN = $1;

-- name: GetSessionsByUserID :many
SELECT * FROM auth.sessions
WHERE user_id = $1;

-- name: UpdateSessionExpiration :exec
UPDATE auth.sessions
SET expires_at = $2
WHERE id = $1;

-- name: DeleteSessionByToken :exec
DELETE FROM auth.sessions
WHERE token = $1;

-- name: DeleteExpiredSession :exec
DELETE FROM auth.sessions
WHERE expires_at < NOW();
