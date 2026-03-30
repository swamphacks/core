-- name: AddEmail :one
-- Adds a new email to the mailing list for a specific user.
-- The unique constraint on `email` will prevent duplicates.
-- Returns the newly created email record.
INSERT INTO interest_submissions (
    email,
    source,
    hackathon_id
) VALUES (
    $1, $2, $3
)
RETURNING *;