import { PageUnderConstruction } from "@/components/PageUnderConstruction";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/events/$eventId/")({
  component: RouteComponent,
});

function RouteComponent() {
  // const { eventId } = Route.useParams();

  return <PageUnderConstruction />;
}
