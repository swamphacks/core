import { createFileRoute, notFound } from "@tanstack/react-router";
import { Heading } from "react-aria-components";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/teams-explorer",
)({
  component: RouteComponent,
  beforeLoad: async ({ context }) => {
    if (
      !context.eventRole ||
      !["applicant", "attendee"].includes(context.eventRole)
    ) {
      return notFound();
    }
    return {};
  },
});

function RouteComponent() {
  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
        Explore Teams
      </Heading>

      <h2>Hey, we are working on this page still....</h2>
    </main>
  );
}
