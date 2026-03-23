-- +goose Up
ALTER TABLE event_interest_submissions RENAME TO interest_submissions;

-- +goose Down
ALTER TABLE interest_submissions RENAME TO event_interest_submissions;

