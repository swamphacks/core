import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/events/$eventId/dashboard/_attendee/schedule",
)({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/events/$eventId/dashboard/schedule"!</div>;
}
