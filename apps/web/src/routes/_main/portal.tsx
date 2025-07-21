import { Button } from "@/components/ui/Button";
import { auth } from "@/lib/authClient";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_main/portal")({
  component: RouteComponent,
});

function RouteComponent() {
  const logout = async () => {
    try {
      await auth.logOut();
      window.location.href = "/";
    } catch (error) {
      console.error("Error during logout:", error);
    }
  };

  return (
    <div>
      <p>Event Portal</p>
      <Button onClick={logout} color="danger">
        Logout
      </Button>
    </div>
  );
}
