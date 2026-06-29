import { createFileRoute, notFound, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_staff")({
  component: () => <Outlet />,
  beforeLoad: ({ context }) => {
    if (!["admin", "staff"].includes(context.user.role)) {
      return notFound();
    }
  },
});
