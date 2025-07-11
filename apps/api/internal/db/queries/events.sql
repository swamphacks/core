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

-- name: GetEventByLocation :many
SELECT * FROM events
WHERE location = $1;

-- name: UpdateEvent :exec
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
    id = @id::uuid;
    
-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;
