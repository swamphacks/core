import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export const adminEventsQueryKey = ["admin", "events"] as const;

type Events =
  paths["/events"]["get"]["responses"]["200"]["content"]["application/json"];

export async function fetchEvents(): Promise<Events> {
  const result = await api
    .get<Events>("events?include_unpublished=true")
    .json();
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
