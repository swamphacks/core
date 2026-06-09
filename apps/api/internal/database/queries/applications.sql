-- name: CreateApplication :one
INSERT INTO applications (user_id, hackathon_id, is_early) VALUES ($1, $2, $3) RETURNING *;

-- name: GetApplicationByUserId :one
SELECT * FROM applications WHERE user_id = $1;

-- name: UpdateApplicationByUserId :exec
UPDATE applications
SET
    status = CASE WHEN @status_do_update::boolean THEN @status::application_status ELSE status END,
    application = CASE WHEN @application_do_update::boolean THEN @application::JSONB ELSE application END,
    submitted_at = CASE WHEN @submitted_at_do_update::boolean THEN @submitted_at::timestamptz ELSE submitted_at END,
    saved_at = CASE WHEN @saved_at_do_update::boolean THEN @saved_at::timestamptz ELSE saved_at END,
    is_early = CASE WHEN @is_early_do_update::boolean THEN @is_early::boolean ELSE is_early END
WHERE
    user_id = @user_id;

-- name: DeleteApplicationByUserId :exec
DELETE FROM applications WHERE user_id = $1;

-- name: ListApplicationsUnderReviewWithTeamIds :many
SELECT a.user_id,
    a.application,
    t.id as team_id
FROM applications a
LEFT JOIN team_members tm
    ON tm.user_id = a.user_id
LEFT JOIN teams t
    ON t.id = tm.team_id
WHERE a.status = 'under_review';

-- name: WaitlistApplicationByUserId :exec
UPDATE applications
SET waitlist_join_time = COALESCE(waitlist_join_time, NOW()),
    status = 'waitlisted'
WHERE user_id = $1;

-- name: UpdateApplicationsStatusByUserIds :exec
UPDATE applications
SET status = @status::application_status
WHERE user_id = ANY(@user_ids::uuid[]);

-- name: WaitlistAcceptedApplications :exec
UPDATE applications
SET waitlist_join_time = COALESCE(waitlist_join_time, NOW()),
    status = 'waitlisted'
WHERE status = 'accepted'
  AND user_id IN (
    SELECT id from users
    WHERE role = 'applicant'
);

-- name: AcceptWaitlistedApplications :many
UPDATE applications
SET waitlist_join_time = NULL,
    status = 'accepted'
WHERE user_id IN (
  SELECT user_id FROM applications
  WHERE status = 'waitlisted'
  ORDER BY waitlist_join_time ASC
  LIMIT @acceptanceCount::int
)
RETURNING user_id;
