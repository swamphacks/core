import {
  createFileRoute,
  Outlet,
  redirect,
  useLocation,
  useRouter,
} from "@tanstack/react-router";
import { PageLoading } from "@/components/PageLoading";
import ApplicantNavbar from "@/components/AppShell/ApplicantNavbar";
import { AppShell } from "@/components/AppShell/AppShell";
import AttendeeNavbar from "@/components/AppShell/AttendeeNavbar";
import VisitorNavbar from "@/components/AppShell/VisitorNavbar";
import { Logo } from "@/components/Logo";
import { useSuspenseQuery } from "@tanstack/react-query";
import { userQueryOptions } from "@/lib/auth/hooks/useUser";

// This layout component performs authentication checks before the user can access protected pages
export const Route = createFileRoute("/_protected")({
  beforeLoad: async ({ context, location }) => {
    const { user, error } = await context.userQuery.promise;

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

    // Unauthenticated, return to login page
    if (!user) {
      console.log("User is not authenticated, redirecting to login.");
      throw redirect({
        to: "/",
        search: { redirect: location.pathname },
      });
    }

    // Return user data for use in the route loader or component
    return { user };
  },
  pendingComponent: PageLoading,
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  const pathname = useLocation({ select: (loc) => loc.pathname });
  const { data } = useSuspenseQuery(userQueryOptions());
  const user = data.user!;

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
      case "attendee":
        return <AttendeeNavbar pathname={pathname} />;
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
