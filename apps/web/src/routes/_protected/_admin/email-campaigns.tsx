import EmailCampaignsPage from "@/modules/EmailCampaigns/EmailCampaignsPage";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_admin/email-campaigns")({
  component: RouteComponent,
});

function RouteComponent() {
  return <EmailCampaignsPage />;
}
