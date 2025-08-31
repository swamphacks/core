import { Outlet, notFound, createFileRoute } from "@tanstack/react-router";

// The path here should match the file's location in your routes directory.
export const Route = createFileRoute("/events/$eventId/dashboard/_attendee")({
  beforeLoad: async ({ context }) => {
    if (
      !context.eventRole ||
      !["attendee", "admin"].includes(context.eventRole)
    ) {
      return notFound();
    }
    return {};
  },
  component: () => <Outlet />,
});
