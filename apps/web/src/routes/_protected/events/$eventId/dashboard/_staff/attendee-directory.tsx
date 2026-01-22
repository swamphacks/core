import { PageLoading } from "@/components/PageLoading";
import AttendeeTable from "@/features/CheckIn/components/ManualCheckInTable";
import { useEventUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventUsers";
import { createFileRoute } from "@tanstack/react-router";
import { Heading } from "react-aria-components";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/attendee-directory",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const eventId = Route.useParams().eventId;
  const { data, isLoading, isError, error } = useEventUsers(eventId);

  if (isLoading) {
    return (
      <div className="flex flex-col gap-6 h-full">
        <Heading className="text-2xl lg:text-3xl font-semibold">
          Attendee Management
        </Heading>
        <div className="flex flex-1 justify-center items-center">
          <PageLoading />
        </div>
      </div>
    );
  }

  if (isError || error) {
    return (
      <div className="flex flex-col gap-6">
        <Heading className="text-2xl lg:text-3xl font-semibold">
          Attendee Management
        </Heading>

        <div>Error loading event users.</div>
      </div>
    );
  }

  // Filter to show only users with "attendee" status
  // const attendeeData = data?.filter((user) => user.event_role === "attendee") ?? [];
  const attendeeData = data;
  return (
    <div className="flex flex-col gap-6">
      <Heading className="text-2xl lg:text-3xl font-semibold">
        User Management
      </Heading>

      <AttendeeTable eventId={eventId} data={attendeeData} />
    </div>
  );
}
