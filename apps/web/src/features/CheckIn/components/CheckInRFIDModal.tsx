import { Modal } from "@/components/ui/Modal";
import { Button } from "@/components/ui/Button";
import { TextField } from "@/components/ui/TextField";
import { useState, useEffect } from "react";
import { Avatar } from "@/components/ui/Avatar";
import RoleBadge from "@/features/EventAdmin/components/RoleBadge";
import TablerCheck from "~icons/tabler/check";
import TablerX from "~icons/tabler/x";
import { toast } from "react-toastify";
import { api } from "@/lib/ky";
import { useQueryClient } from "@tanstack/react-query";
import { HTTPError } from "ky";
import { getEventStaffUsersQueryKey } from "@/features/PlatformAdmin/EventManager/hooks/useEventUsers";
import type { EventUser } from "@/features/PlatformAdmin/EventManager/hooks/useEventUsers";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  user: EventUser;
  eventId: string;
}

function CheckListIcon({ isMet }: { isMet: boolean }) {
  return isMet ? (
    <TablerCheck className="inline h-5 w-5 text-green-600" />
  ) : (
    <TablerX className="inline h-5 w-5 text-red-600" />
  );
}

export default function CheckInRFIDModal({
  isOpen,
  onClose,
  user,
  eventId,
}: Props) {
  const [rfid, setRfid] = useState<string>("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const queryClient = useQueryClient();

  // Auto-focus RFID input when modal opens
  useEffect(() => {
    if (isOpen) {
      setRfid("");
      // Small delay to ensure modal is fully rendered
      setTimeout(() => {
        const hiddenInput = document.getElementById("rfid-hidden-input");
        if (hiddenInput) {
          hiddenInput.focus();
        }
      }, 100);
    }
  }, [isOpen]);

  const isCheckedIn =
    user.checked_in_at !== null && user.checked_in_at !== undefined;
  const isAttendee = user.event_role === "attendee";
  const hasRfid = rfid.trim() !== "";

  const isValidForCheckIn = isAttendee && !isCheckedIn && hasRfid;

  const handleCheckIn = async () => {
    if (!isValidForCheckIn) {
      toast.error("Cannot check in: requirements not met.", {
        position: "bottom-right",
      });
      return;
    }

    setIsSubmitting(true);
    try {
      await api.post(`events/${eventId}/checkin`, {
        json: {
          user_id: user.id,
          rfid: rfid.trim(),
        },
      });

      toast.success("User checked in successfully!", {
        position: "bottom-right",
      });

      // Invalidate and refetch the event users query to update the table
      await queryClient.invalidateQueries({
        queryKey: getEventStaffUsersQueryKey(eventId),
      });

      onClose();
      setRfid("");
    } catch (error) {
      console.error("Check-in error:", error);

      // Handle HTTP errors with specific error messages
      if (error instanceof HTTPError) {
        try {
          const errorBody = (await error.response.json()) as {
            error?: string;
            message?: string;
          };
          const errorMessage =
            errorBody.message || "Failed to check in user. Please try again.";

          toast.error(errorMessage, {
            position: "bottom-right",
          });
        } catch {
          // If we can't parse the error response, show a generic message
          toast.error("Error during check-in. Please try again.", {
            position: "bottom-right",
          });
        }
      } else {
        toast.error("Error during check-in. Please try again.", {
          position: "bottom-right",
        });
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && isValidForCheckIn && !isSubmitting) {
      e.preventDefault();
      handleCheckIn();
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onOpenChange={(isOpen) => {
        if (!isOpen) {
          onClose();
        }
      }}
      title="Check In Attendee"
      size="md"
      padding="lg"
      responsive="standard"
    >
      <div className="flex flex-col gap-6">
        {/* User Info Section */}
        <div className="flex flex-col items-center gap-4 pb-4 border-b border-border">
          <Avatar
            src={user.image ?? undefined}
            alt={user.name}
            size="lg"
            shape="circle"
          />
          <div className="flex flex-col items-center gap-2">
            <h3 className="text-lg font-semibold">{user.name}</h3>
            <div className="flex items-center gap-2">
              <span className="text-sm text-text-secondary">
                {user.preferred_email ?? user.email ?? "No email"}
              </span>
            </div>
            <RoleBadge role={user.event_role} />
          </div>
        </div>

        {/* Check-In Status */}
        {isCheckedIn && (
          <div className="p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-md">
            <p className="text-sm text-yellow-800 dark:text-yellow-200">
              ⚠️ This user has already been checked in.
            </p>
          </div>
        )}

        {/* Check-In Requirements */}
        <div className="space-y-2">
          <h4 className="font-medium text-sm">Check-In Requirements:</h4>
          <div className="space-y-1.5 text-sm text-text-secondary flex flex-col">
            <span className="flex flex-row items-center gap-2">
              <CheckListIcon isMet={isAttendee} />
              <p>User must be an attendee</p>
            </span>
            <span className="flex flex-row items-center gap-2">
              <CheckListIcon isMet={!isCheckedIn} />
              <p>User has not checked in yet</p>
            </span>
            <span className="flex flex-row items-center gap-2">
              <CheckListIcon isMet={hasRfid} />
              {hasRfid ? (
                <p>
                  RFID scanned:{" "}
                  <span className="font-mono font-medium">{rfid}</span>
                </p>
              ) : (
                <p>Scan or enter RFID tag</p>
              )}
            </span>
          </div>
        </div>

        {/* RFID Input */}
        <div className="space-y-2">
          <label htmlFor="rfid-input" className="text-sm font-medium">
            RFID Tag
          </label>
          <TextField
            id="rfid-input"
            value={rfid}
            onChange={setRfid}
            onKeyDown={handleKeyDown}
            placeholder="Scan RFID tag or enter manually..."
            aria-label="RFID tag input"
            className="font-mono"
          />
          <p className="text-xs text-text-secondary">
            Scan the RFID tag or type the RFID value manually
          </p>
        </div>

        {/* Hidden input for RFID scanner - captures scanner input automatically */}
        <input
          type="text"
          id="rfid-hidden-input"
          onChange={(e) => setRfid(e.target.value)}
          inputMode="none"
          autoFocus
          className="opacity-0 absolute h-0 w-0 overflow-hidden"
          aria-hidden="true"
        />

        {/* Action Buttons */}
        <div className="flex flex-col gap-2 pt-2">
          <Button
            variant="primary"
            onPress={handleCheckIn}
            isDisabled={!isValidForCheckIn || isSubmitting}
            className="w-full"
          >
            {isSubmitting ? "Checking In..." : "Check In"}
          </Button>
          <Button
            variant="secondary"
            onPress={onClose}
            isDisabled={isSubmitting}
            className="w-full"
          >
            Cancel
          </Button>
        </div>
      </div>
    </Modal>
  );
}
