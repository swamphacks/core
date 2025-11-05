import { createFileRoute } from "@tanstack/react-router";
import { useEventUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventUsers";
import StaffTable from "@/features/EventAdmin/components/StaffTable";
import { DialogTrigger, Heading } from "react-aria-components";
import { Button } from "@/components/ui/Button";
import AddStaffModal from "@/features/EventAdmin/components/AddStaffModal";
import { PageLoading } from "@/components/PageLoading";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_admin/user-management",
)({
  component: RouteComponent,
});

// TODO:
/*
- Create a component that can accept a column definition,
- Takes in Data
- Takes in filter data
- Takes in "Actions Components?"
- Create side popup
- Other random things
*/
function RouteComponent() {
  // Make request to get user data
  const eventId = Route.useParams().eventId;
  const { data, isLoading, isError } = useEventUsers(eventId);
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
        User Management
      </Heading>

      <div>Error loading event users.</div>
    </div>;
  }

  return (
    <div className="flex flex-col gap-6">
      <Heading className="text-2xl lg:text-3xl font-semibold">
        User Management
      </Heading>

      <div className="flex flex-row w-full justify-end">
        {eventRole === "admin" && (
          <DialogTrigger>
            <Button variant="primary">Actions</Button>

            <AddStaffModal eventId={eventId} />
          </DialogTrigger>
        )}
      </div>

      <StaffTable eventId={eventId} data={data} />
    </div>
  );
}
