-- +goose Up
create type application_auto_decision_type as enum (
    'auto_accept',
    'auto_reject'
);

create table application_auto_decision_requests (
    id uuid default gen_random_uuid() not null primary key,
    application_id uuid not null references applications(user_id) on delete cascade,
    reviewer_user_id uuid references users(id) on delete set null,
    requested_decision application_auto_decision_type not null,
    justification text,
    approved boolean,
    approved_or_denied_by uuid references users(id) on delete set null,
    updated_at timestamptz not null default now(),
    created_at timestamptz not null default now()
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