-- name: GetRedeemables :many
-- total_redeemed counts NFC scans (one nfc_tags_redeemables row per item handed
-- out). COUNT(ntr.id), not COUNT(*): an item with no scans must count 0, not 1.
SELECT r.id,
    r.name,
    r.max_user_amount,
    r.created_at,
    r.updated_at,
    COUNT(ntr.id) AS total_redeemed
FROM redeemables r
LEFT JOIN nfc_tags_redeemables ntr ON ntr.redeemable_id = r.id
GROUP BY r.id
ORDER BY r.name, r.id;

-- name: CreateRedeemable :one
-- Inserts into the active hackathon. With no active hackathon the SELECT
-- yields no rows, so nothing is inserted and pgx returns ErrNoRows.
INSERT INTO redeemables (name, max_user_amount, type, hackathon_id)
SELECT @name, @max_user_amount, @type, h.id
FROM hackathons h
WHERE h.is_active
RETURNING *;

-- name: UpdateRedeemable :one
UPDATE redeemables
SET 
    name = COALESCE(sqlc.narg('name'), name),
    max_user_amount = COALESCE(sqlc.narg('max_user_amount'), max_user_amount),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteRedeemable :exec
DELETE FROM redeemables
WHERE id = $1;
