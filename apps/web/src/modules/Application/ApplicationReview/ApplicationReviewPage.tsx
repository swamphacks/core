import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import { Sheet } from "@/components/ui/Sheet";
import { Table } from "@/components/ui/Table";
import type { UserContext } from "@/lib/auth/types";
import type {
  ApplicationAutoDecisionRequest,
  StaffHackathon,
} from "@/lib/openapi/types";
import ReviewerAssignmentModal from "@/modules/Application/ApplicationReview/ReviewersAssignmentModal";
import ReviewNotStartedAdmin from "@/modules/Application/ApplicationReview/ReviewNotStartedAdmin";
import { useApplicationReviewAdminActions } from "@/modules/Application/hooks/useApplicationReviewAdminActions";
import {
  useAutoDecisionRequests,
  useUpdateAutoDecisionRequest,
} from "@/modules/Application/hooks/useAutoDecisionRequests";
import { Link } from "@tanstack/react-router";
import { useMemo, useState } from "react";
import { DialogTrigger } from "react-aria-components";
import TablerBrowserShare from "~icons/tabler/browser-share";
import TablerUserEdit from "~icons/tabler/user-edit";
import { useApplicationForReview } from "../hooks/useApplicationForReview";
import { Popover } from "@/components/ui/Popover";
import {
  getCoreRowModel,
  getFilteredRowModel,
  useReactTable,
  type ColumnDef,
  type ColumnFiltersState,
} from "@tanstack/react-table";
import { useReviewersProgress } from "@/modules/Application/hooks/useReviewersProgress";
import { Input } from "@/components/ui/Field";
import { Select } from "@/components/ui/Select";
import TablerSearch from "~icons/tabler/search";
import useParsedForm from "@/modules/Application/hooks/useParsedForm";
import ApplicationResponsesViewer from "@/modules/Application/ApplicationResponsesViewer";

interface ApplicationReviewPageProps {
  hackathon: StaffHackathon;
  user: UserContext;
}

