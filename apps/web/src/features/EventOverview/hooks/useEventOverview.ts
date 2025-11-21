import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type EventOverview =
  paths["/events/{eventId}/overview"]["get"]["responses"]["200"]["content"]["application/json"];

const fetchEventOverview = async (eventId: string) => {
  const response = await api
    .get<EventOverview>(`events/${eventId}/overview`)
    .json();

  console.log("Fetched event overview:", response);

  return response;
};

export const useEventOverview = (eventId: string) => {
  return useQuery({
    queryKey: ["eventOverview", eventId],
    queryFn: () => fetchEventOverview(eventId),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};
