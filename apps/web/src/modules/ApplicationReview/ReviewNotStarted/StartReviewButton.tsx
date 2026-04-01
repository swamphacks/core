import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { DialogTrigger } from "react-aria-components";
import ReviewerAssignmentModal from "./ReviewerAssignmentModal";
import type { Event } from "@/modules/Event/schemas/event";
import type { ApplicationStatistics } from "@/modules/Application/hooks/useApplicationStatistics";
import type { StaffUsers } from "@/modules/PlatformAdmin/EventManager/hooks/useEventStaffUsers";

interface Props {
  event: Event;
  stats: ApplicationStatistics;
  staff: StaffUsers;
  validNumOfApplicants: boolean;
  validEventPhase: boolean;
}

export default function StartReviewButton({
  event,
  stats,
  staff,
  validNumOfApplicants,
  validEventPhase,
}: Props) {
  const [open, setOpen] = useState(false);

  return (
    <DialogTrigger onOpenChange={setOpen}>
      <Button
        variant="primary"
        className="w-fit"
        isDisabled={!(validNumOfApplicants && validEventPhase)}
      >
        Start Now
      </Button>

      {!validNumOfApplicants && (
        <p className="text-text-secondary text-sm">
          There are no submitted applications to review yet.
        </p>
      )}

      {!validEventPhase && (
        <p className="text-text-secondary text-sm">
          The application period is still ongoing.
        </p>
      )}

      {open && (
        <ReviewerAssignmentModal
          event={event}
          stats={stats}
          staff={staff}
          onClose={() => setOpen(false)}
        />
      )}
    </DialogTrigger>
  );
}
