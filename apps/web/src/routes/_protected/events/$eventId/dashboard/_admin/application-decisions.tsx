import { Button } from "@/components/ui/Button";
import { useApplicationReviewComplete } from "@/features/ApplicationDecisions/hooks/useApplicationReviewStatus";
import { useDecisionRuns } from "@/features/ApplicationDecisions/hooks/useDecisionRuns";
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
  const appReviewComplete = useApplicationReviewComplete(eventId);
  const runs = useDecisionRuns(eventId);

  const loading =
    !user ||
    appReviewComplete.isLoading ||
    runs.isLoading ||
    !appReviewComplete.data ||
    !runs.data;

  if (loading) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Decision Runs
        </Heading>
        <div className="h-84 w-full max-w-xl bg-neutral-300 dark:bg-neutral-800 rounded animate-pulse" />
      </main>
    );
  }

  if (
    appReviewComplete.isError ||
    appReviewComplete.error ||
    runs.isError ||
    runs.error
  ) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Decision Runs
        </Heading>

        <p className="text-red-500">
          Something wen&apos;t wrong... please refresh the page.
        </p>
      </main>
    );
  }

  return (
    <main>
      <div className="flex flex-row items-center justify-between">
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Decision Runs
        </Heading>

        <Button>New Run</Button>
      </div>
      <div className="w-full md:max-w-1/3">
        {appReviewComplete.data ? (
          runs.data.length !== 0 ? (
            // Render your runs component or list here
            // <RunsList runs={runs.data} />
            <></>
          ) : (
            <div>
              <p>No decision runs yet. Start a new run!</p>
            </div>
          )
        ) : (
          <div>
            <p>Not all reviews have been completed, come back later...</p>
          </div>
        )}
      </div>
    </main>
  );
}
