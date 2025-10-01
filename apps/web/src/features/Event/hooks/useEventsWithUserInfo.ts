import { api } from "@/lib/ky";
import { useQuery } from "@tanstack/react-query";
import { mapEventsAPIResponseToEventCardProps } from "../utils/mapper";
import type { EventWithUserInfo } from "@/lib/openapi/types";

export const eventsQueryKey = ["events", "published"] as const;

//TODO: Use the operations from openapi types when the schema is updated
export type EventsWithUserInfo = EventWithUserInfo[];

export async function fetchEvents(): Promise<EventsWithUserInfo> {
  const result = await api
    .get<EventsWithUserInfo>("events?include_unpublished=false")
    .json();

  console.log("Fetched events:", result);
  return result;
}

export function useEventsWithUserInfo() {
  return useQuery({
    queryKey: eventsQueryKey,
    queryFn: fetchEvents,
    staleTime: 1000 * 60 * 5, // 5 minutes,
    select: (data) => data.map(mapEventsAPIResponseToEventCardProps),
  });
}
