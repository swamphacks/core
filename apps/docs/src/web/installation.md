# Installation & Setup

The web app is a React + Vite SPA. It can run via Docker (part of the full stack) or directly on the host with Node.js.

## Prerequisites

- Node.js 22.16 via nvm (see [Getting Started](../getting-started.md))
- pnpm: `npm install -g pnpm`

## Environment

Copy the example file:

```bash
cp apps/web/.env.example apps/web/.env
```

| Variable | Default | Description |
|---|---|---|
| `VITE_BASE_API_URL` | `https://api.swamphacks.com` | API base URL. Use `http://localhost:8080` for local development |
| `VITE_DISCORD_OAUTH_CLIENT_ID` | — | Discord OAuth application client ID (must match the API's) |
| `VITE_ALLOWED_HOSTS` | `[""]` | JSON array of allowed host origins |

For local development, set `VITE_BASE_API_URL=http://localhost:8080`.

## Running locally

```bash
cd apps/web
nvm use          # picks up .nvmrc (Node 22.16)
pnpm install
pnpm dev
```

The app is available at **http://localhost:5173**.

## Running via Docker

```bash
make local       # full stack including web
# or
docker compose up web
```

## Generating API types

The web app uses auto-generated TypeScript types from the API's OpenAPI spec:

```bash
cd apps/web
pnpm generate:openapi
```

Run this whenever the API schema changes. The output is written to `src/lib/openapi/schema.d.ts`.
