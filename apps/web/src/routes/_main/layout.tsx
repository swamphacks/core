import { NavLink } from "@/components/AppShell/NavLink";
import {
  createFileRoute,
  Outlet,
  redirect,
  useLocation,
  useRouter,
} from "@tanstack/react-router";
import TablerLayoutCollage from "~icons/tabler/layout-collage";
import TablerBooks from "~icons/tabler/books";
import TablerSocial from "~icons/tabler/social";
import { AppShell } from "@/components/AppShell/AppShell";
import { Link } from "react-aria-components";

// This layout component performs authentication checks before the user can access protected pages
export const Route = createFileRoute("/_main")({
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
  const { user } = Route.useRouteContext();
  const pathname = useLocation({ select: (loc) => loc.pathname });
  const router = useRouter();

  return (
    <AppShell>
      <AppShell.Header>
        <div className="w-full px-4 flex flex-row justify-between h-full items-center">
          <h1 className=" text-2xl font-bold">SwampHacks</h1>

          <div className="flex items-center h-full flex-row gap-6">
            {user?.role === "superuser" && (
              <Link
                onClick={() => router.navigate({ to: "/admin" })}
                className="text-sm text-blue-500 select-none cursor-pointer hover:underline"
              >
                To Admin Portal
              </Link>
            )}
            <img
              src="https://i.pinimg.com/736x/8b/d2/f6/8bd2f653f38322972e404925ab67294a.jpg"
              className="h-5/8 aspect-square rounded-full"
            />
          </div>
        </div>
      </AppShell.Header>

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
