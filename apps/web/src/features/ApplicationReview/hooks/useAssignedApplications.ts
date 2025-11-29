import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type AssignedApplications =
  paths["/events/{eventId}/application/assigned"]["get"]["responses"]["200"]["content"]["application/json"];

export const fetchAssignedApplications = async (eventId: string) => {
  const data = await api
    .get<AssignedApplications>(`events/${eventId}/application/assigned`)
    .json();

  return data;
};

export function useAssignedApplications(eventId: string) {
  return useQuery({
    queryKey: getAssignedApplicationsQueryKey(eventId),
    queryFn: () => fetchAssignedApplications(eventId),
    staleTime: 1000 * 60 * 15, // 15 minutes,
  });
}

export function getAssignedApplicationsQueryKey(eventId: string) {
  return ["assignedApplications", eventId] as const;
}
