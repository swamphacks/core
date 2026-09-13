import { createFileRoute, redirect } from "@tanstack/react-router";
import { hackathonQueryOptions } from "@/modules/Hackathon/hooks/useHackathon";
import { useSuspenseQuery } from "@tanstack/react-query";
import { PageLoading } from "@/components/PageLoading";
import ApplicationPage from "@/modules/Application/ApplicationPage";

export const Route = createFileRoute("/_protected/application")({
  component: RouteComponent,
  beforeLoad: ({ context }) => {
    if (context.user.role === "attendee") {
      throw redirect({
        to: "/hacker-portal",
      });
    }
  },
  pendingComponent: PageLoading,
  loader: ({ context }) => {
    return Promise.all([
      context.queryClient.ensureQueryData(hackathonQueryOptions()),
    ]);
  },
});

function RouteComponent() {
  const { user } = Route.useRouteContext();
  const hackathon = useSuspenseQuery(hackathonQueryOptions());

  return <ApplicationPage hackathon={hackathon.data} user={user} />;
}
