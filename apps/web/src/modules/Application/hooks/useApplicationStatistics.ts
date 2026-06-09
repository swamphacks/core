import { api } from "@/lib/ky";
import type { ApplicationStats } from "@/lib/openapi/types";
import { useQuery } from "@tanstack/react-query";

const fetchApplicationStatistics = async () => {
  const response = await api.get<ApplicationStats>(`application/stats`).json();

  return response;
};

export const useApplicationStatistics = () => {
  return useQuery({
    queryKey: ["applicationStatistics"],
    queryFn: () => fetchApplicationStatistics(),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};
