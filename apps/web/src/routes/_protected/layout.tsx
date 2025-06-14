import { auth } from "@/lib/authClient";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/layout")({
  beforeLoad: async () => {
    const { user, error } = await auth.getUser();

    if (error) {
      console.error("Auth error in beforeLoad:", error);
      throw error;
    }

    // If no user, redirect to login
    if (!user) {
      console.log("No user found, redirecting to login.");
      throw redirect({
        to: "/",
        search: { redirectTo: "/layout" },
      });
    }

    // Return user data for use in the route loader or component
    return { user };
  },
  pendingMs: 5000,
  pendingComponent: () => <p>Loading...</p>,
  component: RouteComponent,
});

function RouteComponent() {
  return <Outlet />;
}
