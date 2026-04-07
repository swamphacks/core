import { api } from "@/lib/ky";
import type { Application } from "@/lib/openapi/types";
import { useQuery } from "@tanstack/react-query";

export const myApplicationQueryKey = ["my-application"];

export async function fetchMyApplication(): Promise<Application> {
  return await api.get<Application>(`application`).json();
}

export function useMyApplication() {
  return useQuery({
    queryKey: myApplicationQueryKey,
    queryFn: () => fetchMyApplication(),
    staleTime: Infinity,
  });
}
