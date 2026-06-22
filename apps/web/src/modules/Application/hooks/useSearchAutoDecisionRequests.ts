import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type SearchAutoDecisionRequestsResponse =
  operations["search-auto-decision-requests"]["responses"]["200"]["content"]["application/json"];

export const searchAutoDecisionRequestsQueryKey = [
  "searchAutoDecisionRequests",
];

export async function searchAutoDecisionRequests(
  limit: number,
  offset: number,
  search: string,
  approved: string,
  decision: string,
): Promise<SearchAutoDecisionRequestsResponse> {
  const params = new URLSearchParams();
  params.set("limit", limit.toString());
  params.set("offset", offset.toString());
  params.set("search", search);
  params.set("approved", approved);
  params.set("decision", decision);

  return await api
    .get<SearchAutoDecisionRequestsResponse>(
      `application/review/search-auto-decision-requests?${params}`,
    )
    .json();
}

export function useSearchAutoDecisionRequests(
  limit: number,
  offset: number,
  search: string,
  approved: string,
  decision: string,
) {
  return useQuery({
    queryKey: [
      "searchAutoDecisionRequests",
      limit,
      offset,
      search,
      approved,
      decision,
    ],
    queryFn: () =>
      searchAutoDecisionRequests(limit, offset, search, approved, decision),
    staleTime: 1000 * 60 * 5, // 5 minutes,
  });
}
