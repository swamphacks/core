import { createFileRoute, notFound, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_attendee")({
  component: () => <Outlet />,
  beforeLoad: ({ context }) => {
    if (context.user.role !== "attendee") {
      return notFound();
    }
  },
});
