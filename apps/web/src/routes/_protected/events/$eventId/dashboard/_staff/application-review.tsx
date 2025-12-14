import { Button } from "@/components/ui/Button";
import { Tooltip } from "@/components/ui/Tooltip";
import ResetReviewWarningModal from "@/features/ApplicationReview/components/ResetReviewModal";
import ApplicationReviewPage from "@/features/ApplicationReview/components/Review/ApplicationReviewPage";
import ReviewNotStarted from "@/features/ApplicationReview/components/ReviewNotStarted/ReviewNotStarted";
import { useAppReviewAdminActions } from "@/features/ApplicationReview/hooks/useAppReviewAdminActions";
import { useEvent } from "@/features/Event/hooks/useEvent";
import { createFileRoute } from "@tanstack/react-router";
import { DialogTrigger, Heading } from "react-aria-components";
import { toast } from "react-toastify";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/application-review",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId } = Route.useParams();
  const { user, eventRole } = Route.useRouteContext();
  const { reset } = useAppReviewAdminActions(eventId);
  const event = useEvent(eventId);

  const onReset = async () => {
    await reset.mutateAsync(undefined, {
      onSuccess: () => {
        toast.success("Successfully reset application reviews.");
      },
      onError: () => {
        toast.error("Failed to reset application reviews. Please try again.");
      },
    });
  };

  const loading = !user || event.isLoading;

  if (loading) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Application Review
        </Heading>
        <div className="h-84 w-full max-w-xl bg-neutral-300 dark:bg-neutral-800 rounded animate-pulse" />
      </main>
    );
  }

  if (!event.data || event.isError) {
    return <div>Event not found</div>;
  }

  return (
    <main className="h-full">
      <div className="w-full flex flex-row justify-between items-center">
        <Heading className="text-2xl lg:text-3xl font-semibold mb-4">
          Application Review
        </Heading>

        {eventRole === "admin" && event.data.application_review_started && (
          <DialogTrigger>
            <Tooltip
              tooltipProps={{
                label: "[ADMIN] Reset all application reviews for this event",
              }}
              triggerProps={{
                delay: 100,
              }}
            >
              <Button size="sm" className="m-0 h-fit w-fit">
                Reset Reviews
              </Button>
            </Tooltip>

            <ResetReviewWarningModal
              isPending={reset.isPending}
              onReset={onReset}
            />
          </DialogTrigger>
        )}
      </div>

      {event.data.application_review_started ? (
        <div>
          <ApplicationReviewPage eventId={event.data.id} />
        </div>
      ) : (
        <>
          {eventRole === "staff" && (
            <p className="text-text-secondary">
              Application review has not started yet. Come back later!
            </p>
          )}
          {eventRole === "admin" && <ReviewNotStarted event={event.data} />}
        </>
      )}
    </main>
  );
}
