import { PageLoading } from "@/components/PageLoading";
import ApplicationReviewPage from "@/modules/Application/ApplicationReview/ApplicationReviewPage";
import { staffHackathonQueryOptions } from "@/modules/Hackathon/hooks/useHackathon";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_staff/application-review/")({
  component: RouteComponent,
  pendingComponent: PageLoading,
  loader: ({ context }) => {
    return context.queryClient.ensureQueryData(staffHackathonQueryOptions());
  },
});

function RouteComponent() {
  const { user } = Route.useRouteContext();
  const hackathon = useSuspenseQuery(staffHackathonQueryOptions());

  return <ApplicationReviewPage hackathon={hackathon.data} user={user} />;
}
