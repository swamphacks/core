-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
DROP COLUMN application_review_finished
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
ADD COLUMN application_review_finished BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd
