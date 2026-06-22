import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type SearchApplicationsResponse =
  operations["search-applications"]["responses"]["200"]["content"]["application/json"];

async function searchApplications(
  limit: number,
  offset: number,
  search: string,
): Promise<SearchApplicationsResponse> {
  const params = new URLSearchParams();
  params.set("limit", limit.toString());
  params.set("offset", offset.toString());
  params.set("search", search);

  return await api
    .get<SearchApplicationsResponse>(`application/search?${params}`)
    .json();
}

export function useSearchApplications(
  limit: number,
  offset: number,
  search: string,
) {
  return useQuery({
    queryKey: ["search-applications", limit, offset, search],
    queryFn: () => searchApplications(limit, offset, search),
    staleTime: 1000 * 60 * 5, // 5 minutes,
  });
}
