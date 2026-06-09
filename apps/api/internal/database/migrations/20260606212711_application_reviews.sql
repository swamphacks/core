-- +goose Up
alter table applications drop column if exists experience_rating;
alter table applications drop column if exists passion_rating;
alter table applications drop column if exists assigned_reviewer_id;

create table application_reviews (
    application_id uuid not null references applications(user_id) on delete cascade,
    reviewer_user_id uuid references users(id) on delete set null,
    experience_rating integer,
	passion_rating integer,
    notes text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    primary key (application_id, reviewer_user_id)
);

-- +goose Down
drop table application_reviews;

alter table applications
    add column experience_rating integer;

alter table applications
    add column passion_rating integer;

alter table applications 
    add column assigned_reviewer_id uuid references users on delete set null;