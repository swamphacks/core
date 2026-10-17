import { useState } from "react";
import QrWorkshopScanner from "./QrWorkshopScanner";

export default function WorkshopAttendancePage() {
  const [open, setOpen] = useState(false);
  return (
    <main className="max-w-4xl space-y-4 p-6">
      <h1 className="text-2xl font-bold">Workshop & social attendance</h1>
      <p>
        Select a session, then scan each hacker’s portal QR code. Hackers must
        already be confirmed and checked in. NFC and QR scans count toward the
        same attendance record.
      </p>
      <button
        type="button"
        className="rounded-lg bg-blue-600 px-4 py-2 font-medium text-white"
        onClick={() => setOpen(true)}
      >
        Scan QR
      </button>
      <QrWorkshopScanner isOpen={open} onClose={() => setOpen(false)} />
    </main>
  );
}
