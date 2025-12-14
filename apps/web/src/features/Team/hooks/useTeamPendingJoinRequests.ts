import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type JoinRequestArray =
  paths["/teams/{teamId}/pending-joins"]["get"]["responses"]["200"]["content"]["application/json"];

export function getTeamPendingJoinRequestsQueryKey(teamId: string) {
  return ["teamPendingJoinRequests", teamId] as const;
}

export function useTeamPendingJoinRequests(teamId: string) {
  async function fetchTeamPendingJoinRequests() {
    return api.get<JoinRequestArray>(`teams/${teamId}/pending-joins`).json();
  }

  return useQuery({
    queryKey: getTeamPendingJoinRequestsQueryKey(teamId),
    queryFn: fetchTeamPendingJoinRequests,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
