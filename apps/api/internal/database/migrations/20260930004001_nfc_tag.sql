-- +goose Up

-- TABLES
drop table if exists nfc_tags_user cascade;
drop table if exists nfc_tags_redeemables cascade;
drop table if exists nfc_tags_workshops cascade;


ALTER TABLE redeemables ADD COLUMN type text NOT NULL CHECK (type IN ('meal', 'tshirt'));

ALTER TABLE workshops ADD COLUMN type text NOT NULL CHECK (type IN ('workshop', 'social'));

create table nfc_tags_user
(
    id uuid default gen_random_uuid() not null primary key,
    tag_id text not null unique,
    user_id uuid unique references users(id) on delete set null,
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
    updated_at timestamptz default now() not null,
	unique (tag_id, workshop_id)
);
-- TRIGGERS
-- +goose StatementBegin
create or replace function check_redeemable_count()
returns trigger as $$
declare
    current_amt int;
    max_amt int;
begin
    perform 1 from nfc_tags_user where tag_id = new.tag_id for update;

    select max_user_amount into max_amt
    from redeemables where id = new.redeemable_id;

    select count(*) into current_amt
    from nfc_tags_redeemables
    where redeemable_id = new.redeemable_id
      and tag_id = new.tag_id;

    if current_amt >= max_amt then
        raise exception 'Max redeems reached for this item. max: (%), current: (%)', max_amt, current_amt;
    end if;

    return new;
end;
$$ language plpgsql;
-- +goose StatementEnd

create trigger check_redeemable_count_trigger
before insert on nfc_tags_redeemables
for each row execute function check_redeemable_count();

-- +goose Down
drop table nfc_tags_redeemables;
drop table nfc_tags_workshops;
drop table nfc_tags_user;
drop table redeemables;
drop table workshops;
create table redeemables
(
	id uuid default gen_random_uuid() not null primary key,
	name varchar(255) not null,
	amount integer not null constraint redeemables_amount_check check (amount >= 0),
	max_user_amount integer not null constraint redeemables_max_user_amount_check check (max_user_amount >= 1),
	created_at timestamptz default now() not null,
	updated_at timestamptz default now() not null,
	hackathon_id text not null references hackathons(id)
);

create table workshops
(
	id uuid default gen_random_uuid() not null primary key,
	title text not null,
	description text,
	start_time timestamptz not null,
	end_time timestamptz not null,
	num_attendees integer default 0 not null,
	location text,
	presenter text,
	created_at timestamptz default now() not null,
	updated_at timestamptz default now() not null
);

drop function if exists check_redeemable_count();
