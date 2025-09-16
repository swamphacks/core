import { SettingsPage } from "@/features/Settings/components/SettingsPage";
import { auth } from "@/lib/authClient";
import { createFileRoute, useRouter } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/settings")({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();

  const logout = async () => {
    await auth.logOut();
    await router.navigate({ to: "/", replace: true, reloadDocument: true }); // Force reload to reset auth state (Not ideal, but works for now)
  };

  return (
    <div className="flex items-center justify-center w-full h-full overflow-y-auto">
      <SettingsPage logout={logout} />
    </div>
  );
}
