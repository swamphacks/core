-- name: AddBatRun :one
INSERT INTO bat_runs DEFAULT VALUES RETURNING *;

-- name: GetBatRunById :one
SELECT *
FROM bat_runs
WHERE id = $1;

-- name: GetBatRuns :many
SELECT
    *
FROM bat_runs
ORDER BY created_at DESC;

-- name: UpdateBatRunById :exec
UPDATE bat_runs
SET
    accepted_applicants = CASE WHEN @accepted_applicants_do_update::boolean THEN @accepted_applicants ELSE accepted_applicants END,
    rejected_applicants = CASE WHEN @rejected_applicants_do_update::boolean THEN @rejected_applicants ELSE rejected_applicants END,
    status = CASE WHEN @status_do_update::boolean THEN @status ELSE status END,
    created_at = CASE WHEN @created_at_do_update::boolean THEN @created_at ELSE created_at END
WHERE
    id = @id::uuid
RETURNING *;

-- name: DeleteBatRunById :execrows
DELETE FROM bat_runs
WHERE id = $1;
