# Architecture & State Management

## Directory Layout

```
src/
├── components/        # Shared, reusable UI components
├── config/            # Runtime and environment configuration
├── features/          # Feature modules (see below)
├── forms/             # Shared form definitions and field presets
├── integrations/      # Third-party library bootstrapping (e.g. TanStack Query)
├── lib/               # Core internal libraries (auth, HTTP client, OpenAPI types)
├── routes/            # TanStack Router file-based route tree
├── utils/             # Pure utility functions (cn, date formatting, etc.)
├── main.tsx           # Application entry point
└── routeTree.gen.ts   # Auto-generated route tree (do not edit manually)
```

| Directory | Purpose |
|---|---|
| `components/` | Application-wide presentational components: `AppShell`, `ThemeProvider`, form primitives, `ui/` (React Aria-based design system components), icon wrappers, and loading states |
| `config/` | Zod-validated environment schema. Supports both Vite build-time env vars (`import.meta.env`) and a runtime `window.ENV` object for Docker deployments |
| `features/` | Self-contained feature modules — each owns its components, hooks, API calls, schemas, and utilities |
| `forms/` | Shared form field definitions and validation presets reused across features |
| `integrations/` | Library-specific provider bootstrapping isolated from app code. Currently contains the TanStack Query `QueryClient` factory and `Provider` wrapper |
| `lib/` | Internal libraries: `auth/` (custom OAuth2 client), `ky.ts` (HTTP client instance), `openapi/` (generated types + Zod schemas), `qr-intents/`, and `toast/` |
| `routes/` | File-based routing via TanStack Router. Includes `__root.tsx`, public pages, and the `_protected/` subtree which enforces authentication |
| `utils/` | Stateless helpers: `cn.ts` (class merging), `date.ts`, `object.ts`, `formHelper.ts` |

---

## State Management

The app uses **TanStack Query** exclusively for server state. There is no global Redux or Zustand store — component-local `useState`/`useReducer` handles ephemeral UI state, and TanStack Query owns all remote data.

### QueryClient setup

The `QueryClient` is instantiated inside `src/integrations/tanstack-query/root-provider.tsx` and injected into both the React tree and the TanStack Router context:

```ts title="src/integrations/tanstack-query/root-provider.tsx"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

export function getContext() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        experimental_prefetchInRender: true, // React 19 Suspense-compatible prefetching
      },
    },
  });

  return { queryClient };
}

export function Provider({ children, queryClient }) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
```

A second `queryClient` singleton (`src/lib/query.ts`) is used by non-React code — specifically by the auth client to `invalidateQueries` after logout:

```ts title="src/lib/authClient.ts"
export const auth = Auth({
  providers: [Discord],
  redirectUri: authConfig.OAUTH_REDIRECT_URL,
  hooks: {
    afterLogout: async () => {
      await queryClient.invalidateQueries({ queryKey: useUserQueryKey });
    },
  },
});
```

### Query conventions

Feature hooks follow a consistent pattern: an `async` fetch function calls the `api` client, and a `useQuery` or `useMutation` wrapper is exported alongside an explicit `queryKey` factory:

```ts title="src/features/Event/hooks/useEvent.ts"
export function getEventQueryKey(eventId: string) {
  return ["event", eventId] as const;
}

export function useEvent(eventId: string) {
  return useQuery({
    queryKey: getEventQueryKey(eventId),
    queryFn: () => fetchEvent(eventId),
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
```

Mutations use `useQueryClient()` directly to update or invalidate related queries on success:

```ts title="src/features/Event/hooks/useUpdateEvent.ts"
export const useUpdateEvent = (eventId: string) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data) => updateEventById(eventId, data),
    onSuccess: (updatedEvent) => {
      queryClient.setQueryData<Event>(["event", eventId], () => updatedEvent);
    },
  });
};
```

Cache invalidation is preferred over optimistic updates for write-heavy operations (e.g. creating redeemables, redeeming items).

---

## Component Model

### `features/` — feature-based colocation

Each feature directory owns everything it needs:

