-- name: CreateEvent :one
INSERT INTO events (
    id, name, description,
    location, loction_url, max_attendees,
    application_open, application_close, rsvp_deadline, decision_release,
    start_time, end_time,
    website_url 
) VALUES (
    $1, $2, $3,
    $4, $5, $6,
    $7, $8, $9, $10,
    $11, $12,
    $13
)
RETURNING *;

-- name: GetEventByID: one
SELECT * FROM events
WHERE id = $1;

-- name: GetEventByLocation: many
SELECT * FROM events
WHERE location = $1;

-- name: UpdateEventName: exec
UPDATE events
SET name = $2
WHERE id = $1;

-- name: UpdateEventDescription: exec
UPDATE events
SET description = $2
WHERE id = $1;

-- name: UpdateEventLocation: exec
UPDATE events
SET location = $2
WHERE id = $1;

-- name: UpdateEventLocationUrl: exec
UPDATE events
SET location_url = $2
WHERE id = $1;

-- name: UpdateEventMaxAttendees: exec
UPDATE events
SET max_attendees = $2
WHERE id = $1;

-- name: UpdateEventApplicationOpen: exec
UPDATE events
SET application_open = $2
WHERE id = $1;

-- name: UpdateEventApplicationClose: exec
UPDATE events
SET application_close = $2
WHERE id = $1;

-- name: UpdateEventRsvpDeadline: exec
UPDATE events
SET rsvp_deadline = $2
WHERE id = $1;

-- name: UpdateEventDecisionRelease: exec
UPDATE events
SET decision_release = $2
WHERE id = $1;

-- name: UpdateEventStartTime: exec
UPDATE events
SET start_time = $2
WHERE id = $1;

-- name: UpdateEventEndTime: exec
UPDATE events
SET end_time = $2
WHERE id = $1;

-- name: UpdateEventIsPublished: exec
UPDATE events
SET is_published = $2
WHERE id = $1;

-- name: UpdateEventWebsiteUrl: exec
UPDATE events
SET website_url = $2
WHERE id = $1;
 
-- name: DeleteEvent: exec
DELETE FROM events
WHERE id = $1;
