import { Modal } from "@/components/ui/Modal";
import { auth } from "@/lib/authClient";
import { Button } from "@/components/ui/Button";
import { api } from "@/lib/ky";
import { showToast } from "@/lib/toast/toast";
import { useQueryClient } from "@tanstack/react-query";
import { eventsQueryKey } from "../hooks/useEventsWithUserInfo";
import { myApplicationBaseKey } from "@/features/Application/hooks/useMyApplication";

interface EventWaitlistModalProps {
  eventId: string;
}

function EventWaitlistModal({ eventId }: EventWaitlistModalProps) {
  const queryClient = useQueryClient();

  const { data: userData } = auth.useUser();
  if (!userData?.user) {
    return <div>Loading...</div>;
  }

  const { user } = userData;

  const handleJoinWaitlist = async (userId: string, eventId: string) => {
    try {
      await api.patch(
        `events/${eventId}/application/join-waitlist?userId=${userId}`,
      );
      showToast({
        title: "Waitlist Joined",
        message: "Successfully Joined Waitlist.",
        type: "success",
      });
      await queryClient.invalidateQueries({
        queryKey: eventsQueryKey,
      });
      await queryClient.invalidateQueries({
        queryKey: [...myApplicationBaseKey, eventId],
      });
    } catch (error) {
      console.error("Failed to join waitlist", error);
      showToast({
        title: "Join Failed",
        message: "Failed to Join Waitlist. Please try again.",
        type: "error",
      });
    }
  };

  return (
    <Modal
      isDismissible
      className="bg-grey p-6 rounded-xl shadow-lg max-w-sm w-full"
    >
      <div className="text-center">
        <h3 className="text-lg font-semibold dark:text-gray-100">
          Join Waitlist?
        </h3>
        <p className="mt-1 text-sm text-gray-500">
          Unfortunately, we weren't able to accomodate everybody who applied
          this year. Join the waitlist, and we will let you know if a spot opens
          up!
        </p>
      </div>
      <div className="mt-3 flex justify-center gap-3">
        <Button
          variant="success"
          onPress={() => handleJoinWaitlist(user.userId, eventId)}
        >
          Join Waitlist
        </Button>
      </div>
    </Modal>
  );
}

export { EventWaitlistModal };
