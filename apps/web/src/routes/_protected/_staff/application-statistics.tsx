import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_protected/_staff/application-statistics",
)({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_protected/_staff/application-statistics"!</div>;
}
