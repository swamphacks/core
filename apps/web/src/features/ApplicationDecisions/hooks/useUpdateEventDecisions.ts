import { useUpdateEvent } from "@/features/Event/hooks/useUpdateEvent";
import { type Event } from "@/features/Event/schemas/event";
import { api } from "@/lib/ky";
import { toast } from "react-toastify";

interface checkAppReviewStatus {
  checkAppReviewStatus: (event: Event) => Promise<boolean>;
}

export const useCheckAppReviewStatus = (
  eventId: string,
): checkAppReviewStatus => {
  const { mutateAsync: updateEvent } = useUpdateEvent(eventId);

  const checkAppReviewStatus = async (event: Event) => {
    const { id, application_review_finished } = event;

    await api.post(`events/${eventId}/app-review-decision-status`, {
      json: {
        eventId: id,
      },
    });

    await updateEvent(
      { application_review_finished: application_review_finished },
      {
        onError: () =>
          toast.error(
            "Failed to check if application reviews are finished for this event.",
            { position: "bottom-right" },
          ),
      },
    );

    return application_review_finished;
  };

  return { checkAppReviewStatus };
};
