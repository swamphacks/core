import { useEffect, useRef, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Scanner, type IDetectedBarcode } from "@yudiel/react-qr-scanner";
import { Modal } from "@/components/ui/Modal";
import { parseQrIntent } from "@/lib/qr-intents/parse";
import TablerCircleCheckFilled from "~icons/tabler/circle-check-filled";
import {
  redeemablesQueryKey,
  redeemForUser,
  useRedeemables,
  type Redemption,
} from "./hooks/useRedeemables";

type ScanOutcome =
  | { kind: "redeemed"; redemption: Redemption }
  | { kind: "error"; message: string };
type ScanResult = ScanOutcome & { id: number };

// The camera reads a code about twice a second while it is in view. A code is
// only acted on again once it has been out of view this long, so a hacker
// holding up their phone gets one item, and can step back and show it again
// for a second one.
const SAME_CODE_COOLDOWN_MS = 2000;
const HISTORY_LENGTH = 5;
// How long the green flash holds before it fades out.
const FLASH_HOLD_MS = 250;

const fieldClass =
  "rounded-lg border border-border bg-background px-3 py-2 text-sm";
const actionClass = `${fieldClass} whitespace-nowrap font-medium`;

// Item first, like the NFC mobile app: staff pick what they are handing out
// once, then scan hackers one after another without pressing anything. Each
// scan redeems one item.
export default function QrRedeemScanner({
  isOpen,
  onClose,
}: {
  isOpen: boolean;
  onClose: () => void;
}) {
  const [redeemableId, setRedeemableId] = useState("");
  const [history, setHistory] = useState<ScanResult[]>([]);
  const [busy, setBusy] = useState(false);
  const [cameraError, setCameraError] = useState("");
  const [flashing, setFlashing] = useState(false);
  const flashTimer = useRef<ReturnType<typeof setTimeout>>(undefined);
  // Refs, not state: the camera can report again before React re-renders, and
  // these must already reflect the previous report when it does.
  const inFlight = useRef(false);
  const lastSeen = useRef(new Map<string, number>());
  const nextId = useRef(0);
  const queryClient = useQueryClient();
  const redeemables = useRedeemables({ enabled: isOpen });
  const selected = redeemables.data?.find((r) => r.id === redeemableId);

  useEffect(() => () => clearTimeout(flashTimer.current), []);

  // Covers the camera in green so staff can see from arm's length that the
  // scan went through and they shouldn't scan again. Phones that support it
  // also buzz.
  function flashSuccess() {
    clearTimeout(flashTimer.current);
    setFlashing(true);
    flashTimer.current = setTimeout(() => setFlashing(false), FLASH_HOLD_MS);
    navigator.vibrate?.(120);
  }

  function reset() {
    setHistory([]);
    lastSeen.current.clear();
  }

  function close() {
    reset();
    setCameraError("");
    onClose();
  }

  function record(outcome: ScanOutcome) {
    const entry = { ...outcome, id: nextId.current++ };
    setHistory((h) => [entry, ...h].slice(0, HISTORY_LENGTH));
  }

  async function onScan(codes: IDetectedBarcode[]) {
    const raw = codes[0]?.rawValue;
    if (!raw || !redeemableId || inFlight.current) return;

    const now = Date.now();
    const previous = lastSeen.current.get(raw);
    lastSeen.current.set(raw, now);
    if (previous !== undefined && now - previous < SAME_CODE_COOLDOWN_MS) {
      return;
    }

    const parsed = parseQrIntent(raw);
    if (!parsed.ok) {
      record({
        kind: "error",
        message: "This is not a SwampHacks hacker QR code.",
      });
      return;
    }

    inFlight.current = true;
    setBusy(true);
    const response = await redeemForUser(redeemableId, parsed.value.user_id);
    record(
      response.ok
        ? { kind: "redeemed", redemption: response.value }
        : { kind: "error", message: response.message },
    );
    // The code stayed in view during the request; count from now so it is not
    // read as shown again.
    lastSeen.current.set(raw, Date.now());
    inFlight.current = false;
    setBusy(false);
    if (response.ok) {
      flashSuccess();
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["check-in-dashboard"] }),
        queryClient.invalidateQueries({ queryKey: redeemablesQueryKey }),
      ]);
    }
  }

  return (
    <Modal
      title="Scan to redeem"
      dialogAriaLabel="Scan QR code to redeem an item"
      responsive="adaptive"
      isOpen={isOpen}
      onOpenChange={(open) => {
        if (!open) close();
      }}
    >
      <div className="mt-4 flex flex-col gap-4">
        <label className="flex flex-col gap-1 text-sm">
          Item to hand out
          <select
            className={fieldClass}
            value={redeemableId}
            disabled={busy}
            onChange={(e) => {
              setRedeemableId(e.target.value);
              reset();
            }}
          >
            <option value="">Pick an item</option>
            {redeemables.data?.map((r) => (
              <option key={r.id} value={r.id}>
                {r.name} · {r.max_user_amount} per hacker
              </option>
            ))}
          </select>
        </label>

        {redeemables.isError && (
          <p role="alert" className="text-sm text-red-600">
            Unable to load items: {redeemables.error.message}
          </p>
        )}
        {redeemables.data?.length === 0 && (
          <p className="text-sm text-text-secondary">
            No items yet. Create one with New Item first.
          </p>
        )}

        {selected && (
          <>
            <div className="relative mx-auto w-full max-w-sm overflow-hidden rounded-xl bg-black aspect-square">
              <Scanner
                onScan={onScan}
                onError={() =>
                  setCameraError(
                    "Camera unavailable. Allow camera access for this site and make sure no other app is using it.",
                  )
                }
                paused={!isOpen}
                // Report every read, not only new codes; onScan decides what
                // counts as a fresh scan.
                allowMultiple
                constraints={{ facingMode: "environment" }}
                formats={["qr_code"]}
                sound={false}
              />
              {/* Appears instantly, then fades out once flashing ends. */}
              <div
                aria-hidden
                className={`pointer-events-none absolute inset-0 z-10 flex items-center justify-center bg-green-500/80 ${
                  flashing
                    ? "opacity-100"
                    : "opacity-0 transition-opacity duration-500"
                }`}
              >
                <TablerCircleCheckFilled className="h-24 w-24 text-white" />
              </div>
            </div>
            <p className="text-center text-sm text-text-secondary">
              {busy
                ? "Redeeming…"
                : `Scanning for ${selected.name}. Hold up each hacker's QR code.`}
            </p>
            {cameraError && (
              <p role="alert" className="text-sm text-red-600">
                {cameraError}
              </p>
            )}
          </>
        )}

        {history.length > 0 && (
          <ul className="flex flex-col gap-2" aria-live="polite">
            {history.map((entry, i) => (
              <li
                key={entry.id}
                className={`rounded-xl px-4 ${i === 0 ? "py-3" : "py-2 text-sm opacity-60"} ${
                  entry.kind === "redeemed"
                    ? "bg-green-100 text-green-800 dark:bg-green-950 dark:text-green-200"
                    : "bg-red-100 text-red-800 dark:bg-red-950 dark:text-red-200"
                }`}
              >
                {entry.kind === "redeemed" ? (
                  <>
                    <p
                      className={
                        i === 0 ? "text-lg font-semibold" : "font-medium"
                      }
                    >
                      {entry.redemption.redeemableName} for{" "}
                      {entry.redemption.name}
                    </p>
                    {i === 0 && (
                      <p className="mt-1 text-sm">
                        {entry.redemption.itemsRemaining === 0
                          ? "That was their last one."
                          : `${entry.redemption.itemsRemaining} of ${entry.redemption.maxUserAmount} left for this hacker.`}
                      </p>
                    )}
                  </>
                ) : (
                  <p className="font-medium">{entry.message}</p>
                )}
              </li>
            ))}
          </ul>
        )}

        <button type="button" className={actionClass} onClick={close}>
          Close
        </button>
      </div>
    </Modal>
  );
}
