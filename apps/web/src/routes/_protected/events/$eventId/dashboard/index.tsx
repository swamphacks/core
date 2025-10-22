import AttendeeOverview from "@/features/EventOverview/components/AttendeeOverview";
import StaffOverview from "@/features/EventOverview/components/StaffOverview";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/events/$eventId/dashboard/")({
  component: RouteComponent,
  beforeLoad: ({ context, params }) => {
    if (context.eventRole === "applicant") {
      throw redirect({
        to: `/events/$eventId/dashboard/application-status`,
        params: {
          eventId: params.eventId,
        },
      });
    }
  },
});

function RouteComponent() {
  const { eventRole } = Route.useRouteContext();
  const { eventId } = Route.useParams()

  // Both staff and admin are technically "staff" for dashboard purposes
  const isStaff = eventRole === "staff" || eventRole === "admin";

  if (isStaff) {
    return <StaffOverview eventId={eventId} />;
  }

  if (eventRole === "attendee") {
    return <AttendeeOverview />;
  }

  // Should never reach here due to the redirect in beforeLoad
  return (
    <div>
      <p>Are you sure you&apos;re supposed to be here?</p>
    </div>
  );
}
