# Infrastructure Overview

SwampHacks Core runs on two DigitalOcean droplets — one for the API stack, one for the web stack — each hosting Docker containers behind a Caddy reverse proxy. Secrets are managed by Infisical and injected at deploy time. All deployments are driven by GitHub Actions.

---

## Components

| Component | Role |
|-----------|------|
| **DigitalOcean** | Hosts the API droplet and the web droplet |
| **Docker / Docker Compose** | Runs all services as containers on each droplet |
| **Caddy** | TLS termination and reverse proxy; certificates issued via Cloudflare DNS challenge |
| **Infisical** | Secrets store; secrets are pulled at deploy time and written to `.env` files |
| **GitHub Actions** | CI/CD pipelines that build images, run migrations, and deploy to each droplet |
| **GHCR** | Docker images are published to GitHub Container Registry (`ghcr.io/swamphacks/`) |

---

## Services and Hosts

### API Droplet (`docker-compose.api.yml`)

| Service | Image tag | Internal port | Public host |
|---------|-----------|---------------|-------------|
| `api` (prod) | `core-api:latest` | 8080 | `api.swamphacks.com` |
| `api-dev` (dev) | `core-api:dev` | 8081 | `dev-api.swamphacks.com` |
| `email-worker` (prod) | `core-email-worker:latest` | — | — |
| `email-worker-dev` | `core-email-worker:dev` | — | — |
| `bat-worker` (prod) | `core-bat-worker:latest` | — | — |
| `bat-worker-dev` | `core-bat-worker:dev` | — | — |
| `redis` (prod) | `redis:8.2.1-alpine` | 6379 | — |
| `redis-dev` | `redis:8.2.1-alpine` | 6380 | — |
| `asynqmon` (prod) | `hibiken/asynqmon` | 6767 | `asynqmon.swamphacks.com` |
| `asynqmon-dev` | `hibiken/asynqmon` | 6768 | `dev-asynqmon.swamphacks.com` |
| `caddy` | `caddy-cloudflare` | 80 / 443 | all of the above |

### Web Droplet (`docker-compose.web.yml`)

| Service | Image tag | Public host |
|---------|-----------|-------------|
| `web` (prod) | `core-web:latest` | `app.swamphacks.com` |
| `web-dev` | `core-web:dev` | `dev-app.swamphacks.com` |
| `discord` | `core-discord:latest` | — |
| `caddy` | `caddy-cloudflare` | all of the above |

---

## Environments

Both droplets run **prod** and **dev** side-by-side:

- **Production** — triggered by pushes to `master`. Images are tagged `:latest`. Secrets come from the `prod` Infisical environment.
- **Development** — triggered by pushes to `dev`. Images are tagged `:dev`. Secrets come from the `dev` Infisical environment.

The local development environment (`docker-compose.yml` at the repo root) is separate — it builds images from source and uses a local Postgres instance. It is not used on the droplets.

---

## CI/CD Workflow Structure

Each service has a dedicated GitHub Actions workflow file. The naming convention reflects the environment and service:

```
dev-build-deploy-<service>.yml   # pushes to dev branch
prod-build-deploy-<service>.yml  # pushes to master branch
deploy-caddy.yml                 # manual dispatch only
```

Every deploy workflow follows the same three-step pattern:

1. **Build & push** — Docker image built and pushed to GHCR.
2. **Migrate** — Goose migrations run against the target database (API workflows only).
3. **Deploy** — SSH into the droplet, pull secrets from Infisical, pull the new image, recreate the container.

---

## Further Reading

- [Docker](docker.md) — Compose file structure, image naming, and how services are organized.
- [Secrets / Infisical](secrets.md) — How secrets are stored, pulled, and injected into containers.
- [DigitalOcean](digitalocean.md) — Droplet setup and SSH access.
- [CI/CD](cicd.md) — Workflow details, required GitHub secrets, and the deploy pipeline.
