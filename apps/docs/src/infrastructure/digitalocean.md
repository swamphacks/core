# DigitalOcean

SwampHacks Core runs on two plain DigitalOcean droplets — no managed databases, no App Platform, no Kubernetes. Each droplet runs Docker Compose directly. This setup is an amalgamation of decisions made over time; it works, but it is more manually operated than a fully managed solution would be.

---

## Droplets

| Droplet | Secret ref | What runs on it |
|---------|-----------|-----------------|
| **API droplet** | `API_HOST` | API (prod + dev), both background workers (prod + dev), Redis (prod + dev), Asynqmon (prod + dev), Caddy |
| **Web droplet** | `WEB_HOST` | Web frontend (prod + dev), Discord bot, Caddy |

SSH access uses `root` with a password stored in GitHub Actions secrets (`API_PASSWORD` / `WEB_PASSWORD`). There is no key-based auth in the current workflows.

---

## Environments

Both droplets run **prod** and **dev** side-by-side as separate containers within the same Compose stack:

- **Production** — triggered by pushes to `master`. Images tagged `:latest`. Secrets pulled from the `prod` Infisical environment.
- **Development** — triggered by pushes to `dev`. Images tagged `:dev`. Secrets pulled from the `dev` Infisical environment.

The root `docker-compose.yml` and `infra/docker-compose.dev.yml` are **not used on the droplets**. Those exist for local development only (bind mounts, local Postgres, source-built images). The droplet-specific files are `infra/docker-compose.api.yml` and `infra/docker-compose.web.yml`.

---

## What Runs Where

### API Droplet — `docker-compose.api.yml`

| Service | Image | Port |
|---------|-------|------|
| `api` | `ghcr.io/swamphacks/core-api:latest` | 8080 (host) |
| `api-dev` | `ghcr.io/swamphacks/core-api:dev` | 8081 (host) |
| `email-worker` | `ghcr.io/swamphacks/core-email-worker:latest` | — |
| `email-worker-dev` | `ghcr.io/swamphacks/core-email-worker:dev` | — |
| `bat-worker` | `ghcr.io/swamphacks/core-bat-worker:latest` | — |
| `bat-worker-dev` | `ghcr.io/swamphacks/core-bat-worker:dev` | — |
| `redis` | `redis:8.2.1-alpine` | 6379 (host) |
| `redis-dev` | `redis:8.2.1-alpine` | 6380 (host) |
| `asynqmon` | `hibiken/asynqmon:latest` | 6767 (host) |
| `asynqmon-dev` | `hibiken/asynqmon:latest` | 6768 (host) |
| `caddy` | `ghcr.io/caddybuilds/caddy-cloudflare:latest` | 80, 443 |

Redis prod is capped at 500 MB; Redis dev at 200 MB. Both use `allkeys-lru` eviction.

### Web Droplet — `docker-compose.web.yml`

| Service | Image | Network |
|---------|-------|---------|
| `web` | `ghcr.io/swamphacks/core-web:latest` | `caddy_net` (internal only, port 80) |
| `web-dev` | `ghcr.io/swamphacks/core-web:dev` | `caddy_net` (internal only, port 80) |
| `discord` | `ghcr.io/swamphacks/core-discord:latest` | default |
| `caddy` | `ghcr.io/caddybuilds/caddy-cloudflare:latest` | `caddy_net`, ports 80 / 443 |

Web containers are not bound to host ports. Caddy reaches them over the shared `caddy_net` Docker bridge network.

---

## Docker Images and Registry

All application images are built by GitHub Actions and pushed to the **GitHub Container Registry** (GHCR) under `ghcr.io/swamphacks/`:

| Image | Tags |
|-------|------|
| `ghcr.io/swamphacks/core-api` | `:latest`, `:dev` |
| `ghcr.io/swamphacks/core-web` | `:latest`, `:dev` |
| `ghcr.io/swamphacks/core-email-worker` | `:latest`, `:dev` |
| `ghcr.io/swamphacks/core-bat-worker` | `:latest`, `:dev` |
| `ghcr.io/swamphacks/core-discord` | `:latest`, `:dev` |

Caddy (`ghcr.io/caddybuilds/caddy-cloudflare`) and Redis (`redis:8.2.1-alpine`) are pulled directly from their public registries. GHCR authentication in workflows uses `GITHUB_TOKEN` — no additional credentials are required.

---

## Deployment Process

Each service has a dedicated workflow. The naming convention is:

