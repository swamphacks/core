import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

export type ReviewAssignments =
  operations["get-review-assignments"]["responses"]["200"]["content"]["application/json"];

export const reviewAssignmentsQueryKey = ["reviewAssignments"];

export const fetchReviewAssignments = async () => {
  const data = await api
    .get<ReviewAssignments>(`application/review/assignments`)
    .json();

  return data;
};

export function useReviewAssignments() {
  return useQuery({
    queryKey: reviewAssignmentsQueryKey,
    queryFn: () => fetchReviewAssignments(),
    staleTime: 1000 * 60 * 15, // 15 minutes,
  });
}
