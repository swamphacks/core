import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type TeamMembersResponse =
  operations["get-team-members"]["responses"]["200"]["content"]["application/json"];

export async function fetchTeamMembers(
  teamId: string,
): Promise<TeamMembersResponse> {
  return api
    .get<TeamMembersResponse>(`team/${teamId}/members`, { retry: 0 })
    .json();
}

export const teamMembersQueryKey = (teamId: string) =>
  ["team-members", teamId] as const;

export function useTeamMembers(teamId?: string) {
  return useQuery({
    queryKey: teamId
      ? teamMembersQueryKey(teamId)
      : (["team-members", "unknown"] as const),
    queryFn: () => fetchTeamMembers(teamId!),
    staleTime: 1000 * 60 * 15, // 15 minutes,
    retry: 1,
    enabled: Boolean(teamId),
  });
}
