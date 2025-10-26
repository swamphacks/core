-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_applications ON applications;
DROP FUNCTION IF EXISTS update_application_modified_column;

-- +goose StatementEnd
