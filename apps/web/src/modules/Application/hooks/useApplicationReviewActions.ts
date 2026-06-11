import { api } from "@/lib/ky";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { assignedApplicationsQueryKey } from "./useAssignedApplications";
import type {
  ApplicationReviewDetails,
  AssignedApplications,
} from "@/lib/openapi/types";
import { applicationReviewDetailsQueryKey } from "./useApplicationForReview";

const requestAutoDecisionFn = async (
  applicationId: string,
  accept: boolean,
  justification: string,
) => {
  await api.post(`application/review/auto-decision`, {
    json: {
      applicationId,
      decision: accept ? "auto_accept" : "auto_reject",
      justification,
    },
  });
};

const deleteAutoDecisionRequestFn = async (requestId: string) => {
  await api.delete(`application/review/auto-decision`, {
    json: {
      requestId,
    },
  });
};

const submitReview = async (
  applicationId: string,
  experienceRating: number,
  passionRating: number,
  notes: string,
) => {
  await api.post(`application/review/${applicationId}`, {
    json: {
      experienceRating,
      passionRating,
      notes,
    },
  });

  return {
    experienceRating,
    passionRating,
    notes,
  };
};

export const useApplicationReviewActions = (applicationId: string) => {
  const queryClient = useQueryClient();

  const review = useMutation({
    mutationFn: ({
      experienceRating,
      passionRating,
      notes,
    }: {
      experienceRating: number;
      passionRating: number;
      notes: string;
    }) => submitReview(applicationId, experienceRating, passionRating, notes),
    onSuccess: ({ experienceRating, passionRating, notes }) => {
      queryClient.setQueryData<AssignedApplications>(
        assignedApplicationsQueryKey,
        (oldData) =>
          oldData?.map((app) =>
            app.applicationId === applicationId
              ? { ...app, status: "completed" }
              : app,
          ),
      );

      queryClient.setQueryData<ApplicationReviewDetails>(
        applicationReviewDetailsQueryKey(applicationId),
        (oldData) => {
          if (!oldData) return oldData;

          return {
            ...oldData,
            experienceRating,
            passionRating,
            notes,
          };
        },
      );
    },
  });

  const requestAutoDecision = useMutation({
    mutationFn: ({
      applicationId,
      accept,
      justification,
    }: {
      applicationId: string;
      accept: boolean;
      justification: string;
    }) => requestAutoDecisionFn(applicationId, accept, justification),
    onSuccess: (_, { applicationId, accept, justification }) => {
      queryClient.setQueryData<ApplicationReviewDetails>(
        applicationReviewDetailsQueryKey(applicationId),
        (oldData) => {
          if (!oldData) return oldData;

          return {
            ...oldData,
            autoDecision: accept ? "auto_accept" : "auto_reject",
            justification,
          };
        },
      );
    },
  });

  const deleteAutoDecisionRequest = useMutation({
    mutationFn: ({ requestId }: { requestId: string }) =>
      deleteAutoDecisionRequestFn(requestId),
    onSuccess: (_) => {
      queryClient.setQueryData<ApplicationReviewDetails>(
        applicationReviewDetailsQueryKey(applicationId),
        (oldData) => {
          if (!oldData) return oldData;

          return {
            ...oldData,
            autoDecision: null,
            autoDecisionRequestId: null,
          };
        },
      );
    },
  });

  return { review, requestAutoDecision, deleteAutoDecisionRequest };
};
