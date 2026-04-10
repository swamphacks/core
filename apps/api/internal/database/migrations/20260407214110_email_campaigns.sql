-- +goose Up
-- +goose StatementBegin

create type email_campaign_status as enum (
    'draft',
    'scheduled',
    'sending',
    'sent',
    'failed'
);

create type email_campaign_format as enum (
    'text',
    'html'
);

create type email_recipient_type as enum (
    'admins',
    'staff',
    'accepted_applicants',
    'rejected_applicants',
    'waitlisted_applicants',
    'visitors',
    'interest_subscribers'
);

create table email_campaigns (
    id uuid default gen_random_uuid() not null primary key,
    hackathon_id text not null references hackathons(id) on delete cascade,

    title text not null,
    description text,
    subject text not null,
    body text not null,
    format email_campaign_format default 'text'::email_campaign_format not null,

    recipient_types email_recipient_type[] not null,
    status email_campaign_status default 'draft'::email_campaign_status not null,

    scheduled_at timestamptz,
    sent_at timestamptz,
    last_error text,

    created_by_user_id uuid references users(id) on delete set null,
    updated_by_user_id uuid references users(id) on delete set null,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null,

    constraint email_campaigns_recipient_types_nonempty
        check (cardinality(recipient_types) > 0),

    constraint email_campaigns_scheduled_at_required
        check (
            status != 'scheduled'::email_campaign_status
            or scheduled_at is not null
        ),

    constraint email_campaigns_sent_at_required
        check (
            status != 'sent'::email_campaign_status
            or sent_at is not null
        )
);

create index idx_email_campaigns_hackathon_id
    on email_campaigns (hackathon_id);

create index idx_email_campaigns_status
    on email_campaigns (status);

create index idx_email_campaigns_created_at
    on email_campaigns (created_at desc);

create index idx_email_campaigns_scheduled_at
    on email_campaigns (scheduled_at)
    where status = 'scheduled'::email_campaign_status;

create trigger set_updated_at_email_campaigns
    before update on email_campaigns
    for each row
    execute procedure update_modified_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop trigger if exists set_updated_at_email_campaigns on email_campaigns;
drop table if exists email_campaigns;
drop type if exists email_recipient_type;
drop type if exists email_campaign_format;
drop type if exists email_campaign_status;

-- +goose StatementEnd
