-- +goose Up
alter table nfc_tags_redeemables
drop column tag_id;

alter table nfc_tags_workshops
drop column tag_id;

alter table nfc_tags_redeemables
add column user_id uuid not null references users(id);

alter table nfc_tags_workshops
add column user_id uuid not null references users(id);

alter table nfc_tags_redeemables
add column tag_id text;

alter table nfc_tags_workshops
add column tag_id text;

drop trigger if exists check_redeemable_count_trigger on nfc_tags_redeemables;
drop function if exists check_redeemable_count();

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
      and user_id = new.user_id;

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
alter table nfc_tags_redeemables
drop column tag_id;

alter table nfc_tags_workshops
drop column tag_id;

alter table nfc_tags_redeemables
add column tag_id text not null references nfc_tags_user(tag_id) on delete cascade;

alter table nfc_tags_workshops
add column tag_id text not null references nfc_tags_user(tag_id) on delete cascade;

alter table nfc_tags_redeemables
drop column user_id;

alter table nfc_tags_workshops
drop column user_id;

drop trigger if exists check_redeemable_count_trigger on nfc_tags_redeemables;
drop function if exists check_redeemable_count();

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
