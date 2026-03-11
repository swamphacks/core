# CI/CD

All automation runs on **GitHub Actions**. Workflows live in `.github/workflows/` and fall into three categories: CI checks, development deployments, and production deployments.

---

## Overview

| Trigger | What runs |
|---------|-----------|
| `pull_request` → `master` | Lint and test workflows (API, web, discord-bot) |
| `push` → `dev` | Dev build-and-deploy workflows for changed services |
| `push` → `master` | Prod build-and-deploy workflows for changed services; docs deploy |
| `push` (any branch, `apps/api/**`) | `sqlc` CI checks |
| `workflow_dispatch` | Any workflow can be triggered manually |

All workflows run on `ubuntu-latest` runners. Docker images are published to **GitHub Container Registry** (`ghcr.io`) using `GITHUB_TOKEN` for authentication — no separate registry secret is required.

---

## CI Workflows

| Workflow file | Job name(s) | Trigger | What it does |
|---|---|---|---|
| `lint-api.yml` | `API Lint` | `pull_request` → `master`, paths `apps/api/**` | Runs `go mod tidy` then `golangci-lint` v2.1.6 |
| `lint-web.yml` | `Web Lint & Format` | `pull_request` → `master`, paths `apps/web/**` | Installs pnpm 10 / Node 22, runs ESLint and Prettier check |
| `lint-discord-bot.yml` | `Discord Bot Lint & Format` | `pull_request` → `master`, paths `apps/discord-bot/**` | Placeholder (linting not yet implemented) |
| `test-web.yml` | `Web unit tests` | `pull_request` → `master`, paths `apps/web/**` | Installs pnpm 10 / Node 22 / Playwright, runs `pnpm test` |
| `sqlc_ci.yml` | `diff`, `vet` | `push` (any branch), paths `apps/api/**` | `sqlc diff` to verify generated code is up to date; `sqlc vet` against a live Postgres 17 instance |
| `docs.yml` | `deploy` | `push` → `master`, paths `apps/docs/**`; `workflow_dispatch` | Installs `mkdocs-material`, runs `mkdocs gh-deploy --force` to publish to GitHub Pages |

---

## Development Deployment Workflows

All dev workflows trigger on `push` to the `dev` branch (scoped to relevant paths) and support `workflow_dispatch` for manual runs. Images are tagged `:dev` and pushed to `ghcr.io/<owner>/`.

| Workflow file | Image built | Compose service | Migration | SSH target |
|---|---|---|---|---|
| `dev-build-deploy-api.yml` | `core-api:dev` (from `apps/api/`) | `api-dev` | Yes — Goose against `DEV_DB_URL` | `API_HOST` |
| `dev-build-deploy-web.yml` | `core-web:dev` (`--target prod`, `linux/amd64`) | `web-dev` | No | `WEB_HOST` |
| `dev-build-deploy-bat-worker.yml` | `core-bat-worker:dev` (`apps/api/cmd/BAT_worker/Dockerfile`, `--target prod`) | `bat-worker-dev` | Yes — Goose against `DEV_DB_URL` | `API_HOST` |
| `dev-build-deploy-email-worker.yml` | `core-email-worker:dev` (`apps/api/cmd/email_worker/Dockerfile`, `--target prod`) | `email-worker-dev` | Yes — Goose against `DEV_DB_URL` | `API_HOST` |
| `dev-build-deploy-discord-bot.yml` | `core-discord-bot:dev` (multi-arch: `linux/amd64,linux/arm64`) | _(push only, no SSH deploy step)_ | No | — |
| `dev-deploy-asynqmon.yml` | _(no build — uses upstream `hibiken/asynqmon`)_ | `asynqmon-dev` | No | `API_HOST` |

The bat-worker, email-worker, and api workflows all watch the same paths (`apps/api/**`, `infra/docker-compose.api.yml`), so a single push to `apps/api/` triggers all three in parallel.

---

## Production Deployment Workflows

All prod workflows trigger on `push` to the `master` branch (scoped to relevant paths) and support `workflow_dispatch`. Images are tagged `:latest`.

