import { createFileRoute, redirect } from "@tanstack/react-router";
import { hackathonQueryOptions } from "@/modules/Hackathon/hooks/useHackathon";
import { useSuspenseQuery } from "@tanstack/react-query";
import TablerAlertCircle from "~icons/tabler/alert-circle";
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

  const now = new Date();
  const applicationOpen = new Date(hackathon.data.applicationOpen);
  const applicationClose = new Date(hackathon.data.applicationClose);
  const isApplicationOpen = now >= applicationOpen && now <= applicationClose;

  if (!isApplicationOpen) {
    return (
      <div className="max-w-xs mx-auto h-full flex flex-col justify-center items-center gap-8 text-text-secondary">
        <div className="flex flex-row items-center justify-center gap-2">
          <TablerAlertCircle />
          <p>Applications are currently closed.</p>
        </div>
      </div>
    );
  }

  return <ApplicationPage hackathon={hackathon} user={user} />;
}
