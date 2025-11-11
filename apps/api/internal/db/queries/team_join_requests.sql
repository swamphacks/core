-- name: CreateJoinRequest :one
INSERT INTO team_join_requests (
    team_id, user_id, request_message
) VALUES (
    @team_id, @user_id, @request_message
)
RETURNING *;

-- name: GetJoinRequestByID :one
SELECT * FROM team_join_requests 
WHERE id = @id;

-- name: ListJoinRequestsByTeamIDAndStatus :many
SELECT *
FROM team_join_requests
WHERE team_id = @team_id AND status = @status::join_request_status
ORDER BY created_at DESC;

-- name: ListJoinRequestsByUserID :many
SELECT *
FROM team_join_requests
WHERE user_id = @user_id
ORDER BY created_at DESC;

-- name: UpdateJoinRequest :one
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

-- name: ApproveJoinRequest :one
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

-- name: RejectJoinRequest :one
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