-- +goose Up
ALTER TYPE application_status
ADD VALUE IF NOT EXISTS 'confirmed';

-- +goose Down
