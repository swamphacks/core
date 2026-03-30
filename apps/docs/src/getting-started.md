# Getting Started

This guide covers the prerequisites and steps to get the full SwampHacks stack running locally.

## Prerequisites

Install the following tools before continuing.

### Docker

Docker runs the full stack (API, web, database, Redis) in containers so you don't need to install each runtime manually.

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (macOS / Windows / WSL)
- Linux: install [Docker Engine](https://docs.docker.com/engine/install/) and the [Compose plugin](https://docs.docker.com/compose/install/)

Verify:

```bash
docker --version
docker compose version
```

### Node.js (via nvm)

The web app requires **Node 22.16+**. Use [nvm](https://github.com/nvm-sh/nvm) to manage versions.

You can install [nvm here](https://github.com/nvm-sh/nvm) (MacOS / Linux / WSL).

If you use windows (bad bad bad) you can install the windows version of [nvm here](https://github.com/coreybutler/nvm-windows).

> Note: The windows version `should` work exactly the same as the posix compliant nvm, so no worries between system discrepencies.

You can install `Node 22.16` with the following commands:
```
nvm install 22.16
nvm use 22.16
```

Furthermore, if you aren't sure what version you currently are using, you can navigate to `apps/web` and run
`nvm use` which retrieves the version from the `.nvmrc` file and sets it as your current node version.

### Editor

[VS Code](https://code.visualstudio.com/) is recommended. The repository is set up with extensions and settings that work well out of the box. Any editor works!

---

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/swamphacks/core.git
cd core
```

### 2. Configure environment variables

Each service needs its own `.env` file. Copy the examples to get started:

```bash
cp apps/api/.env.dev.example apps/api/.env.dev
cp apps/web/.env.example     apps/web/.env
cp apps/discord-bot/.env.example apps/discord-bot/.env
```

Fill in the required secrets. Most local values are pre-filled in the examples; secrets (Discord OAuth, Gemini API key, etc.) must be obtained from the team.

See each service's installation guide for a full breakdown of variables:

- [API →](api/installation.md)
- [Web →](web/installation.md)
- [Discord Bot →](discord-bot/installation.md)

### 3. Start the stack

```bash
make local
```

This starts the API, web, background workers, PostgreSQL, Redis, and Asynqmon via Docker Compose.

| Service    | URL                        |
|------------|----------------------------|
| Web        | http://localhost:5173      |
| API        | http://localhost:8080      |
| Asynqmon   | http://localhost:6767      |
| PostgreSQL | `localhost:5432`           |
| Redis      | `localhost:6379`           |

#### Partial startup

If you only need part of the stack:

```bash
make api       # API only (+ postgres + redis)
make storage   # PostgreSQL + Redis only
make backend   # API + background workers + Asynqmon
```
