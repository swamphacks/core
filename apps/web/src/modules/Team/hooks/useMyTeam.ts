import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type MyTeamResponse =
  operations["get-my-team"]["responses"]["200"]["content"]["application/json"];

export async function fetchMyTeam(): Promise<MyTeamResponse> {
  return api.get<MyTeamResponse>("team/me", { retry: 0 }).json();
}

export const myTeamQueryKey = ["my-team"];

export function useMyTeam() {
  return useQuery({
    queryKey: myTeamQueryKey,
    queryFn: fetchMyTeam,
    staleTime: 1000 * 60 * 15, // 15 minutes,
    retry: 1,
  });
}
