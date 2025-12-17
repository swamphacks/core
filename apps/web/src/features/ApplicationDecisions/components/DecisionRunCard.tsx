import type { Run } from "../hooks/useDecisionRuns";

interface Props {
  run: Run;
  onRelease?: (runId: string) => void;
}

export default function DecisionRunCard({ run, onRelease }: Props) {
  const statusConfig = {
    running: {
      text: "var(--badge-text-under-review)",
      label: "In Progress",
    },
    completed: {
      text: "var(--badge-text-accepted)",
      label: "Completed",
    },
    failed: {
      text: "var(--badge-text-rejected)",
      label: "Failed",
    },
  };

  const config =
    statusConfig[run.status.bat_run_status] || statusConfig.running;

  return (
    <div className="group flex flex-col sm:flex-row sm:items-center justify-between w-full min-h-[72px] p-4 sm:px-6 rounded-sm border border-border bg-surface transition-all hover:border-zinc-400 dark:hover:border-zinc-500 shadow-sm gap-4">
      {/* Left Section: Always consistent height */}
      <div className="flex flex-wrap items-center gap-x-6 gap-y-3">
        {/* Status Indicator */}
        <div className="flex items-center gap-2.5 min-w-[110px]">
          <div className="relative flex h-2 w-2">
            {run.status.bat_run_status === "running" && (
              <span
                className="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75"
                style={{ backgroundColor: config.text }}
              ></span>
            )}
            <span
              className="relative inline-flex rounded-full h-2 w-2"
              style={{ backgroundColor: config.text }}
            ></span>
          </div>
          <span
            className="text-xs font-bold tracking-wider uppercase"
            style={{ color: config.text }}
          >
            {config.label}
          </span>
        </div>

        <div className="hidden sm:block h-4 w-[1px] bg-border" />

        {/* Stats */}
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-1.5">
            <span className="text-zinc-400 text-[10px]">✔</span>
            <span className="text-sm font-semibold text-text-main">
              {run.accepted_applicants.length}
            </span>
          </div>
          <div className="flex items-center gap-1.5">
            <span className="text-zinc-400 text-[10px]">✖</span>
            <span className="text-sm font-semibold text-text-main">
              {run.rejected_applicants.length}
            </span>
          </div>
        </div>

        {/* Metadata */}
        <div className="flex items-center gap-3">
          <span className="font-mono text-[10px] text-text-secondary opacity-60 group-hover:opacity-100 transition-opacity">
            {run.id.slice(0, 8)}
          </span>
          <span className="text-[10px] text-text-secondary opacity-40">
            {new Date(run.created_at).toLocaleTimeString([], {
              hour: "2-digit",
              minute: "2-digit",
            })}
          </span>
        </div>
      </div>

      {/* Right Section: Occupies space even if empty to prevent height jitter */}
      <div className="flex items-center justify-end sm:ml-auto min-h-[32px]">
        {run.status.bat_run_status === "completed" && onRelease ? (
          <button
            onClick={() => onRelease(run.id)}
            className="w-full sm:w-auto text-xs font-bold px-5 py-2 rounded-sm bg-button-primary text-white hover:bg-button-primary-hover active:bg-button-primary-pressed transition-colors shadow-sm uppercase tracking-tighter"
          >
            Release Results
          </button>
        ) : (
          /* Placeholder to keep layout spacing balanced on desktop */
          <div className="hidden sm:block w-24" />
        )}
      </div>
    </div>
  );
}
