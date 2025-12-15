import { api } from "@/lib/ky";
import { useQuery } from "@tanstack/react-query";

export const myApplicationBaseKey = ["events", "application"];

export async function fetchMyApplication(eventId: string): Promise<any> {
  const result = await api.get<any>(`events/${eventId}/application`).json();
  return result;
}

export function useMyApplication(eventId: string) {
  return useQuery({
    queryKey: [...myApplicationBaseKey, eventId],
    queryFn: () => fetchMyApplication(eventId),
    staleTime: Infinity,
  });
}
