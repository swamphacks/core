import { api } from "@/lib/ky";
import type { ApplicationReviewDetails } from "@/lib/openapi/types";
import {
  ApplicationFieldsSchema,
  type ApplicationFields,
} from "@/modules/Application/hooks/useApplication";
import { useQuery } from "@tanstack/react-query";

export const applicationReviewDetailsQueryKey = (applicationId: string) => [
  "applicationReviewDetails",
  applicationId,
];

export type ParsedApplicationReviewDetails = Omit<
  ApplicationReviewDetails,
  "application"
> & {
  application: ApplicationFields;
};

export async function fetchApplicationForReview(
  applicationId: string,
): Promise<ParsedApplicationReviewDetails> {
  const result = await api
    .get<ParsedApplicationReviewDetails>(`application/review/${applicationId}`)
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

export function useApplicationForReview(applicationId: string) {
  return useQuery({
    queryKey: applicationReviewDetailsQueryKey(applicationId),
    queryFn: () => fetchApplicationForReview(applicationId),
    staleTime: 1000 * 60 * 15, // 15 minutes,
  });
}
