import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import type { Event } from "@/features/Event/schemas/event";
import type { ApplicationStatistics } from "@/features/Application/hooks/useApplicationStatistics";
import type { StaffUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";
import ReviewerList from "./ReviewerList";
import SummaryFooter from "./SummaryFooter";

export interface AssignedReviewer {
  id: string;
  amount: number | null;
}

interface Props {
  event: Event;
  stats: ApplicationStatistics;
  staff: StaffUsers;
  onClose: () => void;
}

export default function ReviewerAssignmentModal({
  event,
  stats,
  staff,
  onClose,
}: Props) {
  const [assigned, setAssigned] = useState<AssignedReviewer[]>([]);

  const numAssigned = assigned.reduce(
    (acc, curr) => acc + (curr.amount ?? 0),
    0,
  );

  const remaining = stats.status_stats.submitted - numAssigned;

  const handleSubmit = () => {
    const hasNull = assigned.some((r) => r.amount === null);
    const hasZero = assigned.some((r) => r.amount === 0);
    const hasNegative = assigned.some((r) => (r.amount ?? 0) < 0);

    if (hasNegative)
      return alert(
        "Assigned applications cannot be negative for any reviewer.",
      );
    if (hasZero)
      return alert(
        "Assigned applications cannot be zero for any reviewer. Either leave blank or assign a positive number.",
      );
    if (remaining < 0) return alert("Assigned applications exceed total.");
    if (remaining > 0 && !hasNull)
      return alert(
        "Leave at least one reviewer blank to distribute remaining.",
      );
    if (assigned.length === 0) return alert("Select reviewers.");

    // Start the review process with assigned reviewers
    //TODO: Implement the actual submission logic
    console.log("Assigned Reviewers:", assigned);
    console.log("Beginning review process for event:", event.id);
    onClose();
  };

  return (
    <Modal>
      <div className="flex flex-col">
        <h2>Select Reviewers</h2>

        <ReviewerList
          staff={staff}
          assigned={assigned}
          setAssigned={setAssigned}
        />

        <SummaryFooter
          numAssigned={numAssigned}
          total={stats.status_stats.submitted}
        />

        <div className="mt-4 flex justify-end gap-2">
          <Button slot="close" variant="secondary">
            Cancel
          </Button>
          <Button variant="primary" onClick={handleSubmit}>
            Assign Reviewers
          </Button>
        </div>
      </div>
    </Modal>
  );
}
