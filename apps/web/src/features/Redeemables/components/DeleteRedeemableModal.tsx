import { Modal } from "@/components/ui/Modal";
import { Button } from "@/components/ui/Button";
import { api } from "@/lib/ky";
import { showToast } from "@/lib/toast/toast";
import { useQueryClient } from "@tanstack/react-query";
import { useState, useEffect, useContext } from "react";
import TablerAlertTriangle from "~icons/tabler/alert-triangle";
import { OverlayTriggerStateContext } from "react-aria-components";

interface DeleteRedeemableModalProps {
  eventId: string;
  redeemableId: string;
}

function DeleteRedeemableModal({
  eventId,
  redeemableId,
}: DeleteRedeemableModalProps) {
  const queryClient = useQueryClient();
  const [countdown, setCountdown] = useState(3);
  const state = useContext(OverlayTriggerStateContext)!;

  useEffect(() => {
    if (state.isOpen) {
      setCountdown(3);
    } else {
      return;
    }

    // Only start the interval if we are still counting down
    const timer = setInterval(() => {
      setCountdown((prev) => {
        if (prev <= 1) {
          clearInterval(timer);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    // Clean up the timer if the modal is closed or component unmounts
    return () => clearInterval(timer);
  }, [state.isOpen]);

  const handleDeleteRedeemable = async (
    eventId: string,
    redeemableId: string,
  ) => {
    try {
      await api.delete(`events/${eventId}/redeemables/${redeemableId}`);
      showToast({
        title: "Redeemable Deleted",
        message: "Successfully deleted redeemable.",
        type: "success",
      });
      queryClient.invalidateQueries({ queryKey: ["redeemables", eventId] });
    } catch (error) {
      console.error("Failed to delete redeemable", error);
      showToast({
        title: "Delete Failed",
        message: "Failed to delete redeemable. Please try again.",
        type: "error",
      });
    }
  };

  return (
    <Modal isDismissible className="p-6 rounded-xl shadow-lg max-w-sm w-full">
      <div className="flex flex-col items-center text-center">
        <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
          <TablerAlertTriangle className="h-6 w-6 text-red-600" />
        </div>

        <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
          Delete Redeemable?
        </h3>
        <div className="mt-4 rounded-lg bg-red-50 dark:bg-red-900/20 p-4 border border-red-100 dark:border-red-900/30">
          <p className="text-sm leading-relaxed text-red-800 dark:text-red-200 font-medium">
            This action is permanent.
          </p>
          <p className="mt-2 text-xs leading-relaxed text-red-700/80 dark:text-red-300/80">
            Deleting this will erase ALL records for future analytics.
            <span className="block mt-2 font-semibold">
              Do not delete redeemables that have already been used by guests,
              as we may use that data in the future.
            </span>
          </p>
        </div>
      </div>

      <div className="mt-6 flex justify-center">
        <Button
          variant="danger"
          isDisabled={countdown > 0}
          onPress={() => handleDeleteRedeemable(eventId, redeemableId)}
          className="w-full text-sm py-2.5 shadow-sm transition-all"
        >
          {countdown > 0 ? `Wait ${countdown}s...` : "Delete Redeemable"}
        </Button>
      </div>
    </Modal>
  );
}

export { DeleteRedeemableModal };
