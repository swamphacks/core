import type { Event } from "@/lib/openapi/types";
import { useQuery } from "@tanstack/react-query";

export const mockEvents: Event[] = [
  {
    id: "1f6d9f1e-4e0a-4b6b-a5a3-1b2f1cde2a12",
    name: "SwampHacks X",
    description: "University of Florida's annual hackathon.",
    location: "Gainesville, FL",
    location_url: "https://goo.gl/maps/UFHackathon",
    max_attendees: 300,

    application_open: "2025-08-01T00:00:00Z",
    application_close: "2025-08-20T23:59:59Z",
    rsvp_deadline: "2025-08-25T23:59:59Z",
    decision_release: "2025-08-22T12:00:00Z",

    start_time: "2025-09-01T09:00:00Z",
    end_time: "2025-09-03T17:00:00Z",

    website_url: "https://swamphacks.com",
    is_published: true,
    saved_at: "2025-07-30T15:00:00Z",

    created_at: "2025-07-15T12:00:00Z",
    updated_at: "2025-07-30T15:00:00Z",
  },
  {
    id: "2a3d9c42-9633-4df4-881e-6a90e420b61c",
    name: "Art & Code 2025",
    description: "A creative coding retreat in the woods.",
    location: "Catskills, NY",
    location_url: null,
    max_attendees: 50,

    application_open: "2025-09-10T00:00:00Z",
    application_close: "2025-09-25T23:59:59Z",
    rsvp_deadline: "2025-09-28T23:59:59Z",
    decision_release: "2025-09-26T10:00:00Z",

    start_time: "2025-10-10T14:00:00Z",
    end_time: "2025-10-13T18:00:00Z",

    website_url: "https://artandcode.io",
    is_published: false,
    saved_at: "2025-07-29T10:00:00Z",

    created_at: "2025-07-10T08:30:00Z",
    updated_at: "2025-07-29T10:00:00Z",
  },
  {
    id: "3cf6c348-0c1e-472e-a91f-bb6799e10b8a",
    name: "BuildFest 2025",
    description: null,
    location: "San Francisco, CA",
    location_url: "https://example.com/venue",
    max_attendees: null,

    application_open: "2025-08-15T00:00:00Z",
    application_close: "2025-09-01T23:59:59Z",
    rsvp_deadline: null,
    decision_release: null,

    start_time: "2025-09-10T10:00:00Z",
    end_time: "2025-09-12T17:00:00Z",

    website_url: null,
    is_published: true,
    saved_at: "2025-07-28T12:00:00Z",

    created_at: "2025-07-05T11:00:00Z",
    updated_at: "2025-07-28T12:00:00Z",
  },
];

export const adminEventsQueryKey = ["admin", "events"];

export async function fetchEvents(): Promise<Event[]> {
  // Fetch events here

  return mockEvents;
}

export function useAdminEvents() {
  return useQuery({
    queryKey: adminEventsQueryKey,
    queryFn: fetchEvents,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
