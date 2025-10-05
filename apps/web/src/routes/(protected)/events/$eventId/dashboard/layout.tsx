import AttendeeAppShell from "@/features/Dashboard/components/AttendeeAppShell";
import ApplicantAppShell from "@/features/Dashboard/components/ApplicantAppShell";
import StaffAppShell from "@/features/Dashboard/components/StaffAppShell";
import NotFoundPage from "@/features/NotFound/NotFoundPage";
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/(protected)/events/$eventId/dashboard")({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventRole } = Route.useRouteContext();
  const { eventId } = Route.useParams();

  if (!eventRole) {
    return <NotFoundPage />;
  }

  switch (eventRole) {
    case "admin":
      return (
        <StaffAppShell eventRole={eventRole} eventId={eventId}>
          <Outlet />
        </StaffAppShell>
      );
    case "staff":
      return (
        <StaffAppShell eventRole={eventRole} eventId={eventId}>
          <Outlet />
        </StaffAppShell>
      );
    case "attendee":
      return (
        <AttendeeAppShell eventId={eventId}>
          <Outlet />
        </AttendeeAppShell>
      );
    case "applicant":
      return (
        <ApplicantAppShell eventId={eventId}>
          <Outlet />
        </ApplicantAppShell>
      );
    default:
      return <div>Unknown role: {eventRole}</div>;
  }
}
