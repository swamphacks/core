import {
  createFileRoute,
  Outlet,
  useLocation,
  useRouter,
} from "@tanstack/react-router";
import { AppShell } from "@/components/AppShell/AppShell";
import { PageLoading } from "@/components/PageLoading";
import { Logo } from "@/components/Logo";
import { auth } from "@/lib/authClient";
import { useMyApplication } from "@/modules/Application/hooks/useMyApplication";
import ApplicantNavbar from "@/components/AppShell/ApplicantNavbar";
import VisitorNavbar from "@/components/AppShell/VisitorNavbar";

export const Route = createFileRoute("/_protected/_user")({
  pendingComponent: () => PageLoading(),
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const pathname = useLocation({ select: (loc) => loc.pathname });

  // TODO: pass this data down to outlet
  const { data } = auth.useUser();
  const { user } = data!;

  if (!user) {
    throw new Error("An error occurred: user is null");
  }

  const application = useMyApplication(user.role == "applicant");

  if (application.isLoading) {
    return;
  }

  const renderNavbarBasedOnRole = () => {
    switch (user.role) {
      case "visitor":
        return <VisitorNavbar pathname={pathname} />;
      case "applicant":
        return (
          <ApplicantNavbar
            pathname={pathname}
            hasSeenNewApplicationStatus={user.hasSeenNewApplicationStatus}
          />
        );
      case "admin":
        break;
      case "staff":
        break;
      default:
        break;
    }
  };

  return (
    <AppShell>
      <AppShell.Header>
        <div className="items-center gap-2 ml-3 flex">
          <Logo
            onClick={() => router.navigate({ to: "/information" })}
            className="py-2 cursor-pointer"
            label="SwampHacks XII"
          />
        </div>
      </AppShell.Header>

      <AppShell.Navbar>{renderNavbarBasedOnRole()}</AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
