import { api } from "@/lib/ky";
import { useQuery } from "@tanstack/react-query";
import type { paths } from "@/lib/openapi/schema";

export type Application =
  paths["/events/{eventId}/application/assigned"]["get"]["responses"]["200"]["content"]["application/json"];

export async function fetchAssignedApplication(
  eventId: string,
  userId: string,
): Promise<Application> {
  const result = await api
    .get<any>(`events/${eventId}/application/assigned?userId=${userId}`)
    .json();
  return result;
}

export function useAssignedApplication(eventId: string, userId: string) {
  return useQuery({
    queryKey: ["events", "application", "assigned", eventId, userId],
    queryFn: () => fetchAssignedApplication(eventId, userId),
    staleTime: Infinity,
  });
}
