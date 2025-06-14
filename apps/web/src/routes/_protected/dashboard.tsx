import { Button } from "@/components/ui/Button";
import { auth } from "@/lib/authClient";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/dashboard")({
  component: RouteComponent,
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
        search: { redirectTo: "/dashboard" },
      });
    }

    return { user };
  },
});

function RouteComponent() {
  const logout = async () => {
    try {
      await auth.logOut();
    } catch (error) {
      console.error("Error during logout:", error);
    }
  };

  return (
    <div>
      {/* {data.user.role === "admin" ? <Dashboard /> : <HackerDashboard />} */}
      <p>Dashboard</p>
      <Button onClick={logout} color="danger">
        Logout
      </Button>
    </div>
  );
}
