import { PageUnderConstruction } from "@/components/PageUnderConstruction";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/(protected)/admin/overview")({
  component: RouteComponent,
});

function RouteComponent() {
  return <PageUnderConstruction />;
}
