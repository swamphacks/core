import { type Event } from "@/features/Event/schemas/event";
import { api } from "@/lib/ky";
import z from "zod";

interface checkAppReviewStatus {
  checkAppReviewStatus: (event: Event) => Promise<boolean>;
}

const CheckAppReviewStatusReponse = z.boolean();

export const useCheckAppReviewStatus = (
  eventId: string,
): checkAppReviewStatus => {
  const checkAppReviewStatus = async () => {
    const res = await api.post(`events/${eventId}/app-review-decision-status`, {
      json: {
        eventId: eventId,
      },
    });
    return CheckAppReviewStatusReponse.parse(await res.json());
  };

  return { checkAppReviewStatus };
};