| Workflow file | Image built | Compose service | Migration | SSH target |
|---|---|---|---|---|
| `prod-build-deploy-api.yml` | `core-api:latest` | `api` | Yes — Goose against `PROD_DB_URL` | `API_HOST` |
| `prod-build-deploy-web.yml` | `core-web:latest` (`--target prod`, `linux/amd64`) | `web` | No | `WEB_HOST` |
| `prod-build-deploy-bat-worker.yml` | `core-bat-worker:latest` (`apps/api/cmd/BAT_worker/Dockerfile`, `--target prod`) | `bat-worker` | Yes — Goose against `PROD_DB_URL` | `API_HOST` |
| `prod-build-deploy-email-worker.yml` | `core-email-worker:latest` (`apps/api/cmd/email_worker/Dockerfile`, `--target prod`) | `email-worker` | Yes — Goose against `DEV_DB_URL`* | `API_HOST` |
| `prod-deploy-asynqmon.yml` | _(no build — uses upstream `hibiken/asynqmon`)_ | `asynqmon` | No | `API_HOST` |

\* `prod-build-deploy-email-worker.yml` currently references `DEV_DB_URL` in its migration step — this appears to be a bug in the workflow.

---

## Caddy Deployment

`deploy-caddy.yml` is **manual only** (`workflow_dispatch`). It has no build step — Caddy runs from the `caddy-cloudflare` image already present on the droplet. The job SSHs into `API_HOST`, pulls the infra repo to `master`, fetches secrets from Infisical (`--env=master`, path `/api`), then pulls and recreates the `caddy` container via `docker-compose.api.yml`.

Run this workflow whenever the `Caddyfile` or Caddy configuration changes.

---

## Build and Deploy Pattern

Every service that builds a Docker image follows this three-job sequence:

```
build-and-push  →  run-migrations (API services only)  →  deploy
```

**1. `build-and-push`**

Logs in to GHCR using `GITHUB_TOKEN`, builds the image, and pushes it:

```bash
docker build -t ghcr.io/<owner>/core-api:dev ./apps/api
docker push ghcr.io/<owner>/core-api:dev
```

The web image specifies `--platform linux/amd64` and `--target prod`. The discord-bot image uses `docker buildx` for multi-arch (`linux/amd64,linux/arm64`).

**2. `run-migrations`** (API, BAT worker, email worker only)

Installs Goose from the official install script and runs all pending migrations:

```bash
goose -dir ./apps/api/internal/db/migrations postgres "$DB_URL" up
```

This job depends on `build-and-push` completing successfully before it runs.

**3. `deploy`**

Depends on both `build-and-push` and `run-migrations`. Uses `appleboy/ssh-action` to connect to the target droplet as `root`, then:

```bash
cd /root/core/infra
git fetch && git checkout dev && git reset --hard origin/dev && git pull

# Fetch secrets from Infisical and write to .env file
export INFISICAL_TOKEN=$(infisical login --method=universal-auth \
  --client-id='...' --client-secret='...' --silent --plain)

infisical export --token=$INFISICAL_TOKEN --env=dev \
  --format=dotenv --path="/api" --projectId='...' \
  > ./secrets/.env.dev.api

# Pull the new image and recreate the container
docker compose -f docker-compose.api.yml pull api-dev
docker compose -f docker-compose.api.yml up -d --no-deps --force-recreate api-dev
```

For production, `dev` is replaced with `master` (or `main` for the web droplet), `:dev` tags become `:latest`, and the Infisical `--env` is `prod`.

---

## Required GitHub Secrets

| Secret | Used by |
|--------|---------|
| `GITHUB_TOKEN` | All build workflows — authenticates `docker login` to GHCR (automatically provided) |
| `API_HOST` | All workflows that SSH into the API droplet |
| `API_PASSWORD` | SSH password for `root@API_HOST` |
| `WEB_HOST` | Workflows that SSH into the web droplet |
| `WEB_PASSWORD` | SSH password for `root@WEB_HOST` |
| `DEV_DB_URL` | Goose migration connection string for the dev database |
| `PROD_DB_URL` | Goose migration connection string for the production database |
| `INFISICAL_CLIENT_ID` | Infisical universal-auth client ID |
| `INFISICAL_CLIENT_SECRET` | Infisical universal-auth client secret |
| `INFISICAL_PROJECT_ID` | Infisical project ID used when exporting secrets |
