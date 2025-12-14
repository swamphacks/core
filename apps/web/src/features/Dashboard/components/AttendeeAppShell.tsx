import { useLocation, useRouter } from "@tanstack/react-router";
import { AppShell } from "@/components/AppShell/AppShell";
import { NavLink } from "@/components/AppShell/NavLink";
import TablerLayoutCollage from "~icons/tabler/layout-collage";
import TablerUsersGroup from "~icons/tabler/users-group";
import TablerCalendar from "~icons/tabler/calendar";
import { type PropsWithChildren } from "react";
import { Logo } from "@/components/Logo";

interface DashboardAppShellProps {
  eventId: string;
  eventName?: string;
}

export default function AttendeeAppShell({
  eventId,
  children,
  eventName,
}: PropsWithChildren<DashboardAppShellProps>) {
  const router = useRouter();

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
        <div className="items-center gap-2 ml-3 flex">
          <Logo
            onClick={() => router.navigate({ to: "/portal" })}
            className="py-2 cursor-pointer"
            label={eventName || "Event Portal"}
          />
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
