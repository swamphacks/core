import { HTTPError } from "ky";
import { useState } from "react";
import { api } from "@/lib/ky";
import { Modal } from "@/components/ui/Modal";
import QRCode from "react-qr-code";
import { DialogTrigger } from "react-aria-components";
import { Sheet } from "@/components/ui/Sheet";
import ApplicationViewer from "@/modules/Application/ApplicationViewer";
import useParsedForm from "@/modules/Application/hooks/useParsedForm";
import TablerQrcode from "~icons/tabler/qrcode";
import QrCheckInScanner from "./components/QrCheckInScanner";
import { undoCheckIn } from "./checkInActions";
import { useCheckInDashboard } from "./hooks/useCheckInDashboard";
import type { DashboardRow } from "./types";

const fieldClass =
  "rounded-lg border border-border bg-background px-3 py-2 text-sm";
const actionClass = `${fieldClass} whitespace-nowrap font-medium disabled:opacity-50 disabled:cursor-not-allowed`;

const tabs = ["All", "Awaiting arrival", "Standby", "Checked in"] as const;
type Tab = (typeof tabs)[number];

function awaitingArrival(row: DashboardRow) {
  return (
    ["waitlisted", "rejected"].includes(row.status) &&
    !row.standbyArrival &&
    !row.checkedInAt
  );
}

function statusLabel(row: DashboardRow) {
  if (row.checkedInAt) return "Checked in";
  if (row.status === "waitlist_confirmed") return "Standby · arrival recorded";
  if (awaitingArrival(row)) return "Awaiting arrival";
  return row.status.replaceAll("_", " ");
}

function timestamp(value: string | null) {
  return value
    ? new Date(value).toLocaleString(undefined, {
        month: "short",
        day: "numeric",
        hour: "numeric",
        minute: "2-digit",
      })
    : "—";
}

