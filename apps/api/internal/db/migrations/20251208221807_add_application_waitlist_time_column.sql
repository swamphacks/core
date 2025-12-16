-- +goose Up
-- +goose StatementBegin
ALTER TABLE applications
ADD COLUMN waitlist_join_time TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE applications
DROP COLUMN waitlist_join_time;
-- +goose StatementEnd
