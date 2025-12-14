import { PageUnderConstruction } from "@/components/PageUnderConstruction";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/events/$eventId/rejected")({
  component: RouteComponent,
});

function RouteComponent() {
  return <PageUnderConstruction />;
}
