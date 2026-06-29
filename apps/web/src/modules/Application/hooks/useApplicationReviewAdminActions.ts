import { api } from "@/lib/ky";
import { staffHackthonQueryKey } from "@/modules/Hackathon/hooks/useHackathon";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { applicationStatisticsQueryKey } from "./useApplicationStatistics";
import { reviewersProgressQueryKey } from "@/modules/Application/hooks/useReviewersProgress";

export interface AssignedReviewer {
  userId: string;
  amount: number | null;
}

const updateReviewStatus = async (started: boolean) => {
  await api.post(`application/review/update-status`, {
    json: {
      started,
    },
  });
};

const assignReviewers = async (reviewers: AssignedReviewer[]) => {
  await api.post(`application/review/assign`, {
    json: reviewers,
  });
};

const resetReviews = async () => {
  await api.post(`application/review/reset`);
};

export const useApplicationReviewAdminActions = () => {
  const queryClient = useQueryClient();

  const updateStatus = useMutation({
    mutationFn: updateReviewStatus,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: staffHackthonQueryKey,
      });
      queryClient.invalidateQueries({
        queryKey: applicationStatisticsQueryKey,
      });
    },
  });

  const assign = useMutation({
    mutationFn: assignReviewers,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: reviewersProgressQueryKey,
      });
    },
  });

  const reset = useMutation({
    mutationFn: resetReviews,
    onSuccess: () => {},
  });

  return {
    updateStatus,
    assign,
    reset,
  };
};
