# Docker

The project uses Docker Compose for both local development and production deployment. The two environments share the same core services but differ in how images are sourced, how traffic is routed, and how hot-reload is handled.

---

## Development vs Production

| Concern | Development (`docker-compose.yml`) | Production (`infra/docker-compose.api.yml`, `infra/docker-compose.web.yml`) |
|---|---|---|
| Image source | Built locally from Dockerfiles | Pulled from `ghcr.io/swamphacks/` |
| API hot-reload | Air (`Dockerfile.dev`) with bind mount | Pre-built binary in `alpine` image |
| Web hot-reload | Vite dev server (`target: dev`) with bind mount | `serve` serving compiled `dist/` |
| Postgres | `postgres:17.4` container, named volume | Not included — assumed external or managed separately |
| Redis | `redis:8.2.1-alpine`, health-checked | `redis:8.2.1-alpine`, memory-capped, persistence optional |
| TLS | None — direct port exposure | Caddy with Cloudflare DNS challenge |
| Env files | `./apps/api/.env.dev` | `./infra/secrets/.env.api`, `.env.web`, `.env.cf` |

The production side is split across two compose files deployed on separate hosts:

- **`infra/docker-compose.api.yml`** — API server, workers, Redis, Caddy (API), Asynqmon. Runs both a `latest` (production) and `dev`-tagged stack side by side on the same host.
- **`infra/docker-compose.web.yml`** — Web frontend, Discord bot, Caddy (web).

---

## Service Breakdown

### Development (`docker-compose.yml`)

| Service | Image / Build | Ports | Depends On |
|---|---|---|---|
| `api` | `./apps/api` via `Dockerfile.dev` | `8080:8080` | `postgres`, `redis` |
| `bat_worker` | `./apps/api` via `cmd/email_worker/Dockerfile` (`target: dev`) | — | `redis` |
| `email_worker` | `./apps/api` via `cmd/email_worker/Dockerfile` (`target: dev`) | — | `redis` |
| `web` | `./apps/web` via `Dockerfile` (`target: dev`) | `5173:5173` | `api` |
| `asynqmon` | `hibiken/asynqmon:latest` | `6767:6767` | `redis` |
| `postgres` | `postgres:17.4` | `5432:5432` | — |
| `redis` | `redis:8.2.1-alpine` | `6379:6379` | — |

### Production API host (`infra/docker-compose.api.yml`)

Two full stacks run concurrently — one tagged `latest` (production) and one tagged `dev`.

| Service | Image | Ports | Depends On |
|---|---|---|---|
| `api` | `ghcr.io/swamphacks/core-api:latest` | `8080:8080` | `redis` |
| `email-worker` | `ghcr.io/swamphacks/core-email-worker:latest` | — | `redis` |
| `bat-worker` | `ghcr.io/swamphacks/core-bat-worker:latest` | — | `redis` |
| `redis` | `redis:8.2.1-alpine` | `6379:6379` | — |
| `asynqmon` | `hibiken/asynqmon:latest` | `6767:6767` | `redis` |
| `api-dev` | `ghcr.io/swamphacks/core-api:dev` | `8081:8080` | `redis-dev` |
| `email-worker-dev` | `ghcr.io/swamphacks/core-email-worker:dev` | — | `redis-dev` |
| `bat-worker-dev` | `ghcr.io/swamphacks/core-bat-worker:dev` | — | `redis-dev` |
| `redis-dev` | `redis:8.2.1-alpine` | `6380:6379` | — |
| `asynqmon-dev` | `hibiken/asynqmon:latest` | `6768:6767` | `redis-dev` |
| `caddy` | `ghcr.io/caddybuilds/caddy-cloudflare:latest` | `80:80`, `443:443` | `api`, `api-dev` |

### Production web host (`infra/docker-compose.web.yml`)

| Service | Image | Exposed | Depends On |
|---|---|---|---|
| `web` | `ghcr.io/swamphacks/core-web:latest` | `80` (internal) | — |
| `web-dev` | `ghcr.io/swamphacks/core-web:dev` | `80` (internal) | — |
| `discord` | `ghcr.io/swamphacks/core-discord:latest` | — | — |
| `caddy` | `ghcr.io/caddybuilds/caddy-cloudflare:latest` | `80:80`, `443:443` | `web-dev` |

The web containers expose port 80 only to the internal `caddy_net` bridge network — they are never bound to the host directly.

---

## Dockerfile Patterns

### API — production (`apps/api/Dockerfile`)

