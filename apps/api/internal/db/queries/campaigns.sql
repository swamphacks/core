-- name: CreateCampaign :one
INSERT INTO campaigns (
    event_id,
    title, description,
    recipient_roles,
    created_by
) VALUES (
    @event_id, @title,
    coalesce(sqlc.narg(description), null),
    coalesce(sqlc.narg(recipient_roles), null),
    @created_by
)
RETURNING *;

-- name: GetCampaignById :one
SELECT * FROM campaigns 
WHERE id = $1;

-- name: UpdateCampaignById :exec
UPDATE campaigns
SET
    title = CASE WHEN @title_do_update::boolean THEN @title ELSE title END,
    description = CASE WHEN @description_do_update::boolean THEN @description ELSE description END,
    recipient_roles = CASE WHEN @recipient_roles::boolean THEN @recipient_roles ELSE recipient_roles END,
    recipient_addresses = CASE WHEN @recipient_addresses::boolean THEN @recipient_addresses ELSE recipient_addresses END,
    emails = CASE WHEN @emails::boolean THEN @emails ELSE emails END
WHERE
    id = @id::uuid
returning *;

-- name: DeleteCampaignById :execrows
DELETE FROM campaigns
WHERE id = $1;