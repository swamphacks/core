import ApplicationStatistics from "@/features/Application/components/ApplicationStatistics";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/application-statistics",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId } = Route.useParams();

  return (
    <div>
      <ApplicationStatistics eventId={eventId} />
    </div>
  );
}
