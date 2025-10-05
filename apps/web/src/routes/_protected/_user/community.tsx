import { PageUnderConstruction } from "@/components/PageUnderConstruction";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_user/community")({
  component: RouteComponent,
});

function RouteComponent() {
  return <PageUnderConstruction />;
}
