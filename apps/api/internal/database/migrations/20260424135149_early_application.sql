-- +goose Up
alter table applications add is_early boolean default false not null;

alter table hackathons add accept_early_applications boolean default false not null;
alter table hackathons add early_application_open timestamptz;
alter table hackathons add early_application_close timestamptz;

alter table hackathons add constraint early_application_order
    check ((not accept_early_applications) or (early_application_open < early_application_close));

alter table hackathons add constraint early_application_no_overlap
    check ((not accept_early_applications) or (early_application_close <= application_open));

alter table hackathons add constraint early_application_dates_not_null_if_accept_early_applications
    check ((not accept_early_applications) or
        (early_application_open is not null and early_application_close is not null)
    );

alter table hackathons add constraint early_application_both_or_neither
    check (
        (early_application_open is null and early_application_close is null) or
        (early_application_open is not null and early_application_close is not null)
    );

-- +goose Down
alter table applications drop column is_early;

alter table hackathons drop column accept_early_applications;
alter table hackathons drop column early_application_open;
alter table hackathons drop column early_application_close;

alter table hackathons drop constraint if exists early_application_order;
alter table hackathons drop constraint if exists early_application_no_overlap;
alter table hackathons drop constraint if exists early_application_dates_not_null_if_accept_early_applications;
alter table hackathons drop constraint if exists early_application_both_or_neither;