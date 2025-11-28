import ReviewNotStarted from "@/features/ApplicationReview/components/ReviewNotStarted/ReviewNotStarted";
import { useEvent } from "@/features/Event/hooks/useEvent";
import { createFileRoute } from "@tanstack/react-router";
import { Heading } from "react-aria-components";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/application-review",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId } = Route.useParams();
  const { user, eventRole } = Route.useRouteContext();

  const event = useEvent(eventId);

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
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-4">
        Application Review
      </Heading>

      {event.data.application_review_started ? (
        <>
          <p>Something</p>
        </>
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
