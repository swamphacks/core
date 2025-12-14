-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
ADD COLUMN banner TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
DROP COLUMN banner;
-- +goose StatementEnd
