import { api } from "@/lib/ky";
import type { ApplicationAutoDecisionRequest } from "@/lib/openapi/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export const autoDecisionRequestsQueryKey = ["autoDecisionRequests"];

export async function fetchAutoDecisionRequests(): Promise<ApplicationAutoDecisionRequest> {
  return await api
    .get<ApplicationAutoDecisionRequest>(`application/review/auto-decision`)
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
