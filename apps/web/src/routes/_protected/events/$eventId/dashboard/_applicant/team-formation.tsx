import { AvatarStack } from "@/components/ui/AvatarStack";
import { Button } from "@/components/ui/Button";
import { useMyTeam } from "@/features/Team/hooks/useMyTeam";
import { useTeamActions } from "@/features/Team/hooks/useTeamActions";
import { createFileRoute, Link } from "@tanstack/react-router";
import { Heading } from "react-aria-components";
import TablerUsers from "~icons/tabler/users";
export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_applicant/team-formation",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const eventId = Route.useParams().eventId;
  const team = useMyTeam(eventId);
  const { leave } = useTeamActions(eventId);

  const relativeExploreTeamsPath = `/events/${eventId}/dashboard/teams-explorer`;

  if (team.isLoading) {
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
        My Team
      </Heading>
      <div className="flex flex-col gap-4 max-w-xl">
        {/* Event Title */}
        <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

        {/* Description */}
        <div className="h-20 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

        {/* Website URL */}
        <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

        {/* Maximum Attendees */}
        <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

        {/* Venue */}
        <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

        {/* Venue Map URL */}
        <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

        {/* Event Times */}
        <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>

        {/* Application Times */}
        <div className="h-8 bg-neutral-300 dark:bg-neutral-800 rounded w-full animate-pulse"></div>
      </div>
    </main>;
  }

  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-4">
        My Team
      </Heading>

      {team.data ? (
        <div className="border border-input-border rounded-md px-4 py-3 w-full md:w-fit md:min-w-120 flex flex-col gap-4">
          <h3 className="text-text-main text-2xl">{team.data.name}</h3>

          <div className="flex flex-row text-text-secondary items-center gap-2">
            <AvatarStack
              avatars={team.data.members.map((member) => ({
                src: member?.image,
                fallback: member?.name,
              }))}
              max={2}
              size="sm"
            />
            <TablerUsers className="h-6 w-6" />
            <p className="text-lg">{team.data.members.length} / 4 Members</p>
          </div>

          <Button
            variant="danger"
            onPress={() => team.data && leave.mutate(team.data.id)}
          >
            Leave
          </Button>
        </div>
      ) : (
        <div>
          <h2 className="text-text-secondary text-lg">
            It looks like you&apos;re not apart of any teams right now.{" "}
            <Link
              className="text-blue-500 underline underline-offset-4"
              to={relativeExploreTeamsPath}
            >
              Explore Here
            </Link>
          </h2>
        </div>
      )}
    </main>
  );
}
