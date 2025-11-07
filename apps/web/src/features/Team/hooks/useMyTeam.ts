import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";
import { HTTPError } from "ky";

export type TeamWithMembers =
  paths["/events/{eventId}/teams/me"]["get"]["responses"]["200"]["content"]["application/json"];

export function useMyTeam(eventId: string) {
  async function fetchMyTeam(): Promise<TeamWithMembers | null> {
    try {
      const result = await api
        .get<TeamWithMembers>(`events/${eventId}/teams/me`)
        .json();
      return result;
    } catch (err) {
      if (err instanceof HTTPError && err.response.status === 404) {
        return null;
      }
      throw err; // rethrow other errors so Tanstack Query can handle them
    }
  }

  return useQuery({
    queryKey: ["myTeam", eventId],
    queryFn: fetchMyTeam,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
