import { Link, useLocation } from "@tanstack/react-router";
import { AppShell } from "@/components/AppShell/AppShell";
import { NavLink } from "@/components/AppShell/NavLink";
import TablerLayoutCollage from "~icons/tabler/layout-collage";
import TablerUsersGroup from "~icons/tabler/users-group";
import TablerCalendar from "~icons/tabler/calendar";
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
  const scheduleActive = /^\/events\/[^/]+\/dashboard\/schedule\/?$/.test(
    pathname,
  );
  return (
    <AppShell>
      <AppShell.Header>
        <div className="w-full px-4 flex flex-row justify-between h-full items-center">
          {/* This needs to be replaced with the hackathon name */}
          <h1 className="text-2xl font-bold">SwampHacks</h1>

          <Link
            to="/portal"
            className="text-blue-500 underline underline-offset-4"
          >
            Back to portal
          </Link>
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
          label="My Team"
          href={`/events/${eventId}/dashboard/my-team`}
          leftSection={<TablerUsersGroup className="w-5 aspect-square" />}
          active={myTeamActive}
        />

        <NavLink
          label="Schedule"
          href={`/events/${eventId}/dashboard/schedule`}
          leftSection={<TablerCalendar className="w-5 aspect-square" />}
          active={scheduleActive}
        />
      </AppShell.Navbar>

      <AppShell.Main>{children}</AppShell.Main>
    </AppShell>
  );
}
