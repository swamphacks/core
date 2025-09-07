import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_user/settings")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_protected/_user/settings"!</div>;
}
