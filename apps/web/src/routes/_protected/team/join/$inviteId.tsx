import JoinTeam from "@/modules/Team/JoinTeam";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/team/join/$inviteId")({
  component: RouteComponent,
});

function RouteComponent() {
  const { inviteId } = Route.useParams();
  return <JoinTeam inviteId={inviteId} />;
}
