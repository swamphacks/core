import { api } from "@/lib/ky";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  reviewAssignmentsQueryKey,
  type ReviewAssignments,
} from "./useReviewAssignments";
import {
  applicationReviewQueryKey,
  type ParsedApplicationReview,
} from "./useApplicationReview";
import type { operations } from "@/lib/openapi/schema";
import {
  extendedApplicationQueryKey,
  type ParsedExtendedApplicationResponse,
} from "@/modules/Application/hooks/useExtendedApplication";

const requestAutoDecisionFn = async (
  applicationId: string,
  accept: boolean,
  justification: string,
) => {
  return await (
    await api.post<
      operations["request-auto-decision"]["responses"]["200"]["content"]["application/json"]
    >(`application/review/auto-decision`, {
      json: {
        applicationId,
        decision: accept ? "auto_accept" : "auto_reject",
        justification,
      },
    })
  ).json();
};

const deleteAutoDecisionRequestFn = async (requestId: string) => {
  await api.delete(`application/review/auto-decision`, {
    json: {
      requestId,
    },
  });
};

const submitReview = async (
  reviewId: string,
  applicationId: string,
  experienceRating: number,
  passionRating: number,
  notes: string,
) => {
  await api.post(`application/review`, {
    json: {
      id: reviewId,
      applicationId,
      experienceRating,
      passionRating,
      notes,
    },
  });

  return {
    reviewId,
    applicationId,
    experienceRating,
    passionRating,
    notes,
  };
};

export const useApplicationReviewActions = (
  applicationId: string,
  reviewId?: string,
) => {
  const queryClient = useQueryClient();

  const review = useMutation({
    mutationFn: ({
      reviewId,
      experienceRating,
      passionRating,
      notes,
    }: {
      reviewId: string;
      experienceRating: number;
      passionRating: number;
      notes: string;
    }) =>
      submitReview(
        reviewId,
        applicationId,
        experienceRating,
        passionRating,
        notes,
      ),
    onSuccess: ({ reviewId, experienceRating, passionRating, notes }) => {
      queryClient.setQueryData<ReviewAssignments>(
        reviewAssignmentsQueryKey,
        (oldData) =>
          oldData?.map((app) =>
            app.applicationId === applicationId
              ? { ...app, status: "completed" }
              : app,
          ),
      );

      queryClient.setQueryData<ParsedApplicationReview>(
        applicationReviewQueryKey(reviewId),
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
    mutationFn: async ({
      applicationId,
      accept,
      justification,
    }: {
      applicationId: string;
      accept: boolean;
      justification: string;
    }) => await requestAutoDecisionFn(applicationId, accept, justification),
    onSuccess: (res) => {
      if (reviewId) {
        queryClient.setQueryData<ParsedApplicationReview>(
          applicationReviewQueryKey(reviewId),
          (oldData) => {
            if (!oldData) return oldData;

            return {
              ...oldData,
              autoDecisionRequest: {
                applicationId,
                approved: res.approved,
                decidedBy: res.decidedBy,
                createdAt: res.createdAt,
                decision: res.decision,
                id: res.id,
                justification: res.justification,
              },
            };
          },
        );
      }

      queryClient.setQueryData<ParsedExtendedApplicationResponse>(
        extendedApplicationQueryKey(applicationId),
        (oldData) => {
          if (!oldData) return oldData;

          return {
            ...oldData,
            autoDecisionRequest: {
              applicationId,
              approved: res.approved,
              decidedBy: res.decidedBy,
              createdAt: res.createdAt,
              decision: res.decision,
              id: res.id,
              justification: res.justification,
            },
          };
        },
      );
    },
  });

  const deleteAutoDecisionRequest = useMutation({
    mutationFn: ({ requestId }: { requestId: string }) =>
      deleteAutoDecisionRequestFn(requestId),
    onSuccess: (_) => {
      if (reviewId) {
        queryClient.setQueryData<ParsedApplicationReview>(
          applicationReviewQueryKey(reviewId),
          (oldData) => {
            if (!oldData) return oldData;

            return {
              ...oldData,
              autoDecisionRequest: undefined,
            };
          },
        );
      }

      queryClient.setQueryData<ParsedExtendedApplicationResponse>(
        extendedApplicationQueryKey(applicationId),
        (oldData) => {
          if (!oldData) return oldData;

          return {
            ...oldData,
            autoDecisionRequest: undefined,
          };
        },
      );
    },
  });

  return { review, requestAutoDecision, deleteAutoDecisionRequest };
};
