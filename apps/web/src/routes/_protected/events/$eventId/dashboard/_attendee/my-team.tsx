import { PageUnderConstruction } from "@/components/PageUnderConstruction";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_attendee/my-team",
)({
  component: RouteComponent,
});

function RouteComponent() {
  return <PageUnderConstruction />;
}
