import { Intent } from "@/lib/qr-intents/intent";
import {
  parseQrIntent,
  type CheckInIntent,
  type QRIntent,
} from "@/lib/qr-intents/parse";
import { type IDetectedBarcode } from "@yudiel/react-qr-scanner";
import { useState } from "react";

export const useCheckInScanner = (eventId: string) => {
  const [isScannerActive, setIsScannerActive] = useState(true);
  const [isLoadingUserInfo, setIsLoadingUserInfo] = useState(false);

  /** Handles a successful check-in */
  const handleCheckInIntent = async (data: CheckInIntent) => {
    const { user_id, event_id } = data;

    if (event_id !== eventId) {
      console.warn("Scanned QR code for a different event.");
      return;
    }

    console.log(`Checking in user ${user_id} to event ${event_id}`);

    // Simulate async check-in
    await new Promise((resolve) => setTimeout(resolve, 1000));

    console.log(`User ${user_id} checked in successfully.`);
  };

  /** Handles a parsed QR intent based on its type */
  const handleIntent = async (parsedIntent: QRIntent) => {
    switch (parsedIntent.intent) {
      case Intent.CHECK_IN:
        await handleCheckInIntent(parsedIntent);
        break;

      default:
        console.warn("Unhandled intent type:", parsedIntent.intent);
        break;
    }
  };

  /** Main scan handler */
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

    try {
      await handleIntent(res.value);
    } catch (err) {
      console.error("Error handling intent:", err);
    } finally {
      setIsLoadingUserInfo(false);
      setIsScannerActive(true);
    }
  };

  return {
    isScannerActive,
    isLoadingUserInfo,
    onScan,
  };
};
