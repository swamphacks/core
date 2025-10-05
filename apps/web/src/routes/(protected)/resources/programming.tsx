import { PageUnderConstruction } from "@/components/PageUnderConstruction";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/(protected)/resources/programming")({
  component: RouteComponent,
});

function RouteComponent() {
  return <PageUnderConstruction />;
}
