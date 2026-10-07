import { HTTPError } from "ky";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/ky";

interface DashboardRow {
  userId: string;
  name: string;
  email: string;
  status: string;
  checkedInAt: string | null;
  rfid: string | null;
  signupSource: string;
  standbyArrival: string | null;
  redemptions: { redeemableId: string; name: string; amount: number }[];
}

interface Dashboard {
  hackathonId: string;
  confirmed: number;
  confirmedCheckedIn: number;
  dayOfSignups: number;
  standbyWaiting: number;
  rows: DashboardRow[];
}

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
  if (row.status === "waitlist_confirmed") return "Standby";
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

export default function StaffCheckInDashboard() {
  const [search, setSearch] = useState("");
  const [hackathonId, setHackathonId] = useState("");
  const [status, setStatus] = useState("");
  const [redeemableId, setRedeemableId] = useState("");
  const [tab, setTab] = useState<Tab>("All");
  const [showFilters, setShowFilters] = useState(false);
  const [busyUser, setBusyUser] = useState<string | null>(null);
  const [actionError, setActionError] = useState("");

  const dashboard = useQuery({
    queryKey: ["check-in-dashboard", hackathonId, search, status, redeemableId],
    queryFn: ({ signal }) =>
      api
        .get("application/check-in-dashboard", {
          searchParams: { hackathonId, search, status, redeemableId },
          signal,
        })
        .json<Dashboard>(),
    refetchInterval: 15000,
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

  return (
    <main className="w-full max-w-7xl p-4 sm:p-6 pb-24 space-y-5">
      <h1 className="text-2xl font-bold">Check-in Dashboard</h1>

      {data && (
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          {[
            [
              "Confirmed checked in",
              `${data.confirmedCheckedIn} / ${data.confirmed}`,
            ],
            ["Standby waiting", data.standbyWaiting],
            ["Day-of signups", data.dayOfSignups],
          ].map(([label, value]) => (
            <div key={label} className="rounded-xl border border-border p-4">
              <p className="text-sm text-text-secondary">{label}</p>
              <p className="mt-1 text-2xl font-semibold">{value}</p>
            </div>
          ))}
        </div>
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
                {["Hacker", "Status", "Arrival", "Action"].map((label) => (
                  <th key={label} className="p-4 whitespace-nowrap font-medium">
                    {label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr
                  key={row.userId}
                  className="border-b border-border last:border-0"
                >
                  <td className="p-4 align-top">
                    <p className="font-semibold">{row.name}</p>
                    <p className="mt-1 text-text-secondary break-all">
                      {row.email}
                    </p>
                    <span className="inline-block mt-2 rounded-full border border-border px-2 py-0.5 text-xs">
                      {row.signupSource === "day_of"
                        ? "Day-of"
                        : "Preregistered"}
                    </span>
                    <details className="mt-2 text-text-secondary">
                      <summary className="cursor-pointer text-xs">
                        Details
                      </summary>
                      <div className="mt-2 space-y-1 text-xs">
                        <p>Application: {row.status.replaceAll("_", " ")}</p>
                        <p>Check-in: {timestamp(row.checkedInAt)}</p>
                        <p>
                          Redeemables:{" "}
                          {row.redemptions
                            .map((r) => `${r.name} × ${r.amount}`)
                            .join(", ") || "None"}
                        </p>
                      </div>
                    </details>
                  </td>
                  <td className="p-4 align-top">
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
                    {timestamp(row.standbyArrival)}
                  </td>
                  <td className="p-4 align-top min-w-[190px]">
                    {!hackathonId && awaitingArrival(row) && (
                      <button
                        type="button"
                        className={actionClass}
                        disabled={busyUser !== null}
                        onClick={() => updateStandby(row, false)}
                      >
                        {busyUser === row.userId
                          ? "Recording…"
                          : "Record arrival"}
                      </button>
                    )}
                    {row.status === "waitlist_confirmed" && (
                      <button
                        type="button"
                        className="rounded-lg bg-blue-600 text-white px-4 py-2 font-medium whitespace-nowrap disabled:opacity-50"
                        disabled={busyUser !== null}
                        onClick={() => updateStandby(row, true)}
                      >
                        {busyUser === row.userId ? "Accepting…" : "Accept"}
                      </button>
                    )}
                    {row.status === "confirmed" && !row.checkedInAt && (
                      <span className="text-text-secondary whitespace-nowrap">
                        Check in with NFC
                      </span>
                    )}
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
    </main>
  );
}
