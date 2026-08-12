-- name: ListApiKeys :many
SELECT
    id,
    name,
    description,
    role,
    created_at,
    expires_at
FROM api_keys;

-- name: GetApiKeyById :one
SELECT
    id,
    name,
    description,
    role,
    created_at,
    expires_at
FROM api_keys
WHERE id = @id;

-- name: CreateApiKey :one
INSERT INTO api_keys (
    name,
    description,
    role,
    expires_at,
    secret_hash
) VALUES (
    @name,
    @description,
    @role::role,
    sqlc.narg(expires_at),
    @secret_hash
)
RETURNING
    id,
    name,
    description,
    role,
    created_at,
    expires_at;

-- name: DeleteApiKey :exec
DELETE FROM api_keys
WHERE id = @id;