export default function StaffCheckInDashboard({
  isAdmin,
}: {
  isAdmin: boolean;
}) {
  const parsedForm = useParsedForm();
  const [search, setSearch] = useState("");
  const [hackathonId, setHackathonId] = useState("");
  const [status, setStatus] = useState("");
  const [redeemableId, setRedeemableId] = useState("");
  const [tab, setTab] = useState<Tab>("All");
  const [showFilters, setShowFilters] = useState(false);
  const [busyUser, setBusyUser] = useState<string | null>(null);
  const [actionError, setActionError] = useState("");
  const [selectedAttendee, setSelectedAttendee] = useState<DashboardRow | null>(
    null,
  );
  const [scannerOpen, setScannerOpen] = useState(false);

  const dashboard = useCheckInDashboard({
    hackathonId,
    search,
    status,
    redeemableId,
  });

  const data = dashboard.data;
  const rows = (data?.rows ?? []).filter((row) => {
    if (tab === "Awaiting arrival") return awaitingArrival(row);
    if (tab === "Standby") return row.status === "waitlist_confirmed";
    if (tab === "Checked in") return !!row.checkedInAt;
    return true;
  });
  const filterCount = [hackathonId, status, redeemableId].filter(
    Boolean,
  ).length;

  async function undo(row: DashboardRow) {
    if (busyUser) return;
    if (
      !window.confirm(
        `Undo ${row.name}'s check-in? They will need to be checked in again.`,
      )
    )
      return;

    setBusyUser(row.userId);
    setActionError("");
    const error = await undoCheckIn(row.userId);
    if (error) setActionError(error);
    else await dashboard.refetch();
    setBusyUser(null);
  }

  async function updateStandby(row: DashboardRow, accept: boolean) {
    if (!data || busyUser) return;
    if (
      accept &&
      !window.confirm(
        `Accept ${row.name}? Confirm there is space before continuing.`,
      )
    )
      return;

    setBusyUser(row.userId);
    setActionError("");
    try {
      await api.post(
        accept
          ? "application/waitlist/accept-standby"
          : "application/waitlist/in-person",
        {
          json: accept
            ? { hackathonId: data.hackathonId, userId: row.userId }
            : { userId: row.userId },
          retry: 0,
        },
      );
      await dashboard.refetch();
    } catch (error) {
      let message = "Unable to update standby status.";
      if (error instanceof HTTPError) {
        const response = await error.response
          .json<{ detail?: string }>()
          .catch(() => null);
        message = response?.detail ?? message;
      }
      setActionError(message);
    } finally {
      setBusyUser(null);
    }
  }

  async function unlinkTag(row: DashboardRow) {
    if (busyUser) return;
    if (
      !window.confirm(
        `Unlink the NFC tag from ${row.name}? They will need a tag linked again to check in with NFC.`,
      )
    )
      return;

    setBusyUser(row.userId);
    setActionError("");
    try {
      await api.delete(`nfc/checkin/nfc-links/${row.userId}`, { retry: 0 });
      await dashboard.refetch();
    } catch (error) {
      let message = "Unable to unlink NFC tag.";
      if (error instanceof HTTPError) {
        const response = await error.response
          .json<{ detail?: string }>()
          .catch(() => null);
        message = response?.detail ?? message;
      }
      setActionError(message);
    } finally {
      setBusyUser(null);
    }
  }

  return (
    <main className="w-full max-w-7xl p-4 sm:p-6 pb-24 space-y-5">
      <div className="flex items-center justify-between gap-3">
        <h1 className="text-2xl font-bold">Check-in Dashboard</h1>
        <button
          type="button"
          className="flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white whitespace-nowrap"
          onClick={() => setScannerOpen(true)}
        >
          <TablerQrcode className="h-5 w-5" aria-hidden />
          Scan QR
        </button>
      </div>

      {data && (
        <section className="rounded-xl border border-border">
          <div className="grid grid-cols-1 sm:grid-cols-3">
            {[
              [
                "Checked in · confirmed",
                `${data.confirmedCheckedIn} of ${data.confirmed}`,
              ],
              ["Standby waiting", data.standbyWaiting],
              [
                "Unreserved spots",
                data.capacity > 0 ? data.unreservedSpots : "—",
              ],
            ].map(([label, value]) => (
              <div
                key={label}
                className="px-5 py-4 border-b sm:border-b-0 sm:border-r border-border last:border-0"
              >
                <p className="text-xs font-medium text-text-secondary">
                  {label}
                </p>
                <p className="mt-1 text-2xl font-semibold tabular-nums">
                  {value}
                </p>
              </div>
            ))}
          </div>
          <details className="border-t border-border px-5 py-3 text-sm">
            <summary className="cursor-pointer text-text-secondary">
              Capacity and reservation details
            </summary>
            <dl className="mt-3 grid grid-cols-2 sm:grid-cols-4 gap-4">
              {[
                [
                  "Capacity target",
                  data.capacity > 0 ? data.capacity : "Not configured",
                ],
                [
                  "Confirmed, not checked in",
                  Math.max(0, data.confirmed - data.confirmedCheckedIn),
                ],
                ["Outstanding offers", data.outstandingOffers],
                ["Day-of signups", data.dayOfSignups],
              ].map(([label, value]) => (
                <div key={label}>
                  <dt className="text-xs text-text-secondary">{label}</dt>
                  <dd className="mt-1 font-medium tabular-nums">{value}</dd>
                </div>
              ))}
            </dl>
          </details>
          {data.capacity > 0 &&
            data.confirmed + data.outstandingOffers > data.capacity && (
              <p className="border-t border-border px-5 py-3 text-sm text-amber-800 dark:text-amber-200">
                Reservations exceed the target by{" "}
                {data.confirmed + data.outstandingOffers - data.capacity}.
                Existing offers remain valid.
              </p>
            )}
        </section>
      )}

      <div className="flex gap-3">
        <input
          type="search"
          aria-label="Search hackers by name or email"
          placeholder="Search name or email"
          className={`${fieldClass} flex-1 min-w-0`}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <button
          type="button"
          className={actionClass}
          aria-expanded={showFilters}
          aria-controls="dashboard-filters"
          onClick={() => setShowFilters(!showFilters)}
        >
          Filters{filterCount > 0 ? ` (${filterCount})` : ""}
        </button>
      </div>

      {showFilters && (
        <div
          id="dashboard-filters"
          className="rounded-xl border border-border p-4 flex flex-wrap gap-3 items-end"
        >
          <label className="flex flex-col gap-1 text-sm">
            Event ID
            <input
              className={fieldClass}
              value={hackathonId}
              placeholder="Active event"
              onChange={(e) => setHackathonId(e.target.value)}
            />
          </label>
          <label className="flex flex-col gap-1 text-sm">
            Application status
            <select
              className={fieldClass}
              value={status}
              onChange={(e) => setStatus(e.target.value)}
            >
              <option value="">All statuses</option>
              {[
                "accepted",
                "confirmed",
                "waitlisted",
                "waitlist_confirmed",
                "rejected",
                "withdrawn",
                "started",
                "submitted",
                "under_review",
              ].map((value) => (
                <option key={value} value={value}>
                  {value.replaceAll("_", " ")}
                </option>
              ))}
            </select>
          </label>
          <label className="flex flex-col gap-1 text-sm">
            Redeemable ID
            <input
              className={fieldClass}
              value={redeemableId}
              placeholder="All redeemables"
              onChange={(e) => setRedeemableId(e.target.value)}
            />
          </label>
          <button
            type="button"
            className={actionClass}
            onClick={() => {
              setHackathonId("");
              setStatus("");
              setRedeemableId("");
            }}
          >
            Clear filters
          </button>
        </div>
      )}

      <div className="flex flex-wrap gap-2" aria-label="Dashboard views">
        {tabs.map((value) => (
          <button
            key={value}
            type="button"
            aria-pressed={tab === value}
            onClick={() => setTab(value)}
            className={`rounded-lg px-4 py-2 text-sm font-medium ${
              tab === value
                ? "bg-blue-600 text-white"
                : "border border-border text-text-secondary"
            }`}
          >
            {value}
          </button>
        ))}
      </div>

      {tab === "Standby" && (
        <p className="text-sm text-text-secondary">
          Preregistered hackers take priority. Accept people in queue order when
          space is available.
        </p>
      )}

      {actionError && (
        <p role="alert" className="text-red-600">
          {actionError}
        </p>
      )}
      {dashboard.isPending && <p>Loading dashboard…</p>}
      {dashboard.isError && (
        <p role="alert">Unable to load dashboard: {dashboard.error.message}</p>
      )}

      {data && (
        <div className="overflow-x-auto rounded-xl border border-border">
          <table className="w-full min-w-[700px] text-sm text-left">
            <thead>
              <tr className="border-b border-border text-text-secondary">
                {["Hacker", "Status", "Standby arrival", "Next step"].map(
                  (label) => (
                    <th
                      key={label}
                      className="p-4 whitespace-nowrap font-medium"
                    >
                      {label}
                    </th>
                  ),
                )}
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr
                  key={row.userId}
                  className="border-b border-border last:border-0"
                >
                  <td className="px-4 py-3 align-top">
                    {tab === "Standby" && row.standbyPosition != null && (
                      <p className="mb-1.5 text-xs font-semibold text-text-secondary">
                        #{row.standbyPosition}
                        {row.standbyPosition === 1 ? " · Next in line" : ""}
                      </p>
                    )}
                    {isAdmin ? (
                      <DialogTrigger>
                        <button
                          type="button"
                          className="font-semibold text-left underline decoration-current/30 underline-offset-4 hover:decoration-current transition-colors focus-visible:outline-2 focus-visible:outline-offset-4"
                          aria-label={`View application for ${row.name}`}
                        >
                          {row.name}
                        </button>
                        <Sheet sheetClassName="w-full sm:w-160 lg:w-200">
                          {parsedForm && (
                            <ApplicationViewer
                              parsedForm={parsedForm}
                              applicationId={row.applicationId}
                            />
                          )}
                        </Sheet>
                      </DialogTrigger>
                    ) : (
                      <p className="font-semibold">{row.name}</p>
                    )}
                    <p className="mt-1 text-text-secondary break-all">
                      {row.email}
                    </p>
                    <div className="mt-2 flex flex-wrap items-center gap-2">
                      <span className="inline-flex h-7 items-center rounded-md bg-slate-100 px-2.5 text-xs font-medium text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                        {row.signupSource === "day_of"
                          ? "Day-of"
                          : "Preregistered"}
                      </span>
                      <button
                        type="button"
                        className="inline-flex h-7 items-center rounded-md border border-border px-2.5 text-xs font-medium text-text-secondary hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                        onClick={() => setSelectedAttendee(row)}
                      >
                        QR code
                      </button>
                    </div>
                    <details className="mt-1.5 text-text-secondary">
                      <summary className="cursor-pointer text-xs">
                        Details
                      </summary>
                      <div className="mt-2 space-y-1 text-xs">
                        <p>Application: {row.status.replaceAll("_", " ")}</p>
                        <p>
                          NFC check-in:{" "}
                          {row.checkedInAt
                            ? timestamp(row.checkedInAt)
                            : "Not checked in"}
                        </p>
                        <p>
                          Redeemables:{" "}
                          {row.redemptions
                            .map((r) => `${r.name} × ${r.amount}`)
                            .join(", ") || "None"}
                        </p>
                      </div>
                    </details>
                  </td>
                  <td className="px-4 py-3 align-top">
                    <span
                      className={`inline-block rounded-full px-3 py-1 text-xs font-medium capitalize ${
                        row.checkedInAt
                          ? "bg-green-100 text-green-800 dark:bg-green-950 dark:text-green-200"
                          : row.status === "waitlist_confirmed"
                            ? "bg-amber-100 text-amber-800 dark:bg-amber-950 dark:text-amber-200"
                            : "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200"
                      }`}
                    >
                      {statusLabel(row)}
                    </span>
                  </td>
                  <td className="p-4 align-top whitespace-nowrap">
                    {row.standbyArrival
                      ? timestamp(row.standbyArrival)
                      : [
                            "waitlisted",
                            "waitlist_confirmed",
                            "rejected",
                          ].includes(row.status)
                        ? "Not recorded"
                        : "Not applicable"}
                  </td>
                  <td className="px-4 py-3 align-top min-w-[190px]">
                    <div className="flex flex-col items-start gap-2">
                      {!hackathonId && awaitingArrival(row) && (
                        <button
                          type="button"
                          className={actionClass}
                          disabled={busyUser !== null}
                          onClick={() => updateStandby(row, false)}
                        >
                          {busyUser === row.userId
                            ? "Recording…"
                            : "Record standby arrival"}
                        </button>
                      )}
                      {row.status === "waitlist_confirmed" && (
                        <button
                          type="button"
                          className="rounded-lg bg-blue-600 text-white px-4 py-2 font-medium whitespace-nowrap disabled:opacity-50"
                          disabled={busyUser !== null}
                          onClick={() => updateStandby(row, true)}
                        >
                          {busyUser === row.userId
                            ? "Accepting…"
                            : "Offer spot"}
                        </button>
                      )}
                      {row.checkedInAt && (
                        <p className="text-text-secondary">Check-in complete</p>
                      )}
                      {!row.checkedInAt && row.status === "accepted" && (
                        <p className="text-text-secondary">
                          Waiting for portal confirmation
                        </p>
                      )}
                      {!row.checkedInAt &&
                        awaitingArrival(row) &&
                        hackathonId && (
                          <p className="text-text-secondary">
                            Select the active event to record standby arrival
                          </p>
                        )}
                      {!row.checkedInAt &&
                        !awaitingArrival(row) &&
                        ![
                          "accepted",
                          "confirmed",
                          "waitlist_confirmed",
                        ].includes(row.status) && (
                          <p className="text-text-secondary">
                            No check-in action
                          </p>
                        )}
                      {row.status === "confirmed" && !row.checkedInAt && (
                        <span className="text-text-secondary whitespace-nowrap">
                          Check in with Scan QR or NFC
                        </span>
                      )}
                      {row.checkedInAt && !row.nfcLinked && (
                        <button
                          type="button"
                          className={actionClass}
                          disabled={busyUser !== null}
                          onClick={() => undo(row)}
                        >
                          {busyUser === row.userId
                            ? "Undoing…"
                            : "Undo check-in"}
                        </button>
                      )}
                      {row.nfcLinked && (
                        <button
                          type="button"
                          className={actionClass}
                          disabled={busyUser !== null}
                          onClick={() => unlinkTag(row)}
                        >
                          {busyUser === row.userId
                            ? "Unlinking…"
                            : "Unlink NFC tag"}
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr>
                  <td
                    colSpan={4}
                    className="p-8 text-center text-text-secondary"
                  >
                    No matching hackers.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      <Modal
        title="Attendee QR code"
        dialogAriaLabel="Attendee QR code"
        isOpen={selectedAttendee !== null}
        onOpenChange={(isOpen) => {
          if (!isOpen) setSelectedAttendee(null);
        }}
        size="sm"
      >
        {selectedAttendee && (
          <div className="flex flex-col items-center gap-4">
            <div className="text-center">
              <p className="font-semibold">{selectedAttendee.name}</p>
              <p className="text-sm text-text-secondary break-all">
                {selectedAttendee.email}
              </p>
              <p className="mt-2 text-sm capitalize">
                Application status:{" "}
                {selectedAttendee.status.replaceAll("_", " ")}
              </p>
            </div>
            <div className="w-full max-w-64 rounded-2xl bg-white p-4 ring-1 ring-slate-200">
              <QRCode
                size={224}
                style={{ height: "auto", maxWidth: "100%", width: "100%" }}
                value={`IDENT::${selectedAttendee.userId}`}
                viewBox="0 0 256 256"
                bgColor="#FFFFFF"
                fgColor="#000000"
                level="H"
              />
            </div>
            <p className="text-center font-mono text-xs text-text-secondary break-all">
              ID: {selectedAttendee.userId}
            </p>
            <button
              type="button"
              className={`${actionClass} self-stretch`}
              onClick={() => setSelectedAttendee(null)}
            >
              Close
            </button>
          </div>
        )}
      </Modal>

      <QrCheckInScanner
        isOpen={scannerOpen}
        onClose={() => setScannerOpen(false)}
      />
    </main>
  );
}
