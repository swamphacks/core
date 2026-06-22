-- name: AssignReviewerToApplications :exec
INSERT INTO application_reviews (
    application_id,
    reviewer_id
)
SELECT applications.id, @reviewer_id::uuid FROM applications
WHERE applications.id = ANY(@application_ids::uuid[])
ON CONFLICT DO NOTHING;

-- name: ListReviewsByReviewerId :many
SELECT 
    ar.*,
    applications.user_id
FROM application_reviews ar
JOIN applications ON applications.id = ar.application_id
WHERE reviewer_id = @reviewer_id
ORDER BY application_id ASC;

-- name: ListApplicationReviewersById :many
SELECT reviewer_id FROM application_reviews
WHERE application_id = @application_id;

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
  ON reviewer.id = ar.reviewer_id
GROUP BY
  reviewer.id;

-- name: GetReviewById :one
SELECT
    ar.*,
    aadr.requested_decision,
    aadr.id AS decision_request_id,
    aadr.justification AS decision_justification,
    aadr.approved AS decision_approved,
    aadr.approved_or_denied_by AS decision_approved_or_denied_by,
    aadr.created_at AS decision_request_created_at,
    applications.user_id,
    applications.application
FROM application_reviews AS ar
JOIN applications ON applications.id = ar.application_id
LEFT JOIN application_auto_decision_requests AS aadr ON aadr.application_id = ar.application_id
WHERE ar.id = @review_id;

-- name: UpdateApplicationReview :exec
UPDATE application_reviews
SET
    experience_rating = CASE WHEN @experience_rating_do_update::boolean THEN @experience_rating ELSE experience_rating END,
    passion_rating = CASE WHEN @passion_rating_do_update::boolean THEN @passion_rating ELSE passion_rating END,
    notes = CASE WHEN @notes_do_update::boolean THEN @notes ELSE notes END,
    updated_by = CASE WHEN @updated_by_do_update::boolean THEN @updated_by ELSE updated_by END,
    updated_at = NOW()
WHERE
    id = @id AND
    reviewer_id = @reviewer_id;

-- name: DeleteAllApplicationReviews :exec
DELETE FROM application_reviews;

-- name: RequestAutoDecision :one
INSERT INTO application_auto_decision_requests (application_id, reviewer_id, requested_decision, justification, approved, approved_or_denied_by) 
VALUES (@application_id, @reviewer_id, @requested_decision, @justification, @approved, @approved_or_denied_by) RETURNING *;

-- name: ListAutoDecisionRequests :many
SELECT
    aadr.*,
    reviewer.name AS reviewer_name,
    reviewer.image AS reviewer_image,
    approver.id AS approver_id,
    approver.name AS approver_name,
    approver.image AS approver_image,
    applications.user_id
FROM application_auto_decision_requests AS aadr
JOIN users AS reviewer
    ON reviewer.id = reviewer_id
JOIN applications
    ON applications.id = aadr.application_id
LEFT JOIN users AS approver
    ON approver.id = aadr.approved_or_denied_by
ORDER BY aadr.created_at DESC;

-- name: DeleteAutoDecisionRequest :exec
DELETE FROM application_auto_decision_requests
WHERE id = @id AND reviewer_id = @reviewer_id;

-- name: UpdateAutoDecisionRequest :exec
UPDATE application_auto_decision_requests
SET
    requested_decision = CASE WHEN @requested_decision_do_update::boolean AND @requested_decision <> '' THEN @requested_decision::application_auto_decision_type ELSE requested_decision END,
    justification = CASE WHEN @justification_do_update::boolean THEN @justification ELSE justification END,
    approved = CASE WHEN @approved_do_update::boolean THEN @approved ELSE approved END,
    approved_or_denied_by = CASE WHEN @approved_by_do_update::boolean THEN @approved_or_denied_by ELSE approved_or_denied_by END
WHERE id = @id AND reviewer_id = @reviewer_id;

-- name: DeleteAllAutoDecisionRequests :exec
DELETE FROM application_auto_decision_requests;