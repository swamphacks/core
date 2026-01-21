-- name: GetRedeemablesByEventID :many
-- Using event id, return all redeemables associated. Should also return data like how many have been redeemed already
SELECT r.id, 
r.event_id, 
r.name, 
r.amount AS total_stock, 
r.max_user_amount,
r.created_at, 
r.updated_at,
COALESCE(SUM(ur.amount), 0) AS total_redeemed
FROM redeemables r
LEFT JOIN user_redemptions ur ON r.id = ur.redeemable_id
WHERE r.event_id = $1
GROUP BY r.id;

-- name: RedeemRedeemable :one
-- Using user id and redeemable id, attempt! to redeem a redeemable
INSERT INTO user_redemptions (user_id, redeemable_id, amount)
SELECT $1, $2, 1
WHERE (
    SELECT COALESCE(SUM(amount), 0) 
    FROM user_redemptions 
    WHERE redeemable_id = $2
) < (SELECT amount FROM redeemables WHERE id = $2)
ON CONFLICT (user_id, redeemable_id) 
DO UPDATE SET 
    amount = user_redemptions.amount + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE user_redemptions.amount < (SELECT max_user_amount FROM redeemables WHERE id = $2)
RETURNING *;

-- name: GetRedemptionInfoByRedeemableID :many
-- Gather all redemption info for a specific reedeemable (who has redeemed already)
SELECT ur.user_id, ur.redeemable_id, ur.amount, ur.created_at, ur.updated_at
FROM user_redemptions ur
WHERE ur.redeemable_id = $1;

-- name: CreateRedeemable :one
-- Create a new redeemable for an event
INSERT INTO redeemables (event_id, name, amount, max_user_amount)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateRedeemable :one
UPDATE redeemables
SET 
    name = COALESCE(sqlc.narg('name'), name),
    amount = COALESCE(sqlc.narg('amount'), amount),
    max_user_amount = COALESCE(sqlc.narg('max_user_amount'), max_user_amount),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteRedeemable :exec
-- Delete a redeemable by id
DELETE FROM redeemables
WHERE id = $1;

-- name: UpdateRedemption :exec
-- Update a redemption record for a user and redeemable (for removing redemption mostly)
UPDATE user_redemptions
SET
    amount = $1
WHERE user_id = $2 AND redeemable_id = $3;


