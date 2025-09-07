import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_user/resources/programming")(
  {
    component: RouteComponent,
  },
);

function RouteComponent() {
  return <div>Hello "/_protected/resources/programming"!</div>;
}
