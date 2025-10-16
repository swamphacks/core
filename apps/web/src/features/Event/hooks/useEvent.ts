import { useQuery } from "@tanstack/react-query";
import { getEventById } from "../api/getEvent";
import { EventSchema } from "../schemas/event";

export const fetchEvent = async (eventId: string) => {
  const data = await getEventById(eventId);

  return EventSchema.parse(data);
};

export function useEvent(eventId: string) {
  return useQuery({
    queryKey: getEventQueryKey(eventId),
    queryFn: () => fetchEvent(eventId),
    staleTime: 1000 * 60 * 5, // 5 minutes,
  });
}

/**
 *
 * @param eventId - The ID of the event
 * @returns A tuple representing the query key for the event
 */
export function getEventQueryKey(eventId: string) {
  return ["event", eventId] as const;
}
