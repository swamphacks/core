import { PageUnderConstruction } from "@/components/PageUnderConstruction";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/(protected)/admin/logs")({
  component: RouteComponent,
});

function RouteComponent() {
  return <PageUnderConstruction />;
}
