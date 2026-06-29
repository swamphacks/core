import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type ApplicationStats =
  operations["get-application-statistics"]["responses"]["200"]["content"]["application/json"];

export const applicationStatisticsQueryKey = ["applicationStatistics"];

const fetchApplicationStatistics = async () => {
  const response = await api.get<ApplicationStats>(`application/stats`).json();

  return response;
};

export const useApplicationStatistics = () => {
  return useQuery({
    queryKey: applicationStatisticsQueryKey,
    queryFn: () => fetchApplicationStatistics(),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};
