import { createFileRoute, notFound, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_admin")({
  component: () => <Outlet />,
  beforeLoad: ({ context }) => {
    if (context.user.role !== "admin") {
      return notFound();
    }
  },
});
