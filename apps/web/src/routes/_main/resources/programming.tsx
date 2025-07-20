import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_main/resources/programming")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_protected/resources/programming"!</div>;
}
