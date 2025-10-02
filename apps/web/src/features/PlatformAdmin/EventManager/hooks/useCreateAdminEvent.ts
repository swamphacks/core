import { useMutation, useQueryClient } from "@tanstack/react-query";
import { type AddEvent } from "../components/AddEventModal";
import { adminEventsQueryKey } from "./useAdminEvents";
import type { CreateEvent, Event } from "@/lib/openapi/types";
import { api } from "@/lib/ky";

export function useCreateAdminEvent() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: AddEvent) => {
      // Prep data shape
      const body: Pick<
        CreateEvent,
        | "name"
        | "start_time"
        | "end_time"
        | "application_close"
        | "application_open"
      > = {
        name: data.eventName,
        start_time: data.eventDateRange.start.toISOString(),
        end_time: data.eventDateRange.end.toISOString(),
        application_open: data.applicationDateRange.start.toISOString(),
        application_close: data.applicationDateRange.end.toISOString(),
      };

      const result = await api.post<Event>("events", { json: body }).json();
      return result;
    },
    onSuccess: async (newEventData) => {
      await queryClient.cancelQueries({ queryKey: adminEventsQueryKey });

      const previousEvents =
        queryClient.getQueryData<Event[]>(adminEventsQueryKey) || [];

      queryClient.setQueryData<Event[]>(adminEventsQueryKey, [
        ...previousEvents,
        newEventData,
      ]);
    },
  });
}
