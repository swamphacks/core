-- +goose Up

-- Lets staff replace a hacker's lost tag in place:
--   UPDATE nfc_tags_user SET tag_id = '<new tag id>' WHERE tag_id = '<old tag id>';
-- The hacker's scans and workshop check-ins move to the new tag, so their
-- per-item limits and dashboard history are kept. Deleting the old
-- nfc_tags_user row instead would cascade those rows away and hand the
-- stock back.

-- Take every lock up front, children before nfc_tags_user, the same order
-- scans and workshop check-ins use. Dropping a foreign key also locks the
-- table it references, so altering one child at a time could deadlock with
-- a workshop scan that arrives between the two ALTERs. This also makes a
-- second, concurrent `goose up` (the API and email-worker prod workflows
-- both run one) wait here until this transaction commits; it then
-- re-applies harmlessly.
LOCK TABLE nfc_tags_redeemables, nfc_tags_workshops, nfc_tags_user IN ACCESS EXCLUSIVE MODE;

ALTER TABLE nfc_tags_redeemables
    DROP CONSTRAINT nfc_tags_redeemables_tag_id_fkey,
    ADD CONSTRAINT nfc_tags_redeemables_tag_id_fkey
        FOREIGN KEY (tag_id) REFERENCES nfc_tags_user (tag_id)
        ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE nfc_tags_workshops
    DROP CONSTRAINT nfc_tags_workshops_tag_id_fkey,
    ADD CONSTRAINT nfc_tags_workshops_tag_id_fkey
        FOREIGN KEY (tag_id) REFERENCES nfc_tags_user (tag_id)
        ON UPDATE CASCADE ON DELETE CASCADE;

-- +goose Down

LOCK TABLE nfc_tags_redeemables, nfc_tags_workshops, nfc_tags_user IN ACCESS EXCLUSIVE MODE;

ALTER TABLE nfc_tags_workshops
    DROP CONSTRAINT nfc_tags_workshops_tag_id_fkey,
    ADD CONSTRAINT nfc_tags_workshops_tag_id_fkey
        FOREIGN KEY (tag_id) REFERENCES nfc_tags_user (tag_id)
        ON DELETE CASCADE;

ALTER TABLE nfc_tags_redeemables
    DROP CONSTRAINT nfc_tags_redeemables_tag_id_fkey,
    ADD CONSTRAINT nfc_tags_redeemables_tag_id_fkey
        FOREIGN KEY (tag_id) REFERENCES nfc_tags_user (tag_id)
        ON DELETE CASCADE;
