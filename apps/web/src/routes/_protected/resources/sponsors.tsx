import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/resources/sponsors")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_protected/resources/sponsors"!</div>;
}
