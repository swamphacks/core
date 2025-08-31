import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/events/$eventId/dashboard/_staff/team-management",
)({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/events/$eventId/dashboard/_staff/team-management"!</div>;
}
