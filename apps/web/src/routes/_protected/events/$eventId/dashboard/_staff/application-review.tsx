import { Button } from "@/components/ui/Button";
import { Checkbox, CheckboxGroup } from "@/components/ui/Checkbox";
import { Modal } from "@/components/ui/Modal";
import { Tooltip } from "@/components/ui/Tooltip";
import { useEventStaffUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import {
  DialogTrigger,
  Heading,
  Input,
  NumberField,
  Separator,
} from "react-aria-components";
import TablerHelpCircle from "~icons/tabler/help-circle";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/application-review",
)({
  component: RouteComponent,
});

interface AssignedReviewer {
  id: string;
  amount: number | null;
}

function RouteComponent() {
  const { user, eventRole } = Route.useRouteContext();
  const { eventId } = Route.useParams();
  const eventStaff = useEventStaffUsers(eventId);

  const [assignedReviewer, setAssignedReviewer] = useState<AssignedReviewer[]>(
    [],
  );

  const numOfApplicationsAssigned = assignedReviewer.reduce((acc, curr) => {
    return acc + (curr.amount ?? 0);
  }, 0);
  const mockNumOfApplications = 1000;
  const remainingApplications =
    mockNumOfApplications - numOfApplicationsAssigned;

  const onHandleSelectionChange = (selected: string[]) => {
    // Each string in selected is a user ID
    // Update the assignedReviewer state accordingly

    const updatedAssignedReviewers = selected.map((id) => {
      const existing = assignedReviewer.find((rev) => rev.id === id);
      return existing ? existing : { id, amount: null };
    });

    setAssignedReviewer(updatedAssignedReviewers);
  };

  const handleSubmit = () => {
    const hasNullAmount = assignedReviewer.some((rev) => rev.amount === null);
    if (remainingApplications < 0) {
      alert("Assigned applications exceed total number of applications.");
      return;
    }

    if (remainingApplications > 0 && !hasNullAmount) {
      // Check if there are any reviewers with null amount
      alert(
        "Please leave at least one reviewer with no specified amount to distribute remaining applications.",
      );
      return;
    }

    if (assignedReviewer.length === 0) {
      alert("R U DEADAHH RN?");
      return;
    }

    console.log("Assigned Reviewers:", assignedReviewer);
  };

  if (!user || eventStaff.isLoading) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          My Team
        </Heading>
        <div className="flex flex-col gap-4 max-w-xl">
          <div className="h-84 w-full md:w-120 bg-neutral-300 dark:bg-neutral-800 rounded animate-pulse"></div>
        </div>
      </main>
    );
  }

  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-4">
        Application Review
      </Heading>

      {eventRole === "staff" ? (
        <div>
          <p className="text-text-secondary">
            Application review has not started yet. Check back later!
          </p>
        </div>
      ) : (
        <div className="flex flex-col gap-4">
          <p className="text-text-secondary">
            Application review has not started.
          </p>
          <DialogTrigger
            onOpenChange={(open) => !open && setAssignedReviewer([])}
          >
            <Button variant="primary" className="w-fit">
              Start Now
            </Button>

            <Modal>
              <div className="flex flex-col">
                <h2>Select Reviewers</h2>
                <div className="flex flex-row w-full items-center justify-between text-text-secondary text-sm mt-2">
                  <p>Reviewer</p>

                  <div className="flex flex-row gap-1 items-center justify-center">
                    <Tooltip
                      tooltipProps={{
                        label:
                          "Number of applications to assign to this reviewer. If left blank, we will automatically distribute remaining applications among all blank selected reviewers.",
                      }}
                      triggerProps={{
                        delay: 200,
                      }}
                    >
                      <Button variant="unstyled" className="p-0 m-0">
                        <TablerHelpCircle className="inline-block h-4 w-4 mt-0.5" />
                      </Button>
                    </Tooltip>

                    <p># of Apps</p>
                  </div>
                </div>
                <CheckboxGroup onChange={onHandleSelectionChange}>
                  {eventStaff.data?.map((staffUser) => (
                    <div
                      className="flex flex-row items-center justify-between w-full"
                      key={staffUser.id}
                    >
                      <Checkbox value={staffUser.id} className="flex-1">
                        <div className="flex flex-row items-center gap-2">
                          <span>{staffUser.name || "Anonymous"}</span>
                          <span className="text-text-secondary text-xs">
                            {staffUser.email}
                          </span>
                        </div>
                      </Checkbox>

                      <NumberField className="ml-auto">
                        <Input
                          onChange={(e) => {
                            const value = e.target.value;
                            setAssignedReviewer((prev) =>
                              prev.map((rev) =>
                                rev.id === staffUser.id
                                  ? {
                                      ...rev,
                                      amount:
                                        value === "" ? null : Number(value),
                                    }
                                  : rev,
                              ),
                            );
                          }}
                          className="border-border border-2 w-12 h-6 text-sm disabled:bg-neutral-200 dark:disabled:bg-neutral-800 disabled:hover:cursor-not-allowed px-1"
                          disabled={
                            !assignedReviewer.find(
                              (rev) => rev.id === staffUser.id,
                            )
                          }
                        />
                      </NumberField>
                    </div>
                  ))}
                </CheckboxGroup>

                <Separator className="my-4 text-neutral-700" />

                <div className="flex flex-col items-end text-text-secondary text-sm">
                  <p>Assigned Applications: {numOfApplicationsAssigned}</p>
                  <p>Total Applications: {mockNumOfApplications}</p>
                  <p>Remaining: {remainingApplications}</p>
                </div>

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
          </DialogTrigger>
        </div>
      )}
    </main>
  );
}
