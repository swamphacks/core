-- +goose Up
alter table users add column is_fake boolean not null default false;
alter table applications add column is_fake boolean not null default false;

-- +goose Down
alter table users drop column is_fake;
alter table applications drop column is_fake;