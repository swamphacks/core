import { createFileRoute } from "@tanstack/react-router";
import { Heading, Link } from "react-aria-components";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_attendee/schedule",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const notionUrl =
    "https://swamphack.notion.site/29a3b41de22f8061a741c84448a3f5ce?pvs=25#2ee3b41de22f80cbb2cee020b7ad4be3";

  return (
    <main className="max-w-4xl p-6">
      <Heading className="text-xl lg:text-2xl font-semibold mb-6 text-slate-900 dark:text-slate-50">
        Event Schedule
      </Heading>

      <div className="grid gap-6">
        <Link
          href={notionUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="group relative flex flex-col md:flex-row items-start md:items-center justify-between gap-6 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-sm transition-all hover:border-blue-500 dark:hover:border-blue-400 hover:shadow-md focus-visible:outline focus-visible:outline-2 focus-visible:outline-blue-600 focus-visible:outline-offset-2 no-underline"
        >
          <div className="flex items-center gap-5">
            {/* Themed Icon Container */}
            <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-xl bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 group-hover:bg-blue-600 dark:group-hover:bg-blue-500 group-hover:text-white transition-colors">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="28"
                height="28"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                <line x1="16" y1="2" x2="16" y2="6"></line>
                <line x1="8" y1="2" x2="8" y2="6"></line>
                <line x1="3" y1="10" x2="21" y2="10"></line>
              </svg>
            </div>

            <div>
              <div className="flex items-center gap-2 mb-1">
                <h3 className="text-lg font-bold text-slate-900 dark:text-slate-100">
                  View Full Schedule
                </h3>
                <svg
                  className="text-slate-400 dark:text-slate-500 group-hover:text-blue-500 dark:group-hover:text-blue-400 transition-colors"
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
              </div>
              <p className="text-slate-500 dark:text-slate-400 text-sm md:text-base leading-relaxed max-w-md">
                Find the workshop timings, meal breaks, and hacker activities on
                our live Notion page.
              </p>
            </div>
          </div>

          {/* Themed Action Badge */}
          <div className="flex items-center gap-1.5 font-semibold text-sm text-blue-600 dark:text-blue-400 md:bg-blue-50 md:dark:bg-blue-900/20 md:px-5 md:py-2.5 md:rounded-full group-hover:bg-blue-100 dark:group-hover:bg-blue-900/40 transition-colors">
            Open Notion
            <svg
              width="16"
              height="16"
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
      </div>
    </main>
  );
}
