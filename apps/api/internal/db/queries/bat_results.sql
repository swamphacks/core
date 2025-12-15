-- name: AddResult :one
INSERT INTO bat_results (
    event_id,
    accepted_applicants,
    rejected_applicants
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetResultsByEventId :many
SELECT
    accepted_applicants,
    rejected_applicants,
    created_at
FROM bat_results
WHERE event_id = @event_id::uuid
ORDER BY created_at DESC;
