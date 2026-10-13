import { createFileRoute } from "@tanstack/react-router";
import StaffCheckInDashboard from "@/modules/CheckIn/StaffCheckInDashboard";

export const Route = createFileRoute("/_protected/_staff/check-in")({
  component: CheckInRoute,
});

function CheckInRoute() {
  const { user } = Route.useRouteContext();
  return <StaffCheckInDashboard isAdmin={user.role === "admin"} />;
}
