-- name: AddEmail :one
-- Adds a new email to the mailing list for a specific user and event.
-- The unique constraint on (event_id, user_id) will prevent duplicates.
-- Returns the newly created email record.
INSERT INTO mailing_list_emails (
    event_id,
    user_id,
    email
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateEmailByID :one
-- Updates the email address for a specific mailing list entry by its unique ID.
-- The 'updated_at' field will be automatically updated by the database trigger.
-- This is useful when you have the specific record ID, e.g., from a list view.
-- Returns the updated email record.
UPDATE mailing_list_emails
SET email = $2
WHERE id = $1
RETURNING *;

-- name: DeleteEmailByID :exec
-- Deletes an email from the mailing list by its unique primary key ID.
-- Use this when the client has fetched a list of emails and knows the specific ID to delete.
-- :exec specifies that this query does not return any rows.
DELETE FROM mailing_list_emails
WHERE id = $1;

-- name: DeleteEmailByUserAndEvent :exec
-- Deletes an email from the mailing list using the logical composite key (user_id and event_id).
-- This is often more practical for APIs, as the client is more likely to know these two IDs.
-- For example, an API endpoint could be: DELETE /events/{event_id}/subscribers/{user_id}
DELETE FROM mailing_list_emails
WHERE event_id = $1 AND user_id = $2;

-- name: GetEmailsByEvent :many
-- Retrieves all email records associated with a specific event_id.
-- Results are ordered by creation date.
SELECT * FROM mailing_list_emails
WHERE event_id = $1
ORDER BY created_at DESC;

-- name: GetEmailsByUser :many
-- Retrieves all email records associated with a specific user_id.
-- This would show all the event mailing lists a single user has joined.
SELECT * FROM mailing_list_emails
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetEmailByUserAndEvent :one
-- Retrieves a single, specific email record for a user within an event.
-- Useful for checking if a user is already subscribed before attempting an insert.
SELECT * FROM mailing_list_emails
WHERE event_id = $1 AND user_id = $2
LIMIT 1;
