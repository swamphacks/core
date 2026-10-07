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

const fieldClass = "rounded-md border border-border bg-background px-3 py-2";

export default function StaffCheckInDashboard() {
  const [search, setSearch] = useState("");
  const [hackathonId, setHackathonId] = useState("");
  const [status, setStatus] = useState("");
  const [checkedIn, setCheckedIn] = useState("");
  const [redeemableId, setRedeemableId] = useState("");

  const dashboard = useQuery({
    queryKey: [
      "check-in-dashboard",
      hackathonId,
      search,
      status,
      checkedIn,
      redeemableId,
    ],
    queryFn: ({ signal }) =>
      api
        .get("application/check-in-dashboard", {
          searchParams: {
            hackathonId,
            search,
            status,
            checkedIn,
            redeemableId,
          },
          signal,
        })
        .json<Dashboard>(),
    refetchInterval: 15000,
  });

  const data = dashboard.data;

  const [busyUser, setBusyUser] = useState<string | null>(null);
  const [actionError, setActionError] = useState("");

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
    <main className="max-w-7xl p-4 sm:p-6 space-y-6">
      <h1 className="text-2xl font-bold">Check-in Dashboard</h1>

      {data && (
        <div className="grid gap-4 sm:grid-cols-3">
          <div className="rounded-lg border border-border p-4">
            <p className="text-sm text-text-secondary">
              Confirmed hackers checked in
            </p>
            <p className="text-2xl font-semibold">
              {data.confirmedCheckedIn} / {data.confirmed}
            </p>
          </div>
          <div className="rounded-lg border border-border p-4">
            <p className="text-sm text-text-secondary">Day-of signups</p>
            <p className="text-2xl font-semibold">{data.dayOfSignups}</p>
          </div>
          <div className="rounded-lg border border-border p-4">
            <p className="text-sm text-text-secondary">Standby waiting</p>
            <p className="text-2xl font-semibold">{data.standbyWaiting}</p>
          </div>
        </div>
      )}

      <div className="flex flex-wrap gap-3">
        <label className="flex flex-col gap-1">
          Search by name or email
          <input
            className={fieldClass}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            type="search"
          />
        </label>
        <label className="flex flex-col gap-1">
          Event ID
          <input
            className={fieldClass}
            value={hackathonId}
            placeholder="Active event"
            onChange={(e) => setHackathonId(e.target.value)}
          />
        </label>
        <label className="flex flex-col gap-1">
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
        <label className="flex flex-col gap-1">
          Check-in
          <select
            className={fieldClass}
            value={checkedIn}
            onChange={(e) => setCheckedIn(e.target.value)}
          >
            <option value="">Everyone</option>
            <option value="yes">Checked in</option>
            <option value="no">Not checked in</option>
          </select>
        </label>
        <label className="flex flex-col gap-1">
          Redeemable ID
          <input
            className={fieldClass}
            value={redeemableId}
            placeholder="All redeemables"
            onChange={(e) => setRedeemableId(e.target.value)}
          />
        </label>
      </div>

      <p className="text-sm text-text-secondary">
        Previously registered hackers appear before day-of signups in the
        standby queue.
      </p>

      {actionError && <p role="alert">{actionError}</p>}

      {dashboard.isPending && <p>Loading dashboard…</p>}
      {dashboard.isError && (
        <p role="alert">Unable to load dashboard: {dashboard.error.message}</p>
      )}

      {data && (
        <div className="overflow-x-auto">
          <table className="w-full text-sm text-left">
            <thead>
              <tr className="border-b border-border">
                {[
                  "Name",
                  "Email",
                  "Status",
                  "Registration",
                  "Standby arrival",
                  "Check-in",
                  "Redeemables",
                  "Action",
                ].map((label) => (
                  <th key={label} className="p-3 whitespace-nowrap">
                    {label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {data.rows.map((row) => (
                <tr key={row.userId} className="border-b border-border">
                  <td className="p-3">{row.name}</td>
                  <td className="p-3">{row.email}</td>
                  <td className="p-3">{row.status.replaceAll("_", " ")}</td>
                  <td className="p-3">
                    {row.signupSource === "day_of"
                      ? "Day-of signup"
                      : "Preregistered"}
                  </td>
                  <td className="p-3">
                    {row.standbyArrival
                      ? new Date(row.standbyArrival).toLocaleString()
                      : "—"}
                  </td>
                  <td className="p-3">
                    {row.checkedInAt
                      ? new Date(row.checkedInAt).toLocaleString()
                      : "Not checked in"}
                  </td>
                  <td className="p-3">
                    {row.redemptions
                      .map((r) => `${r.name} × ${r.amount}`)
                      .join(", ") || "—"}
                  </td>
                  <td className="p-3">
                    <div className="flex flex-col gap-2">
                      {data.hackathonId === "xii" &&
                        ["waitlisted", "rejected"].includes(row.status) && (
                          <button
                            className={fieldClass}
                            disabled={busyUser !== null}
                            onClick={() => updateStandby(row, false)}
                          >
                            Record standby arrival
                          </button>
                        )}
                      {row.status === "waitlist_confirmed" && (
                        <button
                          className={fieldClass}
                          disabled={busyUser !== null}
                          onClick={() => updateStandby(row, true)}
                        >
                          Accept
                        </button>
                      )}
                      {row.status === "confirmed" && !row.checkedInAt && (
                        <span>Use the NFC phone app</span>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
              {data.rows.length === 0 && (
                <tr>
                  <td colSpan={8} className="p-3">
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
