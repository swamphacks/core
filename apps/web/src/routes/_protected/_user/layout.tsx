import { NavLink } from "@/components/AppShell/NavLink";
import {
  createFileRoute,
  Outlet,
  useLocation,
  useRouter,
} from "@tanstack/react-router";
import TablerLayoutCollage from "~icons/tabler/layout-collage";
import TablerBooks from "~icons/tabler/books";
import TablerSocial from "~icons/tabler/social";
import { AppShell } from "@/components/AppShell/AppShell";
import { PageLoading } from "@/components/PageLoading";
import { Logo } from "@/components/Logo";

export const Route = createFileRoute("/_protected/_user")({
  pendingComponent: () => PageLoading(),
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const pathname = useLocation({ select: (loc) => loc.pathname });

  return (
    <AppShell>
      <AppShell.Header>
        <div className="items-center gap-2 ml-3 flex">
          <Logo
            onClick={() => router.navigate({ to: "/portal" })}
            className="py-2 cursor-pointer"
            label="SwampHacks"
          />
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
