import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type TeamsWithMembers =
  paths["/events/{eventId}/teams"]["get"]["responses"]["200"]["content"]["application/json"];

export function useEventTeams(eventId: string, limit: number, offset: number) {
  async function fetchEventTeams(): Promise<TeamsWithMembers | null> {
    const result = await api
      .get<TeamsWithMembers>(`events/${eventId}/teams`)
      .json();
    return result ?? null;
  }

  return useQuery({
    queryKey: ["teams", eventId, limit, offset],
    queryFn: fetchEventTeams,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
