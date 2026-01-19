import { createFileRoute } from "@tanstack/react-router";
import { Heading } from "react-aria-components";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_attendee/schedule",
)({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <main>
      <Heading className="text-1xl lg:text-2xl font-semibold mb-6">
        Event Schedule
      </Heading>

      <p className="text-text-secondary text-md">To be announced...</p>
    </main>
  );
}
