# Getting Started

### Setup with Docker Compose (main setup)
1. Navigate to `core/apps/web`

2. **Set up environment variables**:
``` bash
cp .env.example .env
```

Fill in the required keys and tokens in your new `.env` file. For `VITE_DISCORD_OAUTH_CLIENT_ID`, retrive the token via the instructions given in [the API installation page](../api/installation.md).

`VITE_` must be prefixed to all environment variables in order for them to be accessible.

3. Continue with the [main setup instructions](../getting-started.md)

### Setup without Docker Compose

1. Make sure [pnpm](https://pnpm.io/) is installed on your system.

2. Navigate to `core/apps/web`

3. Install dependencies

```bash
pnpm install
```

4. Configure environment variables:

```bash
cp .env.example .env
```

Fill in the required keys and tokens in your new `.env` file.

`VITE_` must be prefixed to all environment variables in order for them to be accessible.

5. Finally, launch the app

```bash
pnpm run dev
```
