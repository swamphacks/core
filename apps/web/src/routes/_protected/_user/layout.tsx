import { NavLink } from "@/components/AppShell/NavLink";
import {
  createFileRoute,
  Outlet,
  useLocation,
  useRouter,
} from "@tanstack/react-router";
import TablerCode from "~icons/tabler/code";
import TablerClipboard from "~icons/tabler/clipboard";
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
            onClick={() => router.navigate({ to: "/information" })}
            className="py-2 cursor-pointer"
            label="SwampHacks"
          />
        </div>
      </AppShell.Header>

      <AppShell.Navbar>
        <NavLink
          label="Information"
          href="/information"
          leftSection={<TablerCode className="w-5 aspect-square" />}
          active={pathname.startsWith("/information")}
        />

        <NavLink
          label="Application"
          href="/application"
          leftSection={<TablerClipboard className="w-5 aspect-square" />}
          active={pathname.startsWith("/application")}
        />
      </AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
