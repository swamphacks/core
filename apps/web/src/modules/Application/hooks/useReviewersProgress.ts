import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export const reviewersProgressQueryKey = ["reviewers-progress"];

export type ReviewersProgress =
  paths["/application/review/progress"]["get"]["responses"]["200"]["content"]["application/json"];

export async function fetchReviewersProgress(): Promise<ReviewersProgress> {
  return await api.get(`application/review/progress`).json();
}

export function useReviewersProgress() {
  return useQuery({
    queryKey: reviewersProgressQueryKey,
    queryFn: fetchReviewersProgress,
  });
}
