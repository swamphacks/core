# Secrets Management (Infisical)

[Infisical](https://infisical.com) is the centralised secrets manager for SwampHacks Core. No secret values are stored in the repository. All secrets live in Infisical and are pulled at deploy time.

---

## Why Infisical

- **No secrets in the repo.** `.env` files that contain real values are git-ignored. Only `.env.example` / `.env.dev.example` files (with blank or placeholder values) are committed.
- **Centralised management.** All secrets for all services and environments are in one place. Rotating a secret means updating it in Infisical once — the next deploy picks it up automatically.
- **Per-environment isolation.** Infisical organises secrets by environment (`dev`, `prod`) and by path (`/api`, `/web`). Each deploy workflow targets the correct combination.

---

## How Secrets Flow into the System

```
Infisical (source of truth)
        │
        │  pulled via `infisical export` using a machine-identity token
        ▼
  .env file written to infra/secrets/ on the droplet
        │
        │  mounted into Docker containers via docker-compose env_file
        ▼
  Running container has environment variables available
```

The GitHub Actions workflows drive this process. The deploy step of each workflow:

1. SSHes into the target droplet.
2. Authenticates with Infisical using a machine identity (`universal-auth`).
3. Exports the secrets for the target environment and path to a `.env` file inside `infra/secrets/`.
4. Pulls the new Docker image and recreates the container, which picks up the freshly written `.env` file.

---

## GitHub Actions Secrets

The following secrets must be configured in the GitHub repository (`Settings → Secrets and variables → Actions`) for the deploy workflows to function:

| Secret | Used by | Purpose |
|--------|---------|---------|
| `INFISICAL_CLIENT_ID` | all deploy workflows | Infisical machine identity client ID |
| `INFISICAL_CLIENT_SECRET` | all deploy workflows | Infisical machine identity client secret |
| `INFISICAL_PROJECT_ID` | all deploy workflows | Infisical project identifier |
| `API_HOST` | API workflows | SSH host for the API droplet |
| `API_PASSWORD` | API workflows | SSH password for the API droplet |
| `WEB_HOST` | web workflows | SSH host for the web droplet |
| `WEB_PASSWORD` | web workflows | SSH password for the web droplet |
| `DEV_DB_URL` | `dev-build-deploy-api.yml` | Postgres connection string for running dev migrations |
| `PROD_DB_URL` | `prod-build-deploy-api.yml` | Postgres connection string for running prod migrations |

`INFISICAL_CLIENT_ID`, `INFISICAL_CLIENT_SECRET`, and `INFISICAL_PROJECT_ID` are the only secrets that originate from Infisical itself. All other application secrets are stored in Infisical and never put directly in GitHub.

---

## The Deploy Pattern (Infisical Token Exchange)

Every deploy workflow follows the same two-step Infisical pattern on the server:

```bash
# Step 1 — obtain a short-lived access token using the machine identity
export INFISICAL_TOKEN=$(infisical login \
  --method=universal-auth \
  --client-id='${{ secrets.INFISICAL_CLIENT_ID }}' \
  --client-secret='${{ secrets.INFISICAL_CLIENT_SECRET }}' \
  --silent \
  --plain)

# Step 2 — export secrets for the target environment and path to a dotenv file
infisical export \
  --token=$INFISICAL_TOKEN \
  --env=dev \               # or prod
  --format=dotenv \
  --path="/api" \           # or /web
  --projectId='${{ secrets.INFISICAL_PROJECT_ID }}' \
  > ./secrets/.env.dev.api  # output file name varies per service/environment
```

The four output files written to `infra/secrets/` are:

| File | Environment | Service |
|------|-------------|---------|
| `secrets/.env.dev.api` | dev | API |
| `secrets/.env.api` | prod | API |
| `secrets/.env.dev.web` | dev | web |
| `secrets/.env.web` | prod | web |

---

## The fetch-web-secrets.sh Script

`infra/fetch-web-secrets.sh` is an alternative mechanism for fetching web secrets directly via the Infisical REST API rather than the CLI. It is intended for environments where the Infisical CLI is not installed (e.g. a CI runner that only has `curl` and `jq`).

**What it does:**

1. Requires `INFISICAL_CLIENT_ID` and `INFISICAL_CLIENT_SECRET` to be set in the environment. `WORKSPACE_SLUG` defaults to `swamphacks-core` and `ENVIRONMENT` defaults to `dev`.
2. Calls `POST https://us.infisical.com/api/v1/auth/universal-auth/login` to exchange the client credentials for a short-lived `accessToken`.
3. Calls `GET https://us.infisical.com/api/v3/secrets/raw` with `secretPath=/web` to retrieve all secrets (including any imported secret sets) and pipes them through `jq` to produce `KEY=VALUE` lines.
4. Writes the result to `./secrets/.env.dev.web`.

**When to use it:** The deploy workflows use the Infisical CLI directly (see above). This script is useful for one-off manual secret refreshes or in environments where the CLI cannot be installed.

---

## Local Development

Infisical is **not** required for local development. Each service ships an example env file with safe placeholder values. Copy it and fill in only what you need:

### API (`apps/api`)

```bash
cp apps/api/.env.dev.example apps/api/.env
```

Key variables:

| Variable | Default / Example | Notes |
|----------|-------------------|-------|
| `APP_ENV` | `local` | Set to `local` if running locally |
| `DATABASE_URL` | `postgres://postgres:postgres@postgres:5432/coredb` | Used inside the container |
| `DATABASE_URL_MIGRATION` | `postgres://postgres:postgres@localhost:5432/coredb` | Used from the host when running migrations directly |
| `REDIS_URL` | `redis://redis:6379` | |
| `ALLOWED_ORIGINS` | _(empty)_ | Comma-separated list of allowed CORS origins |
| `AUTH_DISCORD_CLIENT_ID` | _(empty)_ | Discord OAuth application client ID |
| `AUTH_DISCORD_CLIENT_SECRET` | _(empty)_ | Discord OAuth application client secret |
| `AUTH_DISCORD_REDIRECT_URI` | `http://localhost:8080/auth/callback` | |
| `CORE_BUCKETS_USER_QRCODES_BASE_URL` | _(empty)_ | Cloudflare R2 public base URL for QR code assets |
| `COOKIE_DOMAIN` | `localhost` | |
| `COOKIE_SECURE` | `false` | Set to `true` in production |
| `CLIENT_URL` | `http://localhost:5173` | |
| `MAX_ACCEPTED_APPLICATIONS` | `500` | Waitlist configuration |
| `ACCEPT_FROM_WAITLIST_COUNT` | `50` | |
| `ACCEPT_FROM_WAITLIST_PERIOD` | `@every 72h` | Cron-style period |
| `GRAFANA_URL` | `http://grafana:3000` | |
| `MONITORING_DISCORD_WEBHOOK` | _(empty)_ | Discord Webhook used to send Grafana alerts |

### Web (`apps/web`)

```bash
cp apps/web/.env.example apps/web/.env.local
```

| Variable | Example | Notes |
|----------|---------|-------|
| `VITE_BASE_API_URL` | `https://api.swamphacks.com` | Point to `http://localhost:8080` for local API |
| `VITE_DISCORD_OAUTH_CLIENT_ID` | _(empty)_ | Discord OAuth application client ID |
| `VITE_ALLOWED_HOSTS` | `[""]` | JSON array of allowed host strings |

### Discord Bot (`apps/discord-bot`)

```bash
cp apps/discord-bot/.env.example apps/discord-bot/.env
```

| Variable | Notes |
|----------|-------|
| `DISCORD_TOKEN` | Bot token from the Discord Developer Portal |
| `API_KEY` | Internal API key for authenticating against the Core API |
| `GEMINI_API_KEY` | Google Gemini API key |
| `API_URL` | Core API base URL (`https://api.swamphacks.com` or `http://localhost:8080`) |
| `SESSION_COOKIE` | Session cookie value for authenticated API calls |
| `WEBHOOK_URL` | Incoming webhook URL |
| `WEBHOOK_PORT` | Port the webhook listener binds to |
| `EVENT_ID` | Identifier for the active event |

---

## Adding a New Secret

Follow these steps in order:

1. **Add the secret in Infisical.** Navigate to the correct project, environment (`dev` / `prod`), and path (`/api`, `/web`, etc.) and create the secret.

2. **Add the GitHub Actions secret** (only if the secret is needed during the build or migration steps, not just at runtime). Go to `Settings → Secrets and variables → Actions` in the repository and add the secret. If it is only needed at runtime inside the container, skip this step — it will flow through Infisical at deploy time.

3. **Add the variable to the relevant `.env.example` file.** Add the key with an empty or placeholder value and a short comment describing its purpose. This is the in-repo documentation for what variables a service expects.

    ```bash
    # My new secret — obtained from the relevant third-party service dashboard
    MY_NEW_SECRET=
    ```

4. **Document it** in the relevant service's installation or configuration page (e.g. `api/`, `web/`, `discord-bot/` sections of these docs).

5. **Test locally** by adding the value to your local `.env` file. Do not commit the file with a real value.

---

## The `infra/secrets/` Directory

The `infra/secrets/` directory on the server is where all generated `.env` files land. It is git-ignored. The only file committed in that directory is `secrets/README.md`, which describes the expected contents.

At runtime the directory holds:

| File | Contents |
|------|----------|
| `secrets/.env.dev.api` | Dev API secrets (exported from Infisical `/api`, env `dev`) |
| `secrets/.env.api` | Prod API secrets (exported from Infisical `/api`, env `prod`) |
| `secrets/.env.dev.web` | Dev web secrets (exported from Infisical `/web`, env `dev`) |
| `secrets/.env.web` | Prod web secrets (exported from Infisical `/web`, env `prod`) |

These files are referenced by the Docker Compose files via `env_file` directives and are never pushed to the repository.
