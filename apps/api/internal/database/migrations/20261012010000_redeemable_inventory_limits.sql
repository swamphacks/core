-- +goose Up

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION check_redeemable_count()
RETURNS trigger AS $$
DECLARE
    current_user_amount integer;
    current_stock_amount integer;
    stock_amount integer;
    max_usr_amt integer;
BEGIN
    PERFORM 1
    FROM nfc_tags_user
    WHERE tag_id = NEW.tag_id
    FOR UPDATE;

    SELECT amount, max_user_amount
    INTO stock_amount, max_usr_amt
    FROM redeemables
    WHERE id = NEW.redeemable_id
    FOR UPDATE;

    SELECT count(*) INTO current_user_amount
    FROM nfc_tags_redeemables
    WHERE redeemable_id = NEW.redeemable_id
      AND tag_id = NEW.tag_id;

    IF current_user_amount >= max_usr_amt THEN
        RAISE EXCEPTION 'User reached the limit for this redeemable'
            USING ERRCODE = 'P0001';
    END IF;

    SELECT count(*) INTO current_stock_amount
    FROM nfc_tags_redeemables
    WHERE redeemable_id = NEW.redeemable_id;

    IF current_stock_amount >= stock_amount THEN
        RAISE EXCEPTION 'No redeemables left in stock'
            USING ERRCODE = 'P0005';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION check_redeemable_count()
RETURNS trigger AS $$
DECLARE
    current_amt integer;
    max_amt integer;
BEGIN
    PERFORM 1 FROM nfc_tags_user WHERE tag_id = NEW.tag_id FOR UPDATE;

    SELECT max_user_amount INTO max_amt
    FROM redeemables
    WHERE id = NEW.redeemable_id;

    SELECT count(*) INTO current_amt
    FROM nfc_tags_redeemables
    WHERE redeemable_id = NEW.redeemable_id
      AND tag_id = NEW.tag_id;

    IF current_amt >= max_amt THEN
        RAISE EXCEPTION 'Max redeems reached for this item. max: (%), current: (%)', max_amt, current_amt;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
