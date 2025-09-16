import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/events/$eventId/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId } = Route.useParams();

  return <div>This is the info page for {eventId}</div>;
}
