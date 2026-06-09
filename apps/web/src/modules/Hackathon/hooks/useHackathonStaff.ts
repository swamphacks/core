import { api } from "@/lib/ky";
import type { HackathonStaff } from "@/lib/openapi/types";
import { useQuery } from "@tanstack/react-query";

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
