import { AvatarStack } from "@/components/ui/AvatarStack";
import { useMyTeam } from "@/features/Team/hooks/useMyTeam";
import { createFileRoute, Link } from "@tanstack/react-router";
import { Heading } from "react-aria-components";
import TablerUsersGroup from "~icons/tabler/users-group";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_applicant/team-formation",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const eventId = Route.useParams().eventId;
  const team = useMyTeam(eventId);

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
        <div className="border border-input-border rounded-md px-4 py-3 w-fit min-w-128 flex flex-col gap-4">
          <h3 className="text-text-main text-2xl">{team.data.name}</h3>

          <div className="flex flex-row text-text-secondary items-center gap-2">
            <TablerUsersGroup className="h-6 w-6" />
            <p className="text-lg">{team.data.members.length} / 4 Members</p>

            <AvatarStack
              avatars={[
                {
                  src: "https://i.pinimg.com/1200x/b3/33/b9/b333b9da782cb6874da5002f185a80ce.jpg",
                },
                {
                  src: "https://i.pinimg.com/736x/d9/fa/c3/d9fac3d552659e888bbecbd20b12f160.jpg",
                },
                { fallback: "Alexander Wang" },
                { fallback: "Anders Swenson" },
                { fallback: "Hello" },
              ]}
              max={2}
              size="sm"
            />
          </div>
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
