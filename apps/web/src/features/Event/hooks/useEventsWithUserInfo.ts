import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export const eventsQueryKey = ["events", "published"] as const;

export type EventsWithUserInfo =
  operations["get-events"]["responses"]["200"]["content"]["application/json"];

export async function fetchEvents(): Promise<EventsWithUserInfo> {
  const result = await api.get<EventsWithUserInfo>("events?include_unpublished=false").json();
  return result;
}

export function useAdminEvents() {
  return useQuery({
    queryKey: eventsQueryKey,
    queryFn: fetchEvents,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
