-- name: ListActiveAnnouncements :many
-- returns all active announcements for one hackathon
SELECT *
FROM announcements
WHERE hackathon_id = @hackathon_id
    AND (expires_at IS NULL OR expires_at > now())
ORDER BY created_at DESC;

-- name: ListAnnouncements :many
-- returns all announcements for one hackathon
SELECT *
FROM announcements
WHERE hackathon_id = @hackathon_id
ORDER BY created_at DESC;

-- name: CreateAnnouncement :one
-- creates an announcement
INSERT INTO announcements (
    hackathon_id,
    title,
    body,
    source,
    updated_by_user_id,
    expires_at
) VALUES (
    @hackathon_id,
    @title,
    @body,
    @source::announcement_source,
    sqlc.narg(updated_by_user_id)::uuid,
    sqlc.narg(expires_at)
)
RETURNING *;

-- name: UpdateAnnouncement :one
-- edits an announcement
UPDATE announcements
SET
    title =
        CASE WHEN @title_do_update::boolean
        THEN @title
        ELSE title END,
    body =
        CASE WHEN @body_do_update::boolean
        THEN @body
        ELSE body END,
    updated_by_user_id =
        CASE WHEN @updated_by_user_id_do_update::boolean
        THEN sqlc.narg(updated_by_user_id)::uuid
        ELSE updated_by_user_id END,
    expires_at =
        CASE WHEN @expires_at_do_update::boolean
        THEN sqlc.narg(expires_at)
        ELSE expires_at END
WHERE id = @id::uuid
    AND hackathon_id = @hackathon_id
RETURNING *;