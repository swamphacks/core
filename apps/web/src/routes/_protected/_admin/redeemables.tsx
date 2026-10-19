import RedeemablesPage from "@/modules/Redeemables/RedeemablesPage";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_admin/redeemables")({
  component: RouteComponent,
});

function RouteComponent() {
  return <RedeemablesPage />;
}
