import { createFileRoute } from "@tanstack/react-router";
import WorkshopAttendancePage from "@/modules/Workshops/WorkshopAttendancePage";

export const Route = createFileRoute("/_protected/_staff/workshop-attendance")({
  component: WorkshopAttendancePage,
});
