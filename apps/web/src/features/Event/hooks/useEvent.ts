import { useQuery } from "@tanstack/react-query";
import { getEventById } from "../api/getEvent";
import { EventSchema } from "../schemas/event";

export function useEvent(eventId: string) {
  const fetchEvent = async () => {
    const data = await getEventById(eventId);

    return EventSchema.parse(data);
  };

  const eventQueryKey = ["event", eventId] as const;

  return useQuery({
    queryKey: eventQueryKey,
    queryFn: fetchEvent,
    staleTime: 1000 * 60 * 5, // 5 minutes,
  });
}
