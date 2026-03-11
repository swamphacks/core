# API Overview

The SwampHacks API is the central backend for the platform. It handles authentication, event management, applications, teams, email delivery, and check-in operations.

## Stack

| Component | Technology |
|---|---|
| Language | Go 1.24 |
| Router | [chi](https://github.com/go-chi/chi) |
| Database | PostgreSQL 17 (via [pgx](https://github.com/jackc/pgx)) |
| Query layer | [sqlc](https://sqlc.dev) (generated, type-safe) |
| Task queue | [Asynq](https://github.com/hibiken/asynq) + Redis |
| Object storage | Cloudflare R2 (S3-compatible) |
| Email | AWS SES |
| Auth | Discord OAuth2 + session cookies |
| API docs | [Scalar](https://scalar.com) (served at `/docs`) |

## Architecture

The API follows a strict layered architecture:

```
HTTP Request
    └── Middleware (auth, roles, logging)
            └── Handler (parse, validate, respond)
                    └── Service (business logic)
                            └── Repository (data access)
                                    └── Database (PostgreSQL via sqlc)
```

Each layer has a single responsibility. Handlers never touch the database directly; services never parse HTTP requests.

## Background Workers

Two separate processes handle async work and run alongside the API:

- **Email Worker** — processes the email task queue (confirmation emails, welcome emails, decision emails)
- **BAT Worker** — runs the Balanced Admissions Thresher, which calculates accept/reject/waitlist decisions from reviewer scores

Both workers share the same codebase and configuration as the API but are started as separate binaries.

## Key Domains

| Domain | Description |
|---|---|
| Auth | Discord OAuth2 login, session management |
| Users | Profiles, onboarding, email consent |
| Events | Hackathon event lifecycle, banners, scopes |
| Applications | Submission, review assignment, BAT decisions, waitlist |
| Teams | Creation, join requests, membership |
| Redeemables | Prize tracking and redemption |
| Email | Async delivery via SES + task queue |
| Discord | Role lookup and attendee queries for the bot |
| Mobile | RFID-based check-in for the mobile check-in app |

## API Documentation

Interactive API documentation is available at `/docs` when the server is running. The raw OpenAPI spec is at `apps/api/docs/swagger.yaml`.

An external hosted version is available at [core.apidocumentation.com](https://core.apidocumentation.com/guide/swamphacks-core-api).
