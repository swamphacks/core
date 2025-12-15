-- name: AddResult :one
INSERT INTO bat_results (
    event_id,
    accepted_applicants,
    rejected_applicants,
    status,
    
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetResultsByEventId :many
SELECT
    id,
    accepted_applicants,
    rejected_applicants,
    status,
    created_at,
    completed_at
FROM bat_results
WHERE event_id = $1
ORDER BY created_at DESC;

-- name: DeleteResultById :execrows
DELETE FROM bat_results
WHERE id = $1;
