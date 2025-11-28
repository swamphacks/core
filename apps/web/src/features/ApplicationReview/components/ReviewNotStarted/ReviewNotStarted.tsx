import type { Event } from "@/features/Event/schemas/event";
import StartReviewButton from "./StartReviewButton";
import { useApplicationStatistics } from "@/features/Application/hooks/useApplicationStatistics";
import { useEventStaffUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";

interface Props {
  event: Event;
}

export default function ReviewNotStarted({ event }: Props) {
  const stats = useApplicationStatistics(event.id);
  const staff = useEventStaffUsers(event.id);

  if (stats.isLoading || staff.isLoading || !stats.data || !staff.data) {
    return <div>Loading...</div>;
  }

  const validNumOfApplicants = stats.data.status_stats.submitted > 0;
  const validEventPhase = new Date(event.application_close) <= new Date();

  return (
    <div className="flex flex-col gap-4">
      <p className="text-text-secondary">Application review has not started.</p>

      <StartReviewButton
        event={event}
        stats={stats.data}
        staff={staff.data}
        validNumOfApplicants={validNumOfApplicants}
        validEventPhase={validEventPhase}
      />
    </div>
  );
}
