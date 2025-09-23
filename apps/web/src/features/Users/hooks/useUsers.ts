import { api } from "@/lib/ky";
import type { User } from "@/lib/openapi/types";
import { useQuery } from "@tanstack/react-query";

export const useUsers = (
  searchTerm: string | null,
  limit: number = 10,
  offset: number = 0,
) => {
  const fetchUsers = async () => {
    const params = new URLSearchParams();
    if (searchTerm) {
      params.set("search", searchTerm);
    }
    params.set("limit", limit.toString());
    params.set("offset", offset.toString());

    const result = await api.get<User[]>(`users?${params.toString()}`).json();

    return result;
  };

  return useQuery({
    queryKey: ["users", searchTerm, limit, offset],
    queryFn: fetchUsers,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};
