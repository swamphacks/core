import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_protected/events/$eventId/feedback/declined",
)({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/events/$eventId/feedback/declined"!</div>;
}
