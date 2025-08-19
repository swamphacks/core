import type { components, operations } from "./schema";

export type Session = components["schemas"]["Session"];
export type ErrorResponse = components["schemas"]["ErrorResponse"];
export type UserContext = components["schemas"]["UserContext"];
export type PlatformRole = components["schemas"]["PlatformRole"];
export type Event = components["schemas"]["Event"];
export type CreateEvent =
  operations["post-event"]["requestBody"]["content"]["application/json"];
export type User = components["schemas"]["User"];
export type EventWithUserInfo = components["schemas"]["EventWithUserInfo"];
