import { createFileRoute } from "@tanstack/react-router";

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
  // const data = useLoaderData({ from: "/_protected/dashboard" });

  return (
    <div>
      {/* {data.user.role === "admin" ? <Dashboard /> : <HackerDashboard />} */}
      <p>Dashboard</p>
    </div>
  );
}
