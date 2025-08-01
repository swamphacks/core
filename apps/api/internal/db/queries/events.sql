-- name: CreateEvent :one
-- TODO: allow optional parameters
INSERT INTO events (
    name,
    application_open, application_close,
    start_time, end_time,
    description, location, location_url, max_attendees,
    rsvp_deadline, decision_release, 
    website_url, is_published 
    -- saved_at is omitted due to having a default value (NOW) which will likely not be different if it was an optional param
) VALUES (
    @name,
    @application_open, @application_close,
    @start_time, @end_time,
    coalesce(sqlc.narg(description), NULL),
    coalesce(sqlc.narg(location), NULL),
    coalesce(sqlc.narg(location_url), NULL),
    coalesce(sqlc.narg(max_attendees), NULL::INT),     
    coalesce(sqlc.narg(rsvp_deadline), NULL::TIMESTAMPTZ), 
    coalesce(sqlc.narg(decision_release), NULL::TIMESTAMPTZ),
    coalesce(sqlc.narg(website_url), NULL),
    coalesce(sqlc.narg(is_published), NULL::BOOLEAN)
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
    
-- name: DeleteEventById :exec
-- TODO: return error when 0 rows are deleted
DELETE FROM events
WHERE id = $1;
