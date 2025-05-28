import { createFileRoute, redirect, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected")({
  beforeLoad: async ({ context }) => {
    const user = await context.auth.promise;

    if (!user) {
      throw redirect({ to: "/" });
    }

    // Returns the user object to make it available to RouteComponent's context
    return {
      user,
    };
  },
  pendingMs: 5000,
  // TODO: update loading component to an animated Icon
  pendingComponent: () => <p>Loading...</p>,
  component: RouteComponent,
});

function RouteComponent() {
  return <Outlet />;
}
