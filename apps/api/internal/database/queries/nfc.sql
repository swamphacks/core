-- name: getMeals :many
SELECT name, id, max_user_amount FROM redeemables
WHERE type = 'meal'
Order by name;

-- name: getTshirts :many
SELECT name, id, max_user_amount FROM redeemables
where type = 'tshirt'
Order by name;

-- name: tagToRedeemable :exec
INSERT INTO nfc_tags_redeemables (tag_id, redeemable_id) VALUES ($1, $2);

-- name: checkinUser :one
INSERT INTO nfc_tags_user (tag_id, user_id) VALUES
($1, $2)
ON CONFLICT DO NOTHING 
RETURNING *;

-- name: getWorkshops :many
SELECT title as name, id FROM workshops
WHERE type = 'workshop'
ORDER BY title;

-- name: getSocials :many
SELECT title as name, id FROM workshops
WHERE type = 'social'
ORDER BY title;

-- name: tagToWorkshop :execrows
INSERT INTO nfc_tags_workshops (tag_id, workshop_id)
VALUES ($1, $2)
ON CONFLICT (tag_id, workshop_id) DO NOTHING;
