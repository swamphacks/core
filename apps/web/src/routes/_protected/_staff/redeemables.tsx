import RedeemablesPage from "@/modules/Redeemables/RedeemablesPage";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_staff/redeemables")({
  component: RouteComponent,
});

function RouteComponent() {
  return <RedeemablesPage />;
}
