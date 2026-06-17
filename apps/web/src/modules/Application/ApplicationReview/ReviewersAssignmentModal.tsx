import { useState } from "react";
import { Button } from "@/components/ui/Button";
import {
  useApplicationReviewAdminActions,
  type AssignedReviewer,
} from "../hooks/useApplicationReviewAdminActions";
import { toast } from "react-toastify";
import { useApplicationStatistics } from "@/modules/Application/hooks/useApplicationStatistics";
import { useHackathonStaff } from "@/modules/Hackathon/hooks/useHackathonStaff";
import { Checkbox, CheckboxGroup } from "@/components/ui/Checkbox";
import { NumberField, Input } from "react-aria-components";
import { Tooltip } from "@/components/ui/Tooltip";
import TablerHelpCircle from "~icons/tabler/help-circle";
import type { Dispatch, SetStateAction } from "react";
import type { HackathonStaff } from "@/lib/openapi/types";

interface Props {
  onClose: () => void;
  alreadyAssigned: boolean;
}

export default function ReviewerAssignmentModal({
  onClose,
  alreadyAssigned,
}: Props) {
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

    let confirmed;
    if (alreadyAssigned) {
      confirmed = window.confirm(
        "Re-assigning will delete all application reviews and auto decision requests. Are you sure you want to proceed?",
      );
    }

    if (confirmed || !alreadyAssigned) {
      await assign.mutateAsync(assigned, {
        onSuccess: () => {
          toast.success("Reviewers assigned successfully.");
          onClose();
        },
        onError: () => {
          toast.error("Failed to assign reviewers. Please try again.");
        },
      });
    }
  };

  return (
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
        <div className="flex flex-col items-end text-text-secondary text-sm my-4">
          <p>Assigned Applications: {numAssigned}</p>
          <p>Total Applications: {stats.data.statusStats.under_review}</p>
          <p>Remaining: {stats.data.statusStats.under_review - numAssigned}</p>
        </div>
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
  );
}

interface ReviewerListProps {
  staff: HackathonStaff;
  assigned: AssignedReviewer[];
  setAssigned: Dispatch<SetStateAction<AssignedReviewer[]>>;
}

function ReviewerList({ staff, assigned, setAssigned }: ReviewerListProps) {
  const handleSelectionChange = (selected: string[]) => {
    const updated = selected.map((userId) => {
      const existing = assigned.find((r) => r.userId === userId);
      return existing ?? { userId, amount: null };
    });

    setAssigned(updated);
  };

  return (
    <>
      <div className="flex justify-between text-sm text-text-secondary mt-2">
        <p>Reviewer</p>

        <div className="flex gap-1 items-center">
          <Tooltip
            tooltipProps={{
              label:
                "If blank, remaining applications will be auto-distributed.",
            }}
            triggerProps={{
              delay: 200,
            }}
          >
            <Button variant="unstyled" className="p-0 m-0">
              <TablerHelpCircle className="size-5" />
            </Button>
          </Tooltip>
          <p># of Apps</p>
        </div>
      </div>

      <CheckboxGroup onChange={handleSelectionChange}>
        {staff.map((s) => {
          const isSelected = assigned.some((a) => a.userId === s.id);

          return (
            <div className="flex justify-between w-full" key={s.id}>
              <Checkbox value={s.id} className="flex-1">
                <div className="flex flex-row items-center gap-2">
                  <span>{s.name || "Anonymous"}</span>
                  <span className="text-xs text-text-secondary hidden sm:block">
                    {s.email}
                  </span>
                </div>
              </Checkbox>

              <NumberField
                aria-label={`Number of apps for ${s.name ?? "Anonymous"}`}
                className="ml-auto"
              >
                <Input
                  disabled={!isSelected}
                  onChange={(e) => {
                    const value = e.target.value;
                    setAssigned((prev) =>
                      prev.map((r) =>
                        r.userId === s.id
                          ? {
                              ...r,
                              amount: value === "" ? null : Number(value),
                            }
                          : r,
                      ),
                    );
                  }}
                  className="border-border border-1 disabled:border-0 w-12 h-6 text-sm px-1 disabled:bg-input-bg-disbaled dark:disabled:bg-neutral-800 disabled:cursor-not-allowed"
                />
              </NumberField>
            </div>
          );
        })}
      </CheckboxGroup>
    </>
  );
}
