import { Modal } from "@/components/ui/Modal";
import { auth } from "@/lib/authClient";
import { Button } from "@/components/ui/Button";
import { api } from "@/lib/ky";
import { showToast } from "@/lib/toast/toast";
import { useQueryClient } from "@tanstack/react-query";
import { eventsQueryKey } from "../hooks/useEventsWithUserInfo";
import { myApplicationBaseKey } from "@/features/Application/hooks/useMyApplication";

interface EventAcceptanceWithdrawalModalProps {
  eventId: string;
}

function EventAcceptanceWithdrawalModal({
  eventId,
}: EventAcceptanceWithdrawalModalProps) {
  const queryClient = useQueryClient();

  const { data: userData } = auth.useUser();
  if (!userData?.user) {
    return <div>Loading...</div>;
  }

  const { user } = userData;

  const handleWithdrawAcceptance = async (userId: string, eventId: string) => {
    try {
      await api.patch(
        `events/${eventId}/application/withdraw-acceptance?userId=${userId}`,
      );
      showToast({
        title: "Acceptance Withdrawn",
        message: "Successfully Withdrawn Application.",
        type: "success",
      });
      await queryClient.invalidateQueries({
        queryKey: eventsQueryKey,
      });
      await queryClient.invalidateQueries({
        queryKey: [...myApplicationBaseKey, eventId],
      });
    } catch (error) {
      console.error("Failed to withdraw application", error);
      showToast({
        title: "Failed to Withdraw",
        message: "Failed to Withdraw Acceptance. Please try again.",
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
          Withdraw Acceptance?
        </h3>
        <p className="mt-1 text-sm text-gray-500">
          Are you sure? This action cannot be undone.
        </p>
        {/* <p className="mt-1 text-sm text-gray-500">
          You can still join the waitlist after withdrawing.
        </p> */}
      </div>
      <div className="mt-3 flex justify-center gap-3">
        <Button
          variant="danger"
          onPress={() => handleWithdrawAcceptance(user.userId, eventId)}
        >
          Withdraw
        </Button>
      </div>
    </Modal>
  );
}

export { EventAcceptanceWithdrawalModal as EventAcceptanceWithdrawalModal };
