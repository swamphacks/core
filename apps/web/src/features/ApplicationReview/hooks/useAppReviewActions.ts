import { type AssignedApplications } from "@/features/Application/hooks/useAssignedApplication";
import { api } from "@/lib/ky";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { getAssignedApplicationsQueryKey } from "./useAssignedApplications";
import {
  type Application,
  getApplicationQueryKey,
} from "@/features/Application/hooks/useApplication";

const submitReview = async (
  eventId: string,
  applicationId: string,
  experienceRating: number,
  passionRating: number,
) => {
  await api.post(`events/${eventId}/application/${applicationId}/review`, {
    json: {
      experience_rating: experienceRating,
      passion_rating: passionRating,
    },
  });

  return {
    experienceRating,
    passionRating,
  };
};

export const useAppReviewActions = (eventId: string, applicationId: string) => {
  const queryClient = useQueryClient();

  const review = useMutation({
    mutationFn: ({
      experienceRating,
      passionRating,
    }: {
      experienceRating: number;
      passionRating: number;
    }) => submitReview(eventId, applicationId, experienceRating, passionRating),
    onSuccess: ({ experienceRating, passionRating }) => {
      // Mutate the review data in the cache if needed
      queryClient.setQueryData<AssignedApplications>(
        getAssignedApplicationsQueryKey(eventId),
        (oldData) =>
          oldData?.map((app) =>
            app.user_id === applicationId
              ? { ...app, status: "completed" }
              : app,
          ),
      );

      queryClient.setQueryData<Application>(
        getApplicationQueryKey(eventId, applicationId),
        (oldData) => {
          if (!oldData) return oldData;

          return {
            ...oldData,
            experience_rating: experienceRating,
            passion_rating: passionRating,
          };
        },
      );
    },
  });

  return { review };
};