```
features/Event/
├── api/          # Raw fetch functions (no React, no hooks)
│   ├── getEvent.ts
│   └── updateEvent.ts
├── components/   # Feature-specific React components
│   ├── EventCard.tsx
│   └── EventSettingsForm.tsx
├── hooks/        # TanStack Query wrappers and local state hooks
│   ├── useEvent.ts
│   └── useUpdateEvent.ts
├── schemas/      # Zod schemas for the feature's data types
│   └── event.ts
└── utils/        # Pure helpers scoped to this feature
    └── mapper.ts
```

Not all features need every subdirectory — small features (e.g. `Dashboard/`) may consist of just `components/`.

### `components/` — shared infrastructure

`components/` is for UI that is genuinely feature-agnostic:

| Directory / File | Contents |
|---|---|
| `AppShell/` | Sidebar + topbar layout shell, mobile slideout nav, `AppShellContext` |
| `ui/` | Design system primitives built on React Aria Components: `Button`, `Modal`, `ComboBox`, `DatePicker`, `Badge`, etc. |
| `Form/` | Shared form layout wrappers |
| `ThemeProvider.tsx` | Dark/light/system theme context using `localStorage` |
| `Loading.tsx`, `PageLoading.tsx` | Suspense and full-page loading states |
| `ECharts.tsx` | ECharts wrapper component |

Features consume from `components/` but never import across feature boundaries.

---

## Data Flow

```
User action (click / form submit)
        │
        ▼
Feature hook (useQuery / useMutation)
        │
        ▼
api/ fetch function  ──►  ky HTTP client (src/lib/ky.ts)
        │                       │
        │                       ▼
        │               REST API (prefixUrl: BASE_API_URL, credentials: "include")
        │                       │
        ▼                       ▼
TanStack Query cache  ◄──  JSON response / error
        │
        ▼
React component re-renders with updated data
```

### HTTP client

All API calls go through a single `ky` instance configured in `src/lib/ky.ts`:

```ts title="src/lib/ky.ts"
import ky from "ky";
import config from "@/config";

export const api = ky.create({
  prefixUrl: config.BASE_API_URL,
  credentials: "include", // session cookie sent on every request
});
```

Feature `api/` functions call `api.get(...)`, `api.post(...)`, etc. and return typed JSON. They are plain `async` functions with no React dependency, making them independently testable.

### OpenAPI types

`src/lib/openapi/schema.d.ts` is generated from the API's Swagger spec via `openapi-typescript`. Feature code references `operations` and `components` from this file directly for request/response types, keeping API contracts enforced at compile time:

```ts title="src/features/Event/api/getEvent.ts"
import type { operations } from "@/lib/openapi/schema";

type Event =
  operations["get-single-event"]["responses"]["201"]["content"]["application/json"];

export async function getEventById(eventId: string): Promise<Event> {
  return await api.get<Event>(`events/${eventId}`).json();
}
```

---

## Global Providers (`main.tsx`)

The provider stack, outermost to innermost:

```tsx title="src/main.tsx"
<StrictMode>
  <ThemeProvider defaultTheme="system">       // dark/light/system CSS class on <html>
    <TanStackQueryProvider.Provider ...>      // QueryClientProvider
      <InnerApp />                            // resolves auth, then mounts RouterProvider
      <ToastContainer />                      // react-toastify portal
    </TanStackQueryProvider.Provider>
  </ThemeProvider>
</StrictMode>
```

`InnerApp` calls `auth.useUser()` (a TanStack Query hook with a 10-minute stale time and `retry: false`) and passes the resulting query object into the router context. Every route in the `_protected/` subtree awaits `context.userQuery.promise` in `beforeLoad` and redirects unauthenticated users to `/`.

The `QueryClient` and `userQuery` are both part of the TanStack Router context (`RouterContext`), making them available to route loaders and `beforeLoad` guards without prop drilling.

### Auth client

`src/lib/auth/` is a small bespoke OAuth2 library — not a third-party package. It exposes:

| Export | Description |
|---|---|
| `auth.useUser()` | TanStack Query hook; fetches `/auth/me`, returns `{ user, error }` |
| `auth.oauth.signIn(provider, redirect?)` | Initiates Discord OAuth2 redirect with a CSRF nonce cookie |
| `auth.logOut()` | POSTs to `/auth/logout`, then invalidates the `["auth", "me"]` query |
| `auth.getUser()` | Non-hook version of the `getUser` fetch for use outside React |
