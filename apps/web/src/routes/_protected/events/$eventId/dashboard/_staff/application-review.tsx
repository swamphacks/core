import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/application-review",
)({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div>Hello "/events/$eventId/dashboard/_staff/application-review"!</div>
  );
}
