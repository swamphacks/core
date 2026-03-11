# Database Testing

The API does not currently have automated unit or integration tests. This page covers how to work with the database during local development — inspecting data, running ad-hoc queries, and verifying migrations.

## Connecting to the local database

When the stack is running via Docker, PostgreSQL is exposed on `localhost:5432`.

**Connection string:**
```
postgres://postgres:postgres@localhost:5432/coredb
```

Connect with `psql`:

```bash
psql postgres://postgres:postgres@localhost:5432/coredb
```

Or use a GUI client (TablePlus, DBeaver, DataGrip) with the same credentials.

## Useful queries

**Check applied migrations:**
```sql
SELECT version_id, is_applied, tstamp
FROM goose_db_version
ORDER BY id DESC;
```

**Inspect user sessions:**
```sql
SELECT id, user_id, expires_at, updated_at AS last_used_at
FROM auth.sessions
ORDER BY updated_at DESC
LIMIT 20;
```

**View applications by status:**
```sql
SELECT status, COUNT(*)
FROM applications
GROUP BY status
ORDER BY count DESC;
```

**Check BAT run results:**
```sql
SELECT id, status, array_length(accepted_applicants, 1) AS accepted,
       array_length(rejected_applicants, 1) AS rejected, created_at
FROM bat_runs
ORDER BY created_at DESC;
```

## Resetting the database

To wipe and recreate the local database (useful after a bad migration or schema experiment):

```bash
# Stop and remove the postgres container + its volume
docker compose down -v

# Restart — migrations run automatically on startup
make local
```

!!! warning
    `docker compose down -v` deletes all persisted data including the `postgres_data` volume. Only do this locally.



## Inspecting the task queue

[Asynqmon](http://localhost:6767) provides a web UI for inspecting queued, active, and failed tasks when running `make local` or `make backend`. Use it to verify that email and BAT tasks are being enqueued and processed correctly.

## Working with sqlc

All SQL queries are in `internal/db/queries/`. After modifying a query or migration, regenerate the Go bindings:

```bash
cd apps/api
make generate
```

The generated code in `internal/db/sqlc/` reflects the current schema and query set. If `make generate` fails, the SQL query is invalid against the current schema — fix the query or migration before proceeding.
