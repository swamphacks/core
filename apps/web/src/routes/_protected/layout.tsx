import { createFileRoute, redirect, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected")({
  beforeLoad: async ({ context }) => {
    try {
      const user = await context.user.promise;

      console.log(user);

      if (!user) {
        throw new Error();
      }

      // If authenticated successfully, returns the user object to make it available to RouteComponent's context
      return {
        user,
      };
    } catch {
      // TODO: is there a better way to handle errors?
      throw redirect({
        to: "/",
        search: {
          redirectTo: location.href,
        },
      });
    }
  },
  pendingMs: 1000,
  // TODO: update loading component to an animated Icon
  pendingComponent: () => <p>Loading...</p>,
  component: RouteComponent,
});

function RouteComponent() {
  return <Outlet />;
}
