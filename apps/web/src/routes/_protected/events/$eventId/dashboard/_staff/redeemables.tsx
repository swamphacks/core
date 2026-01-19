import { PageUnderConstruction } from "@/components/PageUnderConstruction";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/redeemables",
)({
  component: RouteComponent,
});

// TODO: Implement redeemables dashboard page

function RouteComponent() {
  // TODO: Create Redeemable Card Component (with edit and delete functionality)
  // TODO: Create manager for loading redeemables
  return <PageUnderConstruction />;
}
