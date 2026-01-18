import { Avatar } from "@/components/ui/Avatar";
import { Modal } from "@/components/ui/Modal";
import { CheckInBadge, CheckInEventRoleBadge } from "./CheckInBadge";
import { Button } from "@/components/ui/Button";
import type { UserEventInfo } from "../hooks/useCheckInScanner";
import { useState } from "react";
import { toast } from "react-toastify";

interface Props {
  isModalOpen: boolean;
  setIsModalOpen: (isOpen: boolean) => void;
  resetUserInfo: () => void;
  userEventInfo: UserEventInfo | null;
}

export default function CheckInModal({
  isModalOpen,
  setIsModalOpen,
  userEventInfo,
  resetUserInfo,
}: Props) {
  const [rfid, setRfid] = useState<string | null>(null);

  const onCheckIn = async (userId: string) => {
    if (!rfid || rfid.trim().length !== 10) {
      toast.error("Didn't scan valid RFID.");
      return;
    }

    console.log(`Checking in ${userId}`);

    //
    setIsModalOpen(false);
    setRfid(null);
    resetUserInfo();
  };

  if (userEventInfo === null) {
    return null;
  }

  const invalidCheckin =
    userEventInfo.checked_in_at !== null ||
    userEventInfo.event_role !== "attendee" ||
    !rfid ||
    rfid.trim().length != 10;

  return (
    <Modal
      responsive="sheet"
      dialogAriaLabel="checkin-fields"
      isOpen={isModalOpen}
      isDismissible
      onOpenChange={setIsModalOpen}
    >
      <div className="flex flex-col gap-6">
        <div className="flex items-center gap-4">
          <Avatar
            src={
              userEventInfo.image ||
              "https://static.vecteezy.com/system/resources/thumbnails/009/292/244/small/default-avatar-icon-of-social-media-user-vector.jpg"
            }
            size="lg"
          />

          <div className="flex flex-col flex-1 min-w-0">
            <h1 className="text-md font-medium truncate">
              {userEventInfo.name}
            </h1>
            <p className="text-sm text-text-secondary truncate">
              {userEventInfo.email}
            </p>

            <div className="mt-2 flex flex-wrap gap-2">
              <CheckInEventRoleBadge role={userEventInfo.event_role} />

              <CheckInBadge
                isCheckedIn={userEventInfo.checked_in_at ? true : false}
              />
            </div>
          </div>
        </div>

        <p className="sm:hidden text-sm text-yellow-700 bg-yellow-50 border border-yellow-200 rounded-md px-3 py-2">
          On mobile, check-in requires an attached RFID scanner.
        </p>

        <input
          type="text"
          autoFocus
          onChange={(e) => setRfid(e.target.value)}
          inputMode="none"
          className="opacity-0 h-0 w-0 overflow-hidden"
        />

        {rfid && <p>{rfid}</p>}

        <Button
          variant="primary"
          onPress={() => onCheckIn(userEventInfo.user_id)}
          isDisabled={invalidCheckin}
          className="w-full"
        >
          Check In
        </Button>

        {invalidCheckin && (
          <p className="text-xs text-text-secondary text-center">
            This user cannot be checked in.
          </p>
        )}
      </div>
    </Modal>
  );
}
