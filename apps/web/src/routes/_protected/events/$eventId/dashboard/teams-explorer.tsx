import TeamCard from "@/features/Team/components/TeamCard";
import { useEventTeams } from "@/features/Team/hooks/useEventTeams";
import { useMyTeam } from "@/features/Team/hooks/useMyTeam";
import { createFileRoute, notFound } from "@tanstack/react-router";
import { useState } from "react";
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
  const [page, setPage] = useState(0);
  const teams = useEventTeams(eventId, 10, page);
  const myTeam = useMyTeam(eventId);

  if (teams.isLoading || myTeam.isLoading) {
    return (
      <main>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Explore Teams
        </Heading>

        <p>Loading....</p>
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
            <TeamCard team={team} isCurrentTeam={team.id === myTeam.data?.id} />
          );
        })}
      </section>
    </main>
  );
}
