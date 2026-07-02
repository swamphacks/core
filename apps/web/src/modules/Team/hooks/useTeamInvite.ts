import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type InvitationResponse =
  operations["get-invitation"]["responses"]["200"]["content"]["application/json"];

export async function fetchInvitation(
  teamId: string,
): Promise<InvitationResponse> {
  return api
    .get<InvitationResponse>(`team/${teamId}/invitation`, { retry: 0 })
    .json();
}

export const invitationQueryKey = ["invitation"];

export function useTeamInvitation(teamId: string) {
  return useQuery({
    queryKey: invitationQueryKey,
    queryFn: () => fetchInvitation(teamId),
    staleTime: 1000 * 60 * 15, // 15 minutes,
    retry: 1,
  });
}
