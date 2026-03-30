# Migrations

The API uses [Goose](https://github.com/pressly/goose) for database migrations. Migration files live in `internal/db/migrations/` and are embedded into the binary — migrations run automatically on startup.

## File naming

Migration files follow the Goose timestamp convention:

```
YYYYMMDDHHMMSS_description.sql
```

Each file contains an `Up` and `Down` block:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE ...;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ...;
-- +goose StatementEnd
```

## Creating a migration

Generate a new migration file (requires [goose installed](installation.md#api-development-tools)):

```bash
cd apps/api
goose -dir internal/database/migrations create <description> sql
```

This creates a timestamped file in `internal/database/migrations/`. Write your `Up` and `Down` SQL, then commit the file.

!!! warning
    Never modify an existing migration file. If you need to alter a table, create a new migration.

## Running migrations

Migrations run automatically when the API starts. To run them manually from the host:

```bash
cd apps/api
make migrate-up
```

This uses `DATABASE_URL_MIGRATION` (pointing to `localhost:5432`) rather than `DATABASE_URL` (the Docker-internal hostname), so it works when run outside the container.

## Rolling back

To roll back the last migration:

```bash
cd apps/api
make migrate-down
```

## SQLC regeneration

After modifying a migration or a query file in `internal/database/queries/`, regenerate the Go code:

```bash
cd apps/api
make generate
```

This updates `internal/database/sqlc/` — never edit that directory manually. See `sqlc.yml` for the full codegen configuration.
