import { api } from "@/lib/ky";
import type { Hackathon } from "@/lib/openapi/types";
import { queryOptions, useQuery } from "@tanstack/react-query";

export const hackthonQueryKey = ["hackathon"];

export async function fetchHackathon(): Promise<Hackathon> {
  return await api.get<Hackathon>(`hackathon`).json();
}

export const hackathonQueryOptions = () =>
  queryOptions({
    queryKey: hackthonQueryKey,
    queryFn: fetchHackathon,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });

export function useHackathon() {
  return useQuery({
    queryKey: hackthonQueryKey,
    queryFn: () => fetchHackathon(),
    staleTime: 1000 * 60 * 5, // 5 minutes,
  });
}
