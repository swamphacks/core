import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export function getEventStaffUsersQueryKey(eventId: string) {
  return ["event", eventId, "bat-runs"] as const;
}

export type Runs =
  paths["/events/{eventId}/bat-runs"]["get"]["responses"]["200"]["content"]["application/json"];
export type Run = Runs[number];

export function useEventStaffUsers(eventId: string) {
  async function fetchEventStaffUsers(): Promise<StaffUsers> {
    const result = await api.get<StaffUsers>(`events/${eventId}/staff`).json();
    return result;
  }

  return useQuery({
    queryKey: getEventStaffUsersQueryKey(eventId),
    queryFn: fetchEventStaffUsers,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
