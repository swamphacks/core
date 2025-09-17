import EventSettingsForm from "@/features/Event/components/EventSettingsForm";
import { useEvent } from "@/features/Event/hooks/useEvent";
import { createFileRoute } from "@tanstack/react-router";
import { Heading } from "react-aria-components";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_admin/event-settings",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const eventId = Route.useParams().eventId;
  const { data, isLoading, isError, error } = useEvent(eventId);

  if (isLoading || !data) return <div>Loading...</div>;

  if (isError && !data)
    return <div>Error: {String(error || "Error has occurred")}</div>;

  return (
    <div>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
        Event Settings
      </Heading>

      <EventSettingsForm event={data} />
    </div>
  );
}
