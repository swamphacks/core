-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_applications ON applications;

CREATE TRIGGER set_updated_at_applications
BEFORE UPDATE ON applications
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- +goose StatementEnd
