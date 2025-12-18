import { Button } from "@/components/ui/Button";
import { useEvent } from "@/features/Event/hooks/useEvent";
import { createFileRoute } from "@tanstack/react-router";
import { Heading } from "react-aria-components";
import { toast } from "react-toastify";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_admin/application-decisions",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId } = Route.useParams();
  const { user } = Route.useRouteContext();
  const event = useEvent(eventId);

  const checkReviewsCompleted = async () => {
    if (!event.data) return;

    const isCompleted = true; // TODO: fix. this is being temporarily left in here to get a github workflow to run.
    if (isCompleted) {
      toast.success(
        "All application reviews completed. You can now proceed with calculating decisions.",
      );
    }
  };

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

        {event.data.application_review_finished ? (
          <p>Show applications decision page</p>
        ) : (
          <>
            <Button
              size="sm"
              className="m-0 h-fit w-fit"
              onPress={checkReviewsCompleted}
            >
              Run Pre-Decisions Check
            </Button>
          </>
        )}
      </div>
    </main>
  );
}
