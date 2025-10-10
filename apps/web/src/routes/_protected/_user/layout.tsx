import { NavLink } from "@/components/AppShell/NavLink";
import { createFileRoute, Outlet, useLocation } from "@tanstack/react-router";
import TablerLayoutCollage from "~icons/tabler/layout-collage";
import TablerBooks from "~icons/tabler/books";
import TablerSocial from "~icons/tabler/social";
import { AppShell } from "@/components/AppShell/AppShell";
import { PageLoading } from "@/components/PageLoading";

export const Route = createFileRoute("/_protected/_user")({
  pendingComponent: () => PageLoading(),
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
