-- name: GetRedeemables :many
-- total_redeemed counts NFC scans (one nfc_tags_redeemables row per item handed
-- out). COUNT(ntr.id), not COUNT(*): an item with no scans must count 0, not 1.
SELECT r.id,
    r.name,
    r.type,
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
    type = COALESCE(sqlc.narg('type'), type),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: RedeemForUser :one
-- Records one item handed out after staff scan a hacker's QR code on the web
-- portal. No NFC tag is involved, so tag_id stays NULL. Inserts only for a
-- hacker who is checked in (NFC nametag linked or checked_in_at set) with a
-- confirmed application to the item's hackathon; otherwise nothing is inserted
-- and pgx returns ErrNoRows. check_redeemable_count_trigger enforces the
-- per-hacker limit. The outer SELECT sees the table as it was before the
-- insert, hence the extra - 1.
WITH redeemed AS (
    INSERT INTO nfc_tags_redeemables (user_id, redeemable_id)
    SELECT u.id, r.id
    FROM users u
    JOIN redeemables r ON r.id = @redeemable_id
    JOIN applications a ON a.user_id = u.id AND a.hackathon_id = r.hackathon_id
    WHERE u.id = @user_id
      AND a.status = 'confirmed'
      AND (
          u.checked_in_at IS NOT NULL
          OR EXISTS (SELECT 1 FROM nfc_tags_user ntu WHERE ntu.user_id = u.id)
      )
    RETURNING user_id, redeemable_id
)
SELECT
    COALESCE(
        NULLIF(trim(concat_ws(' ', a.application->>'firstName', a.application->>'lastName')), ''),
        u.name
    )::text AS name,
    r.name AS redeemable_name,
    r.max_user_amount,
    GREATEST(
        r.max_user_amount - (
            SELECT COUNT(*)::integer
            FROM nfc_tags_redeemables ntr
            WHERE ntr.user_id = redeemed.user_id
              AND ntr.redeemable_id = redeemed.redeemable_id
        ) - 1,
        0
    )::integer AS items_remaining
FROM redeemed
JOIN users u ON u.id = redeemed.user_id
JOIN redeemables r ON r.id = redeemed.redeemable_id
JOIN applications a ON a.user_id = u.id AND a.hackathon_id = r.hackathon_id;

-- name: GetRedeemEligibility :one
-- Explains why RedeemForUser inserted nothing. No row means the redeemable
-- does not exist.
SELECT
    r.name AS redeemable_name,
    r.max_user_amount,
    (u.id IS NOT NULL)::boolean AS user_exists,
    COALESCE(a.status::text, '')::text AS application_status,
    COALESCE(
        u.checked_in_at IS NOT NULL
        OR EXISTS (SELECT 1 FROM nfc_tags_user ntu WHERE ntu.user_id = u.id),
        false
    )::boolean AS checked_in
FROM redeemables r
LEFT JOIN users u ON u.id = @user_id
LEFT JOIN applications a ON a.user_id = u.id AND a.hackathon_id = r.hackathon_id
WHERE r.id = @redeemable_id;

-- name: DeleteRedeemable :exec
DELETE FROM redeemables
WHERE id = $1;
