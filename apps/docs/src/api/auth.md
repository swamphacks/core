# Authentication & Roles

## Overview

The API uses **Discord OAuth2** for authentication. On success, the server issues a session cookie. All protected routes validate this cookie on every request.

A separate key-based scheme exists for the mobile check-in app.

---

## OAuth2 Flow

```
Client                          API                         Discord
  │                              │                              │
  │── GET /auth/callback?code ──▶│                              │
  │                              │── exchange code ────────────▶│
  │                              │◀─ access token ─────────────│
  │                              │── GET /users/@me ───────────▶│
  │                              │◀─ Discord user info ─────────│
  │                              │                              │
  │                              │  (upsert user + session)     │
  │                              │                              │
  │◀─ Set-Cookie: sh_session_id ─│                              │
  │◀─ 302 → CLIENT_URL ──────────│                              │
```

1. The frontend initiates the OAuth2 flow by redirecting the user to Discord with a `state` nonce stored in the `sh_auth_nonce` cookie.
2. Discord redirects back to `/auth/callback` with a `code` and `state`.
3. The API validates the nonce, exchanges the code for a Discord access token, and fetches the user's Discord profile.
4. If the user is new, an `auth.users` record and `auth.accounts` record are created in a transaction. Otherwise, a new session is created for the existing user.
5. The session ID is set as the `sh_session_id` cookie and the user is redirected to the frontend.

---

## Session Validation

Every request to a protected route goes through `RequireAuth` middleware:

1. Reads the `sh_session_id` cookie.
2. Looks up the session in `auth.sessions` (must not be expired).
3. Fetches the associated user record.
4. Attaches a `UserContext` to the request context.

**Rolling expiration:** If the session has not been used in the past 24 hours, its expiration is extended by 30 days and the cookie is refreshed.

### UserContext fields

| Field | Type | Description |
|---|---|---|
| `UserID` | UUID | Unique user identifier |
| `Email` | `*string` | Primary email from Discord |
| `PreferredEmail` | `*string` | User-set preferred email |
| `Name` | string | Display name |
| `Onboarded` | bool | Whether onboarding is complete |
| `Image` | `*string` | Profile image URL |
| `Role` | `UserRole` | `admin`, `staff`, `attendee`, `applicant`, `visitor` |
| `EmailConsent` | bool | Whether the user opted into emails |

---

## User Roles

| Role | Description |
|---|---|
| `admin` | Full event management (create/delete/assign roles, release decisions) |
| `staff` | Event operations (check-in, review applications, manage redeemables) |
| `attendee` | Accepted attendee |
| `applicant` | Has submitted an application |
| `visitor` | Has made an account, but never submitted an application |

R oles are enforced by `RequireEventRole(roles)` middleware, which fetches the user's role for the event from the URL path. Superusers bypass event role checks.

---

## Mobile Authentication

The mobile check-in app uses a static key instead of session cookies:

```
Authorization: Key <MOBILE_AUTH_KEY>
```

Routes under `/mobile` require this header. The key is configured via the `MOBILE_AUTH_KEY` environment variable (not in `.env.dev.example` — request it from the team).

---

## Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/auth/callback` | None | OAuth2 callback |
| `GET` | `/users/me` | Session | Get current user |
| `POST` | `/auth/logout` | Session | Invalidate session |
