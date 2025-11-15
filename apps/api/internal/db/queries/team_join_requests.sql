-- name: CreateTeamJoinRequest :one
INSERT INTO team_join_requests (
    team_id, user_id, request_message
) VALUES (
    @team_id, @user_id, @request_message
)
RETURNING *;

-- name: GetTeamJoinRequestByID :one
SELECT * FROM team_join_requests 
WHERE id = @id;

-- name: ListTeamJoinRequestsByTeamIDAndStatus :many
SELECT *
FROM team_join_requests
WHERE team_id = @team_id AND status = @status::join_request_status
ORDER BY created_at DESC;

-- name: ListTeamJoinRequestsByUserID :many
SELECT *
FROM team_join_requests
WHERE user_id = @user_id
ORDER BY created_at DESC;

-- name: UpdateTeamJoinRequest :one
UPDATE team_join_requests
SET
    request_message = CASE WHEN @request_message_do_update::boolean THEN @request_message ELSE request_message END,
    status = CASE WHEN @status_do_update::boolean THEN @status::join_request_status ELSE status END,
    processed_by_user_id = CASE WHEN @processed_by_user_id_do_update::boolean THEN @processed_by_user_id::uuid ELSE processed_by_user_id END,
    processed_at = CASE WHEN @processed_at_do_update::boolean THEN @processed_at::timestamptz ELSE processed_at END,
    updated_at = NOW()
WHERE
    id = @id::uuid
RETURNING *;

-- name: ApproveTeamJoinRequest :one
UPDATE team_join_requests
SET
    status = 'APPROVED'::join_request_status,
    processed_by_user_id = @processed_by_user_id::uuid,
    processed_at = NOW(),
    updated_at = NOW()
WHERE
    id = @id AND 
    status = 'PENDING'::join_request_status
RETURNING *;

-- name: RejectTeamJoinRequest :one
UPDATE team_join_requests
SET
    status = 'REJECTED'::join_request_status,
    processed_by_user_id = @processed_by_user_id::uuid,
    processed_at = NOW(),
    updated_at = NOW()
WHERE
    id = @id AND 
    status = 'PENDING'::join_request_status
RETURNING *;

-- name: DeleteJoinRequest :exec
DELETE FROM team_join_requests 
WHERE id = @id;

-- name: ListTeamJoinRequestsByUserAndEventAndStatus :many
SELECT tjr.*
FROM team_join_requests tjr
WHERE tjr.user_id = $1
  AND tjr.status = $2
  AND EXISTS (
      SELECT 1
      FROM teams t
      WHERE t.id = tjr.team_id
        AND t.event_id = $3
  )
ORDER BY tjr.created_at DESC;
