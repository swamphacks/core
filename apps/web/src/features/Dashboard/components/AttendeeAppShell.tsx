import { useLocation } from "@tanstack/react-router";
import { AppShell } from "@/components/AppShell/AppShell";
import { NavLink } from "@/components/AppShell/NavLink";
import TablerLayoutCollage from "~icons/tabler/layout-collage";
import TablerBooks from "~icons/tabler/books";
import TablerSocial from "~icons/tabler/social";
import { type PropsWithChildren } from "react";

interface DashboardAppShellProps {
  eventId: string;
}

export default function AttendeeAppShell({
  eventId,
  children,
}: PropsWithChildren<DashboardAppShellProps>) {
  const pathname = useLocation({ select: (loc) => loc.pathname });

  const dashboardOverviewActive = /^\/events\/[^/]+\/dashboard\/?$/.test(
    pathname,
  );
  const myTeamActive = /^\/events\/[^/]+\/dashboard\/my-team\/?$/.test(
    pathname,
  );
  const exploreTeamsActive =
    /^\/events\/[^/]+\/dashboard\/explore-teams\/?$/.test(pathname);
  const scheduleActive = /^\/events\/[^/]+\/dashboard\/schedule\/?$/.test(
    pathname,
  );

  const teamsTabActive = myTeamActive || exploreTeamsActive;

  return (
    <AppShell>
      <AppShell.Header>
        <div className="w-full px-4 flex flex-row justify-between h-full items-center">
          <h1 className="text-2xl font-bold">Attendee Dashboard</h1>
        </div>
      </AppShell.Header>

      <AppShell.Navbar>
        <NavLink
          label="Overview"
          href={`/events/${eventId}/dashboard`}
          leftSection={<TablerLayoutCollage className="w-5 aspect-square" />}
          active={dashboardOverviewActive}
        />

        <NavLink
          label="Teams"
          leftSection={<TablerBooks className="w-5 aspect-square" />}
          initialExpanded={teamsTabActive}
        >
          <NavLink
            label="My Team"
            href={`/events/${eventId}/dashboard/my-team`}
            leftSection={<TablerLayoutCollage className="w-5 aspect-square" />}
            active={myTeamActive}
          />
          <NavLink
            label="Explore Teams"
            href={`/events/${eventId}/dashboard/explore-teams`}
            leftSection={<TablerLayoutCollage className="w-5 aspect-square" />}
            active={exploreTeamsActive}
          />
        </NavLink>

        <NavLink
          label="Schedule"
          href={`/events/${eventId}/dashboard/schedule`}
          leftSection={<TablerSocial className="w-5 aspect-square" />}
          active={scheduleActive}
        />
      </AppShell.Navbar>

      <AppShell.Main>{children}</AppShell.Main>
    </AppShell>
  );
}
