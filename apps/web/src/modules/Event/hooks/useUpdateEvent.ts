import { useMutation, useQueryClient } from "@tanstack/react-query";
import { updateEventById } from "../api/updateEvent";
import type { Event } from "../schemas/event";

export const useUpdateEvent = (eventId: string) => {
  const queryClient = useQueryClient();

  const eventQueryKey = ["event", eventId] as const;

  return useMutation<Event, unknown, Partial<Event>, unknown>({
    mutationFn: (data) => updateEventById(eventId, data),
    onSuccess: (updatedEvent) => {
      queryClient.setQueryData<Event>(eventQueryKey, (old) => {
        if (!old) return old;

        return updatedEvent;
      });
    },
  });
};
