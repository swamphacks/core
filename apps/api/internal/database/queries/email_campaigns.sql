-- name: CreateEmailCampaign :one
-- creates a draft campaign. It stores the title, subject, body, format, recipient groups, and optional schedule time.
INSERT INTO email_campaigns (
    hackathon_id,
    title,
    description,
    subject,
    body,
    format,
    recipient_types,
    scheduled_at,
    created_by_user_id,
    updated_by_user_id
) VALUES (
    @hackathon_id,
    @title,
    sqlc.narg(description),
    @subject,
    @body,
    @format::email_campaign_format,
    @recipient_types::email_recipient_type[],
    sqlc.narg(scheduled_at),
    sqlc.narg(created_by_user_id),
    sqlc.narg(updated_by_user_id)
)
RETURNING *;

-- name: GetEmailCampaignByID :one
-- fetches one campaign by ID, scoped to a hackathon so one event cannot accidentally read another event’s campaign.
SELECT *
FROM email_campaigns
WHERE id = @id::uuid
    AND hackathon_id = @hackathon_id;

-- name: ListEmailCampaigns :many
-- returns all campaigns for one hackathon, newest first.
SELECT *
FROM email_campaigns
WHERE hackathon_id = @hackathon_id
ORDER BY created_at DESC;

-- name: UpdateEmailCampaign :one
-- edits draft-like campaign fields: title, description, subject, body, format, recipients, and scheduled time.
UPDATE email_campaigns
SET
    title = 
        CASE WHEN @title_do_update::boolean 
        THEN @title 
        ELSE title END,
    description = 
        CASE WHEN @description_do_update::boolean
        THEN @description  
        ELSE description END,
    subject = 
        CASE WHEN @subject_do_update::boolean 
        THEN @subject 
        ELSE subject END,
    body = 
        CASE WHEN @body_do_update::boolean
        THEN @body
        ELSE body END,
    format = 
        CASE WHEN @format_do_update::boolean
        THEN @format::email_campaign_format
        ELSE format END,
    recipient_types = 
        CASE WHEN @recipient_types_do_update::boolean
        THEN @recipient_types::email_recipient_type[]
        ELSE recipient_types END,
    scheduled_at = 
        CASE WHEN @scheduled_at_do_update::boolean
        THEN sqlc.narg(scheduled_at)
        ELSE scheduled_at END,
    updated_by_user_id = 
        CASE WHEN @updated_by_user_id_do_update::boolean
        THEN @updated_by_user_id::uuid
        ELSE updated_by_user_id END
WHERE id = @id::uuid
    AND hackathon_id = @hackathon_id
RETURNING *;

-- name: UpdateEmailCampaignStatus :one
-- changes lifecycle fields like draft -> scheduled, scheduled -> sending, sending -> sent, or sending -> failed.
UPDATE email_campaigns
SET
    status = @status::email_campaign_status,
    scheduled_at =
        CASE WHEN @scheduled_at_do_update::boolean 
        THEN sqlc.narg(scheduled_at)
        ELSE scheduled_at END,
    sent_at = 
        CASE WHEN @sent_at_do_update::boolean
        THEN sqlc.narg(sent_at)
        ELSE sent_at END,
    last_error = 
        CASE WHEN @last_error_do_update::boolean
        THEN @last_error
        ELSE last_error END,
    updated_by_user_id = 
        CASE WHEN @updated_by_user_id_do_update::boolean
        THEN @updated_by_user_id::uuid
        ELSE updated_by_user_id END
WHERE id = @id::uuid
    AND hackathon_id = @hackathon_id
RETURNING *;