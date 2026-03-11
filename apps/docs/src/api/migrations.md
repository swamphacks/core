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
goose -dir internal/db/migrations create <description> sql
```

This creates a timestamped file in `internal/db/migrations/`. Write your `Up` and `Down` SQL, then commit the file.

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

## Migration history

| File | Description |
|---|---|
| `20250512145328_auth_init` | Auth schema: `users`, `accounts`, `sessions` tables |
| `20250608015747_remove_session_token_and_update_users` | Removed token column, updated users schema |
| `20250608194039_add_roles_to_auth_users` | Added `role` column + `auth_user_role` enum |
| `20250619161938_event_schema` | `events` and `event_roles` tables |
| `20250621222955_create_applications` | `applications` table + `application_status` enum |
| `20250627064521_create_mailing_list` | `event_interest_submissions` table |
| `20250825013123_remove_resume_url_column` | Removed `resume_url` from applications |
| `20250825033231_update_application_trigger` | Updated application `updated_at` trigger |
| `20250905210705_add_preferred_email_column` | Added `preferred_email` to users |
| `20250908055118_add_email_consent_column` | Added `email_consent` to users |
| `20250917044049_add_event_banner` | Added `banner_url` to events |
| `20251002000347_add_get_event_scope_type` | Added event scope enum for query filtering |
| `20251022032904_submitted_by_column_applications` | Added `submitted_by` to applications |
| `20251101211753_teams_table` | `teams` and `team_members` tables |
| `20251106160823_saved_at_trigger` | Added `saved_at` trigger to applications |
| `20251111212010_invitations_and_join_requests` | `team_invitations` and `team_join_requests` tables |
| `20251121165836_add_app_review_columns` | Added reviewer rating columns + `application_review_started` |
| `20251208221807_add_application_waitlist_time_column` | Added `waitlisted_at` to applications |
| `20251215225937_add_bat_runs_schema` | `bat_runs` table + `bat_run_status` enum |
| `20251216200020_add_application_review_finished` | Added review finished flag |
| `20251217065809_remove_application_review_finished_column` | Removed review finished flag |
| `20260116002956_checked_in_time` | Added check-in timestamp to event roles |
| `20260119015108_create_redeemables_tables` | `redeemables` and `user_redemptions` tables |

## SQLc regeneration

After modifying a migration or a query file in `internal/db/queries/`, regenerate the Go code:

```bash
cd apps/api
make generate
```

This updates `internal/db/sqlc/` — never edit that directory manually. See `sqlc.yml` for the full codegen configuration.
