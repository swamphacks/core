import type { Event } from "@/features/Event/schemas/event";
import StartReviewButton from "./StartReviewButton";
import type { ApplicationStatistics } from "@/features/Application/hooks/useApplicationStatistics";
import type { StaffUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";

interface Props {
  event: Event;
  stats: ApplicationStatistics;
  staff: StaffUsers;
}

export default function ReviewNotStarted({ event, stats, staff }: Props) {
  const validNumOfApplicants = stats.status_stats.submitted > 0;
  const validEventPhase = new Date(event.application_close) <= new Date();

  return (
    <div className="flex flex-col gap-4">
      <p className="text-text-secondary">Application review has not started.</p>

      <StartReviewButton
        event={event}
        stats={stats}
        staff={staff}
        validNumOfApplicants={validNumOfApplicants}
        validEventPhase={validEventPhase}
      />
    </div>
  );
}
