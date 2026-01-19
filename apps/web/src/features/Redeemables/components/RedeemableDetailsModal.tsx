import { Modal } from "@/components/ui/Modal";
import { OverlayTriggerStateContext, Button } from "react-aria-components";
import { useContext, useRef, useEffect, useState } from "react";
import { Text } from "react-aria-components";
import { useRedeemRedeemable, getUserByRFID } from "../hooks/useRedeemables";
import { showToast } from "@/lib/toast/toast";
import TablerRefresh from "~icons/tabler/refresh";

export interface RedeemableDetailsModalProps {
  id: string;
  name: string;
  totalStock: number;
  maxUserAmount: number;
  totalRedeemed: number;
  eventId: string;
}

export function RedeemableDetailsModal({
  id,
  name,
  totalStock,
  maxUserAmount,
  totalRedeemed,
  eventId,
}: RedeemableDetailsModalProps) {
  const state = useContext(OverlayTriggerStateContext)!;
  const remaining = totalStock - totalRedeemed;
  const percentageRemaining =
    totalStock > 0 ? ((remaining / totalStock) * 100).toFixed(0) : 0;

  const [rfidInput, setRfidInput] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);
  const { mutateAsync: redeemRedeemable, isPending } = useRedeemRedeemable(
    eventId,
    id,
  );

  // Auto-focus the invisible input when modal opens
  useEffect(() => {
    if (state.isOpen && inputRef.current) {
      inputRef.current.focus();
    }
  }, [state.isOpen]);

  // Handle RFID input - when 10 digits are entered, trigger redemption
  useEffect(() => {
    if (rfidInput.length === 10) {
      handleRedeem();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [rfidInput]);

  const handleRedeem = async () => {
    if (rfidInput.length !== 10) {
      showToast({
        title: "Invalid RFID",
        message: "RFID must be exactly 10 digits",
        type: "error",
      });
      return;
    }

    try {
      // First, convert RFID to UserID
      const userResponse = await getUserByRFID(eventId, rfidInput);

      // Then, redeem using the user ID
      await redeemRedeemable(userResponse.user_id);

      showToast({
        title: "Success",
        message: "Redeemable redeemed successfully!",
        type: "success",
      });
      setRfidInput("");
      if (inputRef.current) {
        inputRef.current.focus();
      }
    } catch (error) {
      showToast({
        title: "Error",
        message: "Failed to redeem. User not found or already redeemed.",
        type: "error",
      });
      setRfidInput("");
      if (inputRef.current) {
        inputRef.current.focus();
      }
      console.log("Redeem error:", error);
    }
  };

  const handleRefresh = () => {
    setRfidInput("");
    if (inputRef.current) {
      inputRef.current.focus();
    }
  };

  return (
    <Modal title={name} responsive="standard" size="lg" isDismissible>
      <div className="flex flex-col gap-6">
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <Text className="text-text-secondary">Total Stock</Text>
            <Text className="text-text-main font-semibold">{totalStock}</Text>
          </div>

          <div className="flex justify-between items-center">
            <Text className="text-text-secondary">Redeemed</Text>
            <Text className="text-text-main font-semibold">
              {totalRedeemed}
            </Text>
          </div>

          <div className="flex justify-between items-center">
            <Text className="text-text-secondary">Remaining</Text>
            <Text className="text-text-main font-semibold">{remaining}</Text>
          </div>

          <div className="flex justify-between items-center">
            <Text className="text-text-secondary">Max per User</Text>
            <Text className="text-text-main font-semibold">
              {maxUserAmount}
            </Text>
          </div>

          <div className="space-y-2">
            <div className="flex justify-between items-center">
              <Text className="text-text-secondary">Stock Status</Text>
              <Text className="text-text-main font-semibold">
                {percentageRemaining}% remaining
              </Text>
            </div>
            <div className="w-full bg-neutral-200 dark:bg-neutral-700 rounded-full h-3">
              <div
                className="bg-violet-600 h-3 rounded-full transition-all"
                style={{ width: `${percentageRemaining}%` }}
              />
            </div>
          </div>
        </div>

        <div className="pt-4 border-t border-neutral-200 dark:border-neutral-700">
          <div className="flex flex-col gap-4">
            <div className="flex items-center justify-between">
              <Text className="text-sm text-text-secondary">
                Redeemable ID: {id}
              </Text>
              <Button
                onPress={handleRefresh}
                className="flex items-center gap-2 px-3 py-1.5 text-sm bg-neutral-100 dark:bg-neutral-800 hover:bg-neutral-200 dark:hover:bg-neutral-700 rounded transition-colors"
              >
                <TablerRefresh className="w-4 h-4" />
                Refresh
              </Button>
            </div>

            {/* Invisible RFID input field */}
            <input
              ref={inputRef}
              type="text"
              value={rfidInput}
              onChange={(e) => {
                // Only allow digits and limit to 10
                const value = e.target.value.replace(/\D/g, "").slice(0, 10);
                setRfidInput(value);
              }}
              className="sr-only"
              aria-label="RFID input"
              disabled={isPending}
              placeholder=""
            />

            {rfidInput.length > 0 && (
              <div className="text-xs text-text-secondary">
                RFID: {"•".repeat(rfidInput.length)}
              </div>
            )}
          </div>
        </div>
      </div>
    </Modal>
  );
}
