-- name: GetTag :one
SELECT * FROM nfc_tags_user WHERE tag_id = @tag_id;

-- name: DeleteTagByUserId :execrows
DELETE FROM nfc_tags_user WHERE user_id = @user_id;

-- name: GetMeals :many
SELECT name, id, max_user_amount FROM redeemables
WHERE type = 'meal'
Order by name;

-- name: GetTshirts :many
SELECT name, id, max_user_amount FROM redeemables
where type = 'tshirt'
Order by name;

-- name: GetRedeemablesByType :many
SELECT name, id, max_user_amount FROM redeemables
WHERE type = @redeemable_type
ORDER BY name, id;

-- name: TagToRedeemable :one
WITH redeemed AS (
    INSERT INTO nfc_tags_redeemables (user_id, tag_id, redeemable_id)
    VALUES ($1, $2, $3)
    RETURNING user_id, tag_id, redeemable_id
)
SELECT
    users.name,
    redeemables.name AS redeemable_name,
    GREATEST(
        redeemables.max_user_amount -
        (
            SELECT COUNT(*)::integer
            FROM nfc_tags_redeemables
            WHERE user_id = redeemed.user_id
              AND redeemable_id = redeemed.redeemable_id
        ) - 1,
        0
    )::integer AS items_remaining
FROM redeemed
JOIN nfc_tags_user ON nfc_tags_user.user_id = redeemed.user_id
JOIN users ON users.id = nfc_tags_user.user_id
JOIN redeemables ON redeemables.id = redeemed.redeemable_id;

-- name: CheckinUser :one
WITH checked_in AS (
    INSERT INTO nfc_tags_user (tag_id, user_id) VALUES
    ($1, $2)
    ON CONFLICT DO NOTHING
    RETURNING *
)
SELECT
    checked_in.id,
    checked_in.tag_id,
    checked_in.user_id,
    checked_in.created_at,
    checked_in.updated_at,
    users.name,
    COALESCE(applications.status::text, 'none') AS application_status
FROM checked_in
JOIN users ON users.id = checked_in.user_id
LEFT JOIN applications ON applications.user_id = users.id;

-- name: GetWorkshops :many
SELECT title as name, id FROM workshops
WHERE type = 'workshop'
ORDER BY title;

-- name: GetSocials :many
SELECT title as name, id FROM workshops
WHERE type = 'social'
ORDER BY title;

-- name: TagToWorkshop :one
WITH tagged AS (
    INSERT INTO nfc_tags_workshops (user_id, tag_id, workshop_id)
    VALUES ($1, $2, $3)
    ON CONFLICT (user_id, workshop_id) DO NOTHING
    RETURNING user_id, workshop_id
)
SELECT users.name, workshops.title AS workshop_name
FROM tagged
JOIN nfc_tags_user ON nfc_tags_user.user_id = tagged.user_id
JOIN users ON users.id = nfc_tags_user.user_id
JOIN workshops ON workshops.id = tagged.workshop_id;
