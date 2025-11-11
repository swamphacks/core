import { DialogTrigger, Link } from "react-aria-components";
import type { ApplicationStatistics } from "../hooks/useApplicationStatistics";
import { Modal } from "@/components/ui/Modal";
import SubmissionsChart from "./SubmissionsChart";
import type { EventOverview } from "@/features/EventOverview/hooks/useEventOverview";

interface ApplicationStatsProps {
  data: ApplicationStatistics;
  eventData: EventOverview;
}

export default function ApplicationStats({
  data,
  eventData,
}: ApplicationStatsProps) {
  return (
    <div className="space-y-2">
      <div className="bg-surface py-2 px-4 rounded-md">
        <p className="text-text-secondary">Total</p>
        <p className="text-xl md:text-2xl">
          {data.status_stats.started + data.status_stats.submitted}
        </p>
      </div>

      <div className="bg-surface py-2 px-4 rounded-md">
        <p className="text-text-secondary">Started</p>
        <p className="text-xl md:text-2xl">{data.status_stats.started}</p>
      </div>

      <div className="bg-surface py-2 px-4 rounded-md">
        <p className="text-text-secondary">Submitted</p>
        <p className="text-xl md:text-2xl">{data.status_stats.submitted}</p>
      </div>

      <div className="flex gap-2 items-center text-text-secondary text-sm">
        <span>Application Open:</span>
        {new Intl.DateTimeFormat("en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
          hour: "numeric",
          minute: "2-digit",
          hour12: true,
        }).format(new Date(eventData.event_details.application_open))}
      </div>

      <div className="flex gap-2 items-center text-text-secondary text-sm">
        <span>Deadline:</span>
        {new Intl.DateTimeFormat("en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
          hour: "numeric",
          minute: "2-digit",
          hour12: true,
        }).format(new Date(eventData.event_details.application_close))}
      </div>
      <div>
        <DialogTrigger>
          <Link className="flex items-center gap-1 underline text-blue-400 text-sm opacity-90">
            View submissions chart
          </Link>

          <Modal size="xl" isDismissible>
            <SubmissionsChart
              submission_stats={eventData.application_submission_stats}
            />
          </Modal>
        </DialogTrigger>
      </div>
    </div>
  );
}
