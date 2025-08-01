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
    resume_url = CASE WHEN @resume_url_do_update::boolean THEN @resume_url ELSE resume_url END
WHERE
    user_id = @user_id AND event_id = @event_id;

-- name: DeleteApplication :exec
DELETE FROM applications
WHERE user_id = $1 AND event_id = $2;
