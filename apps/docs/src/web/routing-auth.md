# Routing & Auth

## Routing Setup

The web app uses **TanStack Router** with **file-based routing**. Route files live under `src/routes/` and the router automatically generates a typed route tree at `src/routeTree.gen.ts` (committed, regenerated on build).

The router is instantiated in `src/main.tsx` with two pieces of context injected at startup — the TanStack Query client and the user query — making both available in every `beforeLoad` and `loader` function across the tree:

```ts
const router = createRouter({
  routeTree,
  context: {
    ...TanStackQueryProviderContext, // { queryClient }
    userQuery: undefined!,           // populated by InnerApp below
  },
});

function InnerApp() {
  const userQuery = auth.useUser();
  return <RouterProvider router={router} context={{ userQuery }} />;
}
```

### Route Tree Overview

```
/                              ← login page (src/routes/index.tsx)
/privacy                       ← public
/terms                         ← public

/_protected                    ← auth guard layout
  /settings
  /_user                       ← app-shell layout (navbar, logo)
    /portal
    /community
    /resources/programming
    /resources/sponsors
  /admin                       ← superuser-only layout
    /overview
    /events-management
    /users-management
    /logs
    /settings
  /events/$eventId
    /                          ← under construction
    /application
    /summary
    /rejected
    /feedback/declined
    /waitlist/info
    /dashboard                 ← event-role-aware layout
      /                        ← role-based overview
      /my-team
      /teams-explorer
      /_admin/…                ← admin-only sub-routes
      /_staff/…                ← staff + admin sub-routes
      /_applicant/…            ← applicant + admin sub-routes
      /_attendee/…             ← attendee + admin sub-routes
```

Pathless layout segments (prefixed with `_`) group routes under a shared layout or guard without contributing a URL segment. For example, `/_protected` and `/_user` are both pathless — they exist only to run `beforeLoad` checks and render a wrapping component.

---

## Protected Routes

A single pathless layout route at `src/routes/_protected/layout.tsx` gates the entire authenticated section of the app. Its `beforeLoad` hook runs before any child route is matched:

```ts
export const Route = createFileRoute("/_protected")({
  beforeLoad: async ({ context, location }) => {
    const { user, error } = await context.userQuery.promise;

    if (!user && !error) {
      throw redirect({
        to: "/",
        search: { redirect: location.pathname },
      });
    }

    if (error) {
      throw redirect({ to: "/" });
    }

    return { user };
  },
  pendingMs: 1000,
  pendingComponent: () => PageLoading(),
});
```

- If the user fetch returns neither a user nor an error (unauthenticated), the visitor is sent to `/` with a `?redirect=` query param preserving the intended destination.
- If the fetch itself errors (network or server fault), the visitor is sent to `/` with no redirect param.
- On success, `user` is forwarded into the route context for all child routes.

The login page (`/`) has a symmetric check — if a user is already present when the root route loads, they are immediately redirected to `/portal`:

```ts
export const Route = createFileRoute("/")({
  validateSearch: z.object({
    redirect: z.string().optional().catch(""),
  }),
  beforeLoad: async ({ context }) => {
    const { user } = await context.userQuery.promise;
    if (user) {
      throw redirect({ to: "/portal" });
    }
  },
});
```

### Admin Guard

`src/routes/_protected/admin/layout.tsx` adds a second layer of enforcement on top of `/_protected`. It reads `user` from the context already resolved by the parent layout, then checks for `role === "superuser"`:

```ts
export const Route = createFileRoute("/_protected/admin")({
  beforeLoad: async ({ context, location }) => {
    const { user } = context;

    if (!user) {
      throw redirect({ to: "/", search: { redirect: location.pathname } });
    }

    if (user.role !== "superuser") {
      throw redirect({ to: "/portal" });
    }

    if (location.pathname === "/admin") {
      throw redirect({ to: "/admin/overview" });
    }
  },
});
```

Non-superusers are silently redirected to `/portal`.

### Event-Role Guards

Inside the event dashboard, four pathless sub-layouts gate access by event role. The parent layout (`/_protected/events/$eventId/dashboard`) resolves the user's event role via `getUserEventRole()` and places it in context as `eventRole`. Each sub-layout checks this value and calls `notFound()` if the user's role is insufficient:

| Layout file | Allowed roles |
|---|---|
| `_admin/layout.tsx` | `admin` |
| `_staff/layout.tsx` | `admin`, `staff` |
| `_applicant/layout.tsx` | `admin`, `applicant` |
| `_attendee/layout.tsx` | `admin`, `attendee` |

```ts
// _staff/layout.tsx — example
beforeLoad: async ({ context }) => {
  if (!context.eventRole || !["admin", "staff"].includes(context.eventRole)) {
    return notFound();
  }
  return {};
},
```

The dashboard index route (`/dashboard/`) also redirects applicants directly to their application status page:

```ts
beforeLoad: ({ context, params }) => {
  if (context.eventRole === "applicant") {
    throw redirect({
      to: `/events/$eventId/dashboard/application-status`,
      params: { eventId: params.eventId },
    });
  }
},
```

---

## Auth Flow

