import { createFileRoute } from "@tanstack/react-router";
import ApplicationStatus from "@/features/Application/components/ApplicationStatus";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_applicant/application-status",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId } = Route.useParams();

  return <ApplicationStatus eventId={eventId} />;
}
