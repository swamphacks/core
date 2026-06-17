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

-- name: ListAllApplications :many
SELECT a.user_id, a.status, a.created_at, a.submitted_at, a.application, a.is_early, u.name, u.image, u.email FROM applications a
LEFT JOIN users u ON u.id = a.user_id
WHERE a.hackathon_id = @hackathon_id AND 
    LOWER(u.name) LIKE LOWER('%' || COALESCE(sqlc.arg('search'), '') || '%')
    OR LOWER(u.email) LIKE LOWER('%' || COALESCE(sqlc.arg('search'), '') || '%')
ORDER BY a.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetApplicationsCount :one
SELECT COUNT(*) FROM applications;

-- name: GetApplicationFullDetailsByUserId :one
SELECT a.*, 
ar.experience_rating, ar.passion_rating, ar.notes, ar.updated_at AS review_updated_at,
reviewer.id AS reviewer_id,
reviewer.name AS reviewer_name, 
reviewer.image AS reviewer_image,
applicant.name AS applicant_name,
applicant.image AS applicant_image,
aadr.requested_decision,
aadr.id as auto_decision_request_id,
aadr.justification AS decision_justification,
aadr.approved AS decision_approved,
aadr.approved_or_denied_by,
aadr.created_at AS decision_request_created_at
FROM applications a
LEFT JOIN application_reviews ar ON ar.application_id = a.user_id
LEFT JOIN users AS reviewer ON ar.reviewer_user_id = reviewer.id
LEFT JOIN users AS applicant ON a.user_id = applicant.id
LEFT JOIN application_auto_decision_requests as aadr ON a.user_id = aadr.application_id
WHERE a.user_id = @user_id;

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