Authentication uses **Discord OAuth2** (authorization code flow). The entire auth implementation lives in `src/lib/auth/` and is exposed through a single configured client at `src/lib/authClient.ts`.

### Initiating Login

Calling `auth.oauth.signIn("discord")` triggers `_oauthSignIn` in `src/lib/auth/services/oauth.ts`:

1. A random UUID nonce is generated with `crypto.randomUUID()`.
2. The nonce is persisted in a `sh_auth_nonce` cookie (`sameSite: lax`; `secure` in production; domain scoped to `localhost` in dev or `.swamphacks.com` in production).
3. An OAuth state object `{ nonce, provider, redirect? }` is base64-encoded (`btoa(JSON.stringify(state))`) and passed as the `state` query parameter.
4. The browser is navigated to Discord's authorization endpoint with these parameters:

```
https://discord.com/oauth2/authorize
  ?response_type=code
  &scope=identify%20email
  &client_id=<DISCORD_OAUTH_CLIENT_ID>
  &redirect_uri=<BASE_API_URL>/auth/callback
  &state=<base64-encoded state>
```

### Callback Handling

The redirect URI is `<BASE_API_URL>/auth/callback` — this is an **API endpoint**, not a frontend route. The API handles the code exchange, validates the nonce, creates or updates the user record, and issues an `sh_session_id` session cookie before redirecting the browser back to the frontend. See the [API auth docs](../api/auth.md) for the server-side flow.

### Session State

After login, all auth state is maintained by the `sh_session_id` HTTP-only cookie sent automatically by the browser. The frontend has no direct access to the session token.

User data is fetched by `_useUser` (`src/lib/auth/hooks/useUser.ts`) using TanStack Query:

```ts
export const queryKey = ["auth", "me"] as const;

export function _useUser() {
  return useQuery({
    queryKey,
    queryFn: async () => await _getUser(),
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    staleTime: 1000 * 60 * 10, // 10 minutes
    retry: false,
  });
}
```

`_getUser` calls `GET /auth/me` with `credentials: "include"`. A `401` response is treated as "not logged in" and returns `{ user: null, error: null }` — this is what the `/_protected` guard reads when deciding to redirect.

The query result is passed into the router as the `userQuery` context value on every render, so all `beforeLoad` hooks can `await context.userQuery.promise` for the resolved value without triggering a separate fetch.

### User Context Shape

The response from `/auth/me` is validated against this Zod schema (`src/lib/auth/types/user.ts`):

| Field | Type | Description |
|---|---|---|
| `userId` | `string` (UUID) | Unique user identifier |
| `email` | `string` | Primary email from Discord |
| `preferredEmail` | `string \| null` | User-set preferred contact email |
| `name` | `string` | Display name |
| `onboarded` | `boolean` | Whether onboarding has been completed |
| `image` | `string \| null` | Profile image URL |
| `role` | `"user" \| "superuser"` | Platform role |
| `emailConsent` | `boolean` | Whether the user has opted into emails |

### Logout

`auth.logOut()` (`src/lib/auth/services/user.ts`) sends `POST /auth/logout` with `credentials: "include"`. On success, the `afterLogout` hook invalidates the `["auth", "me"]` query in TanStack Query. The settings page then forces a full document reload to reset all in-memory router state:

```ts
await auth.logOut();
await router.navigate({ to: "/", replace: true, reloadDocument: true });
```

---

## Onboarding Redirect

There is no hard redirect for un-onboarded users. Instead, the `/portal` route checks the `onboarded` field from user context and the presence of a `welcome-modal-skipped` cookie to decide whether to show an onboarding modal on load:

```ts
// src/routes/_protected/_user/portal.tsx
beforeLoad: (context) => {
  const { user } = context.context;
  const hasSkippedCookie = Cookies.get("welcome-modal-skipped") === "true";
  const showOnboardingModal = !hasSkippedCookie && !!user && !user.onboarded;
  return { showOnboardingModal };
},
```

The `<OnboardingModal>` component is then rendered conditionally based on this context value. Users who dismiss without completing onboarding have the `welcome-modal-skipped` cookie set, suppressing the modal on future visits.

---

## Layout Routes & Nesting

The app uses two distinct patterns for nested layouts:

**Auth/guard layouts** — pathless segments that run `beforeLoad` and render `<Outlet />` without adding UI chrome. Examples: `/_protected`, `/_protected/events/$eventId/dashboard/_staff`.

**Shell layouts** — pathless segments that also render the application shell (navbar, header, sidebar) around `<Outlet />`. Examples:

- `/_protected/_user` — renders `<AppShell>` with the main navigation (Events Portal, Resources, Community).
- `/_protected/admin` — renders `<AppShell>` with the admin sidebar (Overview, Events Management, Users Management, Logs, Settings).
- `/_protected/events/$eventId/dashboard` — renders a role-adaptive shell: `StaffAppShell` for `admin`/`staff`, `AttendeeAppShell` for `attendee`, `ApplicantAppShell` for `applicant`.

The `/_protected/settings` route sits directly under `/_protected` with no shell layout — it renders its own full-page centered layout.
