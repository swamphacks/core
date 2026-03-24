-- +goose Up
ALTER TABLE hackathon_config RENAME TO hackathon;

-- +goose Down
ALTER TABLE hackathon RENAME TO hackathon_config;

