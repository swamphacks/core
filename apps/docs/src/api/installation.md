# Getting Started

### Setup with Docker Compose (main setup)
1. Navigate to `core/apps/api`

2. **Set up environment variables**:
``` bash
cp .env.example .env.local
```

The example file will temporarily be empty while the this project is in early development. However, .env.local still must be made in order for Docker Compose to work.

3. Continue with the [main setup instructions](../getting-started.md)

### Setup without Docker Compose

1. Make sure you have [Go](https://go.dev/) installed on your system.
2. Initialize the Go project]
``` bash
go mod tidy
```

``` bash
go install github.com/air-verse/air@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

3. (Mac OS/Linux) Give yourself permissions to run the project entrypoint:
```bash
chmod +x ./entrypoint.sh
```

4. Run the program with
```bash
./entrypoint.sh
```
