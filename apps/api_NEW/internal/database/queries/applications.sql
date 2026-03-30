-- name: CreateApplication :one
INSERT INTO applications (user_id, hackathon_id) VALUES ($1, $2) RETURNING *;

-- name: GetApplicationByUserId :one
SELECT * FROM applications WHERE user_id = $1;

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
    user_id = @user_id;

-- name: DeleteApplication :exec
DELETE FROM applications WHERE user_id = $1;

-- An application is considered "available" if the application has a status of submitted and has not been reviewed yet.
-- For optimization purposes, we only select the application IDs.

-- name: ListAvailableApplications :many
SELECT user_id FROM applications
WHERE status = 'submitted'
    AND experience_rating IS NULL
    AND passion_rating IS NULL
ORDER BY 
    user_id ASC;

-- name: ListAdmissionCandidates :many
SELECT a.user_id,
    a.passion_rating,
    a.experience_rating,
    a.application,
    t.id as team_id
FROM applications a
LEFT JOIN team_members tm
    ON tm.user_id = a.user_id
LEFT JOIN teams t
    ON t.id = tm.team_id
WHERE a.status = 'under_review'
    AND a.passion_rating IS NOT NULL
    AND a.experience_rating IS NOT NULL;

-- name: AssignApplicationsToReviewer :exec
UPDATE applications
SET assigned_reviewer_id = @reviewer_id::uuid,
    status = 'under_review'
WHERE user_id = ANY(@application_ids::uuid[]);

-- name: ResetApplicationReviews :exec
UPDATE applications
SET assigned_reviewer_id = NULL,
    status = 'submitted',
    experience_rating = NULL,
    passion_rating = NULL
WHERE status NOT IN ('submitted', 'started');

-- name: ListApplicationByReviewer :many
SELECT user_id, passion_rating, experience_rating FROM applications
WHERE assigned_reviewer_id = $1
    AND status IN ('under_review')
ORDER BY user_id ASC;

-- name: ListNonReviewedApplications :many
SELECT user_id
FROM applications
WHERE status = 'under_review'
  AND (passion_rating IS NULL OR experience_rating IS NULL);

-- name: JoinWaitlist :exec
UPDATE applications
SET waitlist_join_time = COALESCE(waitlist_join_time, NOW()),
    status = 'waitlisted'
WHERE user_id = $1;

-- name: UpdateApplicationStatus :exec
UPDATE applications
SET status = @status::application_status
WHERE user_id = ANY(@user_ids::uuid[]);

-- name: TransitionAcceptedApplicationsToWaitlist :exec
UPDATE applications
SET waitlist_join_time = COALESCE(waitlist_join_time, NOW()),
    status = 'waitlisted'
WHERE status = 'accepted'
  AND user_id IN (
    SELECT id from users
    WHERE role = 'applicant'
);

-- name: TransitionWaitlistedApplicationsToAccepted :many
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

