import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export type AutoDecisionRequestResponse =
  operations["get-auto-decision-requests"]["responses"]["200"]["content"]["application/json"];

export const autoDecisionRequestsQueryKey = ["autoDecisionRequests"];

export async function fetchAutoDecisionRequests(): Promise<AutoDecisionRequestResponse> {
  return await api
    .get<AutoDecisionRequestResponse>(
      `application/review/all-auto-decision-requests`,
    )
    .json();
}

export async function updateAutoDecisionRequest(payload: {
  requestId: string;
  approved: boolean;
}) {
  await api.patch(`application/review/auto-decision`, {
    json: payload,
  });
}

export function useAutoDecisionRequests() {
  return useQuery({
    queryKey: autoDecisionRequestsQueryKey,
    queryFn: () => fetchAutoDecisionRequests(),
    staleTime: 1000 * 60 * 5, // 5 minutes,
  });
}

export function useUpdateAutoDecisionRequest() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateAutoDecisionRequest,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: autoDecisionRequestsQueryKey });
    },
  });
}
