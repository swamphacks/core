# Getting Started

### Setup with Docker Compose (main setup)
1. Navigate to `core/apps/api`

1. **Set up environment variables**:
``` bash
cp .env.dev.example .env.dev
```

1. Open `.env.dev`

    1. For `AUTH_DISCORD_CLIENT_ID` and `AUTH_DISCORD_CLIENT_SECRET`, go to the Discord developer portal and create an account. Create a new application and go to the OAuth2 tab in the left sidebar. Copy the Client ID and the Client Secret into their respective environment variables.

    1. While in the OAuth2 menu, copy the `AUTH_DISCORD_REDIRECT_URI` parameter from the example configuration and paste into the box under the *Redirects* header. This is the URL which discord will redirect the user to after Discord authentication has completed.

    1. Fill out any other required keys and tokens, if empty.

1. Continue with the [main setup instructions](../getting-started.md)

### Setup without Docker Compose

1. Make sure you have [Go](https://go.dev/) installed on your system.
1. Initialize the Go project
``` bash
go mod tidy
```
``` bash
go install github.com/air-verse/air@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

1. Run the program with
```bash
air
```