export default function ApplicationReviewPage({
  hackathon,
  user,
}: ApplicationReviewPageProps) {
  const { updateStatus } = useApplicationReviewAdminActions();

  const handleStartApplicationReview = async () => {
    await updateStatus.mutateAsync(true);
  };

  const handleEndApplicationReview = async () => {
    if (window.confirm("Are you sure you want to end application review?")) {
      await updateStatus.mutateAsync(false);
    }
  };

  if (!hackathon.application_review_started) {
    if (user.role === "admin") {
      return (
        <ReviewNotStartedAdmin
          hackathon={hackathon}
          startApplicationReview={handleStartApplicationReview}
        />
      );
    } else {
      return (
        <p className="text-text-secondary">
          Application review has not started.
        </p>
      );
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-2 my-2 w-60">
        <Link to="/application-review/workspace" target="_blank">
          <Button className="w-full">
            <TablerBrowserShare />
            Open Workspace
          </Button>
        </Link>
        {user.role === "admin" && (
          <Button
            onClick={handleEndApplicationReview}
            variant="danger"
            className="w-full"
          >
            End Application Review
          </Button>
        )}
      </div>

      {user.role === "admin" && (
        <div className="flex flex-col lg:flex-row gap-3 lg:items-start">
          <Reviewers />
          <AutoDecisionRequestsTable />
        </div>
      )}
    </div>
  );
}

function AutoDecisionRequestsTable() {
  const { data, isLoading, isError } = useAutoDecisionRequests();
  const updateRequest = useUpdateAutoDecisionRequest();
  const parsedForm = useParsedForm();

  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);

  const requests: ApplicationAutoDecisionRequest = data ?? [];
  type RequestRow = ApplicationAutoDecisionRequest[number];

  const columns: ColumnDef<RequestRow>[] = useMemo(
    () => [
      {
        id: "reviewer_name",
        header: "Reviewer",
        accessorKey: "reviewer_name",
        size: 200,
        cell: ({ row }) => {
          const avatarUrl = row.original.reviewer_image;
          return (
            <div className="flex items-center gap-2">
              {avatarUrl ? (
                <img
                  src={avatarUrl}
                  alt={"user avatar"}
                  className="h-8 w-8 rounded-full object-cover"
                />
              ) : (
                <div className="h-10 w-10 rounded-full bg-gray-300 dark:bg-neutral-700 flex items-center justify-center">
                  <span className="text-gray-600 dark:text-neutral-400">
                    N/A
                  </span>
                </div>
              )}
              <div className="text-sm">
                <div className="font-medium">{row.original.reviewer_name}</div>
              </div>
            </div>
          );
        },
      },
      {
        header: "Application",
        size: 90,
        cell: ({ row }) => (
          <DialogTrigger>
            <Button variant="secondary" size="sm" className="h-8">
              Open
            </Button>
            <Sheet sheetClassName="w-full sm:w-160 lg:w-200">
              <ApplicationResponsesViewer
                parsedForm={parsedForm!}
                applicationId={row.original.application_id}
              />
              {/* <UserApplicationSideDrawer
                applicationId={row.original.application_id}
              /> */}
            </Sheet>
          </DialogTrigger>
        ),
      },
      {
        header: "Decision",
        accessorKey: "requested_decision",
        size: 130,
        cell: ({ row }) =>
          row.original.requested_decision === "auto_accept"
            ? "Auto Accept"
            : "Auto Reject",
      },
      {
        header: "Justification",
        size: 110,
        cell: ({ row }) => (
          <DialogTrigger>
            <Button variant="secondary" size="sm" className="h-8">
              View
            </Button>
            <Popover>
              <div className="p-2">
                <p className="max-w-60 text-wrap">
                  {row.original.justification}
                </p>
              </div>
            </Popover>
          </DialogTrigger>
        ),
      },
      {
        header: "Created At",
        size: 200,
        cell: ({ row }) => new Date(row.original.created_at).toLocaleString(),
      },
      {
        header: "Actions",
        id: "actions",
        size: 150,
        accessorFn: (row) =>
          row.approved_or_denied_by
            ? row.approved
              ? "approved"
              : "denied"
            : "pending",
        cell: (info) => {
          const row = info.row.original as RequestRow;
          const isResolved = row.approved_or_denied_by !== null;

          return isResolved ? (
            <span className="text-sm text-text-secondary">
              {row.approved ? "Approved" : "Denied"}
            </span>
          ) : (
            <div className="flex gap-2">
              <Button
                className="h-8 px-2"
                size="sm"
                onClick={() =>
                  updateRequest.mutateAsync({
                    requestId: row.id,
                    approved: true,
                  })
                }
                isDisabled={updateRequest.isPending}
              >
                Approve
              </Button>
              <Button
                className="h-8 px-2"
                size="sm"
                variant="danger"
                onClick={() =>
                  updateRequest.mutateAsync({
                    requestId: row.id,
                    approved: false,
                  })
                }
                isDisabled={updateRequest.isPending}
              >
                Deny
              </Button>
            </div>
          );
        },
      },
    ],
    [parsedForm],
  );

  const table = useReactTable({
    globalFilterFn: "includesString",
    columns,
    data: data ?? [],
    state: { columnFilters },
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    onColumnFiltersChange: setColumnFilters,
  });

  if (isLoading) {
    return <p>Loading auto decision requests...</p>;
  }

  if (isError) {
    return (
      <p className="text-sm text-rose-600">
        Failed to load auto decision requests.
      </p>
    );
  }

  return (
    <div className="overflow-x-auto rounded-lg bg-surface p-4 min-w-100 max-w-fit">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold">Auto Decision Requests</h2>
        <span className="text-sm text-text-secondary">
          {requests.length} request{requests.length === 1 ? "" : "s"}
        </span>
      </div>

      <div className="grid lg:grid-cols-[2.5fr_3fr_3fr_1fr] gap-2 mb-4 items-stretch">
        <div className="">
          <Input
            type="text"
            className="max-w-50 rounded-md"
            placeholder="Filter reviewers..."
            value={
              (columnFilters.find((f) => f.id === "reviewer_name")?.value as
                | string
                | undefined) ?? ""
            }
            onChange={(e) => {
              const val = e.target.value;
              setColumnFilters((prev) => {
                const next = prev.filter((p) => p.id !== "reviewer_name");
                return val === ""
                  ? next
                  : [...next, { id: "reviewer_name", value: val }];
              });
            }}
            icon={TablerSearch}
          />
        </div>
        <div className="flex items-center gap-2">
          <span className="text-sm">Decision Filter:</span>
          <Select
            className="text-sm max-w-40"
            selectedKey={String(
              (columnFilters.find((f) => f.id === "requested_decision")
                ?.value as string) ?? "",
            )}
            items={[
              { id: "", name: "All" },
              { id: "auto_accept", name: "Auto Accept" },
              { id: "auto_reject", name: "Auto Reject" },
            ]}
            onSelectionChange={(key) => {
              const val = key === "" ? undefined : key;
              setColumnFilters((prev) => {
                const next = prev.filter((p) => p.id !== "requested_decision");
                return val === undefined
                  ? next
                  : [...next, { id: "requested_decision", value: val }];
              });
            }}
            children={null}
          />
        </div>
        <div className="flex items-center gap-2">
          <span className="text-sm">Approval Filter:</span>
          <Select
            className="text-sm max-w-40"
            selectedKey={String(
              (columnFilters.find((f) => f.id === "actions")
                ?.value as string) ?? "",
            )}
            items={[
              { id: "", name: "All" },
              { id: "approved", name: "Approved" },
              { id: "denied", name: "Denied" },
              { id: "pending", name: "Pending" },
            ]}
            onSelectionChange={(key) => {
              const val = key === "" ? undefined : key;
              setColumnFilters((prev) => {
                const next = prev.filter((p) => p.id !== "actions");
                return val === undefined
                  ? next
                  : [...next, { id: "actions", value: val }];
              });
            }}
            children={null}
          />
        </div>
        {/* <div className="flex items-center gap-4 text-text-secondary"></div> */}
      </div>

      {requests.length === 0 ? (
        <p className="text-sm text-text-secondary">
          No auto decision requests found.
        </p>
      ) : (
        <Table
          className="max-h-100 overflow-y-auto"
          headerClassName="text-text-secondary text-sm"
          table={table}
          showPagination={false}
        />
      )}
    </div>
  );
}

