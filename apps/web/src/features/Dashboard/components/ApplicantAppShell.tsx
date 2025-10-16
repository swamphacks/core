import { useLocation, useRouter } from "@tanstack/react-router";
import { AppShell } from "@/components/AppShell/AppShell";
import { NavLink } from "@/components/AppShell/NavLink";
import TablerProgress from "~icons/tabler/progress";
import TablerTransformPointBottomLeft from "~icons/tabler/transform-point-bottom-left";
import { type PropsWithChildren } from "react";
import { Logo } from "@/components/Logo";

interface DashboardAppShellProps {
  eventId: string;
  eventName?: string;
}

export default function ApplicantAppShell({
  eventId,
  children,
  eventName,
}: PropsWithChildren<DashboardAppShellProps>) {
  const router = useRouter();
  const pathname = useLocation({ select: (loc) => loc.pathname });

  const applicationStatusActive =
    /^\/events\/[^/]+\/dashboard\/application-status\/?$/.test(pathname);
  const teamFormationActive =
    /^\/events\/[^/]+\/dashboard\/team-formation\/?$/.test(pathname);

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
