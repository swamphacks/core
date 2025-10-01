import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

type GetUsersResponse =
  paths["/users"]["get"]["responses"]["200"]["content"]["application/json"];

export const useUsers = (
  query: string | null,
  limit: number = 10,
  offset: number = 0,
) => {
  const fetchUsers = async () => {
    const params = new URLSearchParams();
    if (query) {
      params.set("search", query);
    }
    params.set("limit", limit.toString());
    params.set("offset", offset.toString());

    const result = await api
      .get<GetUsersResponse>(`users?${params.toString()}`)
      .json();

    return result;
  };

  return useQuery({
    queryKey: ["users", query, limit, offset],
    queryFn: fetchUsers,
    enabled: !!query, // Only fetch when there's a search term
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};
