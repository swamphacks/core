import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { EventAttendanceWithdrawalModal } from "@/features/Event/components/EventAttendanceWithdrawalModal";
import { generateIdentifyIntent } from "@/lib/qr-intents/generate";
import { DialogTrigger, Heading, Link } from "react-aria-components";
import QRCode from "react-qr-code";

interface Props {
  userId: string;
  eventId: string;
}

export default function AttendeeOverview({ userId, eventId }: Props) {
  const identificationIntentString = generateIdentifyIntent(userId);
  const hackerGuideUrl = "https://swamphack.notion.site/sh-xi-hacker-guide";

  const allowWithdrawal = false;

  return (
    <main className="max-w-5xl p-4 sm:p-6">
      <Heading className="text-2xl lg:text-3xl font-semibold mb-8 text-slate-900 dark:text-slate-50">
        Overview
      </Heading>

      <div className="flex flex-col lg:flex-row gap-8">
        {/* Left Column: QR ID Card */}
        <div className="w-full lg:w-72 shrink-0">
          <Card className="relative overflow-hidden border-none bg-gradient-to-br from-blue-600 to-indigo-700 p-1 shadow-xl">
            <div className="bg-white dark:bg-slate-900 rounded-[calc(var(--radius)-4px)] p-6">
              <div className="flex flex-col items-center">
                <div className="mb-4 text-center">
                  <span className="text-[10px] uppercase tracking-widest font-bold text-slate-400 dark:text-slate-500">
                    Attendee Pass
                  </span>
                </div>

                {/* HIGH CONTRAST QR WRAPPER */}
                {/* We force bg-white and text-black here regardless of theme for scan reliability */}
                <div className="p-4 bg-white rounded-2xl shadow-sm ring-1 ring-slate-100">
                  <QRCode
                    size={200}
                    style={{ height: "auto", maxWidth: "100%", width: "100%" }}
                    value={identificationIntentString}
                    viewBox={`0 0 256 256`}
                    bgColor="#FFFFFF"
                    fgColor="#000000"
                    level="H" // High error correction for better scanning
                  />
                </div>

                <div className="mt-6 text-center">
                  <p className="text-sm font-semibold text-slate-900 dark:text-slate-100">
                    Personal QR Code
                  </p>
                  <p className="font-mono text-[10px] text-slate-500 dark:text-slate-400 mt-1 uppercase">
                    ID: {userId.slice(0, 12)}
                  </p>
                </div>
              </div>
            </div>
          </Card>
        </div>

        {/* Right Column: Info & Links */}
        <div className="flex-1 flex flex-col gap-8">
          {/* QR Explanation */}
          <section className="space-y-3">
            <h3 className="text-lg font-bold text-slate-900 dark:text-slate-50">
              Your Identifier
            </h3>
            <div className="prose prose-slate dark:prose-invert">
              <p className="text-slate-600 dark:text-slate-400 leading-relaxed">
                This QR code is your unique digital badge. It identifies you to
                our staff and is used for check-ins, swag pickups, and meal
                redemptions!
              </p>
            </div>
          </section>

          <hr className="border-slate-200 dark:border-slate-800" />

          {/* Unified Style Resource Links */}
          <section className="grid gap-4">
            <Link
              href={hackerGuideUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="group relative flex items-center justify-between p-5 rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm transition-all hover:border-blue-500 dark:hover:border-blue-400 hover:shadow-md no-underline"
            >
              <div className="flex items-center gap-5">
                {/* Consistent Icon Style */}
                <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-amber-50 dark:bg-amber-900/20 text-amber-600 dark:text-amber-400 group-hover:bg-amber-600 group-hover:text-white transition-colors">
                  <svg
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
                    <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
                  </svg>
                </div>

                <div>
                  <h4 className="font-bold text-slate-900 dark:text-slate-100 flex items-center gap-2">
                    Official Hacker Guide
                    <svg
                      className="text-slate-400 opacity-0 group-hover:opacity-100 transition-opacity"
                      width="14"
                      height="14"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2.5"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    >
                      <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                      <polyline points="15 3 21 3 21 9"></polyline>
                      <line x1="10" y1="14" x2="21" y2="3"></line>
                    </svg>
                  </h4>
                  <p className="text-sm text-slate-500 dark:text-slate-400">
                    Have questions? Look here!
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-1.5 font-semibold text-xs text-blue-600 dark:text-blue-400 md:bg-blue-50 md:dark:bg-blue-900/20 md:px-4 md:py-2 md:rounded-full group-hover:bg-blue-100 dark:group-hover:bg-blue-900/40 transition-colors">
                Open Guide
                <svg
                  width="14"
                  height="14"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <polyline points="9 18 15 12 9 6"></polyline>
                </svg>
              </div>
            </Link>
          </section>

          {/* Management / Withdrawal */}
          {allowWithdrawal && (
            <div className="mt-auto pt-6 border-t border-slate-200 dark:border-slate-800">
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div>
                  <p className="font-medium text-slate-900 dark:text-slate-100">
                    Plans changed?
                  </p>
                  <p className="text-sm text-slate-500 dark:text-slate-400">
                    Withdrawal is permanent for this event.
                  </p>
                </div>
                <DialogTrigger>
                  <Button variant="danger" className="shrink-0">
                    Withdraw Attendance
                  </Button>
                  <EventAttendanceWithdrawalModal eventId={eventId} />
                </DialogTrigger>
              </div>
            </div>
          )}
        </div>
      </div>
    </main>
  );
}
