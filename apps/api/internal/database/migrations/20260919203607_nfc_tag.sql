-- +goose Up

-- TABLES
create table nfc_tags_user
(
    id uuid default gen_random_uuid() not null primary key,
    tag_id text not null unique,
    user_id uuid references users(id) on delete set null,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null
);

create table nfc_tags_redeemables(
    
    id uuid default gen_random_uuid() not null primary key,
    tag_id text not null references nfc_tags_user(tag_id) on delete cascade,
    redeemable_id uuid not null references redeemables(id) on delete cascade,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null
);

create table nfc_tags_workshops(
    
    id uuid default gen_random_uuid() not null primary key,
    tag_id text not null references nfc_tags_user(tag_id) on delete cascade,
    workshop_id uuid not null references workshops(id) on delete cascade,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null
);

-- +goose Down
drop table nfc_tags_user;
drop table nfc_tags_redeemables;
drop table nfc_tags_workshops;
