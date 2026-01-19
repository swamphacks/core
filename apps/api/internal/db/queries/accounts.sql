-- name: CreateAccount :one
INSERT INTO auth.accounts (
    user_id, provider_id, account_id, hashed_password,
    access_token, refresh_token, id_token,
    access_token_expires_at, refresh_token_expires_at, scope
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7,
    $8, $9, $10
)
RETURNING *;

-- name: GetByProviderAndAccountID :one
SELECT * FROM auth.accounts
WHERE provider_id = $1 AND account_id = $2;

-- name: GetByUserID :many
SELECT * FROM auth.accounts
WHERE user_id = $1;

-- name: UpdateTokens :exec
UPDATE auth.accounts
SET access_token = $3,
    refresh_token = $4,
    id_token = $5,
    access_token_expires_at = $6,
    refresh_token_expires_at = $7,
    scope = $8
WHERE provider_id = $1 AND account_id = $2;

-- name: DeleteAccount :exec
DELETE FROM auth.accounts
WHERE provider_id = $1 AND account_id = $2;

-- name: GetUserIDByDiscordAccountID :one
SELECT user_id
FROM auth.accounts
WHERE provider_id = 'discord' AND account_id = $1;
