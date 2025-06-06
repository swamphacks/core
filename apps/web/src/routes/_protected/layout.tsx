import { createFileRoute, redirect, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected")({
  beforeLoad: async ({ context }) => {
    const user = await context.auth.promise;

    // const user = {
    //   role: "admin",
    // };
    if (!user) {
      throw redirect({
        to: "/",
        search: {
          redirect: location.href,
        },
      });
    }

    // If authenticated successfully, returns the user object to make it available to RouteComponent's context
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
