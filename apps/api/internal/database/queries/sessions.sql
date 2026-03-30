-- name: CreateSession :one
INSERT INTO sessions (user_id, expires_at, ip_address, user_agent)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = $1;

-- name: GetSessionsByUserID :many
SELECT * FROM sessions
WHERE user_id = $1;

-- name: UpdateSessionExpiration :exec
UPDATE sessions
SET expires_at = $2
WHERE id = $1;

-- name: DeleteExpiredSession :exec
DELETE FROM sessions
WHERE expires_at < NOW();

-- name: GetActiveSessionUserInfo :one
SELECT u.id AS user_id, u.name, u.email, u.preferred_email, u.onboarded, u.image, u.role, u.email_consent, u.checked_in_at, u.rfid, s.last_used_at
FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.id = $1
    AND (s.expires_at > NOW())
LIMIT 1;

-- name: TouchSession :exec
UPDATE sessions
SET expires_at = $2, last_used_at = NOW()
WHERE id = $1;

-- name: InvalidateSessionByID :exec
UPDATE sessions
SET expires_at = NOW()
WHERE id = $1;
