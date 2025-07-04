-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- This migration creates the 'mailing_list_emails' table.
-- It enforces that the combination of an event_id and a user_id must be unique,
-- preventing duplicate entries for the same user in the same event mailing list.

CREATE TABLE mailing_list_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    user_id UUID NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- This constraint ensures that each user can only be added to an event's mailing list once.
    CONSTRAINT uq_event_user UNIQUE (event_id, user_id)
);

-- Create the trigger to automatically update the 'updated_at' column on row modification.
-- This trigger uses the function created in the previous migration.
CREATE TRIGGER set_updated_at_mailing_list_emails
BEFORE UPDATE ON mailing_list_emails
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
-- This safely drops the table to reverse the 'Up' migration.

-- To reverse, first drop the trigger that depends on the table.
DROP TRIGGER IF EXISTS set_updated_at_mailing_list_emails ON mailing_list_emails;

DROP TABLE IF EXISTS mailing_list_emails;