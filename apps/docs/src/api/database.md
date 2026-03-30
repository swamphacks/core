# Database Schema

The API uses **PostgreSQL 17**. Everything is in the default `public` schema.

All tables include `created_at` and `updated_at` timestamps. `updated_at` is maintained automatically by a `update_modified_column()` trigger.

---

## Public Schema

### `users`

Core user accounts. One record per registered user.

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | Auto-generated |
| `name` | TEXT | Display name from Discord |
| `email` | TEXT UNIQUE | Primary email from Discord |
| `preferred_email` | TEXT | User-set preferred email |
| `email_verified` | BOOLEAN | Default `false` |
| `email_consent` | BOOLEAN | Marketing email opt-in |
| `onboarded` | BOOLEAN | Whether onboarding is complete |
| `image` | TEXT | Profile image URL |
| `rfid` | TEXT | RFID string |
| `role` | `UserRole` | `admin`, `staff`, `attendee`, `applicant`, `visitor` |
| `role_assigned_at` | TIMESTAMPTZ | |
| `checked_in_at` | TIMESTAMPTZ | |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

### `accounts`

OAuth provider associations. A user can have multiple providers (currently only Discord).

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | |
| `user_id` | UUID FK → `users` | |
| `provider_id` | TEXT | e.g., `discord` |
| `account_id` | TEXT | Provider's user ID |
| `access_token` | TEXT | |
| `refresh_token` | TEXT | |
| `scope` | TEXT | |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

Unique constraint on `(provider_id, account_id)`.

### `sessions`

Active user sessions. Sessions expire and use rolling expiration.

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | Stored in `sh_session_id` cookie |
| `user_id` | UUID FK → `users` | |
| `expires_at` | TIMESTAMPTZ | Extended on use after 24h |
| `ip_address` | TEXT | |
| `user_agent` | TEXT | |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | Used as `last_used_at` |

---

TODO: Document other tables