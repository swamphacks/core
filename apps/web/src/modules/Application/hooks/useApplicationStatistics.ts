import { api } from "@/lib/ky";
import type { ApplicationStats } from "@/lib/openapi/types";
import { useQuery } from "@tanstack/react-query";

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
