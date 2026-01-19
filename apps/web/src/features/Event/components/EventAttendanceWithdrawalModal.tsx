import { Modal } from "@/components/ui/Modal";
import { auth } from "@/lib/authClient";
import { Button } from "@/components/ui/Button";
import { api } from "@/lib/ky";
import { showToast } from "@/lib/toast/toast";
import { useQueryClient } from "@tanstack/react-query";
import { eventsQueryKey } from "../hooks/useEventsWithUserInfo";
import { myApplicationBaseKey } from "@/features/Application/hooks/useMyApplication";
import { useNavigate } from "@tanstack/react-router";

interface EventAttendanceWithdrawalModalProps {
  eventId: string;
}

function EventAttendanceWithdrawalModal({
  eventId,
}: EventAttendanceWithdrawalModalProps) {
  const queryClient = useQueryClient();

  const navigate = useNavigate();

  const { data: userData } = auth.useUser();
  if (!userData?.user) {
    return <div>Loading...</div>;
  }

  const { user } = userData;

  const handleWithdrawAttendance = async (userId: string, eventId: string) => {
    try {
      await api.patch(
        `events/${eventId}/application/withdraw-attendance?userId=${userId}`,
      );
      showToast({
        title: "Attendance Withdrawn",
        message: "Successfully Withdrawn Attendance.",
        type: "success",
      });
      await queryClient.invalidateQueries({
        queryKey: eventsQueryKey,
      });
      await queryClient.invalidateQueries({
        queryKey: [...myApplicationBaseKey, eventId],
      });
      navigate({
        to: "/portal",
      });
    } catch (error) {
      console.error("Failed to withdraw attendance", error);
      showToast({
        title: "Failed to Withdraw",
        message: "Failed to Withdraw Attendance. Please try again.",
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
          Withdraw Attendance?
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
          onPress={() => handleWithdrawAttendance(user.userId, eventId)}
        >
          Withdraw
        </Button>
      </div>
    </Modal>
  );
}

export { EventAttendanceWithdrawalModal as EventAttendanceWithdrawalModal };
