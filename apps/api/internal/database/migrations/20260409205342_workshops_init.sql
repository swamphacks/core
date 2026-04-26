-- +goose Up
-- +goose StatementBegin

-- TABLES
create table workshops
(
	id uuid default gen_random_uuid() not null primary key,
	title text not null,
	description text,
	start_time timestamptz not null,
	end_time timestamptz not null,
	curr_attendees integer default 0 not null,
	location text,
	presenter text,
	created_at timestamptz default now() not null,
	updated_at timestamptz default now() not null
);

create table workshop_registrations
(
    user_id uuid not null,
    workshop_id uuid not null,
    created_at timestamptz default now() not null,

    primary key (user_id, workshop_id),

    foreign key (workshop_id)
        references workshops(id)
        on delete cascade
);

-- TRIGGERS
create trigger set_updated_at_workshops
	before update
	on workshops
	for each row
	execute procedure update_modified_column();

-- +goose StatementEnd
-- +goose Down
drop trigger if exists set_updated_at_workshops on workshops;

drop table if exists workshop_registrations;

drop table if exists workshops;
