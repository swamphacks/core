import TablerClipboardData from "~icons/tabler/clipboard-data";
import type { EventOverview } from "@/features/EventOverview/hooks/useEventOverview";
import TablerArrowNarrowRight from "~icons/tabler/arrow-narrow-right";
import { Link } from "react-aria-components";
import { useRouter } from "@tanstack/react-router";
import SubmissionsChart from "@/features/Application/components/SubmissionsChart";

interface ApplicationOverviewProps {
  data: EventOverview;
  eventId: string;
}

export default function ApplicationOverview({
  data,
  eventId,
}: ApplicationOverviewProps) {
  const router = useRouter();

  return (
    <div className="border border-input-border rounded-md px-4 py-3 w-full">
      <div className="space-y-2 md:space-y-0 md:flex justify-between items-start relative">
        <p className="flex items-center gap-1 bg-badge-bg-attending rounded-md px-2 py-1 w-fit">
          <TablerClipboardData /> Applications
        </p>
        <div className="flex gap-2 items-center text-text-secondary text-sm">
          <span>Deadline:</span>
          {new Intl.DateTimeFormat("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "numeric",
            minute: "2-digit",
            hour12: true,
          }).format(new Date(data.event_details.application_close))}
        </div>
      </div>

      <div className="space-y-3 md:space-y-0 md:flex items-end gap-3">
        <div className="flex gap-3 mt-3 items-center min-w-0">
          <div className="bg-surface py-2 px-4 rounded-md">
            <p className="text-text-secondary">Total</p>
            <p className="text-xl md:text-2xl">
              {data.application_status_stats.started +
                data.application_status_stats.submitted}
            </p>
          </div>

          <div>
            <p className="text-xl">=</p>
          </div>

          <div className="flex gap-3 items-center">
            <div className="bg-surface py-2 px-4 rounded-md">
              <p className="text-text-secondary">Started</p>
              <p className="text-xl md:text-2xl">
                {data.application_status_stats.started}
              </p>
            </div>

            <div>
              <p className="text-xl">+</p>
            </div>

            <div className="bg-surface py-2 px-4 rounded-md">
              <p className="text-text-secondary">Submitted</p>
              <p className="text-xl md:text-2xl">
                {data.application_status_stats.submitted}
              </p>
            </div>
          </div>
        </div>
        <Link
          className="flex items-center gap-1 underline text-blue-400 text-sm opacity-90 cursor-pointer"
          onClick={() => {
            router.navigate({
              to: "/events/$eventId/dashboard/application-statistics",
              params: { eventId },
            });
          }}
        >
          View more insights <TablerArrowNarrowRight />
        </Link>
      </div>

      <SubmissionsChart submission_stats={data.application_submission_stats} />
    </div>
  );
}
