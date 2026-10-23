-- +goose Up
-- Abort rather than discard attendance if existing duplicate rows need review.
CREATE UNIQUE INDEX IF NOT EXISTS nfc_tags_workshops_user_workshop_unique
ON nfc_tags_workshops(user_id,workshop_id);

-- +goose Down
DROP INDEX IF EXISTS nfc_tags_workshops_user_workshop_unique;
