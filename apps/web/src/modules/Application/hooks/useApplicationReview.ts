import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import {
  ApplicationFieldsSchema,
  type ApplicationFields,
} from "@/modules/Application/hooks/useApplication";
import { useQuery } from "@tanstack/react-query";

export const applicationReviewQueryKey = (reviewId: string) => [
  "applicationReview",
  reviewId,
];

export type ParsedApplicationReview = Omit<
  operations["get-review-by-id"]["responses"]["200"]["content"]["application/json"],
  "application"
> & {
  application: ApplicationFields;
};

export async function fetchApplicationReview(
  reviewId: string,
): Promise<ParsedApplicationReview> {
  const result = await api
    .get<ParsedApplicationReview>(`application/review/${reviewId}`)
    .json();

  const parsedApplication = ApplicationFieldsSchema.safeParse(
    JSON.parse(atob(result.application as unknown as string)),
  );

  if (!parsedApplication.success) {
    console.error("Failed to parse application data:", parsedApplication.error);
    throw new Error("Invalid application data format");
  }

  return {
    ...result,
    application: parsedApplication.data,
  };
}

export function useApplicationReview(reviewId: string) {
  return useQuery({
    queryKey: applicationReviewQueryKey(reviewId),
    queryFn: () => fetchApplicationReview(reviewId),
    staleTime: 1000 * 60 * 15, // 15 minutes,
  });
}
