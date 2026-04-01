import { type Event } from "@/lib/openapi/types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { adminEventsQueryKey } from "./useAdminEvents";
import { api } from "@/lib/ky";

async function updateAdminEvent() {
  // Update here

  return {} as Event; // Hacky right now
}

async function deleteAdminEvent(eventId: string) {
  await api.delete(`events/${eventId}`);
  return eventId;
}

export function useAdminEventActions() {
  const queryClient = useQueryClient();

  const update = useMutation({
    mutationFn: updateAdminEvent,
    onSuccess: (updatedEvent) => {
      queryClient.setQueryData<Event[]>(adminEventsQueryKey, (old) => {
        if (!old) return old;

        return old.map((event) =>
          event.id === updatedEvent.id ? updatedEvent : event,
        );
      });
    },
  });

  const remove = useMutation({
    mutationFn: deleteAdminEvent,
    onSuccess: (deletedEventId) => {
      queryClient.setQueryData<Event[]>(adminEventsQueryKey, (old) => {
        if (!old) return old;

        return old.filter((event) => event.id !== deletedEventId);
      });
    },
  });

  return { update, remove };
}
