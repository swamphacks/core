import { api } from "@/lib/ky";
import type { Application } from "@/lib/openapi/types";
import { queryOptions, useQuery } from "@tanstack/react-query";

export const myApplicationQueryKey = ["my-application"];

export async function fetchMyApplication(): Promise<Application> {
  return await api.get<Application>(`application`).json();
}

export const myApplicationQueryOptions = () =>
  queryOptions({
    queryKey: myApplicationQueryKey,
    queryFn: fetchMyApplication,
    staleTime: 1000 * 60 * 1, // 1 minutes
  });

export function useMyApplication(enabled: boolean = true) {
  return useQuery({
    queryKey: myApplicationQueryKey,
    queryFn: () => fetchMyApplication(),
    staleTime: Infinity,
    enabled,
  });
}
