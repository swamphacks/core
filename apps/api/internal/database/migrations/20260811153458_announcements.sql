-- +goose Up

-- TYPES
create type announcement_source as enum (
    'discord',
    'manual'
);

-- TABLES
create table announcements
(
	id uuid default gen_random_uuid() not null primary key,
    hackathon_id text not null references hackathons(id) on delete cascade,

	title text not null,
	body text not null,
    source announcement_source not null,

    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null,
    updated_by_user_id uuid references users(id) on delete set null,
    expires_at timestamptz,

    constraint announcements_expires_at_after_creation
        check (
            expires_at is null
            or expires_at >= created_at
        ),

    constraint announcements_discord_source_has_no_user
    check (
        source != 'discord'::announcement_source
        or updated_by_user_id is null
    )
);

create index idx_announcements_hackathon_id
    on announcements (hackathon_id);

create index idx_announcements_created_at
    on announcements (created_at desc);

create table users_dismissed_announcements
(
    user_id uuid not null references users(id) on delete cascade,
    announcement_id uuid not null references announcements(id) on delete cascade,
    dismissed_at timestamptz default now() not null,

    primary key (user_id, announcement_id)
);

-- TRIGGERS
create trigger set_updated_at_announcements
    before update on announcements
    for each row
    execute procedure update_modified_column();

-- +goose Down

drop trigger if exists set_updated_at_announcements on announcements;
drop table if exists users_dismissed_announcements;
drop index if exists idx_announcements_created_at;
drop index if exists idx_announcements_hackathon_id;
drop table if exists announcements;
drop type if exists announcement_source;