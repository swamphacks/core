-- +goose Up
-- The limit trigger serialised concurrent redemptions by locking the hacker's
-- nfc_tags_user row through new.tag_id. Redemptions scanned by QR code on the
-- web portal have no tag, so tag_id is NULL, the lock matched nothing, and two
-- quick scans could both pass the count check. Lock the users row instead,
-- which every redemption has. FOR NO KEY UPDATE still lets foreign-key checks
-- against users (FOR KEY SHARE) through.

-- +goose StatementBegin
create or replace function check_redeemable_count()
returns trigger as $$
declare
    current_amt int;
    max_amt int;
begin
    perform 1 from users where id = new.user_id for no key update;

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

-- +goose Down
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
