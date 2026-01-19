import { Avatar } from "@/components/ui/Avatar";
import { Modal } from "@/components/ui/Modal";
import { CheckInBadge } from "./CheckInBadge";
import { Button } from "@/components/ui/Button";
import { useState } from "react";
import { useUserEventInfo } from "../hooks/useUserEventInfo";
import RoleBadge from "@/features/EventAdmin/components/RoleBadge";
import TablerCheck from "~icons/tabler/check";
import TablerX from "~icons/tabler/x";
import { toast } from "react-toastify";
import { api } from "@/lib/ky";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  userId: string | null; // only needed after QR scan
  eventId: string;
}

function CheckListIcon({ isMet }: { isMet: boolean }) {
  return isMet ? (
    <TablerCheck className="inline h-5 w-5 text-green-600" />
  ) : (
    <TablerX className="inline h-5 w-5 text-red-600" />
  );
}

export default function CheckInModal({
  isOpen,
  onClose,
  userId,
  eventId,
}: Props) {
  const userInfo = useUserEventInfo(eventId, userId);
  const [rfid, setRfid] = useState<string | null>(null);

  const isLoading = !userInfo;

  const isValidForCheckIn =
    userInfo.data &&
    userInfo.data.event_role === "attendee" &&
    !userInfo.data.checked_in_at &&
    rfid &&
    rfid.trim() !== "";

  const onCheckIn = async () => {
    // validate all conditions are met
    if (!isValidForCheckIn || !userId) {
      toast.error("Cannot check in: requirements not met.", {
        position: "bottom-right",
      });
      return;
    }

    try {
      const res = await api.post(`events/${eventId}/checkin`, {
        json: {
          user_id: userId,
          rfid: rfid,
        },
      });

      if (res.ok) {
        toast.success("User checked in successfully!", {
          position: "bottom-right",
        });
        onClose();
        setRfid(null);
        userInfo.refetch();
      } else {
        toast.error("Failed to check in user. Please try again.", {
          position: "bottom-right",
        });
      }
    } catch (error) {
      console.error(`For support: ${error}`);

      toast.error("Error during check-in. Please try again.", {
        position: "bottom-right",
      });
    }
  };

  return (
    <Modal
      responsive="sheet"
      dialogAriaLabel="checkin-fields"
      isOpen={isOpen}
      isDismissible
      onOpenChange={(isOpen) => {
        if (!isOpen) {
          onClose();
          setRfid(null);
        }
      }}
    >
      {isLoading || (!userInfo.data && !userInfo.isError) ? (
        <p>Loading...</p>
      ) : userInfo.isError || userInfo.error || !userInfo.data ? (
        <p>Error loading user info: {userInfo.error?.message}</p>
      ) : (
        <div className="flex flex-col gap-6">
          <div className="flex items-center gap-4">
            <Avatar
              src={
                userInfo.data.image ||
                "https://static.vecteezy.com/system/resources/thumbnails/009/292/244/small/default-avatar-icon-of-social-media-user-vector.jpg"
              }
              size="lg"
            />

            <div className="flex flex-col flex-1 min-w-0">
              <h1 className="text-md font-medium truncate">
                {userInfo.data.name}
              </h1>
              <p className="text-sm text-text-secondary truncate">
                {userInfo.data.email}
              </p>

              <div className="mt-2 flex flex-wrap gap-2">
                {userInfo.data.event_role && (
                  <RoleBadge role={userInfo.data.event_role} />
                )}

                <CheckInBadge
                  isCheckedIn={userInfo.data.checked_in_at ? true : false}
                />
              </div>
            </div>
          </div>

          <p className="sm:hidden text-sm text-yellow-700 bg-yellow-50 border border-yellow-200 rounded-md px-3 py-2">
            On mobile, check-in requires an attached RFID scanner.
          </p>

          {/* Checklist for checkin */}
          <div>
            <h2 className="font-medium mb-2">Check-In Requirements:</h2>
            <div className="space-y-1 text-md text-text-secondary flex flex-col">
              <span className="flex flex-row items-center gap-1">
                <CheckListIcon
                  isMet={userInfo.data.event_role === "attendee"}
                />
                <p>User must be attendee</p>
              </span>
              <span className="flex flex-row items-center gap-1">
                <CheckListIcon isMet={!userInfo.data.checked_in_at} />
                <p>User has not checked in yet.</p>
              </span>
              <span className="flex flex-row items-center gap-1">
                <CheckListIcon isMet={rfid !== null && rfid.trim() !== ""} />
                {rfid ? (
                  <p>Scanned RFID: {rfid}</p>
                ) : (
                  <p>Have scanned RFID (Go ahead and tap RFID now)</p>
                )}
              </span>
            </div>
          </div>

          <div className="flex flex-col w-full gap-3">
            <Button
              variant="primary"
              onPress={onCheckIn}
              isDisabled={!isValidForCheckIn}
              className="w-full"
            >
              Check In
            </Button>

            {!isValidForCheckIn && (
              <p className="text-xs text-text-secondary text-center">
                Please ensure all check-in requirements are met.
              </p>
            )}
          </div>
          <input
            type="text"
            autoFocus
            onChange={(e) => setRfid(e.target.value)}
            inputMode="none"
            className="opacity-0 absolute h-0 w-0 overflow-hidden"
          />
        </div>
      )}
    </Modal>
  );
}