interface UserApplicationSideDrawerProps {
  applicationId: string;
}

export function UserApplicationSideDrawer({
  applicationId,
}: UserApplicationSideDrawerProps) {
  const applicationReviewDetails = useApplicationForReview(applicationId);

  if (!applicationReviewDetails.data || applicationReviewDetails.isLoading) {
    return <p>Loading...</p>;
  }

  const getHackathonExperienceText = (experience: string) => {
    switch (experience) {
      case "first_time":
        return "Swamphacks would be my first!";
      case "one":
        return "1";
      case "two":
        return "2";
      case "three":
        return "3";
      case "four_or_more":
        return "4+";
      default:
        return "";
    }
  };

  const getProjectExperienceText = (experience: string) => {
    switch (experience) {
      case "no_experience":
        return "Swamphacks would be my first!";
      case "course_experience":
        return "From courses";
      case "independent_project":
        return "Yes";
      default:
        return "";
    }
  };

  const appFields = applicationReviewDetails.data.application;
  const resume = applicationReviewDetails.data.resumeUrl;

  function ApplicantInfo() {
    return (
      <div className="p-4 rounded-md border border-input-border bg-card">
        <h2 className="block mb-3 text-lg font-medium">
          Applicant Information
        </h2>
        <div className="flex gap-15 items-start">
          <div className="grid grid-cols-1 gap-4">
            <div>
              <div className="text-text-secondary">Name</div>
              <div className="font-medium">
                {appFields.firstName + " " + appFields.lastName}
              </div>
            </div>

            <div>
              <div className="text-text-secondary">Major(s)</div>
              <div className="font-medium">{appFields.majors}</div>
            </div>

            <div>
              <div className="text-text-secondary">School</div>
              <div className="truncate max-w-60">{appFields.school}</div>
            </div>

            <div>
              <div className="text-text-secondary">Graduation Year</div>
              <div className="font-medium">{appFields.graduationYear}</div>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4">
            <div>
              <div className="text-text-secondary">
                # of Hackathons Attended
              </div>
              <div className="font-medium">
                {getHackathonExperienceText(appFields.experience)}
              </div>
            </div>
            <div>
              <div className="text-text-secondary">Project Experience</div>
              <div className="font-medium">
                {getProjectExperienceText(appFields.projectExperience)}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  function Essays() {
    return (
      <div className="space-y-4">
        <div className="p-4 rounded-md border border-input-border bg-card">
          <h2 className="block mb-3 text-lg font-medium">Essay Responses</h2>
          <div className="space-y-6">
            <div>
              <div className="text-text-secondary mb-2">
                What is your most memorable experience working in a group? What
                did you learn and accomplish?
              </div>
              <div className="whitespace-pre-wrap bg-surface p-3 rounded-md">
                {appFields.essay1}
              </div>
            </div>

            <div>
              <div className="text-text-secondary mb-2">
                Tell us about a project you are most proud of.
              </div>
              <div className="whitespace-pre-wrap bg-surface p-3 rounded-md">
                {appFields.essay2}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="grid gap-6 max-h-[99vh] overflow-y-auto">
      <div className="space-y-3">
        <ApplicantInfo />
        <Essays />
      </div>

      <div className="space-y-4">
        <div className="p-2 rounded-md border border-input-border h-200">
          {resume === "" ? (
            <p>No resume provided.</p>
          ) : (
            <object
              className="w-full h-full"
              type="application/pdf"
              data={resume}
            >
              <p>
                Your browser does not support PDFs.{" "}
                <a href={resume}>Download the PDF</a>.
              </p>
            </object>
          )}
        </div>
      </div>
    </div>
  );
}

function Reviewers() {
  const { data, isLoading, isError } = useReviewersProgress();

  const reviewers = data ?? [];
  const columns: ColumnDef<(typeof reviewers)[number]>[] = useMemo(
    () => [
      {
        header: "Reviewer",
        size: 190,
        cell: ({ row }) => {
          const avatarUrl = row.original.image;
          return (
            <div className="flex items-center gap-2">
              {avatarUrl ? (
                <img
                  src={avatarUrl}
                  alt={"user avatar"}
                  className="h-8 w-8 rounded-full object-cover"
                />
              ) : (
                <div className="h-10 w-10 rounded-full bg-gray-300 dark:bg-neutral-700 flex items-center justify-center">
                  <span className="text-gray-600 dark:text-neutral-400">
                    N/A
                  </span>
                </div>
              )}
              <span className="text-sm inline-block max-w-40 font-medium truncate">
                {row.original.name}
              </span>
            </div>
          );
        },
      },
      {
        header: "Assigned",
        size: 80,
        cell: ({ row }) => row.original.total_assigned,
      },
      {
        header: "Completed",
        size: 60,
        cell: ({ row }) => row.original.completed_count,
      },
    ],
    [],
  );

  const table = useReactTable({
    columns,
    data: reviewers,
    getCoreRowModel: getCoreRowModel(),
  });

  if (isLoading) {
    return <p>Loading reviewers progress...</p>;
  }

  if (isError) {
    return <p>Unable to load reviewers progress data.</p>;
  }

  const totalApps = reviewers.reduce((prev, curr) => {
    return curr.total_assigned + prev;
  }, 0);

  const completed = reviewers.reduce((prev, curr) => {
    return curr.completed_count + prev;
  }, 0);

  return (
    <div className="rounded-lg bg-surface p-4 min-w-100 max-w-200">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold">Reviewers</h2>
        <span className="text-sm text-text-secondary">
          {reviewers.length} reviewer{reviewers.length === 1 ? "" : "s"}
        </span>
      </div>

      {reviewers.length === 0 ? (
        <p className="text-sm text-text-secondary">
          No reviewers were assigned.
        </p>
      ) : (
        <Table
          className="max-h-100 overflow-y-auto"
          headerClassName="text-text-secondary text-sm"
          table={table}
          showPagination={false}
        />
      )}

      <div className="mt-3">
        <div className="text-text-secondary text-sm my-1">
          <p>Applications under review: {totalApps}</p>
          <p>Reviews completed: {completed}</p>
        </div>
        <AssignReviewers
          text={
            reviewers.length === 0 ? "Assign Reviewers" : "Re-assign Reviewers"
          }
          alreadyAssigned={reviewers.length > 0}
        />
      </div>
    </div>
  );
}

function AssignReviewers({
  alreadyAssigned,
  text = "Assign Reviewers",
}: {
  text?: string;
  alreadyAssigned: boolean;
}) {
  const [open, setOpen] = useState(false);

  return (
    <DialogTrigger onOpenChange={setOpen}>
      <Button variant="secondary" className="w-full">
        <TablerUserEdit />
        {text}
      </Button>
      {open && (
        <Modal>
          <ReviewerAssignmentModal
            alreadyAssigned={alreadyAssigned}
            onClose={() => setOpen(false)}
          />
        </Modal>
      )}
    </DialogTrigger>
  );
}