```
.github/workflows/
  dev-build-deploy-api.yml
  dev-build-deploy-web.yml
  prod-build-deploy-api.yml
  prod-build-deploy-web.yml
  prod-build-deploy-bat-worker.yml
  prod-build-deploy-email-worker.yml
  deploy-caddy.yml
```

### Standard deploy steps (API/worker workflows)

1. **Build and push** — The image is built on the Actions runner and pushed to GHCR.
2. **Run migrations** — Goose migrations are applied against the target database from the runner using `PROD_DB_URL` or `DEV_DB_URL`.
3. **Deploy** — The runner SSHes into the droplet via `appleboy/ssh-action` and runs:

```bash
cd /root/core/infra
git fetch
git checkout <branch>        # master or dev, depending on workflow
git reset --hard origin/<branch>
git pull

export INFISICAL_TOKEN=$(infisical login \
  --method=universal-auth \
  --client-id='...' \
  --client-secret='...' \
  --silent \
  --plain)

infisical export \
  --token=$INFISICAL_TOKEN \
  --env=<prod|dev> \
  --format=dotenv \
  --path="/api" \
  --projectId='...' \
  > ./secrets/.env.api          # or .env.dev.api / .env.web / .env.dev.web

docker compose -f docker-compose.api.yml pull <service>
docker compose -f docker-compose.api.yml up -d --no-deps --force-recreate <service>
```

The `--no-deps --force-recreate` flags ensure only the targeted service is restarted; other running containers are left alone.

### Web deploy

The web workflows follow the same pattern but target `WEB_HOST`/`WEB_PASSWORD` and use `docker-compose.web.yml`.

### Key differences between infra Compose files and the local dev compose

| Aspect | Local (`docker-compose.yml`) | Droplet (`docker-compose.api.yml` / `.web.yml`) |
|--------|------------------------------|-------------------------------------------------|
| Images | Built from source | Pre-built, pulled from GHCR |
| Env vars | Local `.env` files or defaults | Written to `infra/secrets/` at deploy time by Infisical |
| Bind mounts | Yes (source code) | None — containers are entirely image-based |
| Postgres | Local container | External managed database (URL from secrets) |
| Caddy | Not included | Runs on the droplet, manages TLS |

---

## Caddy

Caddy handles TLS termination and reverse proxying on both droplets. It uses the `caddy-cloudflare` image, which bundles the Cloudflare DNS plugin for ACME DNS-01 challenges. The Cloudflare API token is loaded from `infra/secrets/.env.cf`.

### API droplet — `Caddyfile.api`

| Host | Upstream |
|------|----------|
| `api.swamphacks.com` | `api:8080` |
| `dev-api.swamphacks.com` | `api-dev:8080` |
| `asynqmon.swamphacks.com` | `asynqmon:6767` |
| `dev-asynqmon.swamphacks.com` | `asynqmon-dev:6767` |

### Web droplet — `Caddyfile.web`

| Host | Upstream |
|------|----------|
| `app.swamphacks.com` | `web:80` |
| `dev-app.swamphacks.com` | `web-dev:80` |

All virtual hosts set `Strict-Transport-Security`, `X-Content-Type-Options`, `X-Frame-Options`, and `Referrer-Policy` headers, and enable gzip/zstd compression.

### Updating Caddy

The `deploy-caddy.yml` workflow is **manual dispatch only** (`workflow_dispatch`). It SSHes into the API droplet, pulls the latest `infra` branch, refreshes secrets from Infisical, then runs:

```bash
docker compose -f docker-compose.api.yml pull caddy
docker compose -f docker-compose.api.yml up -d --no-deps --force-recreate caddy
```

The `Caddyfile.api` is mounted into the container as a bind mount (`./Caddyfile.api:/etc/caddy/Caddyfile`), so updating the file on the droplet (via the git pull) and recreating the container is all that is needed.

There is no equivalent manual workflow for Caddy on the web droplet — that container is brought up as part of the web Compose stack.

---

## Secrets on the Droplet

Secrets are not stored on the droplet long-term. At every deploy, the workflow authenticates to Infisical using `INFISICAL_CLIENT_ID` and `INFISICAL_CLIENT_SECRET` (stored in GitHub Actions), exports the relevant secret path to a dotenv file under `infra/secrets/`, and that file is consumed by Docker Compose via `env_file`.

For the web stack, a standalone helper script (`infra/fetch-web-secrets.sh`) can be run manually on the droplet to regenerate `secrets/.env.dev.web` without triggering a full deploy. It authenticates against Infisical's REST API directly using `curl` and `jq`.

See [Secrets / Infisical](secrets.md) for full details on the secret layout and path conventions.
