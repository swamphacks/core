import { api } from "@/lib/ky";
import { useQuery } from "@tanstack/react-query";

export const fetchApplicationResume = async (
  eventId: string,
  applicationId: string,
) => {
  const data = await api
    .get<string>(`events/${eventId}/application/${applicationId}/resume`)
    .json();

  return data;
};

export function useApplicationResume(eventId: string, applicationId: string) {
  return useQuery({
    queryKey: ["applicationResume", eventId, applicationId],
    queryFn: () => fetchApplicationResume(eventId, applicationId),
    staleTime: 1000 * 60 * 5, // 5 minutes,
  });
}
