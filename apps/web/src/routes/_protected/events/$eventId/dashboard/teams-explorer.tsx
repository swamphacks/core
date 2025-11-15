import TeamCard from "@/features/Team/components/TeamCard";
import { useEventTeams } from "@/features/Team/hooks/useEventTeams";
import { useMyPendingJoinRequests } from "@/features/Team/hooks/useMyPendingJoinRequests";
import { useMyTeam } from "@/features/Team/hooks/useMyTeam";
import { createFileRoute, notFound } from "@tanstack/react-router";
import { Heading } from "react-aria-components";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/teams-explorer",
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
  const teams = useEventTeams(eventId, 9999, 0);
  const myTeam = useMyTeam(eventId);
  const pendingJoinRequests = useMyPendingJoinRequests(eventId);

  if (teams.isLoading || myTeam.isLoading) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Explore Teams
        </Heading>

        <p>Loading teams....</p>
      </main>
    );
  }

  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
        Explore Teams
      </Heading>

      <section className="w-full flex flex-row flex-wrap gap-4">
        {teams.data?.map((team) => {
          return (
            <TeamCard
              key={team.id}
              team={team}
              isCurrentTeam={team.id === myTeam.data?.id}
              alreadyRequested={
                !!pendingJoinRequests.data?.find(
                  (request) => request.team_id === team.id,
                )
              }
            />
          );
        })}

        {!teams.data && <p>No teams have been created for this event yet.</p>}
      </section>
    </main>
  );
}
