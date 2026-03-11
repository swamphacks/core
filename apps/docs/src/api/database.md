# Database Schema

The API uses **PostgreSQL 17**. Auth-related tables live in the `auth` schema; everything else is in the default `public` schema.

All tables include `created_at` and `updated_at` timestamps. `updated_at` is maintained automatically by a `update_modified_column()` trigger.

---

## Auth Schema

### `auth.users`

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
| `role` | `auth_user_role` | `user` or `superuser` |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

### `auth.accounts`

OAuth provider associations. A user can have multiple providers (currently only Discord).

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | |
| `user_id` | UUID FK → `auth.users` | |
| `provider_id` | TEXT | e.g., `discord` |
| `account_id` | TEXT | Provider's user ID |
| `access_token` | TEXT | |
| `refresh_token` | TEXT | |
| `scope` | TEXT | |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

Unique constraint on `(provider_id, account_id)`.

### `auth.sessions`

Active user sessions. Sessions expire and use rolling expiration.

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | Stored in `sh_session_id` cookie |
| `user_id` | UUID FK → `auth.users` | |
| `expires_at` | TIMESTAMPTZ | Extended on use after 24h |
| `ip_address` | TEXT | |
| `user_agent` | TEXT | |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | Used as `last_used_at` |

---

## Public Schema

### `events`

Hackathon events. Drives the application and attendee lifecycle.

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | |
| `name` | TEXT | |
| `description` | TEXT | |
| `location` | TEXT | |
| `location_url` | TEXT | |
| `max_attendees` | INT | Optional cap |
| `application_open` | TIMESTAMPTZ | Applications open |
| `application_close` | TIMESTAMPTZ | Applications close |
| `rsvp_deadline` | TIMESTAMPTZ | |
| `decision_release` | TIMESTAMPTZ | When decisions are released to applicants |
| `start_time` | TIMESTAMPTZ | |
| `end_time` | TIMESTAMPTZ | |
| `website_url` | TEXT | |
| `banner_url` | TEXT | R2 object key for the banner image |
| `is_published` | BOOLEAN | `false` = draft (only staff+ can see) |
| `application_review_started` | BOOLEAN | Locks in reviewer assignments |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

### `event_roles`

Maps users to their role within a specific event.

| Column | Type | Notes |
|---|---|---|
| `user_id` | UUID FK → `auth.users` | |
| `event_id` | UUID FK → `events` | |
| `role` | `event_role_type` | `admin`, `staff`, `attendee`, `applicant` |
| `assigned_at` | TIMESTAMPTZ | |

Primary key: `(user_id, event_id)`.

### `applications`

One application per user per event. Stores the form data as JSONB.

| Column | Type | Notes |
|---|---|---|
| `user_id` | UUID FK → `auth.users` | |
| `event_id` | UUID FK → `events` | |
| `status` | `application_status` | See statuses below |
| `application` | JSONB | Form field data |
| `experience_rating` | INTEGER | Reviewer score (1–5) |
| `passion_rating` | INTEGER | Reviewer score (1–5) |
| `assigned_reviewer_id` | UUID FK → `auth.users` | Nullable |
| `submitted_by` | UUID | User ID at time of submission |
| `waitlisted_at` | TIMESTAMPTZ | When the applicant was waitlisted |
| `saved_at` | TIMESTAMPTZ | Last draft save |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

Primary key: `(user_id, event_id)`.

**`application_status` enum:**

| Value | Description |
|---|---|
| `started` | Draft — created but not submitted |
| `submitted` | Submitted, awaiting review |
| `under_review` | Assigned to a reviewer |
| `accepted` | Accepted by BAT run |
| `rejected` | Rejected by BAT run |
| `waitlisted` | On waitlist |
| `withdrawn` | Withdrawn by applicant |

### `bat_runs`

Records of BAT (Balanced Admissions Thresher) execution results.

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | |
| `event_id` | UUID FK → `events` | |
| `accepted_applicants` | UUID[] | Array of user IDs |
| `rejected_applicants` | UUID[] | Array of user IDs |
| `status` | `bat_run_status` | `running`, `completed`, `failed` |
| `created_at` | TIMESTAMPTZ | |
| `completed_at` | TIMESTAMPTZ | Nullable |

### `teams`

Teams within an event.

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | |
| `name` | TEXT | |
| `owner_id` | UUID FK → `auth.users` | Nullable (SET NULL on delete) |
| `event_id` | UUID FK → `events` | |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

### `team_members`

Many-to-many membership join table.

| Column | Type | Notes |
|---|---|---|
| `user_id` | UUID FK → `auth.users` | |
| `team_id` | UUID FK → `teams` | |
| `joined_at` | TIMESTAMPTZ | |

Primary key: `(user_id, team_id)`.

### `team_join_requests`

Requests from users to join a team.

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | |
| `team_id` | UUID FK → `teams` | |
| `user_id` | UUID FK → `auth.users` | |
| `request_message` | TEXT | Optional message |
| `status` | `join_request_status` | `PENDING`, `APPROVED`, `REJECTED` |
| `processed_by_user_id` | UUID FK → `auth.users` | Nullable |
| `processed_at` | TIMESTAMPTZ | Nullable |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

Unique partial index: one `PENDING` request per `(team_id, user_id)`.

### `event_interest_submissions`

Mailing list for event interest (pre-registration).

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | |
| `event_id` | UUID FK → `events` | |
| `email` | TEXT | |
| `created_at` | TIMESTAMPTZ | |

### `redeemables`

Prize or reward items associated with an event.

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | |
| `event_id` | UUID FK → `events` | |
| `name` | VARCHAR(255) | |
| `amount` | INT | Total available (≥ 0) |
| `max_user_amount` | INT | Per-user limit (≥ 1) |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

### `user_redemptions`

Tracks how many times a user has redeemed a specific redeemable.

| Column | Type | Notes |
|---|---|---|
| `user_id` | UUID FK → `auth.users` | |
| `redeemable_id` | UUID FK → `redeemables` | |
| `amount` | INT | Times redeemed (≥ 0) |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

Primary key: `(user_id, redeemable_id)`.
