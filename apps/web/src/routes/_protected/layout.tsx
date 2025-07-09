import { NavLink } from "@/components/AppShell/NavLink";
import {
  createFileRoute,
  Outlet,
  redirect,
  useLocation,
} from "@tanstack/react-router";
import TablerLayoutCollage from "~icons/tabler/layout-collage";
import TablerBooks from "~icons/tabler/books";
import TablerSocial from "~icons/tabler/social";

// This layout component performs authentication checks before the user can access protected pages
export const Route = createFileRoute("/_protected")({
  beforeLoad: async ({ context }) => {
    const { user, error } = await context.userQuery.promise;

    // Unauthenticated, return to login page
    if (!user && !error) {
      console.log("User is not authenticated, redirecting to login.");
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
      });
    }

    // Return user data for use in the route loader or component
    return { user };
  },
  pendingMs: 1000,
  pendingComponent: () => <p>Loading...</p>,
  component: RouteComponent,
});

function RouteComponent() {
  const pathname = useLocation({ select: (loc) => loc.pathname });

  return (
    <div className="flex h-screen w-screen flex-col">
      {/* Top bar */}
      <header className="h-14 w-full bg-gray-800 text-white flex items-center px-4 shadow">
        <h1 className="text-lg font-semibold">SwampHacks</h1>
      </header>

      <div className="flex flex-1">
        {/* Sidebar */}
        <aside className="w-64 p-4">
          {/* Replace with your nav links */}
          <nav className="flex flex-col gap-2">
            <NavLink
              label="Dashboard"
              href="/_protected/dashboard"
              leftSection={
                <TablerLayoutCollage className="w-4 aspect-square" />
              }
              active={pathname.startsWith("/dashboard")}
            />

            <NavLink
              label="Resources"
              href="/_protected/resources"
              leftSection={<TablerBooks className="w-4 aspect-square" />}
              active={pathname.startsWith("/resources")}
            >
              <NavLink
                label="Programming"
                href="/_protected/resources/programming"
                active={pathname.startsWith("/resources/programming")}
              />
              <NavLink
                label="Sponsors"
                href="/_protected/resources/sponsors"
                active={pathname.startsWith("/resources/sponsors")}
              />
            </NavLink>

            <NavLink
              label="Community"
              href="/_protected/community"
              leftSection={<TablerSocial className="w-4 aspect-square" />}
              active={pathname.startsWith("/community")}
            />
          </nav>
        </aside>

        {/* Main content */}
        <main className="flex-1 p-6 overflow-y-auto">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
