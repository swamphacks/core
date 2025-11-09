import { createFileRoute } from "@tanstack/react-router";
import { useEventUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventUsers";
import UserTable from "@/features/EventAdmin/components/UserTable";
import { Heading } from "react-aria-components";
import { PageLoading } from "@/components/PageLoading";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_admin/user-management",
)({
  component: RouteComponent,
  validateSearch: (search) => {
    return {
      tableState: search.tableState ? String(search.tableState) : undefined,
    };
  },
});

function RouteComponent() {
  const eventId = Route.useParams().eventId;
  const { data, isLoading, isError } = useEventUsers(eventId);

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

      <UserTable eventId={eventId} data={data} />
    </div>
  );
}
