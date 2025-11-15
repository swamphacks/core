import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type JoinRequestArray =
  paths["/events/{eventId}/teams/me/pending-joins"]["get"]["responses"]["200"]["content"]["application/json"];

export function getMyPendingJoinRequestsQueryKey(eventId: string) {
  return ["myPendingJoinRequests", eventId] as const;
}

export function useMyPendingJoinRequests(eventId: string) {
  async function fetchMyPendingJoinRequests() {
    return api
      .get<JoinRequestArray>(`events/${eventId}/teams/me/pending-joins`)
      .json();
  }

  return useQuery({
    queryKey: getMyPendingJoinRequestsQueryKey(eventId),
    queryFn: fetchMyPendingJoinRequests,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
