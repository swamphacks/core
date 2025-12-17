-- name: AddRun :one
INSERT INTO bat_runs (
    event_id
) VALUES (
    $1
) RETURNING *;

-- name: GetRunById :one
SELECT *
FROM bat_runs
WHERE id = $1;

-- name: GetRunsByEventId :many
SELECT
    id,
    accepted_applicants,
    rejected_applicants,
    status,
    created_at,
    completed_at
FROM bat_runs
WHERE event_id = $1
ORDER BY created_at DESC;

-- name: UpdateRunById :exec
UPDATE bat_runs
SET
    accepted_applicants = CASE WHEN @accepted_applicants_do_update::boolean THEN @accepted_applicants ELSE accepted_applicants END,
    rejected_applicants = CASE WHEN @rejected_applicants_do_update::boolean THEN @rejected_applicants ELSE rejected_applicants END,
    status = CASE WHEN @status_do_update::boolean THEN @status ELSE status END,
    created_at = CASE WHEN @created_at_do_update::boolean THEN @created_at ELSE created_at END
WHERE
    id = @id::uuid
RETURNING *;

-- name: DeleteRunById :execrows
DELETE FROM bat_runs
WHERE id = $1;
