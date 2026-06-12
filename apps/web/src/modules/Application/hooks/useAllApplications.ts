import { api } from "@/lib/ky";
import type { components } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type AllApplicationsData = components["schemas"]["AllApplications"];

async function fetchAllApplications(
  limit: number,
  offset: number,
  search: string,
): Promise<AllApplicationsData> {
  const params = new URLSearchParams();
  params.set("limit", limit.toString());
  params.set("offset", offset.toString());
  params.set("search", search);

  return await api.get<AllApplicationsData>(`application/all?${params}`).json();
}

export function useAllApplications(
  limit: number,
  offset: number,
  search: string,
) {
  return useQuery({
    queryKey: ["all-applications", limit, offset, search],
    queryFn: () => fetchAllApplications(limit, offset, search),
    staleTime: 1000 * 60 * 5, // 5 minutes,
  });
}
