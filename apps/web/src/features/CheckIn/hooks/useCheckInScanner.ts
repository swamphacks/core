import { Intent } from "@/lib/qr-intents/intent";
import { parseQrIntent } from "@/lib/qr-intents/parse";
import { type IDetectedBarcode } from "@yudiel/react-qr-scanner";
import { useState } from "react";

export const useCheckInScanner = () => {
  const [isScannerActive, setIsScannerActive] = useState(true);
  const [isLoadingUserInfo, setIsLoadingUserInfo] = useState(false);
  const [scannedUserId, setScannedUserId] = useState<string | null>(null);

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

    switch (res.value.intent) {
      case Intent.CHECK_IN:
        setScannedUserId(res.value.user_id);
        break;

      default:
        reset();
        console.warn("Unhandled intent type:", res.value.intent);
        break;
    }
  };

  const reset = () => {
    setScannedUserId(null);
    setIsLoadingUserInfo(false);
    setIsScannerActive(true);
  };

  return {
    isScannerActive,
    isLoadingUserInfo,
    scannedUserId,
    onScan,
    reset,
  };
};
