import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

type ApplicationStatistics =
  paths["/events/{eventId}/application/stats"]["get"]["responses"]["200"]["content"]["application/json"];

const fetchApplicationStatistics = async (eventId: string) => {
  const response = await api
    .get<ApplicationStatistics>(`events/${eventId}/application/stats`)
    .json();

  console.log("Fetched application statistics:", response);

  return response;
};

export const useApplicationStatistics = (eventId: string) => {
  return useQuery({
    queryKey: ["applicationStatistics", eventId],
    queryFn: () => fetchApplicationStatistics(eventId),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};
