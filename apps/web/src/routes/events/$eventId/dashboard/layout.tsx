import AttendeeAppShell from "@/features/Dashboard/components/AttendeeAppShell";
import { getUserEventRole } from "@/features/Event/api/getUserEventRole";
import NotFoundPage from "@/features/NotFound/NotFoundPage";
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/events/$eventId/dashboard")({
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
      eventRole: data.role,
      roleAssignedAt: parsedAssignedAt ?? null,
    };
  },
});

function RouteComponent() {
  const { eventRole } = Route.useRouteContext();
  const { eventId } = Route.useParams();

  if (!eventRole) {
    return <NotFoundPage />;
  }

  switch (eventRole) {
    case "staff":
      return <div>Staff Dashboard</div>;
    case "attendee":
      return (
        <AttendeeAppShell eventId={eventId}>
          <Outlet />
        </AttendeeAppShell>
      );
    case "applicant":
      return <div>Applicant Dashboard</div>;
    default:
      return <div>Unknown role: {eventRole}</div>;
  }
}
