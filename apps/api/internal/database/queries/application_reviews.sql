-- name: AssignReviewerToApplications :exec
INSERT INTO application_reviews (
    application_id,
    reviewer_user_id
)
SELECT applications.user_id, @reviewer_id::uuid
FROM applications
WHERE applications.user_id = ANY(@application_ids::uuid[])
ON CONFLICT DO NOTHING;

-- name: ListReviewsByReviewerId :many
SELECT application_id, experience_rating, passion_rating FROM application_reviews
WHERE reviewer_user_id = @reviewer_id::uuid
ORDER BY application_id ASC;

-- name: ListReviewersByApplicationId :many
SELECT reviewer_user_id FROM application_reviews
WHERE application_id = @application_id::uuid;

-- name: ListUnderReviewApplicationIds :many
SELECT user_id FROM applications
WHERE status = 'under_review'
ORDER BY user_id ASC;

-- name: ListReviewersAndProgress :many
SELECT
  reviewer.id,
  reviewer.name,
  reviewer.image,
  COUNT(*) AS total_assigned,
  COUNT(*) FILTER (
    WHERE ar.experience_rating IS NOT NULL AND ar.passion_rating IS NOT NULL
  ) AS completed_count
FROM application_reviews AS ar
LEFT JOIN users AS reviewer
  ON reviewer.id = ar.reviewer_user_id
GROUP BY
  reviewer.id;

-- name: GetApplicationReviewDetails :one
-- SELECT ar.application_id, ar.passion_rating, ar.experience_rating, 
-- ar.notes, a.application, aadr.requested_decision, aadr.id as auto_decision_request_id FROM application_reviews as ar
-- LEFT JOIN applications as a ON ar.application_id = a.user_id
-- LEFT JOIN application_auto_decision_requests as aadr ON ar.application_id = aadr.application_id
-- WHERE ar.application_id = @application_id::uuid;
SELECT ar.application_id, ar.passion_rating, ar.experience_rating, 
ar.notes, a.application, aadr.requested_decision, aadr.id as auto_decision_request_id FROM applications as a
LEFT JOIN application_reviews as ar ON ar.application_id = a.user_id
LEFT JOIN application_auto_decision_requests as aadr ON aadr.application_id = a.user_id
WHERE a.user_id = @application_id::uuid;

-- name: UpdateApplicationReview :exec
UPDATE application_reviews
SET
    experience_rating = CASE WHEN @experience_rating_do_update::boolean THEN @experience_rating ELSE experience_rating END,
    passion_rating = CASE WHEN @passion_rating_do_update::boolean THEN @passion_rating ELSE passion_rating END,
    notes = CASE WHEN @notes_do_update::boolean THEN @notes ELSE notes END,
    updated_at = NOW()
WHERE
    reviewer_user_id = @reviewer_id::uuid AND 
    application_id = @application_id::uuid;

-- name: DeleteAllApplicationReviews :exec
DELETE FROM application_reviews;

-- name: ResetApplicationsToSubmitted :exec
UPDATE applications 
SET status = 'submitted'
WHERE status = 'under_review';

-- name: MarkSubmittedApplicationAsUnderReview :exec
UPDATE applications 
SET status = 'under_review'
WHERE status = 'submitted';

-- name: RequestAutoDecision :exec
INSERT INTO application_auto_decision_requests 
(application_id, reviewer_user_id, requested_decision, justification, approved, approved_or_denied_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: ListAutoDecisionRequests :many
SELECT
  aadr.*,
  reviewer.id AS reviewer_id,
  reviewer.name AS reviewer_name,
  reviewer.image AS reviewer_image,
  approver.id AS approver_id,
  approver.name AS approver_name
FROM application_auto_decision_requests AS aadr
LEFT JOIN users AS reviewer
  ON reviewer.id = aadr.reviewer_user_id
LEFT JOIN users AS approver
  ON approver.id = aadr.approved_or_denied_by
ORDER BY aadr.created_at;

-- name: DeleteAutoDecisionRequest :exec
DELETE FROM application_auto_decision_requests 
WHERE id = @request_id::uuid AND reviewer_user_id = @reviewer_id::uuid;

-- name: UpdateAutoDecisionRequest :exec
UPDATE application_auto_decision_requests
SET
    requested_decision = CASE WHEN @requested_decision_do_update::boolean AND @requested_decision <> '' THEN @requested_decision::application_auto_decision_type ELSE requested_decision END,
    justification = CASE WHEN @justification_do_update::boolean THEN @justification ELSE justification END,
    approved = CASE WHEN @approved_do_update::boolean THEN @approved ELSE approved END,
    approved_or_denied_by = CASE WHEN @approved_by_do_update::boolean THEN @approved_or_denied_by ELSE approved_or_denied_by END
WHERE
    id = @request_id::uuid;

-- name: DeleteAllAutoDecisionRequests :exec
DELETE FROM application_auto_decision_requests;