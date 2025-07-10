import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/community")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello this is the community tab!</div>;
}
