-- name: CreateEvent :one
INSERT INTO events (
    name,
    application_open, application_close,
    start_time, end_time,
    description, location, location_url, max_attendees,
    rsvp_deadline, decision_release, 
    website_url, is_published 
) VALUES (
    -- FIXME: The second parameter in coalesce MUST be the default value created in the schema. I have not found a more automated way to insert the default value.
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
    coalesce(sqlc.narg(is_published), FALSE)
) 
RETURNING *;

-- name: GetEventByID :one
SELECT * FROM events
WHERE id = $1;

-- name: UpdateEventById :exec
UPDATE events
SET
    name = CASE WHEN @name_do_update::boolean THEN @name ELSE name END,
    description = CASE WHEN @description_do_update::boolean THEN @description ELSE description END,
    location = CASE WHEN @location_do_update::boolean THEN @location ELSE location END,
    location_url = CASE WHEN @location_url_do_update::boolean THEN @location_url ELSE location_url END,
    max_attendees = CASE WHEN @max_attendees_do_update::boolean THEN @max_attendees ELSE max_attendees END,
    application_open = CASE WHEN @application_open_do_update::boolean THEN @application_open ELSE application_open END,
    application_close = CASE WHEN @application_close_do_update::boolean THEN @application_close ELSE application_close END,
    rsvp_deadline = CASE WHEN @rsvp_deadline_do_update::boolean THEN @rsvp_deadline ELSE rsvp_deadline END,
    decision_release = CASE WHEN @decision_release_do_update::boolean THEN @decision_release ELSE decision_release END,
    start_time = CASE WHEN @start_time_do_update::boolean THEN @start_time ELSE start_time END,
    end_time = CASE WHEN @end_time_do_update::boolean THEN @end_time ELSE end_time END,
    website_url = CASE WHEN @website_url_do_update::boolean THEN @website_url ELSE website_url END,
    is_published = CASE WHEN @is_published_do_update::boolean THEN @is_published ELSE is_published END
WHERE
    id = @id::uuid
RETURNING *;
    
-- name: DeleteEventById :execrows
-- execrows returns affect row count on top of an error
DELETE FROM events
WHERE id = $1;

-- name: GetEventRoleByIds :one
SELECT * FROM event_roles
WHERE user_id = @user_id::uuid AND event_id = @event_id::uuid;

-- name: GetPublishedEvents :many
SELECT * FROM events
WHERE is_published = TRUE;

-- name: GetAllEvents :many
SELECT * FROM events;
