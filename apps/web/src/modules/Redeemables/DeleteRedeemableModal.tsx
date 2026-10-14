import { Modal } from "@/components/ui/Modal";
import { Button } from "@/components/ui/Button";
import { showToast } from "@/lib/toast/toast";
import { useState, useEffect, useContext } from "react";
import TablerAlertTriangle from "~icons/tabler/alert-triangle";
import { OverlayTriggerStateContext } from "react-aria-components";
import {
  apiErrorMessage,
  useDeleteRedeemable,
  type Redeemable,
} from "./hooks/useRedeemables";

interface DeleteRedeemableModalProps {
  redeemable: Redeemable;
}

function DeleteRedeemableModal({ redeemable }: DeleteRedeemableModalProps) {
  const [countdown, setCountdown] = useState(3);
  const state = useContext(OverlayTriggerStateContext)!;
  const { mutateAsync: deleteRedeemable, isPending } = useDeleteRedeemable();

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

  const handleDelete = async () => {
    try {
      await deleteRedeemable(redeemable.id);
      showToast({
        title: "Item Deleted",
        message: `Deleted ${redeemable.name}.`,
        type: "success",
      });
      state.close();
    } catch (error) {
      showToast({
        title: "Delete Failed",
        message: await apiErrorMessage(
          error,
          "Failed to delete the item. Please try again.",
        ),
        type: "error",
      });
    }
  };

  const redeemed = redeemable.total_redeemed;

  return (
    <Modal isDismissible className="p-6 rounded-xl shadow-lg max-w-sm w-full">
      <div className="flex flex-col items-center text-center">
        <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
          <TablerAlertTriangle className="h-6 w-6 text-red-600" />
        </div>

        <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">
          Delete {redeemable.name}?
        </h3>
        <div className="mt-4 rounded-lg bg-red-50 dark:bg-red-900/20 p-4 border border-red-100 dark:border-red-900/30">
          <p className="text-sm leading-relaxed text-red-800 dark:text-red-200 font-medium">
            This action is permanent.
          </p>
          {redeemed > 0 ? (
            <p className="mt-2 text-xs leading-relaxed text-red-700/80 dark:text-red-300/80">
              It erases {redeemed}{" "}
              {redeemed === 1 ? "redemption" : "redemptions"} of this item. If
              you recreate it, everyone who already got one can redeem it again.
              <span className="block mt-2 font-semibold">
                To change the name, type or limit, edit the item instead.
              </span>
            </p>
          ) : (
            <p className="mt-2 text-xs leading-relaxed text-red-700/80 dark:text-red-300/80">
              Nobody has redeemed this item yet. Staff phones keep showing it
              until the app is reopened.
            </p>
          )}
        </div>
      </div>

      <div className="mt-6 flex justify-center">
        <Button
          variant="danger"
          isDisabled={countdown > 0}
          isPending={isPending}
          onPress={handleDelete}
          className="w-full text-sm py-2.5 shadow-sm transition-all"
        >
          {countdown > 0 ? `Wait ${countdown}s...` : "Delete Item"}
        </Button>
      </div>
    </Modal>
  );
}

export { DeleteRedeemableModal };
