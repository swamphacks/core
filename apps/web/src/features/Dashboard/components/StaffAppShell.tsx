import { useLocation, useRouteContext } from "@tanstack/react-router";
import { AppShell } from "@/components/AppShell/AppShell";
import { NavLink } from "@/components/AppShell/NavLink";
import TablerLayoutDashboard from "~icons/tabler/layout-dashboard";
import TablerUsersGroup from "~icons/tabler/users-group";
import TablerClipboardCheck from "~icons/tabler/clipboard-check";
import TablerMail from "~icons/tabler/mail";
import TablerSettings from "~icons/tabler/settings";
import { type PropsWithChildren } from "react";

interface DashboardAppShellProps {
  eventId: string;
}

export default function AttendeeAppShell({
  eventId,
  children,
}: PropsWithChildren<DashboardAppShellProps>) {
  const pathname = useLocation({ select: (loc) => loc.pathname });
  const { eventRole } = useRouteContext(pathname);
  const dashboardOverviewActive = /^\/events\/[^/]+\/dashboard\/?$/.test(
    pathname,
  );

  const applicationReviewActive =
    /^\/events\/[^/]+\/dashboard\/application-review\/?$/.test(pathname);

  const emailsActive = /^\/events\/[^/]+\/dashboard\/emails\/?$/.test(pathname);

  const staffOverviewActive =
    /^\/events\/[^/]+\/dashboard\/staff-overview\/?$/.test(pathname);

  const configsActive = /^\/events\/[^/]+\/dashboard\/configs\/?$/.test(
    pathname,
  );
  return (
    <AppShell>
      <AppShell.Header>
        <div className="w-full px-4 flex flex-row justify-between h-full items-center">
          <h1 className="text-2xl font-bold">Staff Dashboard</h1>
        </div>
      </AppShell.Header>

      <AppShell.Navbar>
        <NavLink
          label="Overview"
          href={`/events/${eventId}/dashboard`}
          leftSection={<TablerLayoutDashboard className="w-5 aspect-square" />}
          active={dashboardOverviewActive}
        />
        <NavLink
          label="Staff Overview"
          href={`/events/${eventId}/dashboard/staff-overview`}
          leftSection={<TablerUsersGroup className="w-5 aspect-square" />}
          active={staffOverviewActive}
        />
        <NavLink
          label="Review Apps"
          href={`/events/${eventId}/dashboard/application-review`}
          leftSection={<TablerClipboardCheck className="w-5 aspect-square" />}
          active={applicationReviewActive}
        />
        <NavLink
          label="Send Emails"
          href={`/events/${eventId}/dashboard/emails`}
          leftSection={<TablerMail className="w-5 aspect-square" />}
          active={emailsActive}
        />
        {eventRole == "admin" && (
          <NavLink
            label="Configs"
            href={`/events/${eventId}/dashboard/configs`}
            leftSection={<TablerSettings className="w-5 aspect-square" />}
            active={configsActive}
          />
        )}
      </AppShell.Navbar>

      <AppShell.Main>{children}</AppShell.Main>
    </AppShell>
  );
}
