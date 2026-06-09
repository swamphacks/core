import { Button } from "@/components/ui/Button";
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
    await updateStatus.mutateAsync(false);
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

      {user.role === "admin" && <AutoDecisionRequestsTable />}
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
        {open && <ReviewerAssignmentModal onClose={() => setOpen(false)} />}
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

  const requests: ApplicationAutoDecisionRequest = data ?? [];

  return (
    <div className="overflow-x-auto rounded-lg border border-border bg-surface p-4">
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
        <table className="min-w-full border-collapse text-sm">
          <thead>
            <tr className="border-b border-border text-left text-xs uppercase tracking-wide text-text-secondary">
              {/* <th className="px-3 py-2">Request ID</th> */}
              <th className="px-3 py-2">Application</th>
              <th className="px-3 py-2">Reviewer</th>
              <th className="px-3 py-2">Decision</th>
              <th className="px-3 py-2">Approved</th>
              <th className="px-3 py-2">Updated By</th>
              {/* <th className="px-3 py-2">Justification</th> */}
              <th className="px-3 py-2">Created</th>
              <th className="px-3 py-2">Updated</th>
              <th className="px-3 py-2">Actions</th>
            </tr>
          </thead>
          <tbody>
            {requests.map((request) => {
              const isResolved = request.approved_by !== null;
              const isPending = request.approved_by === null;

              return (
                <tr
                  key={request.id}
                  className="border-b border-border even:bg-surface"
                >
                  <td className="px-3 py-2 align-top">
                    {request.application_id}
                  </td>
                  <td className="px-3 py-2 align-top">
                    {request.reviewer_name}
                  </td>
                  <td className="px-3 py-2 align-top">
                    {request.requested_decision === "auto_accept"
                      ? "Auto Accept"
                      : "Auto Reject"}
                  </td>
                  <td className="px-3 py-2 align-top">
                    {isPending ? "Pending" : request.approved ? "Yes" : "No"}
                  </td>
                  <td className="px-3 py-2 align-top">
                    {request.approver_name || "—"}
                  </td>
                  {/* <td className="px-3 py-2 align-top max-w-80 overflow-hidden text-ellipsis whitespace-nowrap">
                    {request.justification || "—"}
                  </td> */}
                  <td className="px-3 py-2 align-top">
                    {new Date(request.created_at).toLocaleString()}
                  </td>
                  <td className="px-3 py-2 align-top">
                    {new Date(request.updated_at).toLocaleString()}
                  </td>
                  <td className="px-3 py-2 align-top">
                    {isResolved ? (
                      <span className="text-sm text-text-secondary">
                        {request.approved ? "Approved" : "Denied"}
                      </span>
                    ) : (
                      <div className="flex gap-2">
                        <Button
                          className="h-8 px-2"
                          size="sm"
                          onClick={() =>
                            updateRequest.mutateAsync({
                              requestId: request.id,
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
                              requestId: request.id,
                              approved: false,
                            })
                          }
                          isDisabled={updateRequest.isPending}
                        >
                          Deny
                        </Button>
                      </div>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      )}
    </div>
  );
}
