import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type HackathonStaff =
  operations["get-hackathon-staff"]["responses"]["200"]["content"]["application/json"];

export const hackathonStaffQueryKey = ["hackathon-staff"];

export async function fetchHackathonStaff(): Promise<HackathonStaff> {
  return await api.get<HackathonStaff>(`hackathon/staff`).json();
}

export function useHackathonStaff() {
  return useQuery({
    queryKey: hackathonStaffQueryKey,
    queryFn: () => fetchHackathonStaff(),
    staleTime: 1000 * 60 * 5, // 5 minutes,
  });
}
