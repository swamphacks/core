import { useApplicationStatistics } from "@/features/Application/hooks/useApplicationStatistics";
import ReviewNotStarted from "@/features/ApplicationReview/components/ReviewNotStarted/ReviewNotStarted";
import { useEvent } from "@/features/Event/hooks/useEvent";
import { useEventStaffUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";
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
  const stats = useApplicationStatistics(eventId);
  const staff = useEventStaffUsers(eventId);

  const loading =
    !user ||
    event.isLoading ||
    stats.isLoading ||
    staff.isLoading ||
    !event.data ||
    !stats.data ||
    !staff.data;

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

  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-4">
        Application Review
      </Heading>

      {eventRole === "staff" && (
        <p className="text-text-secondary">
          Application review has not started yet. Come back later!
        </p>
      )}

      {eventRole === "admin" && (
        <ReviewNotStarted
          event={event.data}
          stats={stats.data}
          staff={staff.data}
        />
      )}
    </main>
  );
}
