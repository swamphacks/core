-- name: CreateApplication :one
INSERT INTO applications (
    user_id, event_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetApplicationByUserAndEventID :one
SELECT * FROM applications
WHERE user_id = $1 AND event_id = $2;

-- name: UpdateApplication :exec
UPDATE applications
SET
    status = CASE WHEN @status_do_update::boolean THEN @status::application_status ELSE status END,
    application = CASE WHEN @application_do_update::boolean THEN @application::JSONB ELSE application END,
    submitted_at = CASE WHEN @submitted_at_do_update::boolean THEN @submitted_at::timestamptz ELSE submitted_at END,
    saved_at = CASE WHEN @saved_at_do_update::boolean THEN @saved_at::timestamptz ELSE saved_at END,
    assigned_reviewer_id = CASE WHEN @assigned_reviewer_id_do_update::boolean THEN @assigned_reviewer_id::UUID ELSE assigned_reviewer_id END,
    experience_rating = CASE WHEN @experience_rating_do_update::boolean THEN @experience_rating::INT ELSE experience_rating END,
    passion_rating = CASE WHEN @passion_rating_do_update::boolean THEN @passion_rating::INT ELSE passion_rating END
WHERE
    user_id = @user_id AND event_id = @event_id;

-- name: DeleteApplication :exec
DELETE FROM applications
WHERE user_id = $1 AND event_id = $2;

-- An application is considered "available" for an event if the application has a status of submitted and has not been reviewed yet.
-- For optimization purposes, we only select the application IDs.

-- name: ListAvailableApplicationsForEvent :many
SELECT user_id FROM applications
WHERE event_id = $1 
    AND status = 'submitted'
    AND experience_rating IS NULL
    AND passion_rating IS NULL
ORDER BY 
    user_id ASC;

-- name: AssignApplicationsToReviewer :exec
UPDATE applications
SET assigned_reviewer_id = @reviewer_id::uuid,
    status = 'under_review'
WHERE user_id = ANY(@application_ids::uuid[])
  AND event_id = @event_id::uuid;

-- name: ResetApplicationReviews :exec
UPDATE applications
SET assigned_reviewer_id = NULL,
    status = 'submitted',
    experience_rating = NULL,
    passion_rating = NULL
WHERE status != 'submitted'
    AND event_id = $1;