import AttendeeAppShell from "@/features/Dashboard/components/AttendeeAppShell";
import ApplicantAppShell from "@/features/Dashboard/components/ApplicantAppShell";
import StaffAppShell from "@/features/Dashboard/components/StaffAppShell";
import { getUserEventRole } from "@/features/Event/api/getUserEventRole";
import NotFoundPage from "@/features/NotFound/NotFoundPage";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { fetchEvent, getEventQueryKey } from "@/features/Event/hooks/useEvent";

export const Route = createFileRoute("/_protected/events/$eventId/dashboard")({
  component: RouteComponent,
  beforeLoad: async ({ context, params }) => {
    // Get event roles, permissions, etc.

    const data = await context.queryClient.fetchQuery({
      queryKey: ["eventRole", params.eventId],
      queryFn: () => getUserEventRole(params.eventId),
    });

    // Parse assigned at string to Date object
    let parsedAssignedAt: Date | null = null;
    if (data?.assigned_at) {
      parsedAssignedAt = new Date(data.assigned_at);
    }

    return {
      eventRole: data?.role,
      roleAssignedAt: parsedAssignedAt ?? null,
    };
  },
  loader: async ({ context, params }) => {
    const data = await context.queryClient.fetchQuery({
      queryKey: getEventQueryKey(params.eventId),
      queryFn: () => fetchEvent(params.eventId),
    });

    return data;
  },
});

function RouteComponent() {
  const { eventRole } = Route.useRouteContext();
  const { eventId } = Route.useParams();
  const { name } = Route.useLoaderData();

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
        <AttendeeAppShell eventId={eventId} eventName={name}>
          <Outlet />
        </AttendeeAppShell>
      );
    case "applicant":
      return (
        <ApplicantAppShell eventId={eventId} eventName={name}>
          <Outlet />
        </ApplicantAppShell>
      );
    default:
      return <div>Unknown role: {eventRole}</div>;
  }
}
