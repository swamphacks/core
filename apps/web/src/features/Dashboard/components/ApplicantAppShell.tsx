import { Link, useLocation } from "@tanstack/react-router";
import { AppShell } from "@/components/AppShell/AppShell";
import { NavLink } from "@/components/AppShell/NavLink";
import TablerProgress from "~icons/tabler/progress";
import TablerTransformPointBottomLeft from "~icons/tabler/transform-point-bottom-left";
import { type PropsWithChildren } from "react";

interface DashboardAppShellProps {
  eventId: string;
}

export default function ApplicantAppShell({
  eventId,
  children,
}: PropsWithChildren<DashboardAppShellProps>) {
  const pathname = useLocation({ select: (loc) => loc.pathname });

  const applicationStatusActive =
    /^\/events\/[^/]+\/dashboard\/application-status\/?$/.test(pathname);
  const teamFormationActive =
    /^\/events\/[^/]+\/dashboard\/team-formation\/?$/.test(pathname);

  return (
    <AppShell>
      <AppShell.Header>
        <div className="w-full px-4 flex flex-row justify-between h-full items-center">
          <h1 className="text-2xl font-bold">Applicant Dashboard</h1>
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
          label="Application Status"
          href={`/events/${eventId}/dashboard/application-status`}
          leftSection={<TablerProgress className="w-5 aspect-square" />}
          active={applicationStatusActive}
        />
        <NavLink
          label="My Team"
          href={`/events/${eventId}/dashboard/team-formation`}
          leftSection={
            <TablerTransformPointBottomLeft className="w-5 aspect-square" />
          }
          active={teamFormationActive}
        />
      </AppShell.Navbar>

      <AppShell.Main>{children}</AppShell.Main>
    </AppShell>
  );
}
