-- +goose Up
ALTER TABLE interest_submissions ADD CONSTRAINT unique_emails UNIQUE (email);


-- +goose Down
ALTER TABLE interest_submissions DROP CONSTRAINT unique_emails;
