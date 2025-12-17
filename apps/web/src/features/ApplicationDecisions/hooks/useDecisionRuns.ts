import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type Run =
  paths["/events/{eventId}/bat-runs"]["get"]["responses"]["200"]["content"]["application/json"][number];

const fetchBatRuns = async (eventId: string) => {
  const data = await api.get<Run[]>(`events/${eventId}/bat-runs`).json();
  return data;
};

export const useDecisionRuns = (eventId: string) => {
  return useQuery({
    queryKey: ["batRuns", eventId],
    queryFn: () => fetchBatRuns(eventId),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};
