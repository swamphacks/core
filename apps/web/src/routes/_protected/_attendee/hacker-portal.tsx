import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_attendee/hacker-portal")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>TODO</div>;
}
