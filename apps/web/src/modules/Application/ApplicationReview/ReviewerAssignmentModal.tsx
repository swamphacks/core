import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import ReviewerList from "./ReviewerList";
import SummaryFooter from "./SummaryFooter";
import {
  useApplicationReviewAdminActions,
  type AssignedReviewer,
} from "../hooks/useApplicationReviewAdminActions";
import { toast } from "react-toastify";
import { useApplicationStatistics } from "@/modules/Application/hooks/useApplicationStatistics";
import { useHackathonStaff } from "@/modules/Hackathon/hooks/useHackathonStaff";

interface Props {
  onClose: () => void;
}

export default function ReviewerAssignmentModal({ onClose }: Props) {
  const [assigned, setAssigned] = useState<AssignedReviewer[]>([]);
  const { assign } = useApplicationReviewAdminActions();
  const staff = useHackathonStaff();
  const stats = useApplicationStatistics();

  if (staff.isLoading || !staff.data) {
    return <p>Loading staff...</p>;
  }

  if (stats.isLoading || !stats.data) {
    return <div>Loading application stats...</div>;
  }

  const numAssigned = assigned.reduce(
    (acc, curr) => acc + (curr.amount ?? 0),
    0,
  );

  const remaining = stats.data.statusStats.under_review - numAssigned;

  const handleSubmit = async () => {
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

    await assign.mutateAsync(assigned, {
      onSuccess: () => {
        toast.success("Reviewers assigned successfully.");
        onClose();
      },
      onError: () => {
        toast.error("Failed to assign reviewers. Please try again.");
      },
    });
  };

  return (
    <Modal>
      <div className="flex flex-col h-full">
        <h2>Select Reviewers</h2>
        <div className="flex-3 overflow-y-auto mt-2 max-h-[500px] px-1">
          <ReviewerList
            staff={staff.data}
            assigned={assigned}
            setAssigned={setAssigned}
          />
        </div>
        <div className="mt-4 flex gap-2 flex-col flex-1">
          <SummaryFooter
            numAssigned={numAssigned}
            total={stats.data.statusStats.under_review}
          />

          <div className="flex flex-row w-full justify-end gap-2">
            <Button className="w-fit" slot="close" variant="secondary">
              Cancel
            </Button>
            <Button className="w-fit" variant="primary" onClick={handleSubmit}>
              Assign Reviewers
            </Button>
          </div>
        </div>
      </div>
    </Modal>
  );
}
