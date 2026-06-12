import ApplicationSearchPage from "@/modules/Application/ApplicationSearch/ApplicationSearchPage";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_admin/application-search")({
  component: RouteComponent,
});

function RouteComponent() {
  return <ApplicationSearchPage />;
}
