import { createFileRoute } from "@tanstack/react-router";
import StaffCheckInDashboard from "@/modules/CheckIn/StaffCheckInDashboard";

export const Route = createFileRoute("/_protected/_staff/check-in")({
  component: StaffCheckInDashboard,
});
