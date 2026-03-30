# Installation & Setup

The API runs as a Docker service alongside PostgreSQL and Redis. No local Go installation is required to run the stack.

## Prerequisites

- Docker (see [Getting Started](../getting-started.md))
- A Discord application for OAuth ([Discord Developer Portal](https://discord.com/developers/applications))

### Go and API Development Tools

If you are writing migrations, modifying database queries, or updating the OpenAPI spec, you'll need Go and three CLI tools.

**1. Install the latest Go** from [go.dev/dl](https://go.dev/dl/). Verify:

```bash
go version
```

**2. Install Huma:**

```bash
go get github.com/danielgtaylor/huma/v2
```

**2. Install goose and sqlc:**

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

All three are invoked via `make` targets in `apps/api/`. See [Migrations](migrations.md) and [OpenAPI](openapi.md) for usage.

## Environment

Copy the example file:

```bash
cp apps/api/.env.dev.example apps/api/.env.dev
```

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://postgres:postgres@postgres:5432/coredb` | PostgreSQL connection string (Docker network) |
| `DATABASE_URL_MIGRATION` | `postgres://postgres:postgres@localhost:5432/coredb` | Connection string used when running migrations from the host |
| `REDIS_URL` | `redis://redis:6379` | Redis connection string |
| `ALLOWED_ORIGINS` | — | Comma-separated list of allowed CORS origins |
| `AUTH_DISCORD_CLIENT_ID` | — | Discord OAuth application client ID |
| `AUTH_DISCORD_CLIENT_SECRET` | — | Discord OAuth application client secret |
| `AUTH_DISCORD_REDIRECT_URI` | `http://localhost:8080/auth/callback` | OAuth callback URL (must match Discord app settings) |
| `CORE_BUCKETS_USER_QRCODES_BASE_URL` | — | Base URL for the Cloudflare R2 QR code bucket |
| `COOKIE_DOMAIN` | `localhost` | Domain for session cookies |
| `COOKIE_SECURE` | `false` | Set to `true` in production (requires HTTPS) |
| `CLIENT_URL` | `http://localhost:5173` | Frontend origin, used for redirects |
| `MAX_ACCEPTED_APPLICATIONS` | `500` | Hard cap on accepted applications |
| `ACCEPT_FROM_WAITLIST_COUNT` | `50` | Number of applicants to pull from the waitlist per cycle |
| `ACCEPT_FROM_WAITLIST_PERIOD` | `@every 72h` | Cron-style interval for waitlist processing |

## Running

Start the API with its dependencies:

```bash
make api
# or
docker compose up api
```

The API is available at **http://localhost:8080**.

To also run the background workers:

```bash
make backend
```

## Migrations

Database migrations must be run manually on the development database.

```bash
cd apps/api
make migrate-up
```

> The `DATABASE_URL_MIGRATION` variable points to `localhost:5432` (host network) rather than the Docker internal hostname, so migrations work when run outside the container.
