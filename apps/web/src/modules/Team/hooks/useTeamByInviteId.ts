import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type TeamResponse =
  operations["get-team-by-invite-id"]["responses"]["200"]["content"]["application/json"];

export async function fetchTeamByInvitationId(
  inviteId: string,
): Promise<TeamResponse> {
  return api
    .get<TeamResponse>(`team/invitation/${inviteId}/team`, { retry: 0 })
    .json();
}

export const teamByInviteIdQueryKey = ["team-by-invite-id"];

export function useTeamByInviteId(inviteId: string) {
  return useQuery({
    queryKey: teamByInviteIdQueryKey,
    queryFn: () => fetchTeamByInvitationId(inviteId),
    staleTime: 1000 * 60 * 15, // 15 minutes,
    retry: 1,
  });
}
