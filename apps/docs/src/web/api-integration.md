# API Integration

The web app communicates with the backend through a typed HTTP client built on [ky](https://github.com/sindresorhus/ky), with [TanStack Query](https://tanstack.com/query) handling caching, invalidation, and async state. Types are generated directly from the backend's OpenAPI spec.

---

## HTTP Client

**File:** `src/lib/ky.ts`

```ts
import config from "@/config";
import ky from "ky";

export const api = ky.create({
  prefixUrl: config.BASE_API_URL,
  credentials: "include",
});
```

`api` is a preconfigured ky instance used throughout every feature. Two things to note:

- **`prefixUrl`** is set to `config.BASE_API_URL`, which resolves from the `VITE_BASE_API_URL` environment variable (see `src/config/env.ts`). In development this is read from `import.meta.env`; in production it is injected at runtime via `window.ENV`.
- **`credentials: "include"`** ensures the browser sends the `sh_session_id` cookie on every request, which is how the API authenticates the user.

All feature-level API calls import `api` from `@/lib/ky` and call methods like `api.get(...)`, `api.post(...)`, `api.patch(...)`, and `api.delete(...)`.

---

## OpenAPI TypeScript Types

### Where they live

```
src/lib/openapi/
  schema.d.ts   ← auto-generated, do not edit
  types.ts      ← hand-maintained convenience re-exports
  zodSchemas.ts ← Zod schemas for a subset of schema enums
```

`schema.d.ts` is generated from the backend's OpenAPI spec (`apps/api/docs/swagger.yaml`) by [openapi-typescript](https://openapi-ts.dev/). It exports two root interfaces:

- **`paths`** — every API route, its parameters, request body, and response shapes.
- **`components`** — shared schema objects (models, error envelopes, etc.).

### How they're used

Feature modules import directly from `schema.d.ts` to derive precise request/response types without any manual duplication:

```ts
// src/features/Event/api/getEvent.ts
import type { operations } from "@/lib/openapi/schema";

type Event =
  operations["get-single-event"]["responses"]["201"]["content"]["application/json"];

export async function getEventById(eventId: string): Promise<Event> {
  return api.get<Event>(`events/${eventId}`).json();
}
```

```ts
// src/features/Team/hooks/useMyTeam.ts
import type { paths } from "@/lib/openapi/schema";

type TeamWithMembers =
  paths["/events/{eventId}/teams/me"]["get"]["responses"]["200"]["content"]["application/json"];
```

```ts
// src/features/EventAdmin/hooks/useStaffActions.ts
import type { components } from "@/lib/openapi/schema";

type AddRoleFields = {
  assignments: components["schemas"]["handlers.AssignRoleFields"][];
};
```

`src/lib/openapi/types.ts` re-exports commonly used types under shorter aliases:

```ts
export type ErrorResponse = components["schemas"]["response.ErrorResponse"];
export type UserContext   = components["schemas"]["middleware.UserContext"];
export type PlatformRole  = components["schemas"]["sqlc.AuthUserRole"];
export type Event         = components["schemas"]["sqlc.Event"];
export type User          = components["schemas"]["sqlc.AuthUser"];
```

### Regenerating types

Run this command from `apps/web` whenever the API spec changes:

```bash
pnpm generate:openapi
```

This executes:

```
openapi-typescript ../api/docs/swagger.yaml -o ./src/lib/openapi/schema.d.ts
```

The output file is committed to the repository. After regenerating, update `types.ts` and `zodSchemas.ts` if any referenced schemas changed.

---

## TanStack Query

**Setup:** `src/integrations/tanstack-query/root-provider.tsx` and `src/lib/query.ts`

```ts
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      experimental_prefetchInRender: true,
    },
  },
});
```

The `QueryClientProvider` wraps the entire application. `experimental_prefetchInRender` allows data fetching to begin during render for routes that use `useSuspenseQuery`.

### Query key pattern

Query keys are exported as named constants or factory functions alongside their hooks so callers can target them for invalidation:

```ts
// src/features/Event/hooks/useEvent.ts
export function getEventQueryKey(eventId: string) {
  return ["event", eventId] as const;
}

// src/features/Application/hooks/useApplication.ts
export function getApplicationQueryKey(eventId: string, userId: string) {
  return ["events", eventId, "application", userId] as const;
}

// src/lib/auth/hooks/useUser.ts
export const queryKey = ["auth", "me"] as const;
```

### Query example

```ts
// src/features/Team/hooks/useMyTeam.ts
export function useMyTeam(eventId: string) {
  return useQuery({
    queryKey: ["myTeam", eventId],
    queryFn: async () => {
      try {
        return await api.get<TeamWithMembers>(`events/${eventId}/teams/me`).json();
      } catch (err) {
        if (err instanceof HTTPError && err.response.status === 404) {
          return null; // user has no team — not an error
        }
        throw err;
      }
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
```

### Mutation example

Mutations invalidate or directly update the cache in `onSuccess`. When the response is sufficient to update local state immediately, `setQueryData` avoids an extra refetch:

```ts
// src/features/Team/hooks/useTeamActions.ts
export function useTeamActions(eventId: string) {
  const queryClient = useQueryClient();

  const create = useMutation({
    mutationFn: (data: NewTeam) =>
      api.post(`events/${eventId}/teams`, { json: data }).json(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["myTeam", eventId] });
    },
  });

  const leave = useMutation({
    mutationFn: (teamId: string) => api.delete(`teams/${teamId}/members/me`),
    onSuccess: () => {
      queryClient.setQueryData(["myTeam", eventId], null);
    },
  });

  return { create, leave };
}
```

When the server response contains the full updated resource, `setQueryData` with the returned value is preferred:

```ts
// src/features/Event/hooks/useUpdateEvent.ts
return useMutation<Event, unknown, Partial<Event>>({
  mutationFn: (data) => updateEventById(eventId, data),
  onSuccess: (updatedEvent) => {
    queryClient.setQueryData<Event>(["event", eventId], () => updatedEvent);
  },
});
```

---

## Error Handling

ky throws an `HTTPError` for any non-2xx response. Errors are handled at two levels:

**1. Inside the API function** — for status codes that require specific handling before the error propagates:

```ts
// src/features/Team/hooks/useJoinRequestActions.ts
import { HTTPError } from "ky";
import type { ErrorResponse } from "@/lib/openapi/types";

async function acceptJoinRequest(requestId: string, teamId: string) {
  try {
    await api.post(`teams/join/${requestId}/accept`);
  } catch (err) {
    if (err instanceof HTTPError) {
      const errorBody = await err.response.json<ErrorResponse>();
      toast.error(errorBody.message || "An error occurred.");
    } else {
      toast.error("An error occurred.");
    }
    throw err; // re-throw so TanStack Query marks the mutation as failed
  }
}
```

The `ErrorResponse` type (`{ error: string; message: string }`) matches the backend's standard error envelope at `components["schemas"]["response.ErrorResponse"]`.

**2. Expected non-error 4xx responses** — a 404 that represents a valid empty state (e.g., user has no team) is caught and converted to `null` instead of being thrown, so the query does not enter an error state.

Toast notifications are rendered via `react-toastify`. A custom `showToast` helper in `src/lib/toast/toast.tsx` supports structured title/message toasts; direct `toast.error()` calls are used in mutation error handlers for brevity.

---

## Session Cookie Flow

Authentication is cookie-based. The backend sets an `sh_session_id` HttpOnly cookie after a successful OAuth login. Because `credentials: "include"` is set on both the `api` ky instance and every bare `fetch` call in the auth service, the browser attaches this cookie automatically to every API request.

**Login flow:**

1. The user initiates sign-in. `auth.oauth.signIn("discord")` (`src/lib/auth/services/oauth.ts`) stores a CSRF nonce in an `sh_auth_nonce` cookie and redirects the browser to Discord's OAuth authorization URL.
2. Discord redirects to `VITE_BASE_API_URL/auth/callback` with `code` and `state` params. The API validates the nonce, exchanges the code for a Discord token, creates or retrieves the user account, and sets `sh_session_id` in the `Set-Cookie` response header.
3. The API redirects the browser back to the app. The session cookie is now present on all subsequent requests.

**Session verification:**

On mount, `auth.useUser()` (backed by `_useUser` in `src/lib/auth/hooks/useUser.ts`) issues `GET /auth/me` with `credentials: "include"`. A 200 response means the session is valid; a 401 means the user is unauthenticated. The result is cached for 10 minutes.

```ts
// src/lib/auth/hooks/useUser.ts
export const queryKey = ["auth", "me"] as const;

export function _useUser() {
  return useQuery({
    queryKey,
    queryFn: async () => await _getUser(),
    refetchOnWindowFocus: false,
    staleTime: 1000 * 60 * 10, // 10 minutes
    retry: false,
  });
}
```

**Logout:**

`auth.logOut()` sends `POST /auth/logout` with `credentials: "include"`. The API clears the session server-side. The `afterLogout` hook (configured in `src/lib/authClient.ts`) then calls `queryClient.invalidateQueries({ queryKey: ["auth", "me"] })`, which causes `useUser` to re-fetch and return an unauthenticated state, triggering a redirect to the login page.
