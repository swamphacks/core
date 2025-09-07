import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_attendee/my-team",
)({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/events/$eventId/dashboard/my-team"!</div>;
}
