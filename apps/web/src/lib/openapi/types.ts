import type { components, paths } from "./schema";

export type ErrorResponse = components["schemas"]["response.ErrorResponse"];
export type UserContext = components["schemas"]["middleware.UserContext"];
export type PlatformRole = components["schemas"]["sqlc.AuthUserRole"];
export type Event = components["schemas"]["sqlc.Event"];
export type CreateEvent =
  paths["/events"]["post"]["requestBody"]["content"]["application/json"];
export type User = components["schemas"]["middleware.UserContext"];

//TODO: Remove extension once OpenAPI is updated
export type EventWithUserInfo =
  components["schemas"]["sqlc.GetEventsWithUserInfoRow"] & {
    banner: string | null;
  };