Two-stage build. The `builder` stage uses `golang:1.25-alpine` with `build-base` and `ca-certificates` to compile a statically linked binary:

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/api
```

The final image is bare `alpine:latest` with only `ca-certificates` and the compiled binary. No Go toolchain ships to production.

### API — development (`apps/api/Dockerfile.dev`)

Single stage using `golang:latest`. Installs [Air](https://github.com/air-verse/air) and sets it as the entrypoint. The entire `./apps/api` directory is bind-mounted at `/app`, so Air watches for source changes and rebuilds in place without restarting the container.

### Workers — BAT and email (`cmd/BAT_worker/Dockerfile`, `cmd/email_worker/Dockerfile`)

Both share the same three-stage pattern:

- **`base`** — `golang:1.25-alpine`, downloads modules, copies source.
- **`dev`** — installs Air, runs with the appropriate Air config. Used by the root `docker-compose.yml` via `target: dev`.
- **`prod`** — compiles a statically linked binary with `CGO_ENABLED=0`, adds `ca-certificates`, runs the binary directly.

### Web (`apps/web/Dockerfile`)

Four stages using `node:22.16.0-slim` with `pnpm`:

- **`base`** — installs dependencies via `pnpm install --frozen-lockfile`.
- **`dev`** — exposes `5173`, runs `pnpm run dev --host 0.0.0.0`. Used in the root compose with a bind mount.
- **`build`** — runs `pnpm run build`, producing `/app/dist`.
- **`prod`** — copies `dist/` into a fresh `node:22.16.0-slim` image, serves it with `serve` on port 80. An `entrypoint.sh` script runs first to inject runtime configuration.

---

## Caddy as Reverse Proxy

Production uses a custom Caddy image with the [Cloudflare DNS plugin](https://github.com/caddy-dns/cloudflare) (`ghcr.io/caddybuilds/caddy-cloudflare`) so that TLS certificates are issued via DNS-01 challenge without requiring inbound port 80 to be reachable by Let's Encrypt.

### `Caddyfile.api`

Routes four domains on the API host. All blocks set HSTS, `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, and enable gzip/zstd compression.

| Domain | Upstream | Purpose |
|---|---|---|
| `api.swamphacks.com` | `api:8080` | Production API |
| `dev-api.swamphacks.com` | `api-dev:8080` | Dev-tagged API |
| `asynqmon.swamphacks.com` | `asynqmon:6767` | Production queue dashboard |
| `dev-asynqmon.swamphacks.com` | `asynqmon-dev:6767` | Dev queue dashboard |

Each `reverse_proxy` block forwards `X-Real-IP`, `X-Forwarded-For`, `X-Forwarded-Port`, and `X-Forwarded-Proto` headers upstream.

### `Caddyfile.web`

Routes two domains on the web host. Web containers are only reachable through the internal `caddy_net` Docker network.

| Domain | Upstream | Purpose |
|---|---|---|
| `app.swamphacks.com` | `web:80` | Production frontend |
| `dev-app.swamphacks.com` | `web-dev:80` | Dev-tagged frontend |

---

## Makefile Targets

All targets run `docker compose` against the root `docker-compose.yml`. Run from the repository root.

| Target | Command | What it starts |
|---|---|---|
| `make local` | `docker compose up` | All services — full local stack |
| `make api` | `docker compose up api` | API only |
| `make bat` | `docker compose up api bat_worker asynqmon` | API + BAT worker + Asynqmon dashboard |
| `make backend` | `docker compose up api email_worker bat_worker asynqmon` | API + both workers + Asynqmon |
| `make storage` | `docker compose up postgres redis` | Postgres + Redis only |

`make storage` is useful when running the API from the host with `go run` directly, keeping only the infrastructure containers managed by Docker.

---

## Volume Strategy

| Volume | Used by | Purpose |
|---|---|---|
| `postgres_data` | `postgres` (dev) | Persists database across container restarts |
| `redis_data` | `redis` (prod) | Persists Redis AOF/RDB data in production |
| `redis_data_dev` | `redis-dev` (prod dev stack) | Separate persistence for the dev Redis instance |
| `caddy_data` | `caddy` | Stores TLS certificates issued by Let's Encrypt |
| `caddy_config` | `caddy` | Caddy runtime configuration |

### Bind mounts in development

| Mount | Service | Effect |
|---|---|---|
| `./apps/api:/app` | `api`, `bat_worker`, `email_worker` | Source changes immediately visible to Air — no rebuild required |
| `./apps/web:/app:cached` | `web` | Source changes picked up by Vite HMR |
| `/app/node_modules` | `web` | Anonymous volume prevents the host `node_modules` from shadowing the container's installed packages |

In production, no bind mounts are used. All application code is baked into the image at build time.
