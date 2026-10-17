import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Scanner, type IDetectedBarcode } from "@yudiel/react-qr-scanner";
import { Modal } from "@/components/ui/Modal";
import { parseQrIntent } from "@/lib/qr-intents/parse";
import TablerNfc from "~icons/tabler/nfc";
import { checkInWithoutNfc, undoCheckIn } from "../checkInActions";
import { useCheckInDashboard } from "../hooks/useCheckInDashboard";
import type { DashboardRow } from "../types";

type ScanResult =
  | { kind: "found"; userId: string }
  | { kind: "invalid"; message: string };

interface Verdict {
  label: string;
  tone: "ready" | "done" | "blocked";
  detail: string;
}

const actionClass =
  "rounded-lg border border-border bg-background px-3 py-2 text-sm font-medium whitespace-nowrap";

const toneClass: Record<Verdict["tone"], string> = {
  ready: "bg-green-100 text-green-800 dark:bg-green-950 dark:text-green-200",
  done: "bg-blue-100 text-blue-800 dark:bg-blue-950 dark:text-blue-200",
  blocked: "bg-red-100 text-red-800 dark:bg-red-950 dark:text-red-200",
};

// Reads the same checkedInAt the dashboard shows (NFC nametag linked or
// checked in here), so the scanner never disagrees with the table behind it.
function verdictFor(row: DashboardRow): Verdict {
  if (row.checkedInAt) {
    return {
      label: "Already checked in",
      tone: "done",
      detail: `Checked in ${new Date(row.checkedInAt).toLocaleString(
        undefined,
        {
          month: "short",
          day: "numeric",
          hour: "numeric",
          minute: "2-digit",
        },
      )}${row.nfcLinked ? " with their NFC nametag" : " without an NFC nametag"}.`,
    };
  }
  if (row.status === "confirmed") {
    return {
      label: "Ready to check in",
      tone: "ready",
      detail: "Confirmed attendee. Check them in below.",
    };
  }
  if (row.status === "waitlist_confirmed") {
    return {
      label: "On standby",
      tone: "blocked",
      detail: "Accept this hacker from the Standby tab before checking in.",
    };
  }
  return {
    label: "Not eligible",
    tone: "blocked",
    detail: `Application status is "${row.status.replaceAll("_", " ")}". Only confirmed hackers can check in.`,
  };
}

