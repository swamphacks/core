import { Modal } from "@/components/ui/Modal";
import {
  OverlayTriggerStateContext,
  DialogTrigger,
} from "react-aria-components";
import { Button } from "@/components/ui/Button";
import TablerTrash from "~icons/tabler/trash";
import { useContext, useRef, useEffect, useState } from "react";
import { Text } from "react-aria-components";
import {
  useRedeemRedeemable,
  getUserByRFID,
  useUpdateRedeemable,
} from "../hooks/useRedeemables";
import { DeleteRedeemableModal } from "./DeleteRedeemableModal";
import { Scanner, type IDetectedBarcode } from "@yudiel/react-qr-scanner";
import { showToast } from "@/lib/toast/toast";
import TablerRefresh from "~icons/tabler/refresh";
import { parseQrIntent } from "@/lib/qr-intents/parse";
import { Intent } from "@/lib/qr-intents/intent";

export interface RedeemableDetailsModalProps {
  id: string;
  name: string;
  totalStock: number;
  maxUserAmount: number;
  totalRedeemed: number;
  eventId: string;
}

type ModalMode = "details" | "qr" | "edit";

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
  const [mode, setMode] = useState<ModalMode>("details");
  const inputRef = useRef<HTMLInputElement>(null);
  const { mutateAsync: redeemRedeemable, isPending } = useRedeemRedeemable(
    eventId,
    id,
  );
  const [formData, setFormData] = useState({
    name: name,
    totalStock: totalStock,
    maxUserAmount: maxUserAmount,
  });
  const [isUpdating, setIsUpdating] = useState(false);
  const { mutateAsync: updateRedeemable } = useUpdateRedeemable(eventId, id);

  useEffect(() => {
    setFormData({ name, totalStock, maxUserAmount });
  }, [name, totalStock, maxUserAmount]);

  const handleReset = () => {
    setFormData({ name, totalStock, maxUserAmount });
    showToast({ title: "Reset", message: "Changes discarded", type: "info" });
  };

  const handleConfirmEdit = async () => {
    setIsUpdating(true);
    try {
      await updateRedeemable({
        name: formData.name,
        total_stock: formData.totalStock,
        max_user_amount: formData.maxUserAmount,
      });

      console.log("Updating redeemable:", id, formData);
      showToast({
        title: "Success",
        message: "Redeemable updated successfully!",
        type: "success",
      });
      setMode("details");
    } catch (error) {
      showToast({ title: "Error", message: "Failed to update", type: "error" });
      console.log("Update redeemable error:", error);
    } finally {
      setIsUpdating(false);
    }
  };
  // Auto-focus the invisible input when modal opens
  useEffect(() => {
    if (state.isOpen) {
      setMode("details");

      setRfidInput("");

      if (inputRef.current) {
        inputRef.current.focus();
      }
    }
  }, [state.isOpen]);

  // Handle RFID input - when 10 digits are entered, trigger redemption
  useEffect(() => {
    if (rfidInput.length === 10) {
      handleRedeem();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [rfidInput]);

  const onScan = async (scannedData: IDetectedBarcode[]) => {
    const res = parseQrIntent(scannedData[0].rawValue);
    if (res.ok && res.value.intent === Intent.IDENT) {
      try {
        await redeemRedeemable(res.value.user_id);

        showToast({
          title: "Success",
          message: "Redeemable redeemed successfully!",
          type: "success",
        });
      } catch (error) {
        showToast({
          title: "Error",
          message:
            "Failed to redeem with QR. User not found or already redeemed.",
          type: "error",
        });
        console.log("Redeem error:", error);
      }
    }
  };

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
        message:
          "Failed to redeem with RFID. User not found or already redeemed.",
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
      {mode === "details" && (
        <div className="flex flex-col gap-4">
          <div className="flex justify-start gap-4">
            <button
              onClick={() => setMode("qr")}
              className="text-purple-600 underline decoration-purple-600/30 underline-offset-4 transition-colors hover:text-purple-800 hover:decoration-purple-800"
            >
              Use QR
            </button>

            <button
              onClick={() => setMode("edit")}
              className="text-purple-600 underline decoration-purple-600/30 underline-offset-4 transition-colors hover:text-purple-800 hover:decoration-purple-800"
            >
              Edit
            </button>
          </div>
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
      )}
      {mode === "qr" && (
        <div className="flex flex-col gap-6">
          <div className="justify-start">
            <button
              onClick={() => setMode("details")}
              className="text-purple-600 underline decoration-purple-600/30 underline-offset-4 transition-colors hover:text-purple-800 hover:decoration-purple-800"
            >
              Back to Details
            </button>
          </div>
          <section className="w-full">
            <div className="mx-auto overflow-hidden rounded-xl bg-black aspect-square w-full max-w-md shadow-lg">
              <Scanner
                classNames={{ container: "w-full h-full" }}
                sound={false}
                formats={["qr_code"]}
                onScan={onScan}
                scanDelay={500}
                onError={(err) => console.log(err)}
              />
            </div>
            <p className="text-sm text-text-secondary my-6">
              Point the camera at a user&apos;s QR code to redeem.
            </p>
          </section>
        </div>
      )}
      {mode === "edit" && (
        <div className="flex flex-col gap-6">
          <div className="flex justify-between">
            <button
              onClick={() => setMode("details")}
              className="text-purple-600 underline decoration-purple-600/30 underline-offset-4 transition-colors hover:text-purple-800"
            >
              Back to Details
            </button>
            <DialogTrigger>
              <Button variant="danger" className="aspect-square p-2">
                <TablerTrash className="h-4 w-4" />
              </Button>
              <DeleteRedeemableModal eventId={eventId} redeemableId={id} />
            </DialogTrigger>
          </div>

          <div className="space-y-5">
            {/* Name Input */}
            <div className="flex flex-col gap-2">
              <label className="text-sm font-medium text-text-secondary">
                Redeemable Name
              </label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) =>
                  setFormData({ ...formData, name: e.target.value })
                }
                className="w-full px-3 py-2 bg-neutral-50 dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-700 rounded-lg focus:ring-2 focus:ring-purple-500 outline-none transition-all"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              {/* Total Stock Input */}
              <div className="flex flex-col gap-2">
                <label className="text-sm font-medium text-text-secondary">
                  Total Stock
                </label>
                <input
                  type="number"
                  value={formData.totalStock}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      totalStock: Number(e.target.value),
                    })
                  }
                  className="w-full px-3 py-2 bg-neutral-50 dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-700 rounded-lg focus:ring-2 focus:ring-purple-500 outline-none transition-all"
                />
              </div>

              {/* Max Per User Input */}
              <div className="flex flex-col gap-2">
                <label className="text-sm font-medium text-text-secondary">
                  Max Per User
                </label>
                <input
                  type="number"
                  value={formData.maxUserAmount}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      maxUserAmount: Number(e.target.value),
                    })
                  }
                  className="w-full px-3 py-2 bg-neutral-50 dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-700 rounded-lg focus:ring-2 focus:ring-purple-500 outline-none transition-all"
                />
              </div>
            </div>
          </div>
          <div className="flex items-center gap-3 pt-4 border-t border-neutral-200 dark:border-neutral-700">
            <Button
              onPress={handleReset}
              className="flex-1 px-4 py-2 font-medium text-sm text-neutral-700 dark:text-neutral-200 bg-neutral-100 dark:bg-neutral-800 hover:bg-neutral-200 dark:hover:bg-neutral-700 rounded transition-colors cursor-pointer"
            >
              Reset Changes
            </Button>
            <Button
              onPress={handleConfirmEdit}
              isDisabled={
                isUpdating ||
                (formData.name === name &&
                  formData.totalStock === totalStock &&
                  formData.maxUserAmount === maxUserAmount)
              }
              className="flex-1 px-4 py-2 text-white text-sm bg-violet-600 hover:bg-violet-700 disabled:opacity-50 rounded font-medium cursor-pointer"
            >
              {isUpdating ? "Saving..." : "Confirm Changes"}
            </Button>
          </div>
        </div>
      )}
    </Modal>
  );
}
