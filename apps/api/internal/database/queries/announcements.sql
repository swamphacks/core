-- name: ListAnnouncements :many
-- returns all announcements for one hackathon
SELECT *
FROM announcements
WHERE hackathon_id = @hackathon_id
ORDER BY created_at DESC;

-- name: ListActiveAnnouncements :many
-- returns all active announcements for one hackathon
SELECT *
FROM announcements
WHERE hackathon_id = @hackathon_id
    AND (expires_at IS NULL OR expires_at > now())
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

-- name: DeleteAnnouncement :exec
-- deletes an announcement
DELETE FROM announcements WHERE id = @id;

-- name: DismissAnnouncement :one
-- dismisses an announcement
INSERT INTO users_dismissed_announcements (
    user_id,
    announcement_id
) VALUES (
    @user_id::uuid,
    @id::uuid
)
ON CONFLICT (user_id, announcement_id)
DO UPDATE SET
    dismissed_at = now()
RETURNING announcement_id;

-- name: ListDismissedAnnouncements :many
-- returns all dismissed announcements for a user
SELECT uda.announcement_id
FROM users_dismissed_announcements uda
JOIN announcements a
    ON a.id = uda.announcement_id
WHERE uda.user_id = @user_id
    AND (uda.dismissed_at >= a.updated_at);