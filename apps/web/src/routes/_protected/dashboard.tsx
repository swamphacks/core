import { createFileRoute, useLoaderData } from "@tanstack/react-router";
// import { Test as Admin } from "@/features/admin/Dashboard/components/Dashboard";
import Dashboard from "@/features/admin/Dashboard";
import HackerDashboard from "@/features/hacker/Dashboard";

export const Route = createFileRoute("/_protected/dashboard")({
  component: RouteComponent,
  loader: ({ context }) => {
    // TODO: fetch approriate data based on user role
    return {
      user: context.user,
    };
  },
});

function RouteComponent() {
  const data = useLoaderData({ from: "/_protected/dashboard" });

  return (
    <div>
      {data.user.role === "admin" ? <Dashboard /> : <HackerDashboard />}
    </div>
  );
}
