import { AppShell } from "@/components/AppShell/AppShell";
import { NavLink } from "@/components/AppShell/NavLink";
import {
  createFileRoute,
  Outlet,
  redirect,
  useLocation,
} from "@tanstack/react-router";
import TablerDashboard from "~icons/tabler/dashboard";
import TablerCalendarSearch from "~icons/tabler/calendar-search";
import TablerUsers from "~icons/tabler/users";
import TablerListDetails from "~icons/tabler/list-details";
import TablerSettings2 from "~icons/tabler/settings-2";

export const Route = createFileRoute("/_protected/admin")({
  beforeLoad: async ({ context, location }) => {
    const { user } = context;

    if (!user) {
      throw redirect({
        to: "/",
        search: { redirect: location.pathname },
      });
    }

    // If not superuser
    if (user.role !== "superuser") {
      throw redirect({
        to: "/portal",
      });
    }

    if (location.pathname === "/admin") {
      throw redirect({ to: "/admin/overview" });
    }
  },
  component: RouteComponent,
});

function RouteComponent() {
  const pathname = useLocation({ select: (loc) => loc.pathname });

  return (
    <AppShell>
      <AppShell.Navbar>
        <NavLink
          label="Overview"
          href="/admin/overview"
          leftSection={<TablerDashboard className="w-5 aspect-square" />}
          active={pathname.startsWith("/admin/overview")}
        />

        <NavLink
          label="Events Management"
          href="/admin/events-management"
          leftSection={<TablerCalendarSearch className="w-5 aspect-square" />}
          active={pathname.startsWith("/admin/events-management")}
        />

        <NavLink
          label="Users Management"
          href="/admin/users-management"
          leftSection={<TablerUsers className="w-5 aspect-square" />}
          active={pathname.startsWith("/admin/users-management")}
        />

        <NavLink
          label="Platform Logs"
          href="/admin/logs"
          leftSection={<TablerListDetails className="w-5 aspect-square" />}
          active={pathname.startsWith("/admin/logs")}
        />

        <NavLink
          label="Platform Settings"
          href="/admin/settings"
          leftSection={<TablerSettings2 className="w-5 aspect-square" />}
          active={pathname.startsWith("/admin/settings")}
        />
      </AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
