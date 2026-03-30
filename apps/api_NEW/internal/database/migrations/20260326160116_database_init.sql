-- +goose Up
-- +goose StatementBegin

-- TYPES
create type application_status as enum ('started', 'submitted', 'under_review', 'accepted', 'rejected', 'waitlisted', 'withdrawn');

create type team_invitation_status as enum ('pending', 'accepted', 'expired', 'rejected');

create type team_join_request_status as enum ('pending', 'approved', 'rejected');

create type bat_run_status as enum ('running', 'completed', 'failed');

create type user_role as enum ('admin', 'staff', 'attendee', 'applicant', 'visitor');

-- TABLES
create table users
(
	id uuid default gen_random_uuid() not null primary key,
    name text not null,
	email text unique,
	email_verified boolean default false not null,
	onboarded boolean default false not null,
	image text,
	created_at timestamptz default now() not null,
	updated_at timestamptz default now() not null,
	preferred_email text,
	email_consent boolean default false not null,
	checked_in_at timestamptz,
	rfid text,
	role_assigned_at timestamptz,
	role user_role default 'visitor'::user_role not null
);

create table accounts
(
	id uuid default gen_random_uuid() not null primary key,
	user_id uuid not null references users on delete cascade,
	provider_id text not null,
	account_id text not null,
	hashed_password text,
	access_token text,
	refresh_token text,
	id_token text,
	access_token_expires_at timestamptz,
	refresh_token_expires_at timestamptz,
	scope text,
	created_at timestamptz default now() not null,
	updated_at timestamptz default now() not null,
	unique (provider_id, account_id)
);

create table sessions
(
	id uuid default gen_random_uuid() not null primary key,
	user_id uuid not null references users on delete cascade,
	expires_at timestamptz not null,
	ip_address text,
	user_agent text,
	created_at timestamptz default now() not null,
	updated_at timestamptz default now() not null,
	last_used_at timestamptz default now() not null
);

create table hackathons
(
	id text not null primary key,
	name text not null,
	description text,
	location text,
	location_url text,
	max_attendees integer,
	application_open timestamptz not null,
	application_close timestamptz not null,
	rsvp_deadline timestamptz,
	decision_release timestamptz,
	start_time timestamptz not null,
	end_time timestamptz not null,
	is_active boolean default false not null,
	created_at timestamptz default now() not null,
	updated_at timestamptz default now() not null,
	banner text,
	application_review_started boolean default false not null
);

create table applications
(
	user_id uuid not null primary key references users on delete cascade,
	status application_status default 'started'::application_status not null,
	application jsonb default '{}'::jsonb not null,
	created_at timestamptz default now() not null,
	saved_at timestamptz default now() not null,
	updated_at timestamptz default now() not null,
	submitted_at timestamptz,
	experience_rating integer,
	passion_rating integer,
	assigned_reviewer_id uuid references users on delete set null,
	waitlist_join_time timestamptz,
	hackathon_id text not null references hackathons(id)
);

create table bat_runs
(
	id uuid default gen_random_uuid() not null primary key,
	accepted_applicants uuid[] default '{}'::uuid[],
	rejected_applicants uuid[] default '{}'::uuid[],
	status bat_run_status default 'running'::bat_run_status not null,
	created_at timestamptz default now() not null,
	completed_at timestamptz,
	hackathon_id text not null references hackathons(id)
);

create table interest_submissions
(
	id uuid default gen_random_uuid() not null primary key,
	email text unique not null,
	created_at timestamptz default now() not null,
	source text,
	hackathon_id text not null references hackathons(id)
);

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

create table teams
(
	id uuid default gen_random_uuid() not null primary key,
	name text not null,
	owner_id uuid references users on delete set null,
	created_at timestamptz default now() not null,
	updated_at timestamptz default now() not null,
	hackathon_id text not null references hackathons(id)
);

create table team_invitations
(
	id uuid default gen_random_uuid() not null primary key,
	team_id uuid not null references teams on delete cascade,
	invited_by_user_id uuid not null references users on delete cascade,
	invited_email text not null,
	invited_user_id uuid references users on delete cascade,
	status team_invitation_status default 'pending'::team_invitation_status not null,
	expires_at timestamptz,
	created_at timestamptz default now() not null,
	updated_at timestamptz default now() not null
);

