import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/attendee-directory",
)({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div>Hello "/events/$eventId/dashboard/_staff/attendee-directory"!</div>
  );
}
