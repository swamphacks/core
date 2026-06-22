import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { queryOptions, useQuery } from "@tanstack/react-query";

export type Hackathon =
  operations["get-hackathon"]["responses"]["200"]["content"]["application/json"];
export type StaffHackathon =
  operations["get-hackathon-for-staff"]["responses"]["200"]["content"]["application/json"];

export const hackthonQueryKey = ["hackathon"];
export const staffHackthonQueryKey = ["staff-hackathon"];

export async function fetchHackathon(): Promise<Hackathon> {
  return await api.get<Hackathon>(`hackathon`).json();
}

export async function fetchStaffHackathon(): Promise<StaffHackathon> {
  return await api.get<Hackathon>(`hackathon/detailed`).json();
}

export const staffHackathonQueryOptions = () =>
  queryOptions({
    queryKey: staffHackthonQueryKey,
    queryFn: fetchStaffHackathon,
    staleTime: 1000 * 60 * 1, // 1 minutes
  });

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
