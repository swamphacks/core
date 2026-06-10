import { PageLoading } from "@/components/PageLoading";
import ApplicationReviewWorkspace from "@/modules/Application/ApplicationReview/Workspace";
import { staffHackathonQueryOptions } from "@/modules/Hackathon/hooks/useHackathon";
import { createFileRoute, HeadContent, redirect } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_protected/_staff/application-review/workspace",
)({
  head: () => ({
    meta: [{ title: "SwampHacks Application Review" }],
  }),
  component: RouteComponent,
  pendingComponent: PageLoading,
  beforeLoad: async ({ context }) => {
    const hackathon = await context.queryClient.ensureQueryData(
      staffHackathonQueryOptions(),
    );
    if (!hackathon.application_review_started) {
      throw redirect({
        to: "/application-review",
      });
    }
  },
});

function RouteComponent() {
  const { user } = Route.useRouteContext();

  return (
    <>
      <HeadContent />
      <ApplicationReviewWorkspace user={user} />
    </>
  );
}
