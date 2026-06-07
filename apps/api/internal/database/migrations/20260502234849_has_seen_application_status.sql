-- +goose Up
alter table users add has_seen_new_application_status boolean;

-- +goose Down
alter table users drop column has_seen_new_application_status;
