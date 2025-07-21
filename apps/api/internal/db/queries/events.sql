-- name: CreateEvent :one
INSERT INTO events (
    name,
    application_open, application_close,
    start_time, end_time
) VALUES (
    $1,
    $2, $3,
    $4, $5
)
RETURNING *;

-- name: GetEventByID :one
SELECT * FROM events
WHERE id = $1;

-- name: UpdateEventById :exec
UPDATE events
SET
    name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    location = coalesce(sqlc.narg('location'), location),
    location_url = coalesce(sqlc.narg('location_url'), location_url),
    max_attendees = coalesce(sqlc.narg('max_attendees'), max_attendees),
    application_open = coalesce(sqlc.narg('application_open'), application_open),
    application_close = coalesce(sqlc.narg('application_close'), application_close),
    rsvp_deadline = coalesce(sqlc.narg('rsvp_deadline'), rsvp_deadline),
    decision_release = coalesce(sqlc.narg('decision_release'), decision_release),
    start_time = coalesce(sqlc.narg('start_time'), start_time),
    end_time = coalesce(sqlc.narg('end_time'), end_time),
    website_url = coalesce(sqlc.narg('website_url'), website_url),
    is_published = coalesce(sqlc.narg('is_published'), is_published),
    saved_at = coalesce(sqlc.narg('saved_at'), is_published)
WHERE
    id = @id::uuid;
    
-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;
