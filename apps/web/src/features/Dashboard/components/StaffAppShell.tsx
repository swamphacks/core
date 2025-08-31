import { Link, useLocation } from "@tanstack/react-router";
import { AppShell } from "@/components/AppShell/AppShell";
import { NavLink } from "@/components/AppShell/NavLink";
import TablerLayoutDashboard from "~icons/tabler/layout-dashboard";
import TablerUsersGroup from "~icons/tabler/users-group";
import TablerClipboardCheck from "~icons/tabler/clipboard-check";
import TablerMail from "~icons/tabler/mail";
import TablerSettings from "~icons/tabler/settings";
import TablerFileText from "~icons/tabler/file-text";
import TablerAlarm from "~icons/tabler/alarm";
import TablerBodyScan from "~icons/tabler/body-scan";
import TablerUserSearch from "~icons/tabler/user-search";
import TablerFlower from "~icons/tabler/flower";
import TablerTicket from "~icons/tabler/ticket";
import TablerAdjustmentsHorizontal from "~icons/tabler/adjustments-horizontal";
import TablerShieldHalfFilled from "~icons/tabler/shield-half-filled";
import { type PropsWithChildren } from "react";

interface DashboardAppShellProps {
  eventId: string;
  eventRole: "admin" | "staff";
}

export default function StaffDashboardShell({
  eventId,
  eventRole,
  children,
}: PropsWithChildren<DashboardAppShellProps>) {
  const pathname = useLocation({ select: (loc) => loc.pathname });
  const dashboardOverviewActive = /^\/events\/[^/]+\/dashboard\/?$/.test(
    pathname,
  );

  const applicationReviewActive =
    /^\/events\/[^/]+\/dashboard\/application-review\/?$/.test(pathname);
  const waitlistActive = /^\/events\/[^/]+\/dashboard\/waitlist\/?$/.test(
    pathname,
  );
  const checkInActive = /^\/events\/[^/]+\/dashboard\/check-in\/?$/.test(
    pathname,
  );
  const attendeeDirectoryActive =
    /^\/events\/[^/]+\/dashboard\/attendee-directory\/?$/.test(pathname);
  const teamManagementActive =
    /^\/events\/[^/]+\/dashboard\/team-management\/?$/.test(pathname);
  const redeemablesActive = /^\/events\/[^/]+\/dashboard\/redeemables\/?$/.test(
    pathname,
  );
  const emailsActive = /^\/events\/[^/]+\/dashboard\/emails\/?$/.test(pathname);

  const eventSettingsActive =
    /^\/events\/[^/]+\/dashboard\/event-settings\/?$/.test(pathname);
  const staffManagementActive =
    /^\/events\/[^/]+\/dashboard\/staff-management\/?$/.test(pathname);

  return (
    <AppShell>
      <AppShell.Header>
        <div className="w-full px-4 flex flex-row justify-between h-full items-center">
          <h1 className="text-2xl font-bold">Staff Dashboard</h1>
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
          leftSection={<TablerLayoutDashboard className="w-5 aspect-square" />}
          active={dashboardOverviewActive}
        />
        <NavLink
          label="Applications"
          leftSection={<TablerFileText className="w-5 aspect-square" />}
          initialExpanded={waitlistActive || applicationReviewActive}
        >
          <NavLink
            label="Application Review"
            href={`/events/${eventId}/dashboard/application-review`}
            leftSection={<TablerClipboardCheck className="w-5 aspect-square" />}
            active={applicationReviewActive}
          />
          <NavLink
            label="Waitlist"
            href={`/events/${eventId}/dashboard/waitlist`}
            leftSection={<TablerAlarm className="w-5 aspect-square" />}
            active={waitlistActive}
          />
        </NavLink>

        <NavLink
          label="Attendees"
          leftSection={<TablerUsersGroup className="w-5 aspect-square" />}
          initialExpanded={
            checkInActive || teamManagementActive || attendeeDirectoryActive
          }
        >
          <NavLink
            label="Check In"
            href={`/events/${eventId}/dashboard/check-in`}
            leftSection={<TablerBodyScan className="w-5 aspect-square" />}
            active={checkInActive}
          />
          <NavLink
            label="Directory"
            href={`/events/${eventId}/dashboard/attendee-directory`}
            leftSection={<TablerUserSearch className="w-5 aspect-square" />}
            active={attendeeDirectoryActive}
          />
          <NavLink
            label="Teams"
            href={`/events/${eventId}/dashboard/team-management`}
            leftSection={<TablerFlower className="w-5 aspect-square" />}
            active={teamManagementActive}
          />
        </NavLink>

        <NavLink
          label="Redeemables"
          href={`/events/${eventId}/dashboard/redeemables`}
          leftSection={<TablerTicket className="w-5 aspect-square" />}
          active={redeemablesActive}
        />

        <NavLink
          label="Email Campaigns"
          href={`/events/${eventId}/dashboard/emails`}
          leftSection={<TablerMail className="w-5 aspect-square" />}
          active={emailsActive}
        />
        {eventRole == "admin" && (
          <NavLink
            label="Configuration"
            leftSection={<TablerSettings className="w-5 aspect-square" />}
            initialExpanded={eventSettingsActive || staffManagementActive}
          >
            <NavLink
              label="Event Settings"
              href={`/events/${eventId}/dashboard/event-settings`}
              leftSection={
                <TablerAdjustmentsHorizontal className="w-5 aspect-square" />
              }
              active={eventSettingsActive}
            />
            <NavLink
              label="Staff & Roles"
              href={`/events/${eventId}/dashboard/staff-management`}
              leftSection={
                <TablerShieldHalfFilled className="w-5 aspect-square" />
              }
              active={staffManagementActive}
            />
          </NavLink>
        )}
      </AppShell.Navbar>

      <AppShell.Main>{children}</AppShell.Main>
    </AppShell>
  );
}
