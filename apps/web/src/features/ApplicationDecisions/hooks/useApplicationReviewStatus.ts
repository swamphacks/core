import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

type ApplicationReviewCompleteResponse =
  paths["/events/{eventId}/review-status"]["get"]["responses"]["200"]["content"]["application/json"];

const fetchApplicationReviewComplete = async (eventId: string) => {
  const data = await api
    .get<ApplicationReviewCompleteResponse>(`events/${eventId}/review-status`)
    .json();
  return data.complete;
};

export const useApplicationReviewComplete = (eventId: string) => {
  return useQuery({
    queryKey: ["appReviewComplete", eventId],
    queryFn: () => fetchApplicationReviewComplete(eventId),
    staleTime: 5 * 60 * 1000, // 5 minutes in milliseconds
  });
};
