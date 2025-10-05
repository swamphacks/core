import { PageUnderConstruction } from "@/components/PageUnderConstruction";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/(protected)/(user)/resources/sponsors")({
  component: RouteComponent,
});

function RouteComponent() {
  return <PageUnderConstruction />;
}
