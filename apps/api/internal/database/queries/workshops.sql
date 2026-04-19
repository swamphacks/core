-- name: GetWorkshop :one
SELECT w.* FROM workshops w
WHERE w.id = @workshop_id;

-- name: GetAllWorkshops :many
SELECT w.* FROM workshops w
WHERE w.start_time > CURRENT_TIMESTAMP
ORDER BY w.start_time ASC;

-- name: DeleteWorkshop :exec
DELETE FROM workshops
WHERE id = @workshop_id;

-- name: CreateWorkshop :one
INSERT INTO workshops (title, description, start_time, end_time, location, presenter)
VALUES (@title, @description, @start_time, @end_time, @location, @presenter)
RETURNING *;

-- name: RegisterUserForWorkshop :exec
INSERT INTO workshop_registrations (user_id, workshop_id)
VALUES (@user_id, @workshop_id);

-- name: UnregisterUserForWorkshop :exec
DELETE FROM workshop_registrations
WHERE user_id = @user_id AND workshop_id = @workshop_id;

-- name: IsUserRegistered :one
SELECT EXISTS (
    SELECT 1
    FROM workshop_registrations
    WHERE user_id = @user_id AND workshop_id = @workshop_id
);

-- name: GetWorkshopRegistrations :many
SELECT *
FROM workshop_registrations
WHERE workshop_id = @workshop_id;

-- name: IncrementWorkshopAttendees :exec
UPDATE workshops
SET curr_attendees = curr_attendees + 1
WHERE id = $1;

-- name: DecrementWorkshopAttendees :exec
UPDATE workshops
SET curr_attendees = GREATEST(curr_attendees - 1, 0)
WHERE id = $1;