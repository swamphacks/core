import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/events/$eventId/dashboard/_admin/event-settings",
)({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/events/$eventId/dashboard/_admin/event-settings"!</div>;
}
