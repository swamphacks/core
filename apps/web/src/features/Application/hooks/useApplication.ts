import { api } from "@/lib/ky";
import { useQuery } from "@tanstack/react-query";

export async function fetchApplication(eventId: string): Promise<any> {
  const result = await api.get<any>(`events/${eventId}/application`).json();
  return result;
}

export function useApplication(eventId: string) {
  return useQuery({
    queryKey: ["events", "application"],
    queryFn: () => fetchApplication(eventId),
    staleTime: Infinity,
  });
}
