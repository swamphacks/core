import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import { Sheet } from "@/components/ui/Sheet";
import { Table } from "@/components/ui/Table";
import type { UserContext } from "@/lib/auth/types";
import type {
  ApplicationAutoDecisionRequest,
  StaffHackathon,
} from "@/lib/openapi/types";
import ReviewerAssignmentModal from "@/modules/Application/ApplicationReview/ReviewerAssignmentModal";
import ReviewNotStartedAdmin from "@/modules/Application/ApplicationReview/ReviewNotStartedAdmin";
import { useApplicationReviewAdminActions } from "@/modules/Application/hooks/useApplicationReviewAdminActions";
import {
  useAutoDecisionRequests,
  useUpdateAutoDecisionRequest,
} from "@/modules/Application/hooks/useAutoDecisionRequests";
import { Link } from "@tanstack/react-router";
import { useState } from "react";
import { DialogTrigger } from "react-aria-components";
import TablerBrowserShare from "~icons/tabler/browser-share";
import TablerUserEdit from "~icons/tabler/user-edit";
import { useApplicationForReview } from "../hooks/useApplicationForReview";
import { Popover } from "@/components/ui/Popover";
import type { ColumnDef } from "@tanstack/react-table";

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
          <AppReviewAdminActions
            endApplicationReview={handleEndApplicationReview}
          />
        )}
      </div>

      <div>{user.role === "admin" && <AutoDecisionRequestsTable />}</div>
    </div>
  );
}

interface AppReviewAdminActionsProps {
  endApplicationReview: () => void;
}

function AppReviewAdminActions({
  endApplicationReview,
}: AppReviewAdminActionsProps) {
  const [open, setOpen] = useState(false);

  return (
    <>
      <DialogTrigger onOpenChange={setOpen}>
        <Button variant="secondary" className="w-full">
          <TablerUserEdit />
          Assign Reviewers
        </Button>
        {open && (
          <Modal>
            <ReviewerAssignmentModal onClose={() => setOpen(false)} />
          </Modal>
        )}
      </DialogTrigger>
      <Button
        onClick={endApplicationReview}
        variant="danger"
        className="w-full"
      >
        End Application Review
      </Button>
    </>
  );
}

function AutoDecisionRequestsTable() {
  const { data, isLoading, isError } = useAutoDecisionRequests();
  const updateRequest = useUpdateAutoDecisionRequest();

  const requests: ApplicationAutoDecisionRequest = data ?? [];
  type RequestRow = ApplicationAutoDecisionRequest[number];

  const columns: ColumnDef<RequestRow>[] = [
    {
      header: "Reviewer",
      maxSize: 120,
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
                <span className="text-gray-600 dark:text-neutral-400">N/A</span>
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
      maxSize: 90,
      cell: ({ row }) => (
        <DialogTrigger>
          <Button variant="secondary" size="sm" className="h-8">
            Open
          </Button>
          <Sheet sheetClassName="w-150">
            <UserApplicationSideDrawer
              applicationId={row.original.application_id}
            />
          </Sheet>
        </DialogTrigger>
      ),
    },
    {
      header: "Decision",
      maxSize: 100,
      cell: ({ row }) =>
        row.original.requested_decision === "auto_accept"
          ? "Auto Accept"
          : "Auto Reject",
    },
    {
      header: "Justification",
      maxSize: 90,
      cell: ({ row }) => (
        <DialogTrigger>
          <Button variant="secondary" size="sm" className="h-8">
            View
          </Button>
          <Popover>
            <div className="p-2">
              <p className="max-w-60 text-wrap">{row.original.justification}</p>
            </div>
          </Popover>
        </DialogTrigger>
      ),
    },
    {
      header: "Created",
      maxSize: 150,
      cell: ({ row }) => new Date(row.original.created_at).toLocaleString(),
    },
    {
      header: "Actions",
      id: "actions",
      maxSize: 100,
      cell: (info) => {
        const row = info.row.original as RequestRow;
        const isResolved = row.approved_by !== null;

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
                updateRequest.mutateAsync({ requestId: row.id, approved: true })
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
  ];

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
    <div className="overflow-x-auto rounded-lg border border-border bg-surface p-4 max-w-230">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold">Auto Decision Requests</h2>
        <span className="text-sm text-text-secondary">
          {requests.length} request{requests.length === 1 ? "" : "s"}
        </span>
      </div>

      {requests.length === 0 ? (
        <p className="text-sm text-text-secondary">
          No auto decision requests found.
        </p>
      ) : (
        <Table
          className="max-h-100 overflow-y-scroll"
          headerClassName="text-text-secondary text-sm"
          data={requests}
          columns={columns}
          showPagination={false}
        />
      )}
    </div>
  );
}

interface UserApplicationSideDrawerProps {
  applicationId: string;
}

function UserApplicationSideDrawer({
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
