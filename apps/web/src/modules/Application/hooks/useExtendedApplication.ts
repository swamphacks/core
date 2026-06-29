import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import {
  ApplicationFieldsSchema,
  type ApplicationFields,
} from "@/modules/Application/hooks/useApplication";
import { useQuery } from "@tanstack/react-query";

type ExtendedApplicationResponse =
  operations["get-extended-application-by-id"]["responses"]["200"]["content"]["application/json"];

export type ParsedExtendedApplicationResponse = Omit<
  ExtendedApplicationResponse,
  "application"
> & {
  application: ApplicationFields;
};

async function fetchExtendedApplication(
  applicationId: string,
): Promise<ParsedExtendedApplicationResponse> {
  const result = await api
    .get<ParsedExtendedApplicationResponse>(
      `application/extended/${applicationId}`,
    )
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

export const extendedApplicationQueryKey = (applicationId: string) => [
  "application-extended",
  applicationId,
];

export function useExtendedApplication(applicationId: string) {
  return useQuery({
    queryKey: extendedApplicationQueryKey(applicationId),
    queryFn: () => fetchExtendedApplication(applicationId),
    staleTime: 1000 * 60 * 5, // 5 minutes,
  });
}
