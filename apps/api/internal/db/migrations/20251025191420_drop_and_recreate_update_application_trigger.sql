-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_applications ON applications;


CREATE OR REPLACE FUNCTION update_application_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = clock_timestamp();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to update application updates
CREATE TRIGGER set_updated_at_applications
BEFORE UPDATE ON applications
FOR EACH ROW
EXECUTE FUNCTION update_application_modified_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_applications ON applications;
DROP FUNCTION IF EXISTS update_application_modified_column;
-- +goose StatementEnd
