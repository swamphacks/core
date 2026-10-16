-- +goose Up

-- Stock counts aren't known for the event, and what matters is the
-- per-hacker limit, so redeemables no longer track stock. Redemptions are
-- still counted (one nfc_tags_redeemables row per scan).

-- Runs first on purpose: the ACCESS EXCLUSIVE lock on redeemables makes a
-- second, concurrent `goose up` (the API and email-worker prod workflows
-- both run one) wait here until this transaction commits. IF EXISTS then
-- lets it skip.
ALTER TABLE redeemables DROP COLUMN IF EXISTS amount;

-- Per-hacker limit only. Locking the tag's nfc_tags_user row serializes
-- scans of the same tag, so two quick taps can't both pass the limit.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION check_redeemable_count()
RETURNS trigger AS $$
DECLARE
    current_user_amount integer;
    max_usr_amt integer;
BEGIN
    PERFORM 1
    FROM nfc_tags_user
    WHERE tag_id = NEW.tag_id
    FOR UPDATE;

    SELECT max_user_amount
    INTO max_usr_amt
    FROM redeemables
    WHERE id = NEW.redeemable_id;

    SELECT count(*) INTO current_user_amount
    FROM nfc_tags_redeemables
    WHERE redeemable_id = NEW.redeemable_id
      AND tag_id = NEW.tag_id;

    IF current_user_amount >= max_usr_amt THEN
        RAISE EXCEPTION 'User reached the limit for this redeemable'
            USING ERRCODE = 'P0001';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down

-- The old stock values are gone. Every item comes back with a stock of
-- 1000000 (effectively unlimited); set real values before relying on it.
ALTER TABLE redeemables
    ADD COLUMN IF NOT EXISTS amount integer NOT NULL DEFAULT 1000000
        CONSTRAINT redeemables_amount_check CHECK (amount >= 0);

ALTER TABLE redeemables ALTER COLUMN amount DROP DEFAULT;

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
