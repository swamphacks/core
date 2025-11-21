import { Button } from "@/components/ui/Button";
import AddStaffModal from "@/features/EventAdmin/components/AddStaffModal";
import StaffTable from "@/features/EventAdmin/components/StaffTable";
import { useEventStaffUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";
import { createFileRoute } from "@tanstack/react-router";
import { DialogTrigger, Heading } from "react-aria-components";
import { PageLoading } from "@/components/PageLoading";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_admin/staff-management",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const eventId = Route.useParams().eventId;
  const { data, isLoading, isError } = useEventStaffUsers(eventId);
  const { eventRole } = Route.useRouteContext();

  if (isLoading) {
    return (
      <div className="flex flex-col gap-6 h-full">
        <Heading className="text-2xl lg:text-3xl font-semibold">
          Staff Management
        </Heading>
        <div className="flex flex-1 justify-center items-center">
          <PageLoading />
        </div>
      </div>
    );
  }

  if (isError) {
    <div className="flex flex-col gap-6">
      <Heading className="text-2xl lg:text-3xl font-semibold">
        Staff Management
      </Heading>

      <div>Error loading staff users.</div>
    </div>;
  }

  return (
    <div className="flex flex-col gap-6">
      <Heading className="text-2xl lg:text-3xl font-semibold">
        Staff Management
      </Heading>

      <div className="flex flex-row w-full justify-end">
        {eventRole === "admin" && (
          <DialogTrigger>
            <Button variant="primary">Add Staff</Button>

            <AddStaffModal eventId={eventId} />
          </DialogTrigger>
        )}
      </div>

      <StaffTable eventId={eventId} data={data} />
    </div>
  );
}
