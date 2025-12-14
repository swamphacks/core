import { api } from "@/lib/ky";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { AssignedReviewer } from "../components/ReviewNotStarted/ReviewerAssignmentModal";
import { getEventQueryKey } from "@/features/Event/hooks/useEvent";
import { getAssignedApplicationsQueryKey } from "./useAssignedApplications";

const assignReviewers = async (
  eventId: string,
  reviewers: AssignedReviewer[],
) => {
  await api.post(`events/${eventId}/application/assign-reviewers`, {
    json: reviewers,
  });
};

const resetReviews = async (eventId: string) => {
  await api.post(`events/${eventId}/application/reset-reviews`);
};

export const useAppReviewAdminActions = (eventId: string) => {
  const queryClient = useQueryClient();

  const assign = useMutation({
    mutationFn: (reviewers: AssignedReviewer[]) =>
      assignReviewers(eventId, reviewers),
    onSuccess: () => {
      queryClient.setQueryData(getEventQueryKey(eventId), (oldData) => {
        if (!oldData) return oldData;

        return {
          ...oldData,
          application_review_started: true,
        };
      });
      queryClient.invalidateQueries({
        queryKey: getAssignedApplicationsQueryKey(eventId),
      });
      queryClient.invalidateQueries({
        queryKey: ["applicationResume"],
        exact: false,
      });
      queryClient.invalidateQueries({
        queryKey: ["events", eventId, "application"],
        exact: false,
      });
    },
  });

  const reset = useMutation({
    mutationFn: () => resetReviews(eventId),
    onSuccess: () => {
      queryClient.setQueryData(getEventQueryKey(eventId), (oldData) => {
        if (!oldData) return oldData;

        return {
          ...oldData,
          application_review_started: false,
        };
      });
      queryClient.invalidateQueries({
        queryKey: getAssignedApplicationsQueryKey(eventId),
      });
      queryClient.invalidateQueries({
        queryKey: ["applicationResume"],
        exact: false,
      });
      queryClient.invalidateQueries({
        queryKey: ["events", eventId, "application"],
        exact: false,
      });
    },
  });

  return {
    assign,
    reset,
  };
};
