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

  if (isLoading || !data)
    return (
      <div>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Event Settings
        </Heading>

        <div className="flex flex-col gap-4 max-w-xl">
          {/* Event Title */}
          <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

          {/* Description */}
          <div className="h-20 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

          {/* Website URL */}
          <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

          {/* Maximum Attendees */}
          <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

          {/* Venue */}
          <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

          {/* Venue Map URL */}
          <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

          {/* Event Times */}
          <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

          {/* Application Times */}
          <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

          {/* Decision Release & RSVP (side by side) */}
          <div className="flex gap-4">
            <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-1/2 animate-pulse"></div>
            <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-1/2 animate-pulse"></div>
          </div>

          {/* Published Switch */}
          <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-1/4 animate-pulse"></div>

          {/* Save Button */}
          <div className="h-10 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse mt-2"></div>
        </div>
      </div>
    );

  if (isError && !data)
    return (
      <div>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Event Settings
        </Heading>

        <div className="text-red-500">
          Error: {error ?? "Something went wrong, try refreshing the page."}
        </div>
      </div>
    );

  return (
    <div>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
        Event Settings
      </Heading>

      <EventSettingsForm event={data} />
    </div>
  );
}
