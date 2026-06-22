import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { queryOptions, useQuery } from "@tanstack/react-query";

export type ApplicationResponse =
  operations["get-my-application"]["responses"]["200"]["content"]["application/json"];

export const myApplicationQueryKey = ["my-application"];

export async function fetchMyApplication(): Promise<ApplicationResponse> {
  return await api.get<ApplicationResponse>(`application`).json();
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
