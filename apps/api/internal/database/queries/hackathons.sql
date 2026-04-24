-- name: CreateHackathon :one
INSERT INTO hackathons (
    id, name,
    application_open, application_close,
    accept_early_applications, early_application_open, early_application_close,
    start_time, end_time,
    description, location, location_url, max_attendees,
    rsvp_deadline, decision_release, is_active
) VALUES (
    -- FIXME: The second parameter in coalesce MUST be the default value created in the schema. I have not found a more automated way to insert the default value.
    @id, @name,
    @application_open, @application_close,
    @accept_early_applications, 
    coalesce(sqlc.narg(early_application_open), NULL::TIMESTAMPTZ), 
    coalesce(sqlc.narg(early_application_close), NULL::TIMESTAMPTZ),
    @start_time, @end_time,
    coalesce(sqlc.narg(description), NULL),
    coalesce(sqlc.narg(location), NULL),
    coalesce(sqlc.narg(location_url), NULL),
    coalesce(sqlc.narg(max_attendees), NULL::INT),     
    coalesce(sqlc.narg(rsvp_deadline), NULL::TIMESTAMPTZ), 
    coalesce(sqlc.narg(decision_release), NULL::TIMESTAMPTZ),
    coalesce(sqlc.narg(is_active), false)
) 
RETURNING *;

-- name: UpdateHackathon :exec
UPDATE hackathons
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
    banner = CASE WHEN @banner_do_update::boolean THEN @banner ELSE banner END,
    application_review_started = CASE WHEN @application_review_started_do_update::boolean THEN @application_review_started ELSE application_review_started END
WHERE is_active = true
RETURNING *;
    
-- name: GetHackathon :one
SELECT * FROM hackathons WHERE is_active = true;

-- name: GetStaff :many
SELECT * FROM users
WHERE role IN ('admin', 'staff');

-- name: GetAttendeesWithDiscord :many
SELECT 
    a.account_id as discord_id,
    u.id as user_id,
    u.name,
    u.email
FROM users u
JOIN accounts a ON u.id = a.user_id
WHERE u.role = 'attendee'
    AND a.provider_id = 'discord';

-- name: GetAttendeeCount :one
SELECT COUNT(*) FROM users
WHERE role = 'attendee';

-- name: GetAttendeeUserIds :many
SELECT id FROM users
WHERE role = 'attendee';