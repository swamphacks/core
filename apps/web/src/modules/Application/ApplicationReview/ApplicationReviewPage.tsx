import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import { Table } from "@/components/ui/Table";
import type { UserContext } from "@/lib/auth/types";
import ReviewerAssignmentModal from "@/modules/Application/ApplicationReview/ReviewersAssignmentModal";
import ReviewNotStartedAdmin from "@/modules/Application/ApplicationReview/ReviewNotStartedAdmin";
import { useApplicationReviewAdminActions } from "@/modules/Application/hooks/useApplicationReviewAdminActions";
import { Link } from "@tanstack/react-router";
import { useMemo, useState } from "react";
import { DialogTrigger } from "react-aria-components";
import TablerBrowserShare from "~icons/tabler/browser-share";
import TablerUserEdit from "~icons/tabler/user-edit";
import {
  getCoreRowModel,
  useReactTable,
  type ColumnDef,
} from "@tanstack/react-table";
import { useReviewersProgress } from "@/modules/Application/hooks/useReviewersProgress";
import { Input } from "@/components/ui/Field";
import { Select } from "@/components/ui/Select";
import TablerSearch from "~icons/tabler/search";
import type { components } from "@/lib/openapi/schema";
import { useDebounce } from "@uidotdev/usehooks";
import AutoDecisionRequestList from "@/modules/Application/ApplicationReview/AutoDecisionRequestList";

interface ApplicationReviewPageProps {
  hackathon: components["schemas"]["Hackathon"];
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
  const [searchInput, setSearchInput] = useState("");
  const debouncedSearchInput = useDebounce(searchInput, 500);

  const [approvalFilter, setApprovalFilter] = useState("all");
  const [decisionFilter, setDecisionFilter] = useState("all");

  return (
    <div className="overflow-x-auto rounded-lg bg-surface p-4 min-w-100 max-w-fit">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold">Auto Decision Requests</h2>
      </div>

      <div className="@container">
        <div className="flex flex-col @2xl:flex-row gap-2 mb-4">
          <div className="grow">
            <Input
              aria-label="search input"
              className="max-w-100 rounded-md"
              placeholder="Search by reviewer or applicant name"
              value={searchInput}
              icon={TablerSearch}
              onChange={(e) => setSearchInput(e.target.value)}
            />
          </div>
          <div className="flex items-center gap-2">
            <span className="text-sm">Decision Filter:</span>
            <Select
              className="text-sm max-w-40"
              value={decisionFilter}
              items={[
                { id: "all", name: "All" },
                { id: "auto_accept", name: "Auto Accept" },
                { id: "auto_reject", name: "Auto Reject" },
              ]}
              onChange={(key) => {
                if (key) setDecisionFilter(key.toString());
              }}
              children={null}
            />
          </div>
          <div className="flex items-center gap-2">
            <span className="text-sm">Approval Filter:</span>
            <Select
              className="text-sm max-w-40"
              value={approvalFilter}
              items={[
                { id: "all", name: "All" },
                { id: "approved", name: "Approved" },
                { id: "denied", name: "Denied" },
                { id: "pending", name: "Pending" },
              ]}
              onChange={(key) => {
                if (key) setApprovalFilter(key.toString());
              }}
              children={null}
            />
          </div>
        </div>
      </div>

      <AutoDecisionRequestList
        searchInput={debouncedSearchInput}
        approvalFilter={approvalFilter}
        decisionFilter={decisionFilter}
      />
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
        cell: ({ row }) => row.original.totalAssigned,
      },
      {
        header: "Completed",
        size: 60,
        cell: ({ row }) => row.original.completedCount,
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
    return curr.totalAssigned + prev;
  }, 0);

  const completed = reviewers.reduce((prev, curr) => {
    return curr.completedCount + prev;
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
