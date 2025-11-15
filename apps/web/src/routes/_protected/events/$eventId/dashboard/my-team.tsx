import MyTeamCard from "@/features/Team/components/MyTeamCard";
import NoTeamCard from "@/features/Team/components/NoTeamCard";
import TeamJoinRequestSection from "@/features/Team/components/TeamJoinRequestSection";
import { useMyTeam } from "@/features/Team/hooks/useMyTeam";
import { createFileRoute, notFound } from "@tanstack/react-router";
import { Heading } from "react-aria-components";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/my-team",
)({
  component: RouteComponent,
  beforeLoad: async ({ context }) => {
    if (
      !context.eventRole ||
      !["applicant", "attendee"].includes(context.eventRole)
    ) {
      return notFound();
    }
    return {};
  },
});

function RouteComponent() {
  const eventId = Route.useParams().eventId;
  const { user } = Route.useRouteContext();
  const team = useMyTeam(eventId);

  const isOwner = team.data?.owner_id === user?.userId;

  if (team.isLoading) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          My Team
        </Heading>
        <div className="flex flex-col gap-4 max-w-xl">
          <div className="h-84 w-full md:w-120 bg-neutral-300 dark:bg-neutral-800 rounded animate-pulse"></div>
        </div>
      </main>
    );
  }

  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-4">
        My Team
      </Heading>

      <div className="flex flex-col gap-8">
        {team.data ? (
          <MyTeamCard eventId={eventId} team={team.data} />
        ) : (
          <NoTeamCard eventId={eventId} />
        )}
        {/* <TeamInvitationSection /> */}

        {isOwner && team.data && (
          <TeamJoinRequestSection teamId={team.data.id} />
        )}
      </div>
    </main>
  );
}
