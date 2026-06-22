-- +goose Up
alter table applications drop column if exists waitlist_join_time;

alter table applications add id uuid not null default gen_random_uuid();

alter table applications drop constraint applications_pkey;

alter table applications add primary key (id);

alter table applications add constraint one_application_per_user unique (user_id, hackathon_id);

-- +goose Down
alter table applications add waitlist_join_time timestamptz;

alter table applications drop constraint one_application_per_user;

alter table applications drop constraint applications_pkey;

alter table applications drop column if exists id;

alter table applications add primary key (user_id);