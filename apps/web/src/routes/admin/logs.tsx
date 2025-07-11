import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/admin/logs")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/admin/logs"!</div>;
}
