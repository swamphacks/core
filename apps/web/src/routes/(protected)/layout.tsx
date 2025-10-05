import { AppShell } from "@/components/AppShell/AppShell";
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
export const Route = createFileRoute("/(protected)")({
  beforeLoad: async ({ context, location }) => {
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
    <AppShell>
      <AppShell.Navbar>
        <NavLink
          label="Events Portal"
          href="/portal"
          leftSection={<TablerLayoutCollage className="w-5 aspect-square" />}
          active={pathname.startsWith("/portal")}
        />

        <NavLink
          label="Resources"
          href="/_protected/resources"
          initialExpanded={pathname.startsWith("/resources")}
          leftSection={<TablerBooks className="w-5 aspect-square" />}
        >
          <NavLink
            label="Programming"
            href="/resources/programming"
            active={pathname.startsWith("/resources/programming")}
          />
          <NavLink
            label="Sponsors"
            href="/resources/sponsors"
            active={pathname.startsWith("/resources/sponsors")}
          />
        </NavLink>

        <NavLink
          label="Community"
          href="/community"
          leftSection={<TablerSocial className="w-5 aspect-square" />}
          active={pathname.startsWith("/community")}
        />
      </AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
