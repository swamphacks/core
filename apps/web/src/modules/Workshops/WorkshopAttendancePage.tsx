import { useState } from "react";
import QrWorkshopScanner from "./QrWorkshopScanner";
import { useWorkshopAttendance } from "./useWorkshopAttendance";

export default function WorkshopAttendancePage() {
  const [open, setOpen] = useState(false);
  const attendance = useWorkshopAttendance();

  return (
    <main className="w-full max-w-6xl space-y-5 p-4 sm:p-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">Workshop & social attendance</h1>
          <p className="mt-1 text-sm text-text-secondary">
            Unique QR and NFC check-ins · refreshes every 10 seconds
          </p>
        </div>
        <button
          type="button"
          className="rounded-lg bg-blue-600 px-4 py-2 font-medium text-white"
          onClick={() => setOpen(true)}
        >
          Scan QR
        </button>
      </div>

      {attendance.isPending && <p role="status">Loading attendance…</p>}
      {attendance.isError && (
        <div role="alert" className="rounded-xl border border-border p-4">
          <p>Unable to load attendance.</p>
          <button
            type="button"
            className="mt-2 underline"
            onClick={() => void attendance.refetch()}
          >
            Try again
          </button>
        </div>
      )}
      {attendance.data?.length === 0 && (
        <p className="rounded-xl border border-border p-6 text-text-secondary">
          No workshops or socials are available.
        </p>
      )}

      <section
        aria-label="Attendance by event"
        className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3"
      >
        {attendance.data?.map((event) => (
          <article
            key={event.id}
            className="flex flex-col rounded-xl border border-border p-5"
          >
            <p className="text-xs font-medium uppercase tracking-wide text-text-secondary">
              {event.type === "social" ? "Social" : "Workshop"}
            </p>
            <h2 className="mt-2 text-lg font-semibold">{event.title}</h2>
            <p className="mt-1 text-sm text-text-secondary">
              {new Date(event.startTime).toLocaleString(undefined, {
                month: "short",
                day: "numeric",
                hour: "numeric",
                minute: "2-digit",
              })}
            </p>
            <div className="mt-auto pt-6">
              <p className="text-3xl font-semibold tabular-nums">
                {event.checkedInCount}
              </p>
              <p className="mt-1 text-sm text-text-secondary">Checked in</p>
            </div>
          </article>
        ))}
      </section>

      <QrWorkshopScanner isOpen={open} onClose={() => setOpen(false)} />
    </main>
  );
}
