import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

// This layout component performs authentication checks before the user can access protected pages
export const Route = createFileRoute("/_protected")({
  beforeLoad: async ({ context }) => {
    const { user, error } = await context.userQuery.promise;

    // Unauthenticated, return to login page
    if (!user && !error) {
      console.log("User is not authenticated, redirecting to login.");
      throw redirect({
        to: "/",
        search: { redirect: location.pathname },
      });
    }

    if (error) {
      // TODO: Display a friendly error to the user?
      console.error("Auth error in beforeLoad in layout.tsx:", error);
      console.log(
        "authentication error occurred while accessing protected page, redirecting to login.",
      );
      throw redirect({
        to: "/",
      });
    }

    // Return user data for use in the route loader or component
    return { user };
  },
  pendingMs: 1000,
  pendingComponent: () => <p>Loading...</p>,
  component: RouteComponent,
});

function RouteComponent() {
  return <Outlet />;
}
