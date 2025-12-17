import ApplicationDecisionsPage from "@/features/ApplicationDecisions/components/Decisions/ApplicationDecisionsPage";
import { useCheckAppReviewStatus } from "@/features/ApplicationDecisions/hooks/useCheckAppReviewStatus";
import { useEvent } from "@/features/Event/hooks/useEvent";
import { createFileRoute } from "@tanstack/react-router";
import { Heading } from "react-aria-components";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_admin/application-decisions",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId } = Route.useParams();
  const { user } = Route.useRouteContext();
  const { checkAppReviewStatus } = useCheckAppReviewStatus(eventId);
  const event = useEvent(eventId);

  const checkReviewsCompleted = async () => {
    if (!event.data) return;

    const isCompleted = await checkAppReviewStatus(event.data);
    console.log(isCompleted);
  };
  checkReviewsCompleted();

  const loading = !user || event.isLoading;

  if (loading) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Application Decisions
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
          Application Decisions
        </Heading>
      </div>

      {event.data.application_review_finished ? (
        <ApplicationDecisionsPage eventId={eventId} />
      ) : (
        <p>Please come back after finishing application reviews!</p>
      )}
    </main>
  );
}
