import { useLocation } from "@tanstack/react-router";
import { AppShell } from "@/components/AppShell/AppShell";
import { NavLink } from "@/components/AppShell/NavLink";
import TablerLayoutCollage from "~icons/tabler/layout-collage";
import { type PropsWithChildren } from "react";

interface DashboardAppShellProps {
  eventId: string;
}

export default function ApplicantAppShell({
  eventId,
  children,
}: PropsWithChildren<DashboardAppShellProps>) {
  const pathname = useLocation({ select: (loc) => loc.pathname });

  const dashboardOverviewActive = /^\/events\/[^/]+\/dashboard\/?$/.test(
    pathname,
  );
  const applicationStatusActive =
    /^\/events\/[^/]+\/dashboard\/application-status\/?$/.test(pathname);

  return (
    <AppShell>
      <AppShell.Header>
        <div className="w-full px-4 flex flex-row justify-between h-full items-center">
          <h1 className="text-2xl font-bold">Applicant Dashboard</h1>
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
          label="Status"
          href={`/events/${eventId}/dashboard/application-status`}
          leftSection={<TablerLayoutCollage className="w-5 aspect-square" />}
          active={applicationStatusActive}
        />
      </AppShell.Navbar>

      <AppShell.Main>{children}</AppShell.Main>
    </AppShell>
  );
}
