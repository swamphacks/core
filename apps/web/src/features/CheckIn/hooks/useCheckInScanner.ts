import { api } from "@/lib/ky";
import { Intent } from "@/lib/qr-intents/intent";
import { parseQrIntent } from "@/lib/qr-intents/parse";
import { type IDetectedBarcode } from "@yudiel/react-qr-scanner";
import { useState } from "react";
import { toast } from "react-toastify";
import { z } from "zod";

const userEventInfoSchema = z.object({
  user_id: z.uuid(),
  name: z.string(),
  email: z.email(),
  image: z.url().nullable().optional(),
  platform_role: z.enum(["user", "superuser"]),
  event_role: z
    .enum(["attendee", "admin", "staff", "applicant"])
    .nullable()
    .optional(),
  checked_in_at: z.date().nullable(),
});

export type UserEventInfo = z.infer<typeof userEventInfoSchema>;

export const useCheckInScanner = (eventId: string) => {
  const [isScannerActive, setIsScannerActive] = useState(true);
  const [isLoadingUserInfo, setIsLoadingUserInfo] = useState(false);
  const [selectedUser, setSelectedUser] = useState<UserEventInfo | null>(null);

  const fetchUserEventInfo = async (userId: string) => {
    try {
      const result = await api.get(`events/${eventId}/users/${userId}`).json();

      return userEventInfoSchema.parse(result);
    } catch {
      toast.error("Could not load user profile for this event.");
    }
  };

  const onScan = async (scannedData: IDetectedBarcode[]) => {
    if (scannedData.length === 0) return;

    setIsScannerActive(false);
    setIsLoadingUserInfo(true);

    const rawVal = scannedData[0].rawValue;
    const res = parseQrIntent(rawVal);

    if (!res.ok) {
      console.error("Error parsing QR intent:", res.error);
      setIsLoadingUserInfo(false);
      setIsScannerActive(true);
      return;
    }

    switch (res.value.intent) {
      case Intent.CHECK_IN: {
        const user = await fetchUserEventInfo(res.value.user_id);

        if (!user) {
          reset();
          break;
        }

        setSelectedUser(user);
        break;
      }

      default: {
        reset();
        console.warn("Unhandled intent type:", res.value.intent);
        break;
      }
    }
  };

  const reset = () => {
    setIsLoadingUserInfo(false);
    setIsScannerActive(true);
    setSelectedUser(null);
  };

  return {
    isScannerActive,
    isLoadingUserInfo,
    selectedUser,
    onScan,
    reset,
  };
};
