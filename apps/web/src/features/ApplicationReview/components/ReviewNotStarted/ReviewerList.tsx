import { Checkbox, CheckboxGroup } from "@/components/ui/Checkbox";
import { NumberField, Input } from "react-aria-components";
import { Tooltip } from "@/components/ui/Tooltip";
import TablerHelpCircle from "~icons/tabler/help-circle";
import type { StaffUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";
import type { AssignedReviewer } from "./ReviewerAssignmentModal";
import type { Dispatch, SetStateAction } from "react";
import { Button } from "@/components/ui/Button";

interface Props {
  staff: StaffUsers;
  assigned: AssignedReviewer[];
  setAssigned: Dispatch<SetStateAction<AssignedReviewer[]>>;
}

export default function ReviewerList({ staff, assigned, setAssigned }: Props) {
  const handleSelectionChange = (selected: string[]) => {
    const updated = selected.map((id) => {
      const existing = assigned.find((r) => r.id === id);
      return existing ?? { id, amount: null };
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
              <TablerHelpCircle className="h-4 w-4" />
            </Button>
          </Tooltip>
          <p># of Apps</p>
        </div>
      </div>

      <CheckboxGroup onChange={handleSelectionChange}>
        {staff.map((s) => {
          const isSelected = assigned.some((a) => a.id === s.id);

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
                        r.id === s.id
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
