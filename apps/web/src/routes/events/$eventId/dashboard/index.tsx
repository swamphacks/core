import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/events/$eventId/dashboard/")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/events/$eventId/dashboard/"!</div>;
}
