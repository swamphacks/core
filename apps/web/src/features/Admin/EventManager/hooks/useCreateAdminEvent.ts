import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { AddEvent } from "../components/AddEventModal";
import { adminEventsQueryKey } from "./useAdminEvents";
import type { Event } from "@/lib/openapi/types";

export function useCreateAdminEvent() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: AddEvent) => {
      // Make request
      return data;
    },
    onSuccess: async (newEventData) => {
      await queryClient.cancelQueries({ queryKey: adminEventsQueryKey });

      const previousEvents =
        queryClient.getQueryData<Event[]>(adminEventsQueryKey) || [];

      const optimisticEvent: Event = {
        id: crypto.randomUUID(),
        name: newEventData.eventName,
        description: null,
        location: null,
        location_url: null,
        max_attendees: null,

        application_open:
          newEventData.applicationDateRange?.start?.toString() || "",
        application_close:
          newEventData.applicationDateRange?.end?.toString() || "",
        rsvp_deadline: null,
        decision_release: null,

        start_time: newEventData.eventDateRange?.start?.toString() || "",
        end_time: newEventData.eventDateRange?.end?.toString() || "",

        website_url: null,
        is_published: false,
        saved_at: new Date().toISOString(),
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      queryClient.setQueryData<Event[]>(adminEventsQueryKey, [
        ...previousEvents,
        optimisticEvent,
      ]);
    },
  });
}
