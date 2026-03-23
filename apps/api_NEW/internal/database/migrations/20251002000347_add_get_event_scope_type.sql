-- +goose Up
-- +goose StatementBegin
CREATE TYPE get_event_scope_type AS ENUM (
    'published',
    'scoped',
    'all'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE get_event_scope_type;
-- +goose StatementEnd