export default function QrCheckInScanner({
  isOpen,
  onClose,
}: {
  isOpen: boolean;
  onClose: () => void;
}) {
  const [result, setResult] = useState<ScanResult | null>(null);
  const [cameraError, setCameraError] = useState("");
  const [busy, setBusy] = useState(false);
  const [checkInError, setCheckInError] = useState("");
  const [justCheckedIn, setJustCheckedIn] = useState(false);
  const queryClient = useQueryClient();
  // Unfiltered, so a scan finds the hacker no matter what the table is
  // filtered to. Shares the cache entry with the dashboard's default view.
  const dashboard = useCheckInDashboard({}, { enabled: isOpen });

  function reset() {
    setResult(null);
    setCheckInError("");
    setJustCheckedIn(false);
  }

  function close() {
    reset();
    setCameraError("");
    onClose();
  }

  async function run(action: "check-in" | "undo", userId: string) {
    if (
      action === "undo" &&
      !window.confirm(
        "Undo this check-in? The hacker will need to be checked in again.",
      )
    )
      return;

    setBusy(true);
    setCheckInError("");
    const error =
      action === "check-in"
        ? await checkInWithoutNfc(userId)
        : await undoCheckIn(userId);
    if (error) {
      setCheckInError(error);
    } else {
      setJustCheckedIn(action === "check-in");
      await queryClient.invalidateQueries({
        queryKey: ["check-in-dashboard"],
      });
    }
    setBusy(false);
  }

  function onScan(codes: IDetectedBarcode[]) {
    const raw = codes[0]?.rawValue;
    if (!raw || result) return;

    const parsed = parseQrIntent(raw);
    if (!parsed.ok) {
      setResult({
        kind: "invalid",
        message: "This is not a SwampHacks check-in QR code.",
      });
      return;
    }
    setResult({ kind: "found", userId: parsed.value.user_id });
  }

  const row =
    result?.kind === "found"
      ? dashboard.data?.rows.find((r) => r.userId === result.userId)
      : undefined;
  const verdict = row ? verdictFor(row) : null;

  return (
    <Modal
      title="Scan to check in"
      dialogAriaLabel="Scan QR code to check in"
      responsive="adaptive"
      isOpen={isOpen}
      onOpenChange={(open) => {
        if (!open) close();
      }}
    >
      <div className="mt-4 flex flex-col gap-4">
        {!result && (
          <>
            <div className="mx-auto w-full max-w-sm overflow-hidden rounded-xl bg-black aspect-square">
              <Scanner
                onScan={onScan}
                onError={() =>
                  setCameraError(
                    "Camera unavailable. Allow camera access for this site and make sure no other app is using it.",
                  )
                }
                constraints={{ facingMode: "environment" }}
                formats={["qr_code"]}
                sound={false}
              />
            </div>
            <p className="text-center text-sm text-text-secondary">
              Point the camera at the hacker's check-in QR code.
            </p>
            {cameraError && (
              <p role="alert" className="text-sm text-red-600">
                {cameraError}
              </p>
            )}
          </>
        )}

        {result?.kind === "invalid" && (
          <p
            role="alert"
            className="rounded-lg bg-red-100 px-3 py-2 text-sm text-red-800 dark:bg-red-950 dark:text-red-200"
          >
            {result.message}
          </p>
        )}

        {result?.kind === "found" && dashboard.isPending && (
          <p>Looking up hacker…</p>
        )}

        {result?.kind === "found" && dashboard.isError && (
          <p role="alert" className="text-sm text-red-600">
            Unable to look up hacker: {dashboard.error.message}
          </p>
        )}

        {result?.kind === "found" && dashboard.data && !row && (
          <p
            role="alert"
            className="rounded-lg bg-red-100 px-3 py-2 text-sm text-red-800 dark:bg-red-950 dark:text-red-200"
          >
            No application for this hackathon matches this QR code.
          </p>
        )}

        {row && verdict && (
          <div className="flex flex-col gap-4">
            <div>
              <p className="text-lg font-semibold">{row.name}</p>
              <p className="text-sm text-text-secondary break-all">
                {row.email}
              </p>
              <span
                className={`mt-2 inline-block rounded-full px-3 py-1 text-xs font-medium ${toneClass[verdict.tone]}`}
              >
                {verdict.label}
              </span>
              <p className="mt-2 text-sm text-text-secondary">
                {verdict.detail}
              </p>
            </div>

            {verdict.tone === "ready" && (
              <div className="flex flex-col gap-2">
                <button
                  type="button"
                  className="w-full rounded-lg bg-green-600 px-4 py-3 font-medium text-white disabled:opacity-50"
                  disabled={busy}
                  onClick={() => run("check-in", row.userId)}
                >
                  {busy ? "Checking in…" : `Check in ${row.name}`}
                </button>
              </div>
            )}

            {verdict.tone === "done" && !row.nfcLinked && (
              <button
                type="button"
                className={`${actionClass} self-start disabled:opacity-50`}
                disabled={busy}
                onClick={() => run("undo", row.userId)}
              >
                {busy ? "Undoing…" : "Undo check-in"}
              </button>
            )}

            {checkInError && (
              <p role="alert" className="text-sm text-red-600">
                {checkInError}
              </p>
            )}

            {(verdict.tone === "ready" || justCheckedIn) && (
              <div className="rounded-xl border-2 border-dashed border-border p-4">
                <div className="flex flex-wrap items-center gap-2 font-medium">
                  <TablerNfc className="h-5 w-5" aria-hidden />
                  {justCheckedIn
                    ? "Next: link their NFC nametag on the mobile app"
                    : "NFC nametag: link it afterwards on the mobile app"}
                  <span className="rounded-full border border-border px-2 py-0.5 text-xs text-text-secondary">
                    NFC not available on web
                  </span>
                </div>
                <p className="mt-2 text-sm text-text-secondary">
                  Checking in here skips the nametag. To use their nametag for
                  meals and swag, link it after check-in on a phone with the
                  SwampHacks mobile app:
                </p>
                <ol className="mt-2 list-decimal space-y-1 pl-5 text-sm text-text-secondary">
                  <li>Give the hacker their NFC nametag.</li>
                  <li>Open the SwampHacks mobile app on a phone.</li>
                  <li>Tap the nametag to link it to {row.name}.</li>
                </ol>
              </div>
            )}
          </div>
        )}

        <div className="flex gap-3">
          {result && (
            <button
              type="button"
              className="flex-1 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white"
              onClick={reset}
            >
              Scan next
            </button>
          )}
          <button
            type="button"
            className={`${actionClass} flex-1`}
            onClick={close}
          >
            Close
          </button>
        </div>
      </div>
    </Modal>
  );
}
