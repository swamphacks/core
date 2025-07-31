import type { components } from "./schema";

export type Session = components["schemas"]["Session"];
export type ErrorResponse = components["schemas"]["ErrorResponse"];
export type UserContext = components["schemas"]["UserContext"];
export type PlatformRole = components["schemas"]["PlatformRole"];

// replace later with the real thing
export type Event = {
  id: string; // UUID
  name: string;
  description?: string | null;
  location?: string | null;
  location_url?: string | null;
  max_attendees?: number | null;

  application_open: string; // ISO 8601 date-time string
  application_close: string;
  rsvp_deadline?: string | null;
  decision_release?: string | null;

  start_time: string;
  end_time: string;

  website_url?: string | null;
  is_published: boolean;
  saved_at: string;

  created_at: string;
  updated_at: string;
};
