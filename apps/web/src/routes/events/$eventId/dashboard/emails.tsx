import { createFileRoute, notFound } from "@tanstack/react-router";

export const Route = createFileRoute("/events/$eventId/dashboard/emails")({
  component: RouteComponent,
  beforeLoad: async ({ context }) => {
    // Only allow admin and staff to access this route
    if (!context.eventRole || !["admin", "staff"].includes(context.eventRole)) {
      return notFound();
    }

    return {};
  },
});

function RouteComponent() {
  return <div>Under construction...</div>;
}
