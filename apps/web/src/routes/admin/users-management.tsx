import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/admin/users-management")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/admin/users-management"!</div>;
}
