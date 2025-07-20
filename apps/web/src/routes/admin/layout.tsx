import { AppShell } from "@/components/AppShell/AppShell";
import { NavLink } from "@/components/AppShell/NavLink";
import {
  createFileRoute,
  Outlet,
  redirect,
  useLocation,
  useRouter,
} from "@tanstack/react-router";
import TablerDashboard from "~icons/tabler/dashboard";
import TablerCalendarSearch from "~icons/tabler/calendar-search";
import TablerUsers from "~icons/tabler/users";
import TablerListDetails from "~icons/tabler/list-details";
import TablerSettings2 from "~icons/tabler/settings-2";
import { Link } from "react-aria-components";

export const Route = createFileRoute("/admin")({
  beforeLoad: async ({ context, location }) => {
    const { user, error } = await context.userQuery.promise;

    if (!user && !user) {
      console.log("You aren't authenticated, redirecting to login.");
      throw redirect({
        to: "/",
        search: { redirect: location.pathname },
      });
    }

    if (error) {
      // TODO: Display a friendly error to the user?
      console.error("Auth error in beforeLoad in layout.tsx:", error);
      console.log(
        "authentication error occurred while accessing protected page, redirecting to login.",
      );
      throw redirect({
        to: "/",
        search: { redirect: location.pathname },
      });
    }

    // If not superuser
    if (user.role !== "superuser") {
      console.log("You aren't a superuser, redirecting to event portal.");
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
  const router = useRouter();

  return (
    <AppShell>
      <AppShell.Header>
        <div className="w-full px-4 flex flex-row justify-between h-full items-center">
          <h1 className=" text-2xl font-bold">Admin Portal</h1>

          <div className="flex items-center h-full flex-row gap-6">
            <Link
              onClick={() => router.navigate({ to: "/portal" })}
              className="text-sm text-blue-500 select-none cursor-pointer hover:underline"
            >
              Back to User Portal
            </Link>
            <img
              src="https://i.pinimg.com/736x/8b/d2/f6/8bd2f653f38322972e404925ab67294a.jpg"
              className="h-5/8 aspect-square rounded-full"
            />
          </div>
        </div>
      </AppShell.Header>

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
