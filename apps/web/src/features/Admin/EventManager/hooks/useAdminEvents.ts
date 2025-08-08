import { api } from "@/lib/ky";
import type { Event } from "@/lib/openapi/types";
import { useQuery } from "@tanstack/react-query";

export const adminEventsQueryKey = ["admin", "events"] as const;

export async function fetchEvents(): Promise<Event[]> {
  const result = await api.get<Event[]>("events").json();
  console.log(result);
  return result;
}

export function useAdminEvents() {
  return useQuery({
    queryKey: adminEventsQueryKey,
    queryFn: fetchEvents,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
