-- +goose Up

-- TYPES
alter type user_role rename to role;

-- TABLES
create table api_keys
(
    id uuid default gen_random_uuid() not null primary key,
    name text not null,
    description text,
    role role default 'visitor'::role not null,
    created_at timestamptz default now() not null,
    expires_at timestamptz,
    secret_hash text not null unique,

    constraint api_key_expires_at_after_creation
        check (
            expires_at is null
            or expires_at >= created_at
        )
);

alter table sessions
    alter column user_id drop not null;
alter table sessions
    add column api_key_id uuid references api_keys(id) on delete cascade;
alter table sessions
    add constraint sessions_exactly_one_principal
        check (
            (user_id is not null and api_key_id is null)
            or (user_id is null and api_key_id is not null)
        );

-- +goose Down
alter table sessions
    drop constraint sessions_exactly_one_principal;

-- explicitly delete current api key sessions
delete from sessions
    where api_key_id is not null;

alter table sessions
    drop column api_key_id;
alter table sessions
    alter column user_id set not null;

drop table api_keys;

alter type role rename to user_role;