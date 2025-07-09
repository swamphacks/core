-- name: AddEmail :one
-- Adds a new email to the mailing list for a specific user and event.
-- The unique constraint on (event_id, user_id) will prevent duplicates.
-- Returns the newly created email record.
INSERT INTO event_interest_submissions (
    event_id,
    email,
    source
) VALUES (
    $1, $2, $3
)
RETURNING *;