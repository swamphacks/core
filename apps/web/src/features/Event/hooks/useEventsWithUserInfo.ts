import { api } from "@/lib/ky";
import { useQuery } from "@tanstack/react-query";
import { mapEventsAPIResponseToEventCardProps } from "../utils/mapper";
import type { paths } from "@/lib/openapi/schema";

export const eventsQueryKey = ["events", "published"] as const;

export type EventsWithUserInfo =
  paths["/events"]["get"]["responses"]["200"]["content"]["application/json"];
export type EventWithUserInfo = EventsWithUserInfo[number];

export async function fetchEvents(): Promise<EventsWithUserInfo> {
  const result = await api
    .get<EventsWithUserInfo>("events?scope=scoped")
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
