-- +goose Up
alter table applications drop column if exists experience_rating;
alter table applications drop column if exists passion_rating;
alter table applications drop column if exists assigned_reviewer_id;

create table application_reviews (
    id uuid default gen_random_uuid() not null primary key,

    application_id uuid not null references applications(id) on delete cascade,
    reviewer_id uuid not null references users(id) on delete cascade,

    experience_rating integer,
    passion_rating integer,
    notes text,

    updated_by uuid references users(id) on delete set null,

    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    UNIQUE(application_id, reviewer_id)
);

create type application_auto_decision_type as enum (
    'auto_accept',
    'auto_reject'
);

create table application_auto_decision_requests (
    id uuid default gen_random_uuid() not null primary key,

    application_id uuid not null references applications(id) on delete cascade,
    reviewer_id uuid not null references users(id) on delete cascade,

    requested_decision application_auto_decision_type not null,
    justification text,
    approved boolean,
    approved_or_denied_by uuid references users(id) on delete set null,

    updated_at timestamptz not null default now(),
    created_at timestamptz not null default now(),

    UNIQUE(application_id, reviewer_id)
);

create trigger application_auto_decision_requests_updated_at
before update on application_auto_decision_requests
for each row
execute function update_modified_column();

-- +goose Down
drop table application_auto_decision_requests;
drop type application_auto_decision_type;
drop trigger if exists application_auto_decision_requests_updated_at
    on application_auto_decision_requests;

drop table application_reviews;

alter table applications
    add column experience_rating integer;

alter table applications
    add column passion_rating integer;

alter table applications 
    add column assigned_reviewer_id uuid references users on delete set null;