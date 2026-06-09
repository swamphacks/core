import { api } from "@/lib/ky";
import { useQuery } from "@tanstack/react-query";

export const fetchApplicationResume = async (userId: string) => {
  const data = await api
    .get<string>(`application/review/${userId}/resume`)
    .json();

  return data;
};

export function useApplicationResume(userId: string) {
  return useQuery({
    queryKey: ["applicationResume", userId],
    queryFn: () => fetchApplicationResume(userId),
    staleTime: 1000 * 60 * 5, // 5 minutes,
  });
}