create table team_join_requests
(
	id uuid default gen_random_uuid() not null primary key,
	team_id uuid not null references teams on delete cascade,
	user_id uuid not null references users on delete cascade,
	request_message text,
	status team_join_request_status default 'pending'::team_join_request_status not null,
	processed_by_user_id uuid references users on delete set null,
	processed_at timestamptz,
	created_at timestamptz default now() not null,
	updated_at timestamptz default now() not null
);

create table team_members
(
	user_id uuid not null references users on delete cascade,
	team_id uuid not null references teams on delete cascade,
	joined_at timestamptz default now() not null,
	primary key (user_id, team_id)
);

create table user_redemptions
(
	user_id uuid not null references users on delete cascade,
	redeemable_id uuid not null references redeemables on delete cascade,
	amount integer not null constraint user_redemptions_amount_check check (amount >= 0),
	created_at timestamptz default now() not null,
	updated_at timestamptz default now() not null,
	hackathon_id text not null references hackathons(id),
	primary key (user_id, redeemable_id)
);

-- INDEXES
create unique index only_one_hackathon_active on hackathons (is_active) where is_active = true;

create unique index idx_unique_pending_request
	on team_join_requests (team_id, user_id)
	where (status = 'pending'::team_join_request_status);

create index idx_applications_status
	on applications (status);

create index idx_accounts_user_id
	on accounts (user_id);

create index idx_accounts_provider_account
	on accounts (provider_id, account_id);

create index idx_sessions_user_id
	on sessions (user_id);

create index idx_sessions_expires_at
	on sessions (expires_at);

-- TRIGGERS
create or replace function update_modified_column()
returns TRIGGER as $$
begin
    NEW.updated_at = clock_timestamp();
    return NEW;
end;
$$ language plpgsql;

create trigger set_updated_at_accounts
	before update
	on accounts
	for each row
	execute procedure update_modified_column();

create trigger set_updated_at_users
	before update
	on users
	for each row
	execute procedure update_modified_column();

create trigger set_updated_at_teams
	before update
	on users
	for each row
	execute procedure update_modified_column();

create trigger set_updated_at_sessions
	before update
	on sessions
	for each row
	execute procedure update_modified_column();

create trigger set_updated_at_hackathon
	before update
	on hackathons
	for each row
	execute procedure update_modified_column();

create trigger set_updated_at_applications
	before update
	on applications
	for each row
	execute procedure update_modified_column();

create trigger set_updated_at_redeemables
	before update
	on redeemables
	for each row
	execute procedure update_modified_column();

create trigger set_updated_at_team_invitations
	before update
	on team_invitations
	for each row
	execute procedure update_modified_column();

create trigger set_updated_at_team_join_requests
	before update
	on team_join_requests
	for each row
	execute procedure update_modified_column();

create trigger set_updated_at_user_redemptions
	before update
	on user_redemptions
	for each row
	execute procedure update_modified_column();

-- +goose StatementEnd
-- +goose Down
drop index if exists only_one_hackathon_active;

drop index if exists idx_unique_pending_request;

drop index if exists idx_applications_status;

drop index if exists idx_accounts_user_id;

drop index if exists idx_accounts_provider_account;

drop index if exists idx_sessions_user_id;

drop index if exists idx_sessions_expires_at;

drop trigger if exists set_updated_at_accounts on accounts;

drop trigger if exists set_updated_at_users on users;

drop trigger if exists set_updated_at_teams on teams;

drop trigger if exists set_updated_at_sessions on sessions;

drop trigger if exists set_updated_at_hackathon on hackathon;

drop trigger if exists set_updated_at_applications on applications;

drop trigger if exists set_updated_at_redeemables on redeemables;

drop trigger if exists set_updated_at_team_invitations on invitations;

drop trigger if exists set_updated_at_team_join_requests on team_join_requests;

drop trigger if exists set_updated_at_user_redemptions on user_redemptions;

drop table user_redemptions;
drop table team_members;
drop table team_join_requests;
drop table team_invitations;
drop table teams;
drop table redeemables;
drop table interest_submissions;
drop table bat_runs;
drop table applications;
drop table hackathons;
drop table sessions;
drop table accounts;
drop table users;

drop type application_status;

drop type team_invitation_status;

drop type team_join_request_status;

drop type bat_run_status;

drop type user_role;

drop function update_modified_column;