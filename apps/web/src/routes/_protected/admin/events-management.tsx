import { createFileRoute } from "@tanstack/react-router";
import { Heading, DialogTrigger, Text } from "react-aria-components";

import { Button } from "@/components/ui/Button";
import { AddEventModal } from "@/features/PlatformAdmin/EventManager/components/AddEventModal";
import { useAdminEvents } from "@/features/PlatformAdmin/EventManager/hooks/useAdminEvents";
import { EventDetailsCard } from "@/features/PlatformAdmin/EventManager/components/EventDetailsCard";

export const Route = createFileRoute("/_protected/admin/events-management")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data, isLoading, isError, error } = useAdminEvents();

  if (isLoading) {
    return (
      <div className="mx-auto p-6">
        <header className="flex justify-between items-center mb-6">
          <Heading className="text-text-main text-2xl">Event Manager</Heading>
          <Button variant="primary" isDisabled aria-details="disabled">
            Loading...
          </Button>
        </header>
        <div className="flex flex-wrap gap-6">
          {[...Array(3)].map((_, i) => (
            <div
              key={i}
              className="w-82 h-52 bg-neutral-200 dark:bg-neutral-800 animate-pulse rounded"
            />
          ))}
        </div>
      </div>
    );
  }

  if (isError || error) {
    return (
      <div className="mx-auto p-6 text-red-600">
        <header className="flex justify-between items-center mb-6">
          <Heading className="text-text-main text-2xl">Event Manager</Heading>
          <Button variant="primary" isDisabled aria-details="disabled">
            Create New Event
          </Button>
        </header>
        <Text>Error loading events: {error?.message || "Unknown error"}</Text>
      </div>
    );
  }

  return (
    <div>
      <div className="mx-auto">
        <header className="flex justify-between items-center mb-6">
          <Heading className="text-text-main text-2xl">Event Manager</Heading>

          <DialogTrigger>
            <Button variant="primary">Create New Event</Button>

            <AddEventModal />
          </DialogTrigger>
        </header>
      </div>

      <div className="flex flex-row flex-wrap gap-6">
        {data?.length === 0 && <Text>Wow, much empty...</Text>}
        {data &&
          data.map((val) => <EventDetailsCard key={val.id} event={val} />)}
      </div>
    </div>
  );
}
