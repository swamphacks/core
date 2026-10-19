import TablerAlertCircle from "~icons/tabler/alert-circle";
import VisitorWaitlistForm from "@/modules/Application/VisitorWaitlistForm";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { hackathonQueryOptions } from "@/modules/Hackathon/hooks/useHackathon";
import { useSuspenseQuery } from "@tanstack/react-query";
import { PageLoading } from "@/components/PageLoading";
import ApplicationPage from "@/modules/Application/ApplicationPage";
import { useMyApplication } from "@/modules/Application/hooks/useMyApplication";

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
  const application = useMyApplication(user.role !== "visitor");

  const now = new Date();
  const applicationOpen = new Date(hackathon.data.applicationOpen);
  const applicationClose = new Date(hackathon.data.applicationClose);

  let isApplicationOpen;
  if (hackathon.data.acceptEarlyApplications) {
    const earlyApplicationOpen = new Date(hackathon.data.earlyApplicationOpen!);
    const earlyApplicationClose = new Date(
      hackathon.data.earlyApplicationClose!,
    );
    isApplicationOpen =
      (now >= earlyApplicationOpen && now <= earlyApplicationClose) ||
      (now >= applicationOpen && now <= applicationClose);
  } else {
    isApplicationOpen = now >= applicationOpen && now <= applicationClose;
  }

  if (user.role === "visitor" && !isApplicationOpen) {
    return <VisitorWaitlistForm />;
  }

  if (application.isLoading) {
    return <PageLoading />;
  }

  const hasDecision = [
    "accepted",
    "confirmed",
    "rejected",
    "waitlisted",
    "waitlist_confirmed",
    "withdrawn",
  ].includes(application.data?.status ?? "");

  if (!isApplicationOpen && !hasDecision) {
    return (
      <div className="max-w-xs mx-auto h-full flex flex-col justify-center items-center gap-8 text-text-secondary">
        <div className="flex flex-row items-center justify-center gap-2">
          <TablerAlertCircle />
          <p>Applications are currently closed.</p>
        </div>
      </div>
    );
  }

  return <ApplicationPage hackathon={hackathon.data} user={user} />;
}
