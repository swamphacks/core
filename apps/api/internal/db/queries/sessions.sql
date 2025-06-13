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

-- name: GetActiveSessionUserInfo :one
SELECT u.id AS user_id, u.name, u.onboarded, u.image, u.role, s.last_used_at
FROM auth.sessions s
JOIN auth.users u ON s.user_id = u.id
WHERE s.id = $1
    AND (s.expires_at > NOW())
LIMIT 1;

-- name: TouchSession :exec
UPDATE auth.sessions
SET expires_at = $2, last_used_at = NOW()
WHERE id = $1;

-- name: InvalidateSessionByID :exec
UPDATE auth.sessions
SET expires_at = NOW()
WHERE id = $1;
