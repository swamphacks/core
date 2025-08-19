import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/events/$eventId/waitlist/info")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/events/$eventId/waitlist/info"!</div>;
}
