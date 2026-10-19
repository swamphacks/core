-- +goose Up

-- The check-in dashboard and check_redeemable_count both look up scans by
-- tag and by item; without these they seq-scan the whole table.

-- SHARE ROW EXCLUSIVE conflicts with itself, so a second, concurrent
-- `goose up` (the API and email-worker prod workflows both run one) waits
-- here until this transaction commits, and IF NOT EXISTS then skips.
-- CREATE INDEX alone only takes SHARE, which does not, and the second
-- run would fail on the duplicate index name.
LOCK TABLE nfc_tags_redeemables IN SHARE ROW EXCLUSIVE MODE;

CREATE INDEX IF NOT EXISTS nfc_tags_redeemables_tag_id_redeemable_id_idx
    ON nfc_tags_redeemables (tag_id, redeemable_id);

CREATE INDEX IF NOT EXISTS nfc_tags_redeemables_redeemable_id_idx
    ON nfc_tags_redeemables (redeemable_id);

-- +goose Down

DROP INDEX IF EXISTS nfc_tags_redeemables_redeemable_id_idx;

DROP INDEX IF EXISTS nfc_tags_redeemables_tag_id_redeemable_id_idx;
