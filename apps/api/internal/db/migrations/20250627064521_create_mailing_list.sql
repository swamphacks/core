-- +goose Up
-- +goose StatementBegin
CREATE TABLE emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Trigger for updated_at
CREATE TRIGGER set_updated_at_emails
BEFORE UPDATE ON emails
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE emails;
-- +goose StatementEnd
